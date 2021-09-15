//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package web

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/rs/xid"
	"log"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"reflect"
	"strings"
	"time"
)

const bindUserEmail = "bindUserEmail.html"

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/20 9:28 上午
 * @version 1.0
 */

type AuthRestApi struct {
	base.BaseController
}

func (c *AuthRestApi) Verify() {
	accessToken := c.GetString(":accessToken")
	userInfo := common.RedisUtil.Get("USER_TOKEN:" + accessToken)
	if userInfo == "" {
		c.ErrorWithMessage("token令牌未被识别")
		return
	}
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(userInfo), &m)
	if err != nil {
		panic(err)
	}
	c.SuccessWithData(m)
}

func (c *AuthRestApi) GetFeedbackList() {
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	if token == "" {
		c.ErrorWithMessage("token令牌未被识别，请重新登录")
		return
	}
	var user models.User
	tokenJson := common.RedisUtil.Get("USER_TOKEN:" + token)
	err := json.Unmarshal([]byte(tokenJson), &user)
	if err != nil {
		panic(err)
	}
	var pageList []models.Feedback
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Feedback{}).Where("status=? and user_uid=?", 1, user.Uid).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where("status=? and user_uid=?", 1, user.Uid).Limit(20).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    20,
		Current: 1,
	}
	c.SuccessWithData(iPage)
}

func (c *AuthRestApi) Delete() {
	accessToken := c.GetString(":accessToken")
	common.RedisUtil.Delete("USER_TOKEN:" + accessToken)
	c.SuccessWithData("注销成功")
}

func (c *AuthRestApi) EditUser() {
	var userVO vo.UserVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userVO)
	if err != nil {
		panic(err)
	}
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	if token == "" {
		c.Result("error", "token令牌未被识别，请重新登录")
		return
	}
	var user models.User
	tokenJson := common.RedisUtil.Get("USER_TOKEN:" + token)
	err1 := json.Unmarshal([]byte(tokenJson), &user)
	if err1 != nil {
		panic(err1)
	}
	if user == (models.User{}) {
		c.Result("error", "用户未找到或者未登陆")
		return
	}
	base.L.Print("获取用户")
	common.DB.Where("uid = ?", user.Uid).Find(&user)
	g := common.InterfaceToInt(userVO.Gender)
	user.NickName = userVO.NickName
	user.Avatar = userVO.Avatar
	user.Birthday = userVO.Birthday
	user.Summary = userVO.Summary
	user.Gender = g
	user.QqNumber = userVO.QqNumber
	user.Occupation = userVO.Occupation
	if userVO.StartEmailNotification == 1 && userVO.Email == "" {
		c.Result("error", "必须填写并绑定邮箱后，才能开启评论邮件通知~")
		return
	}
	user.StartEmailNotification = userVO.StartEmailNotification
	common.DB.Save(&user)
	user.PassWord = ""
	user.PhotoUrl = userVO.PhotoUrl
	userTokenSurvivalTime, _ := beego.AppConfig.Int64("user_token_survival_time")
	if userVO.Email != "" && userVO.Email != user.Email {
		user.Email = userVO.Email
		common.Email.SendRegisterEmail(user, token)
		c.SuccessWithData("您已修改邮箱，请先到邮箱进行确认绑定")
	} else {
		c.SuccessWithData("更新成功")
	}
	b, _ := json.Marshal(user)
	common.RedisUtil.SetEx("USER_TOKEN:"+token, string(b), userTokenSurvivalTime, time.Hour)
}

func (c *AuthRestApi) BindUserEmail() {
	token := c.GetString(":token")
	userInfo := common.RedisUtil.Get("USER_TOKEN:" + token)
	if userInfo == "" {
		c.Data["title"] = "出错了"
		c.Data["subtitle"] = "token令牌未被识别，请重新登录"
		c.TplName = bindUserEmail
		err := c.Render()
		if err != nil {
			panic(err)
		}
		return
	}
	var user models.User
	err := json.Unmarshal([]byte(userInfo), &user)
	if err != nil {
		panic(err)
	}
	common.DB.Model(&user).Select("email").Update("email", user.Email)
	c.Data["title"] = "恭喜你"
	c.Data["subtitle"] = "成功绑定新邮箱"
	c.TplName = bindUserEmail
	err1 := c.Render()
	if err1 != nil {
		panic(err)
	}
}

