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
	"github.com/beego/beego/v2/server/web"
	"github.com/rs/xid"
	"gorm.io/gorm/clause"
	"log"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/10 2:44 下午
 * @version 1.0
 */

type CommentRestApi struct {
	base.BaseController
}

func (c *CommentRestApi) GetList() {
	var commentVO vo.CommentVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &commentVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		panic(err)
	}
	where := "status=? and type=? and source=? and to_uid in (null,'')"
	if commentVO.BlogUid != "" {
		where = where + " and blog_uid='" + commentVO.BlogUid + "'"
	}
	var pageList []models.Comment
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Comment{}).Where(where, 1, 0, commentVO.Source).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1, 0, commentVO.Source).Offset((commentVO.CurrentPage - 1) * commentVO.PageSize).Limit(commentVO.PageSize).Order("create_time desc").Preload(clause.Associations).Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var firstUidList []string
	for _, item := range pageList {
		firstUidList = append(firstUidList, item.Uid)
	}
	if len(firstUidList) > 0 {
		var notFirstList []models.Comment
		common.DB.Where("status=? and first_comment_uid in ?", 1, firstUidList).Preload(clause.Associations).Find(&notFirstList)
		if len(notFirstList) > 0 {
			pageList = append(pageList, notFirstList...)
		}
	}
	var userUidList []string
	for _, item := range pageList {
		userUid := item.UserUid
		toUserUid := item.ToUserUid
		if userUid != "" {
			userUidList = append(userUidList, userUid)
		}
		if toUserUid != "" {
			userUidList = append(userUidList, toUserUid)
		}
	}
	var userList []models.User
	if len(userUidList) > 0 {
		common.DB.Find(&userList, userUidList)
	}
	var filterUserList []models.User
	for _, item := range userList {
		user := models.User{
			Avatar:   item.Avatar,
			Uid:      item.Uid,
			NickName: item.NickName,
			UserTag:  item.UserTag,
		}
		filterUserList = append(filterUserList, user)
	}
	var fileUids strings.Builder
	for _, item := range filterUserList {
		if item.Avatar != "" {
			fileUids.WriteString(item.Avatar + ",")
		}
	}
	pictureList := map[string]interface{}{}
	if !reflect.DeepEqual(fileUids, strings.Builder{}) {
		pictureList = service.FileService.GetPicture(fileUids.String(), ",")
	}
	picList := common.WebUtil.GetPictureMap(pictureList)
	pictureMap := map[string]string{}
	for _, item := range picList {
		pictureMap[item["uid"].(string)] = item["url"].(string)
	}
	userMap := map[string]models.User{}
	for _, item := range filterUserList {
		if item.Avatar != "" && pictureMap[item.Avatar] != "" {
			item.PhotoUrl = pictureMap[item.Avatar]
		}
		userMap[item.Uid] = item
	}
	for i, item := range pageList {
		if item.UserUid != "" {
			pageList[i].User = userMap[item.UserUid]
		}
		if item.ToUserUid != "" {
			pageList[i].ToUser = userMap[item.ToUserUid]
		}
	}
	toCommentListMap := map[string][]models.Comment{}
	var tempList []models.Comment
	for a := 0; a < len(pageList); a++ {
		for b := 0; b < len(pageList); b++ {
			if pageList[a].Uid == pageList[b].ToUid {
				tempList = append(tempList, pageList[b])
			}
		}
		toCommentListMap[pageList[a].Uid] = tempList
		tempList = []models.Comment{}
	}
	var firstComment []models.Comment
	for _, item := range pageList {
		if item.ToUid == "" {
			firstComment = append(firstComment, item)
		}
	}
	iPage := page.IPage{
		Records: getCommentReplys(firstComment, toCommentListMap),
		Total:   total,
		Size:    commentVO.PageSize,
		Current: commentVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func getCommentReplys(list []models.Comment, toCommentListMap map[string][]models.Comment) []models.Comment {
	if list == nil {
		return []models.Comment{}
	}
	for i, item := range list {
		commentUid := item.Uid
		replyCommentList := toCommentListMap[commentUid]
		replyComments := getCommentReplys(replyCommentList, toCommentListMap)
		list[i].ReplyList = replyComments
	}
	return list
}

func (c *CommentRestApi) GetListByUser() {
	var userVO vo.UserVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userVO)
	if err != nil {
		c.ThrowError("00", "传入参数有误！")
		panic(err)
	}
	token := c.Ctx.GetCookie("token")
	if token == "" {
		c.Result("error", "token令牌未被识别，请重新登录")
		return
	}
	userJson := common.RedisUtil.Get("USER_TOKEN:" + token)
	var user models.User
	err1 := json.Unmarshal([]byte(userJson), &user)
	if err1 != nil {
		panic(err1)
	}
	requestUserUid := c.GetSession("userUid").(string)
	var pageList []models.Comment
	common.DB.Where("status=? and type=? and (user_uid=? or to_user_uid=?)", 1, 0, requestUserUid, requestUserUid).Offset((userVO.CurrentPage - 1) * userVO.PageSize).Limit(userVO.PageSize).Order("create_time desc").Find(&pageList)
	var userUidList []string
	for _, item := range pageList {
		userUid := item.UserUid
		toUserUid := item.ToUserUid
		if userUid != "" {
			userUidList = append(userUidList, item.UserUid)
		}
		if toUserUid != "" {
			userUidList = append(userUidList, item.ToUserUid)
		}
	}
	var userList []models.User
	if len(userUidList) > 0 {
		common.DB.Find(&userList, userUidList)
	}
	var fileterUserList []models.User
	for _, item := range userList {
		var user models.User
		user.Avatar = item.Avatar
		user.Uid = item.Uid
		user.NickName = item.NickName
		fileterUserList = append(fileterUserList, user)
	}
	var fileUids strings.Builder
	for _, item := range fileterUserList {
		if item.Avatar != "" {
			fileUids.WriteString(item.Avatar + ",")
		}
	}
	pictureList := map[string]interface{}{}
	if reflect.DeepEqual(fileUids, strings.Builder{}) {
		pictureList = service.FileService.GetPicture(fileUids.String(), ",")
	}
	picList := common.WebUtil.GetPictureMap(pictureList)
	pictureMap := map[string]string{}
	for _, item := range picList {
		pictureMap[item["uid"].(string)] = item["url"].(string)
	}
	userMap := map[string]interface{}{}
	for _, item := range fileterUserList {
		if item.Avatar != "" && pictureMap[item.Avatar] != "" {
			item.PhotoUrl = pictureMap[item.Avatar]
		}
		userMap[item.Uid] = item
	}
	var commentList []models.Comment
	var replayList []models.Comment
	for i, item := range pageList {
		if item.UserUid != "" {
			pageList[i].User = userMap[item.UserUid].(models.User)
		}
		if item.ToUserUid != "" {
			pageList[i].ToUser = userMap[item.ToUserUid].(models.User)
		}
		if item.Source != "" {
			eCommentSource := common.Emu.CommentSourceEmu()
			pageList[i].SourceName = eCommentSource[item.Source]["name"]
		}
		if item.UserUid == requestUserUid {
			commentList = append(commentList, item)
		}
		if item.ToUserUid == requestUserUid {
			replayList = append(replayList, item)
		}
	}
	resultMap := map[string]interface{}{}
	resultMap["commentList"] = commentList
	resultMap["replayList"] = replayList
	c.SuccessWithData(resultMap)
}

