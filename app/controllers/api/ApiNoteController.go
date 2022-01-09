package api

import (
	"github.com/revel/revel"
	//	"encoding/json"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"os"
	"os/exec"
	// "strings"
	"time"
	"regexp"
	//	"github.com/leanote/leanote/app/types"
	//	"io/ioutil"
	//	"fmt"
	//	"bytes"
	//	"os"
)

// 笔记API

type ApiNote struct {
	ApiBaseContrller
}

// 获取同步的笔记
// > afterUsn的笔记
// 无Desc, Abstract, 有Files
/*
  {
    "NoteId": "55195fa199c37b79be000005",
    "NotebookId": "55195fa199c37b79be000002",
    "UserId": "55195fa199c37b79be000001",
    "Title": "Leanote语法Leanote语法Leanote语法Leanote语法Leanote语法",
    "Desc": "",
    "Tags": null,
    "Abstract": "",
    "Content": "",
    "IsMarkdown": true,
    "IsBlog": false,
    "IsTrash": false,
    "Usn": 5,
    "Files": [],
    "CreatedTime": "2015-03-30T22:37:21.695+08:00",
    "UpdatedTime": "2015-03-30T22:37:21.724+08:00",
    "PublicTime": "2015-03-30T22:37:21.695+08:00"
  }
*/
func (c ApiNote) GetSyncNotes(afterUsn, maxEntry int) revel.Result {
	if maxEntry == 0 {
		maxEntry = 100
	}
	notes := noteService.GetSyncNotes(c.getUserId(), afterUsn, maxEntry)
	return c.RenderJSON(notes)
}

// 得到笔记本下的笔记
// [OK]
func (c ApiNote) GetNotes(notebookId string) revel.Result {
	if notebookId != "" && !bson.IsObjectIdHex(notebookId) {
		re := info.NewApiRe()
		re.Msg = "notebookIdInvalid"
		return c.RenderJSON(re)
	}
	_, notes := noteService.ListNotes(c.getUserId(), notebookId, false, c.GetPage(), pageSize, defaultSortField, false, false)
	return c.RenderJSON(noteService.ToApiNotes(notes))
}

// 得到trash
// [OK]
func (c ApiNote) GetTrashNotes() revel.Result {
	_, notes := noteService.ListNotes(c.getUserId(), "", true, c.GetPage(), pageSize, defaultSortField, false, false)
	return c.RenderJSON(noteService.ToApiNotes(notes))

}

// get Note
// [OK]
/*
{
  "NoteId": "550c0bee2ec82a2eb5000000",
  "NotebookId": "54a1676399c37b1c77000004",
  "UserId": "54a1676399c37b1c77000002",
  "Title": "asdfadsf--=",
  "Desc": "",
  "Tags": [
  ],
  "Abstract": "",
  "Content": "",
  "IsMarkdown": false,
  "IsBlog": false,
  "IsTrash": false,
  "Usn": 8,
  "Files": [
    {
      "FileId": "551975d599c37b970f000000",
      "LocalFileId": "",
      "Type": "",
      "Title": "",
      "HasBody": false,
      "IsAttach": false
    },
    {
      "FileId": "551975de99c37b970f000001",
      "LocalFileId": "",
      "Type": "doc",
      "Title": "李铁-print-en.doc",
      "HasBody": false,
      "IsAttach": true
    },
    {
      "FileId": "551975de99c37b970f000002",
      "LocalFileId": "",
      "Type": "doc",
      "Title": "李铁-print.doc",
      "HasBody": false,
      "IsAttach": true
    }
  ],
  "CreatedTime": "2015-03-20T20:00:52.463+08:00",
  "UpdatedTime": "2015-03-31T00:12:44.967+08:00",
  "PublicTime": "2015-03-20T20:00:52.463+08:00"
}
*/
func (c ApiNote) GetNote(noteId string) revel.Result {
	if !bson.IsObjectIdHex(noteId) {
		re := info.NewApiRe()
		re.Msg = "noteIdInvalid"
		return c.RenderJSON(re)
	}

	note := noteService.GetNote(noteId, c.getUserId())
	if note.NoteId == "" {
		re := info.NewApiRe()
		re.Msg = "notExists"
		return c.RenderJSON(re)
	}
	apiNotes := noteService.ToApiNotes([]info.Note{note})
	return c.RenderJSON(apiNotes[0])
}

