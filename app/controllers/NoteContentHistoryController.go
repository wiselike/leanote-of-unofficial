package controllers

import (
	"github.com/wiselike/revel"
)

type NoteContentHistory struct {
	BaseController
}

// 得到list
func (c NoteContentHistory) ListHistories(noteId string) revel.Result {
	histories := noteContentHistoryService.ListHistories(noteId, c.GetUserId())

	return c.RenderJSON(histories)
}

func (c NoteContentHistory) DeleteHistory(noteId, userId, timeToDel string) revel.Result {
	// 因为是删除记录，所以额外添加防攻击检测，
	// 任意一个字符串为空立即返回失败
	if len(noteId)*len(userId)*len(timeToDel) == 0 {
		return c.RenderJSON(false)
	}

	noteContentHistoryService.DeleteHistory(noteId, userId, timeToDel)
	return c.RenderJSON(true)
}