func (c *CommentRestApi) GetPraiseListByUser() {
	currentPage, _ := c.GetInt("currentPage", 1)
	pageSize, _ := c.GetInt("pageSize", 10)
	token := c.Ctx.GetCookie("token")
	if token == "" {
		c.Result("error", "token令牌未被识别，请重新登录")
		return
	}
	userUid := c.GetSession("userUid").(string)
	var pageList []models.Comment
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Comment{}).Where("status=? and type=? and user_uid=?", 1, 1, userUid).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where("status=? and type=? and user_uid=?", 1, 1, userUid).Offset((currentPage - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var blogUids []string
	for _, item := range pageList {
		blogUids = append(blogUids, item.BlogUid)
	}
	blogMap := map[string]models.Blog{}
	if len(blogUids) > 0 {
		var blogList []models.Blog
		common.DB.Find(&blogList, blogUids)
		for _, blog := range blogList {
			blog.Content = ""
			blogMap[blog.Uid] = blog
		}
	}
	for i, item := range pageList {
		if !reflect.DeepEqual(blogMap[item.BlogUid], models.Blog{}) {
			pageList[i].Blog = blogMap[item.BlogUid]
		}
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    pageSize,
		Current: currentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *CommentRestApi) Add() {
	var commentVO vo.CommentVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &commentVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		panic(err)
	}
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	if token == "" {
		c.Result("error", "token令牌未被识别，请重新登录")
		return
	}
	var webConfig models.WebConfig
	common.DB.Where("status=?", 1).First(&webConfig)
	if webConfig.OpenComment == "0" {
		c.Result("error", "网站未开启评论功能")
		return
	}
	var blog models.Blog
	if commentVO.BlogUid != "" {
		common.DB.Where("uid=?", commentVO.BlogUid).Find(&blog)
		if blog.OpenComment == 0 {
			c.Result("error", "博客未开启评论")
			return
		}
	}
	var user models.User
	tokenJson := common.RedisUtil.Get("USER_TOKEN:" + token)
	//common.DB.Where("uid=?", userUid).Find(&user)
	err1 := json.Unmarshal([]byte(tokenJson), &user)
	if err1 != nil {
		panic(err)
	}
	if utf8.RuneCountInString(commentVO.Content) > 1000 {
		c.Result("error", "评论不能大于1000字")
		return
	}
	if user.CommentStatus == 0 {
		c.Result("error", "你没有评论权限")
		return
	}
	jsonResult := common.RedisUtil.Get("USER_PUBLISH_SPAM_COMMENT_COUNT:" + user.Uid)
	var count int
	if jsonResult != "" {
		count, _ = strconv.Atoi(jsonResult)
		if count > 5 {
			c.Result("error", "由于发送过多无意义评论，您已被禁言一小时，请稍后在试~")
			return
		}
	}
	content := commentVO.Content
	if common.StringUtils.IsCommentSpam(content) {
		if len(jsonResult) == 0 {
			common.RedisUtil.SetEx("USER_PUBLISH_SPAM_COMMENT_COUNT:"+user.Uid, "0", 1, time.Hour)
		} else {
			common.RedisUtil.IncrBy("USER_PUBLISH_SPAM_COMMENT_COUNT:"+user.Uid, 1)
		}
		c.Result("error", "请输入有意义的评论内容")
		return
	}
	if commentVO.ToUserUid != "" {
		var toUser models.User
		common.DB.Where("uid=?", commentVO.ToUserUid).Find(&toUser)
		if toUser.StartEmailNotification == 1 {
			var toComment models.Comment
			common.DB.Where("uid=?", commentVO.ToUid).Find(&toComment)
			if !reflect.DeepEqual(toComment, models.Comment{}) && toComment.Content != "" {
				log.Print("发送评论邮件")
				m := map[string]string{}
				m["email"] = toUser.Email
				m["text"] = commentVO.Content
				m["to_text"] = toComment.Content
				m["nickname"] = user.NickName
				m["to_nickname"] = toUser.NickName
				m["user_uid"] = toUser.Uid
				commentSource := toComment.Source
				var url string
				dataWebsiteUrl, _ := web.AppConfig.String("data_website_url")
				switch commentSource {
				case "ABOUT":
					url = dataWebsiteUrl + "about"
				case "BLOG_INFO":
					url = dataWebsiteUrl + "info?blogUid=" + toComment.BlogUid
				case "MESSAGE_BOARD":
					url = dataWebsiteUrl + "messageBoard"
				default:
					c.ErrorWithMessage("跳转到其他链接")
				}
				m["url"] = url
				common.Email.SentCommentEmail(m)

			}
		}
	}
	var comment models.Comment
	comment.Source = commentVO.Source
	comment.BlogUid = commentVO.BlogUid
	comment.Content = commentVO.Content
	comment.ToUserUid = commentVO.ToUserUid
	if commentVO.ToUid != "" {
		var toComment models.Comment
		common.DB.Where("uid=?", commentVO.ToUid).Find(&toComment)
		if !reflect.DeepEqual(toComment, models.Comment{}) && toComment.FirstCommentUid != "" {
			comment.FirstCommentUid = toComment.FirstCommentUid
		} else {
			comment.FirstCommentUid = toComment.Uid
		}
	} else {
		m := service.SystemConfigService.GetConfig()
		systemConfig := m["data"].(models.SystemConfig)
		if systemConfig != (models.SystemConfig{}) && systemConfig.StartEmailNotification == "1" {
			if systemConfig.Email != "" {
				log.Print("发送评论邮件")
				commentSource := common.Emu.CommentSourceEmu()
				sourceName := commentSource[commentVO.Source]["name"]
				linkText := "<a href=\" " + getUrlByCommentSource(commentVO) + "\">" + sourceName + "</a>\n"
				commentContent := linkText + "收到新的评论: " + commentVO.Content
				common.Email.SentSimpleEmail(systemConfig.Email, commentContent)
			} else {
				c.Result("error", "网站没有配置通知接收的邮箱地址！")
			}
		}
	}
	comment.UserUid = commentVO.UserUid
	comment.ToUid = commentVO.ToUid
	comment.Uid = xid.New().String()
	common.DB.Create(&comment)
	if user.Avatar != "" {
		pictureList := service.FileService.GetPicture(user.Avatar, ",")
		u := common.WebUtil.GetPicture(pictureList)
		if len(u) > 0 {
			user.PhotoUrl = u[0]
		}
	}
	comment.User = user

	// 如果是回复某人的评论，那么需要向该用户Redis收件箱中中写入一条记录
	if comment.ToUserUid != "" {
		count := common.RedisUtil.Get("USER_RECEIVE_COMMENT_COUNT:" + comment.ToUserUid)
		if count == "" {
			common.RedisUtil.Set("USER_RECEIVE_COMMENT_COUNT:" + comment.ToUserUid, "1")
		} else {
			countTemp, _ := strconv.Atoi(count)
			countTemp++
			common.RedisUtil.SetEx("USER_RECEIVE_COMMENT_COUNT:" + comment.ToUserUid, strconv.Itoa(countTemp), 10, 7 * 24 * time.Hour)
		}
	}

	c.SuccessWithData(comment)
}

func (c *CommentRestApi) DeleteBatch() {
	var commentVO vo.CommentVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &commentVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		panic(err)
	}
	var comment models.Comment
	common.DB.Where("uid=?", commentVO.Uid).Find(&comment)
	if comment.UserUid != commentVO.UserUid {
		c.Result("error", "没权限删除该数据")
		return
	}
	common.DB.Model(&comment).Select("status").Update("status", 0)
	var commentList []models.Comment
	commentList = append(commentList, comment)
	var firstCommentUid string
	if comment.FirstCommentUid != "" {
		firstCommentUid = comment.FirstCommentUid
	} else {
		firstCommentUid = comment.Uid
	}
	var toCommentList []models.Comment
	common.DB.Where("status=? and first_comment_uid=?", 1, firstCommentUid).Find(&toCommentList)
	var resultList []models.Comment
	getToCommentList(comment, toCommentList, &resultList)
	if len(resultList) > 0 {
		for i := range resultList {
			resultList[i].Status = 0
		}
		common.DB.Model(&resultList).Select("status").Update("status", 0)
	}
	c.SuccessWithData("删除成功")
}

