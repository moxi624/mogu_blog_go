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
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/29 2:21 下午
 * @version 1.0
 */

type SysDictDataRestApi struct {
	base.BaseController
}

func (c *SysDictDataRestApi) GetListByDictType() {
	dictType := c.GetString("dictType")
	if dictType == "" {
		c.Result("error", "操作失败")
		return
	}
	jsonResult := common.RedisUtil.Get("REDIS_DICT_TYPE:" + dictType)
	if jsonResult != "" {
		m := map[string]interface{}{}
		err := json.Unmarshal([]byte(jsonResult), &m)
		if err != nil {
			panic(err)
		}
		c.Result("success", m)
		return
	}
	var sysDictType models.SysDictType
	common.DB.Where("dict_type=? and status=? and is_publish=?", dictType, 1, "1").Limit(1).Find(&sysDictType)
	if sysDictType == (models.SysDictType{}) {
		c.Result("success", map[string]interface{}{})
		return
	}
	var sysDictDataList []models.SysDictData
	common.DB.Where("is_publish=? and status=? and dict_type_uid=?", "1", 1, sysDictType.Uid).Order("sort desc,create_time desc").Find(&sysDictDataList)
	defaultValue := ""
	for _, sysDictDATA := range sysDictDataList {
		if sysDictDATA.IsDefault == 1 {
			defaultValue = sysDictDATA.DictValue
			break
		}
	}
	result := make(map[string]interface{})
	result["defaultValue"] = defaultValue
	result["list"] = sysDictDataList
	b, _ := json.Marshal(result)
	common.RedisUtil.SetEx("REDIS_DICT_TYPE:"+dictType, string(b), 24, time.Hour)
	c.Result("success", result)
}

func (c *SysDictDataRestApi) GetListByDictTypeList() {
	var dictTypeList []string
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &dictTypeList)
	if err != nil {
		panic(err)
	}
	m := map[string]map[string]interface{}{}
	var tempTypeList []string
	for _, item := range dictTypeList {
		jsonResult := common.RedisUtil.Get("REDIS_DICT_TYPE:" + item)
		if jsonResult != "" {
			tempMap := map[string]interface{}{}
			err := json.Unmarshal([]byte(jsonResult), &tempMap)
			if err != nil {
				panic(err)
			}
			m[item] = tempMap
		} else {
			tempTypeList = append(tempTypeList, item)
		}
	}
	if len(tempTypeList) == 0 {
		c.Result("success", m)
		return
	}
	var sysDictTypeList []models.SysDictType
	common.DB.Where("dict_type in ? and status=? and is_publish=?", tempTypeList, 1, "1").Find(&sysDictTypeList)
	for _, item := range sysDictTypeList {
		var list []models.SysDictData
		common.DB.Where("is_publish=? and status=? and dict_type_uid=?", "1", 1, item.Uid).Order("sort desc,create_time desc").Find(&list)
		defaultValue := ""
		for _, sysDictData := range list {
			if sysDictData.IsDefault == 1 {
				defaultValue = sysDictData.DictValue
				break
			}
		}
		result := map[string]interface{}{
			"defaultValue": defaultValue,
			"list":         list,
		}
		m[item.DictType] = result
		b, _ := json.Marshal(result)
		common.RedisUtil.SetEx("REDIS_DICT_TYPE:"+item.DictType, string(b), 24, time.Hour)
	}
	c.Result("success", m)
}

