package admin

import (
	"encoding/json"
	"github.com/rs/xid"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/vo"
	"reflect"
	"sort"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/27 8:33 上午
 * @version 1.0
 */

type CategoryMenuRestApi struct {
	base.BaseController
}

func (c *CategoryMenuRestApi) GetAll() {
	keyword := c.GetString("keyword")
	where := "menu_level=? and status=? and menu_type=?"
	if keyword != "" {
		where += " and uid='" + keyword + "'"
	}
	var list []models.CategoryMenu
	common.DB.Where(where, "1", 1, 0).Order("sort desc").Find(&list)
	var ids []string
	for _, item := range list {
		if item.Uid != "" {
			ids = append(ids, item.Uid)
		}
	}
	var childList []models.CategoryMenu
	common.DB.Where("status=? and parent_uid in ?", 1, ids).Find(&childList)
	var secondMenuUids []string
	for _, item := range childList {
		if item.Uid != "" {
			secondMenuUids = append(secondMenuUids, item.Uid)
		}
	}
	var buttonList []models.CategoryMenu
	common.DB.Where("status=? and parent_uid in ?", 1, secondMenuUids).Find(&buttonList)
	m := map[string][]models.CategoryMenu{}
	for _, item := range buttonList {
		var tempList []models.CategoryMenu
		tempList = append(tempList, item)
		m[item.ParentUid] = tempList
	}
	for i, item := range childList {
		if len(m[item.Uid]) > 0 {
			tempList := m[item.Uid]
			sort.Slice(tempList, func(i, j int) bool {
				return tempList[i].Sort > tempList[j].Sort
			})
			childList[i].ChildCategoryMenu = tempList
		}
	}
	for i, parentItem := range list {
		var tempList []models.CategoryMenu
		for _, item := range childList {
			if item.ParentUid == parentItem.Uid {
				tempList = append(tempList, item)
			}
		}
		sort.SliceStable(tempList, func(i, j int) bool {
			return tempList[i].Sort > tempList[j].Sort
		})
		list[i].ChildCategoryMenu = tempList
	}
	c.SuccessWithData(list)
}

func (c *CategoryMenuRestApi) GetButtonAll() {
	keyword := c.GetString("keyword")
	where := "menu_level=? and status=? and menu_type=?"
	if keyword != "" {
		where += " and uid='" + keyword + "'"
	}
	var list []models.CategoryMenu
	common.DB.Where(where, "2", 1, 0).Order("sort desc").Find(&list)
	var ids []string
	for _, item := range list {
		if item.Uid != "" {
			ids = append(ids, item.Uid)
		}
	}
	var childList []models.CategoryMenu
	common.DB.Where("status=? and parent_uid in ?", 1, ids).Find(&childList)
	var secondUidSet []string
	m := map[string][]models.CategoryMenu{}
	for _, item := range ids {
		var tempList []models.CategoryMenu
		for _, v := range childList {
			if v.ParentUid == item {
				tempList = append(tempList, v)
			}
		}
		m[item] = tempList
	}
	secondUidSet = common.RemoveRepByMap(secondUidSet)
	for i, item := range list {
		if !reflect.DeepEqual(m[item.Uid], models.CategoryMenu{}) {
			tempList := m[item.Uid]
			sort.SliceStable(tempList, func(i, j int) bool {
				return tempList[i].Sort > tempList[j].Sort
			})
			list[i].ChildCategoryMenu = tempList
		}
	}
	c.SuccessWithData(list)
}

func (c *CategoryMenuRestApi) Stick() {
	var categoryMenuVO vo.CategoryMenuVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &categoryMenuVO)
	if err != nil {
		panic(err)
	}
	var categoryMenu models.CategoryMenu
	common.DB.Where("uid=?", categoryMenuVO.Uid).Find(&categoryMenu)
	where := "menu_level=?"
	if categoryMenu.MenuLevel == 2 || categoryMenu.MenuType == 1 {
		where += " and parent_uid='" + categoryMenu.ParentUid + "'"
	}
	var maxSort models.CategoryMenu
	common.DB.Where(where, categoryMenu.MenuLevel).Order("sort desc").Last(&maxSort)
	if maxSort.Uid == "" {
		c.ErrorWithMessage("操作失败")
		return
	}
	sortCount := maxSort.Sort + 1
	categoryMenu.Sort = sortCount
	common.DB.Save(&categoryMenu)
	c.SuccessWithMessage("操作成功")
}

func (c *CategoryMenuRestApi) Edit() {
	var categoryMenuVO vo.CategoryMenuVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &categoryMenuVO)
	if err != nil {
		panic(err)
	}
	var categoryMenu models.CategoryMenu
	common.DB.Where("uid=?", categoryMenuVO.Uid).Find(&categoryMenu)
	categoryMenu.ParentUid = categoryMenuVO.ParentUid
	categoryMenu.Sort = categoryMenuVO.Sort
	categoryMenu.Icon = categoryMenuVO.Icon
	categoryMenu.Summary = categoryMenuVO.Summary
	categoryMenu.MenuLevel = categoryMenuVO.MenuLevel
	categoryMenu.MenuType = categoryMenuVO.MenuType
	categoryMenu.Name = categoryMenuVO.Name
	categoryMenu.Url = categoryMenuVO.Url
	categoryMenu.IsShow = categoryMenuVO.IsShow
	categoryMenu.IsJumpExternalUrl = categoryMenuVO.IsJumpExternalUrl
	common.DB.Save(&categoryMenu)
	deleteAdminVisitUrl()
	c.SuccessWithMessage("更新成功")
}

func (c *CategoryMenuRestApi) Delete() {
	var categoryMenuVO vo.CategoryMenuVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &categoryMenuVO)
	if err != nil {
		panic(err)
	}
	var categoryMenu models.CategoryMenu
	var menuCount int64
	common.DB.Model(&categoryMenu).Where("status=? and parent_uid=?", 1, categoryMenuVO.Uid).Count(&menuCount)
	if menuCount > 0 {
		c.ErrorWithMessage("该菜单下还有子菜单")
		return
	}
	common.DB.Where("uid=?", categoryMenuVO.Uid).Find(&categoryMenu)
	common.DB.Model(&categoryMenu).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}

func (c *CategoryMenuRestApi) Add() {
	var categoryMenuVO vo.CategoryMenuVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &categoryMenuVO)
	if err != nil {
		panic(err)
	}
	if categoryMenuVO.MenuLevel == 1 {
		categoryMenuVO.ParentUid = ""
	}
	categoryMenu := models.CategoryMenu{
		ParentUid:         categoryMenuVO.ParentUid,
		Sort:              categoryMenuVO.Sort,
		Icon:              categoryMenuVO.Icon,
		Summary:           categoryMenuVO.Summary,
		MenuLevel:         categoryMenuVO.MenuLevel,
		MenuType:          categoryMenuVO.MenuType,
		Name:              categoryMenuVO.Name,
		Url:               categoryMenuVO.Url,
		IsShow:            categoryMenuVO.IsShow,
		IsJumpExternalUrl: categoryMenuVO.IsJumpExternalUrl,
		Uid:               xid.New().String(),
	}
	common.DB.Create(&categoryMenu)
	c.SuccessWithMessage("新增成功")
}
