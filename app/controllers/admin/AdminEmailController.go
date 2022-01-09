package admin

import (
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"strconv"
	"strings"
)

// admin 首页

type AdminEmail struct {
	AdminBaseController
}

// email配置
func (c AdminEmail) Email() revel.Result {
	return nil
}

// blog标签设置
func (c AdminEmail) Blog() revel.Result {
	recommendTags := configService.GetGlobalArrayConfig("recommendTags")
	newTags := configService.GetGlobalArrayConfig("newTags")
	c.ViewArgs["recommendTags"] = strings.Join(recommendTags, ",")
	c.ViewArgs["newTags"] = strings.Join(newTags, ",")
	return c.RenderTemplate("admin/setting/blog.html")
}
func (c AdminEmail) DoBlogTag(recommendTags, newTags string) revel.Result {
	re := info.NewRe()

	re.Ok = configService.UpdateGlobalArrayConfig(c.GetUserId(), "recommendTags", strings.Split(recommendTags, ","))
	re.Ok = configService.UpdateGlobalArrayConfig(c.GetUserId(), "newTags", strings.Split(newTags, ","))

	return c.RenderJSON(re)
}

// demo
// blog标签设置
func (c AdminEmail) Demo() revel.Result {
	c.ViewArgs["demoUsername"] = configService.GetGlobalStringConfig("demoUsername")
	c.ViewArgs["demoPassword"] = configService.GetGlobalStringConfig("demoPassword")
	return c.RenderTemplate("admin/setting/demo.html")
}
func (c AdminEmail) DoDemo(demoUsername, demoPassword string) revel.Result {
	re := info.NewRe()

	userInfo, err := authService.Login(demoUsername, demoPassword)
	if err != nil {
		return c.RenderJSON(info.Re{Ok: false})
	}
	if userInfo.UserId == "" {
		re.Msg = "The User is Not Exists"
		return c.RenderJSON(re)
	}

	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "demoUserId", userInfo.UserId.Hex())
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "demoUsername", demoUsername)
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "demoPassword", demoPassword)

	return c.RenderJSON(re)
}

// ToImage
// 长微博的bin路径phantomJs
func (c AdminEmail) ToImage() revel.Result {
	c.ViewArgs["toImageBinPath"] = configService.GetGlobalStringConfig("toImageBinPath")
	return c.RenderTemplate("admin/setting/toImage.html")
}
func (c AdminEmail) DoToImage(toImageBinPath string) revel.Result {
	re := info.NewRe()
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "toImageBinPath", toImageBinPath)
	return c.RenderJSON(re)
}

func (c AdminEmail) Set(emailHost, emailPort, emailUsername, emailPassword, emailSSL string) revel.Result {
	re := info.NewRe()
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "emailHost", emailHost)
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "emailPort", emailPort)
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "emailUsername", emailUsername)
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "emailPassword", emailPassword)
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "emailSSL", emailSSL)

	return c.RenderJSON(re)
}
func (c AdminEmail) Template() revel.Result {
	re := info.NewRe()

	keys := []string{"emailTemplateHeader", "emailTemplateFooter",
		"emailTemplateRegisterSubject",
		"emailTemplateRegister",
		"emailTemplateFindPasswordSubject",
		"emailTemplateFindPassword",
		"emailTemplateUpdateEmailSubject",
		"emailTemplateUpdateEmail",
		"emailTemplateInviteSubject",
		"emailTemplateInvite",
		"emailTemplateCommentSubject",
		"emailTemplateComment",
	}

	userId := c.GetUserId()
	for _, key := range keys {
		v := c.Params.Values.Get(key)
		if v != "" {
			ok, msg := emailService.ValidTpl(v)
			if !ok {
				re.Ok = false
				re.Msg = "Error key: " + key + "<br />" + msg
				return c.RenderJSON(re)
			} else {
				configService.UpdateGlobalStringConfig(userId, key, v)
			}
		}
	}

	re.Ok = true
	return c.RenderJSON(re)
}