func (c *SysDictDataRestApi) GetList() {
	var sysDictDataVO vo.SysDictDataVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysDictDataVO)
	if err != nil {
		panic(err)
	}
	base.L.Print("获取字典数据列表")
	where := "status=?"
	if sysDictDataVO.DictTypeUid != "" {
		where += " and dict_type_uid= '" + sysDictDataVO.DictTypeUid + "'"
	}
	if strings.TrimSpace(sysDictDataVO.DictLabel) != "" {
		where += " and dict_label like '%" + strings.TrimSpace(sysDictDataVO.DictLabel) + "%'"
	}
	var pageList []models.SysDictData
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.SysDictData{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((sysDictDataVO.CurrentPage - 1) * sysDictDataVO.PageSize).Limit(sysDictDataVO.PageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var dictDataUidList []string
	for _, item := range pageList {
		dictDataUidList = append(dictDataUidList, item.DictTypeUid)
	}
	var dictTypeList []models.SysDictType
	if len(dictDataUidList) > 0 {
		dictDataUidList = common.RemoveRepByMap(dictDataUidList)
		common.DB.Find(&dictTypeList, dictDataUidList)
	}
	dictTypeMap := map[string]models.SysDictType{}
	for _, item := range dictTypeList {
		dictTypeMap[item.Uid] = item
	}
	for i, item := range pageList {
		pageList[i].SysDictType = dictTypeMap[item.DictTypeUid]
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    sysDictDataVO.PageSize,
		Current: sysDictDataVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *SysDictDataRestApi) Edit() {
	var sysDictDataVO vo.SysDictDataVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysDictDataVO)
	if err != nil {
		panic(err)
	}
	adminUid := c.GetAdminUid()
	var sysDictData models.SysDictData
	common.DB.Where("uid=?", sysDictDataVO.Uid).Find(&sysDictData)
	var temp models.SysDictData
	if sysDictData.DictLabel != sysDictDataVO.DictLabel {
		common.DB.Where("dict_label=? and dict_type_uid=? and status=?", sysDictDataVO.DictLabel, sysDictDataVO.DictTypeUid, 1).Last(&temp)
		if temp != (models.SysDictData{}) {
			c.ErrorWithMessage("记录已经存在")
			return
		}
	}
	sysDictData.UpdateByUid = adminUid
	sysDictData.Oid = sysDictDataVO.Oid
	sysDictData.DictTypeUid = sysDictDataVO.DictTypeUid
	sysDictData.DictLabel = sysDictDataVO.DictLabel
	sysDictData.DictValue = sysDictDataVO.DictValue
	sysDictData.CssClass = sysDictDataVO.CssClass
	sysDictData.ListClass = sysDictDataVO.ListClass
	sysDictData.IsDefault = sysDictDataVO.IsDefault
	sysDictData.Remark = sysDictDataVO.Remark
	sysDictData.IsPublish = sysDictDataVO.IsPublish
	sysDictData.Sort = sysDictDataVO.Sort
	common.DB.Save(&sysDictData)
	keys := common.RedisUtil.Keys("REDIS_DICT_TYPE:*")
	common.RedisUtil.MultiDelete(keys)
	c.SuccessWithMessage("更新成功")
}

func (c *SysDictDataRestApi) DeleteBatch() {
	var sysDictDataVOList []vo.SysDictDataVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysDictDataVOList)
	if err != nil {
		panic(err)
	}
	adminUid := c.GetAdminUid()
	if len(sysDictDataVOList) <= 0 {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	var uids []string
	for _, item := range sysDictDataVOList {
		uids = append(uids, item.Uid)
	}
	var sysDictDataList []models.SysDictData
	common.DB.Find(&sysDictDataList, uids)
	save := common.DB.Model(&sysDictDataList).Select("status", "update_by_uid").Updates(models.SysDictData{Status: 0, UpdateByUid: adminUid}).Error
	keys := common.RedisUtil.Keys("REDIS_DICT_TYPE:*")
	common.RedisUtil.MultiDelete(keys)
	if save == nil {
		c.SuccessWithMessage("删除成功")
	} else {
		c.ErrorWithMessage("删除失败")
	}
}

func (c *SysDictDataRestApi) Add() {
	var sysDictDataVO vo.SysDictDataVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysDictDataVO)
	if err != nil {
		panic(err)
	}
	adminUid := c.GetAdminUid()
	var temp models.SysDictData
	common.DB.Where("dict_label=? and dict_type_uid=? and status=?", sysDictDataVO.DictLabel, sysDictDataVO.DictTypeUid, 1).Last(&temp)
	if temp != (models.SysDictData{}) {
		c.ErrorWithMessage("记录已经存在")
		return
	}
	sysDictData := models.SysDictData{
		Uid:         xid.New().String(),
		Oid:         sysDictDataVO.Oid,
		DictTypeUid: sysDictDataVO.DictTypeUid,
		DictLabel:   sysDictDataVO.DictLabel,
		DictValue:   sysDictDataVO.DictValue,
		CssClass:    sysDictDataVO.CssClass,
		ListClass:   sysDictDataVO.ListClass,
		IsDefault:   sysDictDataVO.IsDefault,
		CreateByUid: adminUid,
		UpdateByUid: adminUid,
		Remark:      sysDictDataVO.Remark,
		IsPublish:   sysDictDataVO.IsPublish,
		Sort:        sysDictDataVO.Sort,
	}
	common.DB.Create(&sysDictData)
	c.SuccessWithMessage("新增成功")
}
