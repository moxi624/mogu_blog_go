package admin

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/rs/xid"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/26 8:45 上午
 * @version 1.0
 */

type TodoRestApi struct {
	base.BaseController
}

var l = logs.GetLogger()

func (c *TodoRestApi) GetList() {
	var todoVO vo.TodoVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &todoVO)
	if err != nil {
		panic(err)
	}
	l.Println("执行获取代办事项列表")
	//tokenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	if tokenData == "" {
		c.Result("error", "token用户过期")
		return
	}
	tokenHead, _ := beego.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	t := common.Jwx.ParseToken(tokenString)
	a, _ := t.Get("adminUid")
	adminUid := a.(string)
	where := "admin_uid= ? and status=?"
	if todoVO.Keyword != "" || strings.TrimSpace(todoVO.Keyword) != "" {
		where += " AND text like \"%" + strings.TrimSpace(todoVO.Keyword) + "%\""
	}
	var pageList []models.Todo
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Todo{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, adminUid, 1).Offset((todoVO.CurrentPage - 1) * todoVO.PageSize).Limit(todoVO.PageSize).Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    todoVO.PageSize,
		Current: todoVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *TodoRestApi) Add() {
	var todoVO vo.TodoVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &todoVO)
	if err != nil {
		panic(err)
	}
	//tokenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	tokenHead, _ := beego.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	t := common.Jwx.ParseToken(tokenString)
	a, _ := t.Get("adminUid")
	adminUid := a.(string)
	var todo models.Todo
	todo.Text = todoVO.Text
	todo.Done = false
	todo.AdminUid = adminUid
	todo.Uid = xid.New().String()
	common.DB.Create(&todo)
	c.Result("success", "插入成功")
}

func (c *TodoRestApi) Edit() {
	var todoVO vo.TodoVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &todoVO)
	if err != nil {
		panic(err)
	}
	//tokenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	tokenHead, _ := beego.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	t := common.Jwx.ParseToken(tokenString)
	a, _ := t.Get("adminUid")
	adminUid := a.(string)
	var todo models.Todo
	common.DB.Where("uid=?", todoVO.Uid).Find(&todo)
	if todo.AdminUid != adminUid {
		c.Result("error", "该资源无权限访问")
		return
	}
	todo.Text = todoVO.Text
	todo.Done = todoVO.Done
	common.DB.Save(&todo)
	c.Result("success", "更新成功")
}

func (c *TodoRestApi) Delete() {
	var todoVO vo.TodoVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &todoVO)
	if err != nil {
		panic(err)
	}
	//tokenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	tokenHead, _ := beego.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	t := common.Jwx.ParseToken(tokenString)
	a, _ := t.Get("adminUid")
	adminUid := a.(string)
	var todo models.Todo
	common.DB.Where("uid=?", todoVO.Uid).Find(&todo)
	if todo.AdminUid != adminUid {
		c.Result("error", "该数据无权限访问")
		return
	}
	todo.Status = 0
	common.DB.Save(&todo)
	c.Result("success", "删除成功")
}