// 得到note和内容
// [OK]
func (c ApiNote) GetNoteAndContent(noteId string) revel.Result {
	noteAndContent := noteService.GetNoteAndContent(noteId, c.getUserId())

	apiNotes := noteService.ToApiNotes([]info.Note{noteAndContent.Note})
	apiNote := apiNotes[0]
	apiNote.Content = noteService.FixContent(noteAndContent.Content, noteAndContent.IsMarkdown)
	return c.RenderJSON(apiNote)
}

// content里的image, attach链接是
// https://leanote.com/api/file/getImage?fileId=xx
// https://leanote.com/api/file/getAttach?fileId=xx
// 将fileId=映射成ServerFileId, 这里的fileId可能是本地的FileId
func (c ApiNote) fixPostNotecontent(noteOrContent *info.ApiNote) {
	if noteOrContent.Content == "" {
		return
	}

	files := noteOrContent.Files
	if files != nil && len(files) > 0 {
		for _, file := range files {
			if file.LocalFileId != "" {
				LogJ(file)
				if !file.IsAttach {
					// <img src="https://"
					// ![](http://demo.leanote.top/api/file/getImage?fileId=5863219465b68e4fd5000001)
					reg, _ := regexp.Compile(`https*://[^/]*?/api/file/getImage\?fileId=`+file.LocalFileId)
					// Log(reg)
					noteOrContent.Content = reg.ReplaceAllString(noteOrContent.Content, `/api/file/getImage?fileId=`+file.FileId)  

					// // "http://a.com/api/file/getImage?fileId=localId" => /api/file/getImage?fileId=serverId
					// noteOrContent.Content = strings.Replace(noteOrContent.Content, 
					// 	baseUrl + "/api/file/getImage?fileId="+file.LocalFileId, 
					// 	"/api/file/getImage?fileId="+file.FileId, -1)
				} else {
					reg, _ := regexp.Compile(`https*://[^/]*?/api/file/getAttach\?fileId=`+file.LocalFileId)
					// Log(reg)
					noteOrContent.Content = reg.ReplaceAllString(noteOrContent.Content, `/api/file/getAttach?fileId=`+file.FileId)  
					/*
					noteOrContent.Content = strings.Replace(noteOrContent.Content, 
						baseUrl + "/api/file/getAttach?fileId="+file.LocalFileId, 
						"/api/file/getAttach?fileId="+file.FileId, -1)
					*/
				}
			}
		}
	}
}

// 得到内容
func (c ApiNote) GetNoteContent(noteId string) revel.Result {
	userId := c.getUserId()
	note := noteService.GetNote(noteId, userId)
	//	re := info.NewRe()
	noteContent := noteService.GetNoteContent(noteId, userId)
	if noteContent.Content != "" {
		noteContent.Content = noteService.FixContent(noteContent.Content, note.IsMarkdown)
	}

	apiNoteContent := info.ApiNoteContent{
		NoteId:  noteContent.NoteId,
		UserId:  noteContent.UserId,
		Content: noteContent.Content,
	}

	return c.RenderJSON(apiNoteContent)
}