func (c *AuthRestApi) UpdateUserPwd() {
	oldPwd := c.GetString("oldPwd")
	newPwd := c.GetString("newPwd")
	if oldPwd == "" {
		c.Result("error", "原密码不能为空")
		return
	}
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	if token == "" {
		c.Result("error", "token令牌未被识别，请重新登录")
		return
	}
	tokenJson := common.RedisUtil.Get("USER_TOKEN:" + token)
	var user models.User
	err1 := json.Unmarshal([]byte(tokenJson), &user)
	if err1 != nil {
		panic(err1)
	}
	if user.Source != "MOGU" {
		c.Result("error", "第三方登录无法修改密码")
		return
	}
	common.DB.Where("uid = ?", user.Uid).Find(&user)
	if user.PassWord == common.MD5(oldPwd) {
		user.PassWord = common.MD5(newPwd)
		common.DB.Save(&user)
		c.Result("success", "修改成功")
	} else {
		c.Result("error", "密码错误")
	}
}

func (c *AuthRestApi) AddFeedback() {
	var feedbackVO vo.FeedbackVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &feedbackVO)
	if err != nil {
		panic(err)
	}
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	if token == "" {
		c.Result("error", "token令牌未被识别，请重新登录")
		return
	}
	var user models.User
	tokenJson := common.RedisUtil.Get("USER_TOKEN:" + token)
	err1 := json.Unmarshal([]byte(tokenJson), &user)
	if err1 != nil {
		panic(err1)
	}
	if user.CommentStatus == 0 {
		c.Result("error", "你没有反馈权限")
		return
	}
	m := service.SystemConfigService.GetConfig()
	systemConfig := m["data"].(models.SystemConfig)
	if systemConfig != (models.SystemConfig{}) && systemConfig.StartEmailNotification == "1" {
		if systemConfig.Email != "" {
			log.Print("发送反馈邮件通知")
			feedback := "网站收到新的反馈: " + "<br />" + "标题：" + feedbackVO.Title + "<br />" + "<br />" + "内容" + feedbackVO.Content
			common.Email.SentSimpleEmail(systemConfig.Email, feedback)
		} else {
			c.Result("error", "网站没有配置通知接收的邮箱地址！")
		}
	}
	var feedback models.Feedback
	feedback.UserUid = user.Uid
	feedback.Title = feedbackVO.Title
	feedback.Content = feedbackVO.Content
	feedback.Uid = xid.New().String()
	common.DB.Create(&feedback)
	c.Result("success", "新增成功")
}

func (c *AuthRestApi) ReplyBlogLink() {
	var linkVO vo.LinkVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &linkVO)
	if err != nil {
		panic(err)
	}
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	if token == "" {
		c.Result("error", "token令牌未被识别，请重新登录")
		return
	}
	var user models.User
	tokenJson := common.RedisUtil.Get("USER_TOKEN:" + token)
	err1 := json.Unmarshal([]byte(tokenJson), &user)
	if err1 != nil {
		panic(err1)
	}
	if user.CommentStatus == 0 {
		c.Result("error", "你没有申请权限")
		return
	}
	m := service.SystemConfigService.GetConfig()
	systemConfig := m["data"].(models.SystemConfig)
	if systemConfig != (models.SystemConfig{}) && systemConfig.StartEmailNotification == "1" {
		if systemConfig.Email != "" {
			log.Print("发送友链申请邮件通知")
			feedback := "网站收到新的友链申请: " + "<br />" + "名称：" + linkVO.Title + "<br />" + "<br />" + "简介" +
				linkVO.Summary + "<br />" + "<br />" + "地址" + linkVO.Url
			common.Email.SentSimpleEmail(systemConfig.Email, feedback)
		} else {
			c.Result("error", "网站没有配置通知接收的邮箱地址！")
		}
	}
	var existLink models.Link
	common.DB.Where("status = ? and user_uid=? and title=?", 1, user.Uid, linkVO.Title).First(&existLink)
	if !reflect.DeepEqual(existLink, models.Link{}) {
		linkStatus := existLink.LinkStatus
		message := ""
		switch linkStatus {
		case 0:
			message = "您申请的友链，已经在申请列表中！"
		case 1:
			message = "您申请的友链，已经发布!"
		case 2:
			message = "您申请的友链，已经下架！"
		}
		c.Result("error", message)
	}
	var link models.Link
	link.Title = linkVO.Title
	link.Summary = linkVO.Summary
	link.Url = linkVO.Url
	link.FileUid = linkVO.FileUid
	link.Email = linkVO.Email
	link.UserUid = user.Uid
	link.Uid = xid.New().String()
	common.DB.Create(&link)
	c.Result("success", "申请已发送")
}

func (c *AuthRestApi) RenderOauth() {
	source := c.GetString("source")
	isOpenLoginType, _ := service.WebConfigService.IsOpenLoginType(strings.ToUpper(source))
	if !isOpenLoginType {
		c.Result("error", "后台未开启该登录方式!")
	}
	log.Print("进入render")
	c.Result("error", "目前暂时还没有用开通第三方登录")
}

func getAuthRequest(source string) {
	switch source {
	case "wechat":
		break
	case "qq":
	default:
		break
	}
}
