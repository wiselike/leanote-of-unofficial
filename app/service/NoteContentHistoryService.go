package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	//	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"sort"
)

// 历史记录
type NoteContentHistoryService struct {
}

// 新建一个note, 不需要添加历史记录
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

	//注释掉下面这个块，使用mongodb3的块，可以优化速度和效率
	history := info.NoteContentHistory{}
	db.GetByIdAndUserId(db.NoteContentHistories, noteId, userId, &history) // TODO 优化掉, 只获取数字即可
	var historiesLenth int
	if history.NoteId == "" {
		historiesLenth = -1
	} else {
		historiesLenth = len(history.Histories)
	}

	/* mongodb3才支持
	// 优化为只获取数字，不获取所有历史的正文
	historiesLenth := db.GetNoteHistoriesCount(db.NoteContentHistories, noteId, userId)
	*/

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
	sort.Sort(info.EachHistorySlice(histories.Histories)) // TODO 前端倒着展示，就不用排序了
	return histories.Histories
}