func getToCommentList(comment models.Comment, commentList []models.Comment, resultList *[]models.Comment) {
	if reflect.DeepEqual(comment, models.Comment{}) {
		return
	}
	commentUid := comment.Uid
	for _, item := range commentList {
		if commentUid == item.ToUid {
			*resultList = append(*resultList, item)
			getToCommentList(item, commentList, resultList)
		}
	}
}

func getUrlByCommentSource(commentVO vo.CommentVO) string {
	var linkUrl string
	var dataWebsiteUrl, _ = web.AppConfig.String("data_website_url")
	commentSource := commentVO.Source
	switch commentSource {
	case "ABOUT":
		linkUrl = dataWebsiteUrl + "about"
	case "BLOG_INFO":
		linkUrl = dataWebsiteUrl + "info?blogUid=" + commentVO.BlogUid
	case "MESSAGE_BOARD":
		linkUrl = dataWebsiteUrl + "messageBoard"
	default:
		linkUrl = dataWebsiteUrl
		log.Print("跳转到其他链接")
	}
	return linkUrl
}

func (c *CommentRestApi) Report() {
	var commentVO vo.CommentVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &commentVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		panic(err)
	}
	var comment models.Comment
	common.DB.Where("uid=?", commentVO.Uid).Find(&comment)
	if reflect.DeepEqual(comment, models.Comment{}) || comment.Status == 0 {
		c.Result("error", "评论不存在")
		return
	}
	if comment.UserUid == commentVO.UserUid {
		c.Result("error", "不能举报自己的评论")
		return
	}
	var commentReportList []models.CommentReport
	common.DB.Where("user_uid=? and report_comment_uid=?", commentVO.UserUid, commentVO.Uid).Find(&commentReportList)
	if len(commentReportList) > 0 {
		c.Result("error", "不能重复举报")
	}
	var commentReport models.CommentReport
	commentReport.Content = commentVO.Content
	commentReport.Progress = 0
	commentReport.UserUid = commentVO.UserUid
	commentReport.ReportCommentUid = commentVO.Uid
	commentReport.ReportUserUid = comment.UserUid
	commentReport.Uid = xid.New().String()
	common.DB.Create(&commentReport)
	c.Result("success", "举报成功")
}