// 添加笔记
// [OK]
func (c ApiNote) AddNote(noteOrContent info.ApiNote) revel.Result {
	userId := bson.ObjectIdHex(c.getUserId())
	re := info.NewRe()
	myUserId := userId
	// 为共享新建?
	/*
		if noteOrContent.FromUserId != "" {
			userId = bson.ObjectIdHex(noteOrContent.FromUserId)
		}
	*/
	//	Log(noteOrContent.Title)
	//		LogJ(noteOrContent)

	/*
		LogJ(c.Params)
		for name, _ := range c.Params.Files {
			Log(name)
			file, _, _ := c.Request.FormFile(name)
			LogJ(file)
		}
	*/
	//	return c.RenderJSON(re)
	if noteOrContent.NotebookId == "" || !bson.IsObjectIdHex(noteOrContent.NotebookId) {
		re.Msg = "notebookIdNotExists"
		return c.RenderJSON(re)
	}

	noteId := bson.NewObjectId()
	// TODO 先上传图片/附件, 如果不成功, 则返回false
	//
	attachNum := 0
	if noteOrContent.Files != nil && len(noteOrContent.Files) > 0 {
		for i, file := range noteOrContent.Files {
			if file.HasBody {
				if file.LocalFileId != "" {
					// FileDatas[54c7ae27d98d0329dd000000]
					ok, msg, fileId := c.upload("FileDatas["+file.LocalFileId+"]", noteId.Hex(), file.IsAttach)

					if !ok {
						re.Ok = false
						if msg != "" {
							Log(msg)
							Log(file.LocalFileId)
							re.Msg = "fileUploadError"
						}
						// 报不是图片的错误没关系, 证明客户端传来非图片的数据
						if msg != "notImage" {
							return c.RenderJSON(re)
						}
					} else {
						// 建立映射
						file.FileId = fileId
						noteOrContent.Files[i] = file

						if file.IsAttach {
							attachNum++
						}
					}
				} else {
					return c.RenderJSON(re)
				}
			}
		}
	}

	c.fixPostNotecontent(&noteOrContent)

	//	Log("Add")
	//	LogJ(noteOrContent)

	//	return c.RenderJSON(re)

	note := info.Note{UserId: userId,
		NoteId:     noteId,
		NotebookId: bson.ObjectIdHex(noteOrContent.NotebookId),
		Title:      noteOrContent.Title,
		Tags:       noteOrContent.Tags,
		Desc:       noteOrContent.Desc,
		//		ImgSrc:     noteOrContent.ImgSrc,
		IsBlog:      noteOrContent.IsBlog,
		IsMarkdown:  noteOrContent.IsMarkdown,
		AttachNum:   attachNum,
		CreatedTime: noteOrContent.CreatedTime,
		UpdatedTime: noteOrContent.UpdatedTime,
	}
	noteContent := info.NoteContent{NoteId: note.NoteId,
		UserId:      userId,
		IsBlog:      note.IsBlog,
		Content:     noteOrContent.Content,
		Abstract:    noteOrContent.Abstract,
		CreatedTime: noteOrContent.CreatedTime,
		UpdatedTime: noteOrContent.UpdatedTime,
	}

	// 通过内容得到Desc, abstract
	if noteOrContent.Abstract == "" {
		note.Desc = SubStringHTMLToRaw(noteContent.Content, 200)
		noteContent.Abstract = SubStringHTML(noteContent.Content, 200, "")
	} else {
		note.Desc = SubStringHTMLToRaw(noteContent.Abstract, 200)
	}

	note = noteService.AddNoteAndContentApi(note, noteContent, myUserId)

	if note.NoteId == "" {
		re.Ok = false
		return c.RenderJSON(re)
	}

	// 添加需要返回的
	noteOrContent.NoteId = note.NoteId.Hex()
	noteOrContent.Usn = note.Usn
	noteOrContent.CreatedTime = note.CreatedTime
	noteOrContent.UpdatedTime = note.UpdatedTime
	noteOrContent.UserId = c.getUserId()
	noteOrContent.IsMarkdown = note.IsMarkdown
	// 删除一些不要返回的, 删除Desc?
	noteOrContent.Content = ""
	noteOrContent.Abstract = ""
	//	apiNote := info.NoteToApiNote(note, noteOrContent.Files)
	return c.RenderJSON(noteOrContent)
}