// 发送Email
func (c AdminEmail) SendEmailToEmails(sendEmails, latestEmailSubject, latestEmailBody string, verified, saveAsOldEmail bool) revel.Result {
	re := info.NewRe()

	c.updateConfig([]string{"sendEmails", "latestEmailSubject", "latestEmailBody"})

	if latestEmailSubject == "" || latestEmailBody == "" {
		re.Msg = "subject or body is blank"
		return c.RenderJSON(re)
	}

	if saveAsOldEmail {
		oldEmails := configService.GetGlobalMapConfig("oldEmails")
		oldEmails[latestEmailSubject] = latestEmailBody
		configService.UpdateGlobalMapConfig(c.GetUserId(), "oldEmails", oldEmails)
	}

	sendEmails = strings.Replace(sendEmails, "\r", "", -1)
	emails := strings.Split(sendEmails, "\n")

	re.Ok, re.Msg = emailService.SendEmailToEmails(emails, latestEmailSubject, latestEmailBody)
	return c.RenderJSON(re)
}

// 发送Email
func (c AdminEmail) SendToUsers2(emails, latestEmailSubject, latestEmailBody string, verified, saveAsOldEmail bool) revel.Result {
	re := info.NewRe()

	c.updateConfig([]string{"sendEmails", "latestEmailSubject", "latestEmailBody"})

	if latestEmailSubject == "" || latestEmailBody == "" {
		re.Msg = "subject or body is blank"
		return c.RenderJSON(re)
	}

	if saveAsOldEmail {
		oldEmails := configService.GetGlobalMapConfig("oldEmails")
		oldEmails[latestEmailSubject] = latestEmailBody
		configService.UpdateGlobalMapConfig(c.GetUserId(), "oldEmails", oldEmails)
	}

	emails = strings.Replace(emails, "\r", "", -1)
	emailsArr := strings.Split(emails, "\n")

	users := userService.ListUserInfosByEmails(emailsArr)
	LogJ(emailsArr)

	re.Ok, re.Msg = emailService.SendEmailToUsers(users, latestEmailSubject, latestEmailBody)

	return c.RenderJSON(re)
}

// send Email dialog
func (c AdminEmail) SendEmailDialog(emails string) revel.Result {
	emailsArr := strings.Split(emails, ",")
	emailsNl := strings.Join(emailsArr, "\n")

	c.ViewArgs["emailsNl"] = emailsNl
	c.ViewArgs["str"] = configService.GlobalStringConfigs
	c.ViewArgs["map"] = configService.GlobalMapConfigs

	return c.RenderTemplate("admin/email/emailDialog.html")
}

func (c AdminEmail) SendToUsers(userFilterEmail, userFilterWhiteList, userFilterBlackList, latestEmailSubject, latestEmailBody string, verified, saveAsOldEmail bool) revel.Result {
	re := info.NewRe()

	c.updateConfig([]string{"userFilterEmail", "userFilterWhiteList", "userFilterBlackList", "latestEmailSubject", "latestEmailBody"})

	if latestEmailSubject == "" || latestEmailBody == "" {
		re.Msg = "subject or body is blank"
		return c.RenderJSON(re)
	}

	if saveAsOldEmail {
		oldEmails := configService.GetGlobalMapConfig("oldEmails")
		oldEmails[latestEmailSubject] = latestEmailBody
		configService.UpdateGlobalMapConfig(c.GetUserId(), "oldEmails", oldEmails)
	}

	users := userService.GetAllUserByFilter(userFilterEmail, userFilterWhiteList, userFilterBlackList, verified)

	if users == nil || len(users) == 0 {
		re.Ok = false
		re.Msg = "no users"
		return c.RenderJSON(re)
	}

	re.Ok, re.Msg = emailService.SendEmailToUsers(users, latestEmailSubject, latestEmailBody)
	if !re.Ok {
		return c.RenderJSON(re)
	}

	re.Ok = true
	re.Msg = "users:" + strconv.Itoa(len(users))

	return c.RenderJSON(re)
}

// 删除emails
func (c AdminEmail) DeleteEmails(ids string) revel.Result {
	re := info.NewRe()
	re.Ok = emailService.DeleteEmails(strings.Split(ids, ","))
	return c.RenderJSON(re)
}

func (c AdminEmail) List(sorter, keywords string) revel.Result {
	pageNumber := c.GetPage()
	sorterField, isAsc := c.getSorter("CreatedTime", false, []string{"email", "ok", "subject", "createdTime"})
	pageInfo, emails := emailService.ListEmailLogs(pageNumber, userPageSize, sorterField, isAsc, keywords)
	c.ViewArgs["pageInfo"] = pageInfo
	c.ViewArgs["emails"] = emails
	c.ViewArgs["keywords"] = keywords
	return c.RenderTemplate("admin/email/list.html")
}
