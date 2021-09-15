package admin

import (
	"encoding/json"
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
 * @date  2021/1/26 9:14 上午
 * @version 1.0
 */

type RoleRestApi struct {
	base.BaseController
}

func (c *RoleRestApi) GetList() {
	base.L.Print("获取角色信息列表")
	var roleVO vo.RoleVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &roleVO)
	if err != nil {
		panic(err)
	}
	where := "status != ?"
	if strings.TrimSpace(roleVO.Keyword) != "" {
		where += " and role_name like \"%" + strings.TrimSpace(roleVO.Keyword) + "%\""
	}
	var pageList []models.Role
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Role{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 0).Offset((roleVO.CurrentPage - 1) * roleVO.PageSize).Limit(roleVO.PageSize).Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    roleVO.PageSize,
		Current: roleVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *RoleRestApi) Edit() {
	var roleVO vo.RoleVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &roleVO)
	if err != nil {
		panic(err)
	}
	uid := roleVO.Uid
	var getRole models.Role
	common.DB.Where("uid = ?", uid).Find(&getRole)
	if getRole == (models.Role{}) {
		c.ErrorWithMessage("传入参数错误")
		return
	}
	getRole.RoleName = roleVO.RoleName
	getRole.CategoryMenuUids = roleVO.CategoryMenuUids
	getRole.Summary = roleVO.Summary
	common.DB.Save(&getRole)
	deleteAdminVisitUrl()
	c.SuccessWithMessage("更新成功")
}

func deleteAdminVisitUrl() {
	keys := common.RedisUtil.Keys("ADMIN_VISIT_MENU*")
	common.RedisUtil.MultiDelete(keys)
}

func (c *RoleRestApi) Delete() {
	var roleVO vo.RoleVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &roleVO)
	if err != nil {
		panic(err)
	}
	var adminCount int64
	var admin models.Admin
	common.DB.Model(&admin).Where("status=? and role_uid=?", 1, roleVO.Uid).Count(&adminCount)
	if adminCount > 0 {
		c.ErrorWithMessage("该角色下还有管理员")
		return
	}
	var role models.Role
	common.DB.Where("uid = ?", roleVO.Uid).Find(&role)
	common.DB.Model(&role).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}

func (c *RoleRestApi) Add() {
	var roleVO vo.RoleVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &roleVO)
	if err != nil {
		panic(err)
	}
	roleName := roleVO.RoleName
	var getRole models.Role
	common.DB.Where("role_name=?", roleName).First(&getRole)
	if getRole == (models.Role{}) {
		var role models.Role
		role.RoleName = roleVO.RoleName
		role.CategoryMenuUids = roleVO.CategoryMenuUids
		role.Summary = roleVO.Summary
		role.Uid = xid.New().String()
		common.DB.Create(&role)
		c.SuccessWithMessage("新增成功")
		return
	}
	c.ErrorWithMessage("权限已存在")
}