func (c *CommentRestApi) CloseEmailNotification() {
	userUid := c.GetString(":userUid")
	var user models.User
	common.DB.Where("uid=?", userUid).Find(&user)
	if user == (models.User{}) {
		c.Data["title"] = "出错了"
		c.Data["subtitle"] = "用户不存在"
		c.TplName = "bindUserEmail.html"
		err := c.Render()
		if err != nil {
			panic(err)
		}
		return
	}
	common.DB.Model(&user).Select("start_email_notification").Update("start_email_notification", 0)
	if user.ValidCode != "" {
		userInfo := common.RedisUtil.Get("USER_TOKEN:" + user.ValidCode)
		if userInfo != "" {
			m := map[string]interface{}{}
			err := json.Unmarshal([]byte(userInfo), &m)
			if err != nil {
				panic(err)
			}
			m["startEmailNotification"] = 0
			b, _ := json.Marshal(m)
			userTokenSurvivalTime, _ := web.AppConfig.Int64("user_token_survival_time")
			common.RedisUtil.SetEx("USER_TOKEN:"+user.ValidCode, string(b), userTokenSurvivalTime, time.Hour)
		}
	}
	c.Data["title"] = "恭喜你"
	c.Data["subtitle"] = "成功关闭邮件通知"
	c.TplName = "bindUserEmail.html"
	err := c.Render()
	if err != nil {
		panic(err)
	}
}

