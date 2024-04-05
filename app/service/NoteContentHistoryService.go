package service

import (
	"github.com/wiselike/leanote-of-unofficial/app/db"
	"github.com/wiselike/leanote-of-unofficial/app/info"
	//	. "github.com/wiselike/leanote-of-unofficial/app/lea"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"time"
)

// 历史记录
type NoteContentHistoryService struct {
}

var UseMongoVer int

func isMongo2() bool {
	if db.Session != nil {
		if info, err := db.Session.BuildInfo(); err == nil {
			return !info.VersionAtLeast(3)
		}
	}
	return true
}

// 新建一个note, 不添加历史记录
// 添加历史，在数据库中倒序存放：前面的是老的，后面是新的（与原来的顺序相反）
func (this *NoteContentHistoryService) AddHistory(noteId, userId string, oneHistory info.EachHistory) {
	// 检查是否是空
	if oneHistory.Content == "" {
		return
	}

	// 每个历史记录最大值
	maxSize := ConfigS.GlobalAllConfigs["note.history.size"].(int)
	if maxSize < 1 {
		return
	}

	// 判断使用的mongo版本
	if UseMongoVer == 0 {
		if isMongo2() {
			UseMongoVer = 2
		} else {
			UseMongoVer = 3
		}
	}

	var historiesLenth int
	if UseMongoVer == 2 {
		// 使用mongodb2的版本，效率较低，而且耗内存
		history := info.NoteContentHistory{}
		db.GetByIdAndUserId(db.NoteContentHistories, noteId, userId, &history)
		if history.NoteId == "" {
			historiesLenth = -1
		} else {
			historiesLenth = len(history.Histories)
		}
	} else {
		// mongodb3才支持，可以优化速度和效率
		// 只获取数字即可，不获取所有历史的正文内容
		historiesLenth = db.GetNoteHistoriesCount(db.NoteContentHistories, noteId, userId)
	}

	if historiesLenth == -1 {
		this.newHistory(noteId, userId, oneHistory)
	} else {
		// 读取最新的历史记录，判断是否是AutoBackup；
		var lastContentHistory info.NoteContentHistory
		db.GetLastOneInArray(db.NoteContentHistories, noteId, userId, "Histories", &lastContentHistory)
		if len(lastContentHistory.Histories) > 0 && lastContentHistory.Histories[0].IsAutoBackup {
			db.UpdateByIdAndUserIdPop(db.NoteContentHistories, noteId, userId, "Histories", 1)
			historiesLenth--
		}

		// 判断是否超出 maxSize, 如果是则pop掉一个最老的
		if historiesLenth >= maxSize {
			db.UpdateByIdAndUserIdPop(db.NoteContentHistories, noteId, userId, "Histories", -1)
		}

		// 插入一个历史记录，只能后插
		db.UpdateByIdAndUserIdPush(db.NoteContentHistories, noteId, userId, "Histories", oneHistory)
	}

	return
}

// 更新一下最后一条历史记录的状态，由自动历史转为手动历史
func (this *NoteContentHistoryService) UpdateHistoryBackupState(noteId, userId string, isAutoBackup bool) {
	// mongo2没法找到最后数组的最后一个，
	// 所以这里进行了折中，找到第一个IsAutoBackup为true的项
	// 将其替换为isAutoBackup值
	db.UpdateHistoryBackupState(db.NoteContentHistories, noteId, userId, isAutoBackup)
}

// 新建历史
func (this *NoteContentHistoryService) newHistory(noteId, userId string, oneHistory info.EachHistory) {
	history := info.NoteContentHistory{NoteId: bson.ObjectIdHex(noteId),
		UserId:    bson.ObjectIdHex(userId),
		Histories: []info.EachHistory{oneHistory},
	}

	// 保存之
	db.Insert(db.NoteContentHistories, history)
}

// 列表展示
func (this *NoteContentHistoryService) ListHistories(noteId, userId string) []info.EachHistory {
	histories := info.NoteContentHistory{}
	db.GetByIdAndUserId(db.NoteContentHistories, noteId, userId, &histories)
	sort.Sort(info.EachHistorySlice(histories.Histories)) // 前端倒着展示，便于理解和操作
	return histories.Histories
}

// 删除一条历史；
// 使用历史记录的时间戳，作为标志进行查找并删除；
// 实际过程中应该不存在两条时间戳完全相同(时间戳是精确到毫秒级的)历史记录；
// 如果确实存在两条时间戳毫秒级也相同的，则内容肯定也相同，会一起都删除，目前还没遇到此情况
func (this *NoteContentHistoryService) DeleteHistory(noteId, userId, timeToDel string) {
	// 自动解析js返回的RFC 3339格式化时间戳。
	// golang可以自动解析末尾的Z或者时区偏移
	t, err := time.Parse(time.RFC3339, timeToDel)
	if err != nil {
		return
	}
	db.DeleteOneHistory(db.NoteContentHistories, noteId, userId, t)
	return
}
