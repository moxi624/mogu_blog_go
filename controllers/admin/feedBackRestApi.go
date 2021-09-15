package admin

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/22 11:11 上午
 * @version 1.0
 */

type FeedbackRestApi struct {
	base.BaseController
}

func (c *FeedbackRestApi) GetList() {
	base.L.Print("获取反馈列表")
	var feedbackVO vo.FeedbackVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &feedbackVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		panic(err)
	}
	where := "status=?"
	if feedbackVO.Title != "" {
		where += " and title like \"%" + feedbackVO.Title + "%\""
	}
	if feedbackVO.FeedbackStatus != 0 {
		where += " and feedback_status=" + string(feedbackVO.FeedbackStatus)
	}
	var pageList []models.Feedback
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Feedback{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((feedbackVO.CurrentPage - 1) * feedbackVO.PageSize).Limit(feedbackVO.PageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var userUids []string
	for _, item := range pageList {
		if item.UserUid != "" {
			userUids = append(userUids, item.UserUid)
		}
	}
	userList := service.UserService.GetUserListByIds(userUids)
	m := map[string]models.User{}
	for _, item := range userList {
		item.PassWord = ""
		m[item.Uid] = item
	}
	for i, item := range pageList {
		if item.UserUid != "" {
			pageList[i].User = m[item.UserUid]
		}
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    feedbackVO.PageSize,
		Current: feedbackVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *FeedbackRestApi) Edit() {
	base.L.Print("编辑反馈")
	var feedbackVO vo.FeedbackVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &feedbackVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		panic(err)
	}
	//tokenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	tokenHead, _ := beego.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	t := common.Jwx.ParseToken(tokenString)
	a, _ := t.Get("adminUid")
	if a.(string) == "" {
		c.ErrorWithMessage("操作失败,请重新登录")
		return
	}
	var feedback models.Feedback
	common.DB.Where("uid=?", feedbackVO.Uid).Find(&feedback)
	feedback.Title = feedbackVO.Title
	feedback.Content = feedbackVO.Content
	feedback.FeedbackStatus = feedbackVO.FeedbackStatus
	feedback.Reply = feedbackVO.Reply
	common.DB.Save(&feedback)
	c.SuccessWithMessage("更新成功")
}

func (c *FeedbackRestApi) DeleteBatch() {
	base.L.Print("批量删除反馈")
	var feedbackVOList []vo.FeedbackVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &feedbackVOList)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		panic(err)
	}
	//tokenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	tokenHead, _ := beego.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	t := common.Jwx.ParseToken(tokenString)
	a, _ := t.Get("adminUid")
	if a.(string) == "" {
		c.ErrorWithMessage("操作失败,请重新登录")
		return
	}
	adminUid := a.(string)
	if len(feedbackVOList) <= 0 {
		c.ErrorWithMessage("传入参数有误")
		return
	}
	var uids []string
	for _, item := range feedbackVOList {
		uids = append(uids, item.Uid)
	}
	var feedbackList []models.Feedback
	common.DB.Find(&feedbackList, uids)
	save := common.DB.Model(&feedbackList).Select("status", "admin_uid").Updates(models.Feedback{Status: 0, AdminUid: adminUid}).Error
	if save != nil {
		c.ErrorWithMessage("删除失败")
	} else {
		c.SuccessWithMessage("删除成功")
	}

}