// 更新笔记
// [OK]
func (c ApiNote) UpdateNote(noteOrContent info.ApiNote) revel.Result {
	re := info.NewReUpdate()

	noteUpdate := bson.M{}
	needUpdateNote := false

	noteId := noteOrContent.NoteId

	if noteOrContent.NoteId == "" {
		re.Msg = "noteIdNotExists"
		return c.RenderJSON(re)
	}

	if noteOrContent.Usn <= 0 {
		re.Msg = "usnNotExists"
		return c.RenderJSON(re)
	}

	//	Log("_____________")
	//	LogJ(noteOrContent)
	/*
		LogJ(c.Params.Files)
		LogJ(c.Request.Header)
		LogJ(c.Params.Values)
	*/

	// 先判断USN的问题, 因为很可能添加完附件后, 会有USN冲突, 这时附件就添错了
	userId := c.getUserId()
	note := noteService.GetNote(noteId, userId)
	if note.NoteId == "" {
		re.Msg = "notExists"
		return c.RenderJSON(re)
	}
	if note.Usn != noteOrContent.Usn {
		re.Msg = "conflict"
		Log("conflict")
		return c.RenderJSON(re)
	}
	Log("没有冲突")

	// 如果传了files
	// TODO 测试
	/*
		for key, v := range c.Params.Values {
			Log(key)
			Log(v)
		}
	*/
	//	Log(c.Has("Files[0]"))
	if c.Has("Files[0][LocalFileId]") {
		//		LogJ(c.Params.Files)
		if noteOrContent.Files != nil && len(noteOrContent.Files) > 0 {
			for i, file := range noteOrContent.Files {
				if file.HasBody {
					if file.LocalFileId != "" {
						// FileDatas[54c7ae27d98d0329dd000000]
						ok, msg, fileId := c.upload("FileDatas["+file.LocalFileId+"]", noteId, file.IsAttach)
						if !ok {
							Log("upload file error")
							re.Ok = false
							if msg == "" {
								re.Msg = "fileUploadError"
							} else {
								re.Msg = msg
							}
							return c.RenderJSON(re)
						} else {
							// 建立映射
							file.FileId = fileId
							noteOrContent.Files[i] = file
						}
					} else {
						return c.RenderJSON(re)
					}
				}
			}
		}

		//		Log("after upload")
		//		LogJ(noteOrContent.Files)

	}

	// 移到外面来, 删除最后一个file时也要处理, 不然总删不掉
	// 附件问题, 根据Files, 有些要删除的, 只留下这些
	attachService.UpdateOrDeleteAttachApi(noteId, userId, noteOrContent.Files)

	// Desc前台传来
	if c.Has("Desc") {
		needUpdateNote = true
		noteUpdate["Desc"] = noteOrContent.Desc
	}
	/*
		if c.Has("ImgSrc") {
			needUpdateNote = true
			noteUpdate["ImgSrc"] = noteOrContent.ImgSrc
		}
	*/
	if c.Has("Title") {
		needUpdateNote = true
		noteUpdate["Title"] = noteOrContent.Title
	}
	if c.Has("IsTrash") {
		needUpdateNote = true
		noteUpdate["IsTrash"] = noteOrContent.IsTrash
	}

	// 是否是博客
	if c.Has("IsBlog") {
		needUpdateNote = true
		noteUpdate["IsBlog"] = noteOrContent.IsBlog
	}

	/*
		Log(c.Has("tags[0]"))
		Log(c.Has("Tags[]"))
		for key, v := range c.Params.Values {
			Log(key)
			Log(v)
		}
	*/

	if c.Has("Tags[0]") {
		needUpdateNote = true
		noteUpdate["Tags"] = noteOrContent.Tags
	}

	if c.Has("NotebookId") {
		if bson.IsObjectIdHex(noteOrContent.NotebookId) {
			needUpdateNote = true
			noteUpdate["NotebookId"] = bson.ObjectIdHex(noteOrContent.NotebookId)
		}
	}

	if c.Has("Content") {
		// 通过内容得到Desc, 如果有Abstract, 则用Abstract生成Desc
		if noteOrContent.Abstract == "" {
			noteUpdate["Desc"] = SubStringHTMLToRaw(noteOrContent.Content, 200)
		} else {
			noteUpdate["Desc"] = SubStringHTMLToRaw(noteOrContent.Abstract, 200)
		}
	}

	noteUpdate["UpdatedTime"] = noteOrContent.UpdatedTime

	afterNoteUsn := 0
	noteOk := false
	noteMsg := ""
	if needUpdateNote {
		noteOk, noteMsg, afterNoteUsn = noteService.UpdateNote(c.getUserId(), noteOrContent.NoteId, noteUpdate, noteOrContent.Usn)
		if !noteOk {
			re.Ok = false
			re.Msg = noteMsg
			return c.RenderJSON(re)
		}
	}

	//-------------
	afterContentUsn := 0
	contentOk := false
	contentMsg := ""
	if c.Has("Content") {
		// 把fileId替换下
		c.fixPostNotecontent(&noteOrContent)
		// 如果传了Abstract就用之
		if noteOrContent.Abstract == "" {
			noteOrContent.Abstract = SubStringHTML(noteOrContent.Content, 200, "")
		}

		//		Log("--------> afte fixed")
		//		Log(noteOrContent.Content)
		contentOk, contentMsg, afterContentUsn = noteService.UpdateNoteContent(c.getUserId(),
			noteOrContent.NoteId,
			noteOrContent.Content,
			noteOrContent.Abstract,
			needUpdateNote,
			noteOrContent.Usn,
			noteOrContent.UpdatedTime)
	}

	if needUpdateNote {
		re.Ok = noteOk
		re.Msg = noteMsg
		re.Usn = afterNoteUsn
	} else {
		re.Ok = contentOk
		re.Msg = contentMsg
		re.Usn = afterContentUsn
	}

	if !re.Ok {
		return c.RenderJSON(re)
	}

	noteOrContent.Content = ""
	noteOrContent.Usn = re.Usn
	noteOrContent.UpdatedTime = time.Now()

	//	Log("after upload")
	//	LogJ(noteOrContent.Files)
	noteOrContent.UserId = c.getUserId()

	return c.RenderJSON(noteOrContent)
}