// 获取用户收到的评论回复数
func (c *CommentRestApi) GetUserReceiveCommentCount() {
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	if token == "" {
		c.Result("error", "token令牌未被识别，请重新登录")
		return
	}
	userJson := common.RedisUtil.Get("USER_TOKEN:" + token)
	var user models.User
	err1 := json.Unmarshal([]byte(userJson), &user)
	if err1 != nil {
		panic(err1)
	}
	var commentCount = 0
	// 评论数
	count := common.RedisUtil.Get("USER_RECEIVE_COMMENT_COUNT:" + user.Uid)
	if count != "" {
		commentCount, _ = strconv.Atoi(count)
	}
	c.Result("success", commentCount)

}

// 阅读用户接收的评论数
func (c *CommentRestApi) ReadUserReceiveCommentCount() {
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	if token == "" {
		c.Result("error", "token令牌未被识别，请重新登录")
		return
	}
	userJson := common.RedisUtil.Get("USER_TOKEN:" + token)
	var user models.User
	err1 := json.Unmarshal([]byte(userJson), &user)
	if err1 != nil {
		panic(err1)
	}
	common.RedisUtil.Delete("USER_RECEIVE_COMMENT_COUNT:" + user.Uid)
	c.SuccessWithMessage("阅读成功")
}
