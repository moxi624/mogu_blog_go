//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package admin

import (
	"encoding/json"
	"github.com/rs/xid"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/3/11 9:02 上午
 * @version 1.0
 */

type WebNavbarRestApi struct {
	base.BaseController
}

func (c *WebNavbarRestApi) GetAllList() {
	c.SuccessWithData(service.WebNavBarService.GetAllList())
}

func (c *WebNavbarRestApi) Edit() {
	base.L.Print("编辑门户导航栏")
	var webNavbarVO vo.WebNavbarVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &webNavbarVO)
	if err != nil {
		panic(err)
	}
	if webNavbarVO.NavbarLevel == 1 {
		webNavbarVO.ParentUid = ""
	}
	var webNavbar models.WebNavbar
	common.DB.Where("uid=?", webNavbarVO.Uid).Find(&webNavbar)
	webNavbar.ParentUid = webNavbarVO.ParentUid
	webNavbar.NavbarLevel = webNavbarVO.NavbarLevel
	webNavbar.Sort = webNavbarVO.Sort
	webNavbar.Icon = webNavbarVO.Icon
	webNavbar.IsJumpExternalUrl = webNavbarVO.IsJumpExternalUrl
	webNavbar.Name = webNavbarVO.Name
	webNavbar.Url = webNavbarVO.Url
	common.DB.Save(&webNavbar)
	c.SuccessWithData("更新成功")
}

func (c *WebNavbarRestApi) Add() {
	base.L.Print("增加门户导航栏")
	var webNavbarVO vo.WebNavbarVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &webNavbarVO)
	if err != nil {
		panic(err)
	}
	if webNavbarVO.NavbarLevel == 1 {
		webNavbarVO.ParentUid = ""
	}
	webNavbar := models.WebNavbar{
		Uid:               xid.New().String(),
		ParentUid:         webNavbarVO.ParentUid,
		NavbarLevel:       webNavbarVO.NavbarLevel,
		Sort:              webNavbarVO.Sort,
		Icon:              webNavbarVO.Icon,
		IsJumpExternalUrl: webNavbarVO.IsJumpExternalUrl,
		Name:              webNavbarVO.Name,
		Url:               webNavbarVO.Url,
	}
	common.DB.Create(&webNavbar)
	c.SuccessWithMessage("新增成功")
}

func (c *WebNavbarRestApi) Delete() {
	base.L.Print("删除门户导航栏")
	var webNavbarVO vo.WebNavbarVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &webNavbarVO)
	if err != nil {
		panic(err)
	}
	var menuCount int64
	common.DB.Model(&models.WebNavbar{}).Where("status=? and parent_uid = ?", 1, webNavbarVO.Uid).Count(&menuCount)
	if menuCount > 0 {
		c.ErrorWithMessage("该菜单下还有子菜单")
		return
	}
	var webNavbar models.WebNavbar
	common.DB.Where("uid=?", webNavbarVO.Uid).Find(&webNavbar)
	common.DB.Model(&webNavbar).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}
