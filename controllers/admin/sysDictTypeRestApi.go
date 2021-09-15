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
 * @date  2021/1/29 8:31 上午
 * @version 1.0
 */

type SysDictTypeRestApi struct {
	base.BaseController
}

func (c *SysDictTypeRestApi) GetList() {
	var sysDictTypeVO vo.SysDictTypeVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysDictTypeVO)
	if err != nil {
		panic(err)
	}
	base.L.Print("获取字典列表参数")
	where := "status=?"
	if strings.TrimSpace(sysDictTypeVO.DictName) != "" {
		where += " and dict_name like '%" + strings.TrimSpace(sysDictTypeVO.DictName) + "%'"
	}
	if strings.TrimSpace(sysDictTypeVO.DictType) != "" {
		where += " and dict_type like '%" + strings.TrimSpace(sysDictTypeVO.DictType) + "%'"
	}
	var pageList []models.SysDictType
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.SysDictType{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((sysDictTypeVO.CurrentPage - 1) * sysDictTypeVO.PageSize).Limit(sysDictTypeVO.PageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    sysDictTypeVO.PageSize,
		Current: sysDictTypeVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *SysDictTypeRestApi) Edit() {
	var sysDictTypeVO vo.SysDictTypeVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysDictTypeVO)
	if err != nil {
		panic(err)
	}
	var sysDictType models.SysDictType
	common.DB.Where("uid = ?", sysDictTypeVO.Uid).Find(&sysDictType)
	var temp models.SysDictType
	if sysDictType.DictType != sysDictTypeVO.DictType {
		common.DB.Where("dict_type=? and status=?", sysDictTypeVO.DictType, 1).Last(&temp)
		if temp != (models.SysDictType{}) {
			c.ErrorWithMessage("数据已经存在")
			return
		}
	}
	sysDictType.DictName = sysDictTypeVO.DictName
	sysDictType.DictType = sysDictTypeVO.DictType
	sysDictType.Remark = sysDictTypeVO.Remark
	sysDictType.IsPublish = sysDictTypeVO.IsPublish
	sysDictType.Sort = sysDictTypeVO.Sort
	sysDictType.UpdateByUid = c.GetAdminUid()
	common.DB.Save(&sysDictType)
	keys := common.RedisUtil.Keys("REDIS_DICT_TYPE:*")
	common.RedisUtil.MultiDelete(keys)
	c.SuccessWithMessage("更新成功")
}

func (c *SysDictTypeRestApi) DeleteBatch() {
	var sysDictTypeVOList []vo.SysDictTypeVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysDictTypeVOList)
	if err != nil {
		panic(err)
	}
	adminUid := c.GetAdminUid()
	if len(sysDictTypeVOList) == 0 {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	var uids []string
	for _, item := range sysDictTypeVOList {
		uids = append(uids, item.Uid)
	}
	var count int64
	common.DB.Model(&models.SysDictData{}).Where("status=? and dict_type_uid in ?", 1, uids).Count(&count)
	if count > 0 {
		c.ErrorWithMessage("该分类下还有字典数据")
		return
	}
	var sysDictTypeList []models.SysDictType
	common.DB.Find(&sysDictTypeList, uids)
	save := common.DB.Model(&sysDictTypeList).Select("status", "update_by_uid").Updates(models.SysDictData{Status: 0, UpdateByUid: adminUid}).Error
	keys := common.RedisUtil.Keys("REDIS_DICT_TYPE:*")
	common.RedisUtil.MultiDelete(keys)
	if save == nil {
		c.SuccessWithMessage("删除成功")
	} else {
		c.ErrorWithMessage("删除失败")
	}
}

func (c *SysDictTypeRestApi) Add() {
	var sysDictTypeVO vo.SysDictTypeVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysDictTypeVO)
	if err != nil {
		panic(err)
	}
	var temp models.SysDictType
	common.DB.Where("dict_type=? and status=?", sysDictTypeVO.DictType, 1).Last(&temp)
	if temp != (models.SysDictType{}) {
		c.ErrorWithMessage("数据已存在")
		return
	}
	adminUid := c.GetAdminUid()
	sysDictType := models.SysDictType{
		Uid:         xid.New().String(),
		DictName:    sysDictTypeVO.DictName,
		DictType:    sysDictTypeVO.DictType,
		Remark:      sysDictTypeVO.Remark,
		IsPublish:   sysDictTypeVO.IsPublish,
		Sort:        sysDictTypeVO.Sort,
		CreateByUid: adminUid,
		UpdateByUid: adminUid,
	}
	common.DB.Create(&sysDictType)
	c.SuccessWithMessage("新增成功")
}