// 删除trash
func (c ApiNote) DeleteTrash(noteId string, usn int) revel.Result {
	re := info.NewReUpdate()
	re.Ok, re.Msg, re.Usn = trashService.DeleteTrashApi(noteId, c.getUserId(), usn)
	return c.RenderJSON(re)
}

// 得到历史列表
/*
func (c ApiNote) GetHistories(noteId string) revel.Result {
	re := info.NewRe()
	histories := noteContentHistoryService.ListHistories(noteId, c.getUserId())
	if len(histories) > 0 {
		re.Ok = true
		re.Item = histories
	}
	return c.RenderJSON(re)
}
*/

// 0.2 新增
// 导出成PDF
func (c ApiNote) ExportPdf(noteId string) revel.Result {
	re := info.NewApiRe()
	userId := c.getUserId()
	if noteId == "" {
		re.Msg = "noteNotExists"
		return c.RenderJSON(re)
	}

	note := noteService.GetNoteById(noteId)
	if note.NoteId == "" {
		re.Msg = "noteNotExists"
		return c.RenderJSON(re)
	}

	noteUserId := note.UserId.Hex()
	// 是否有权限
	if noteUserId != userId {
		// 是否是有权限协作的
		if !note.IsBlog && !shareService.HasReadPerm(noteUserId, userId, noteId) {
			re.Msg = "noteNotExists"
			return c.RenderJSON(re)
		}
	}

	// path 判断是否需要重新生成之
	guid := NewGuid()
	fileUrlPath := "files/export_pdf"
	dir := revel.BasePath + "/" + fileUrlPath
	if !MkdirAll(dir) {
		re.Msg = "noDir"
		return c.RenderJSON(re)
	}
	filename := guid + ".pdf"
	path := dir + "/" + filename

	appKey, _ := revel.Config.String("app.secretLeanote")
	if appKey == "" {
		appKey, _ = revel.Config.String("app.secret")
	}

	// 生成之
	binPath := configService.GetGlobalStringConfig("exportPdfBinPath")
	// 默认路径
	if binPath == "" {
		binPath = "/usr/local/bin/wkhtmltopdf"
	}

	url := configService.GetSiteUrl() + "/note/toPdf?noteId=" + noteId + "&appKey=" + appKey
	var cc string
	if note.IsMarkdown {
		cc = binPath + " --lowquality --window-status done \"" + url + "\"  \"" + path + "\"" //  \"" + cookieDomain + "\" \"" + cookieName + "\" \"" + cookieValue + "\""
	} else {
		cc = binPath + " --lowquality \"" + url + "\"  \"" + path + "\"" //  \"" + cookieDomain + "\" \"" + cookieName + "\" \"" + cookieValue + "\""
	}

	cmd := exec.Command("/bin/sh", "-c", cc)
	_, err := cmd.Output()
	if err != nil {
		re.Msg = "sysError"
		return c.RenderJSON(re)
	}
	file, err := os.Open(path)
	if err != nil {
		re.Msg = "sysError"
		return c.RenderJSON(re)
	}

	filenameReturn := note.Title
	filenameReturn = FixFilename(filenameReturn)
	if filenameReturn == "" {
		filenameReturn = "Untitled.pdf"
	} else {
		filenameReturn += ".pdf"
	}
	return c.RenderBinary(file, filenameReturn, revel.Attachment, time.Now()) // revel.Attachment
}
