package admin

import (
	"encoding/json"
	"github.com/rs/xid"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/29 3:29 下午
 * @version 1.0
 */
type SysParamsRestApi struct {
	base.BaseController
}

func (c *SysParamsRestApi) GetList() {
	base.L.Print("获取参数配置列表")
	var sysParamsVO vo.SysParamsVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysParamsVO)
	if err != nil {
		panic(err)
	}
	where := "status=?"
	if sysParamsVO.ParamsName != "" {
		where += " and params_name like '%" + sysParamsVO.ParamsName + "%'"
	}
	if sysParamsVO.ParamsKey != "" {
		where += " and params_key like '%" + sysParamsVO.ParamsKey + "%'"
	}
	var pageList []models.SysParams
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.SysParams{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((sysParamsVO.CurrentPage - 1) * sysParamsVO.PageSize).Limit(sysParamsVO.PageSize).Order("sort desc, create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    sysParamsVO.PageSize,
		Current: sysParamsVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *SysParamsRestApi) Edit() {
	var sysParamsVO vo.SysParamsVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysParamsVO)
	if err != nil {
		panic(err)
	}
	var sysParams models.SysParams
	common.DB.Where("uid=?", sysParamsVO.Uid).Find(&sysParams)
	var temp models.SysParams
	if sysParamsVO.ParamsKey != sysParams.ParamsKey {
		common.DB.Where("params_key=? and status=?", sysParamsVO.ParamsKey, 1).Last(&temp)
		if temp != (models.SysParams{}) {
			c.ErrorWithMessage("记录不存在")
			return
		}
	}
	paramsType := common.InterfaceToString(sysParamsVO.ParamsType)
	sysParams.ParamsName = sysParamsVO.ParamsName
	sysParams.ParamsKey = sysParamsVO.ParamsKey
	sysParams.ParamsValue = sysParamsVO.ParamsValue
	sysParams.ParamsType = paramsType
	sysParams.Remark = sysParamsVO.Remark
	sysParams.Sort = sysParamsVO.Sort
	common.DB.Save(&sysParams)
	common.RedisUtil.Delete("SYSTEM_PARAMS:" + sysParamsVO.ParamsKey)
	c.SuccessWithMessage("更新成功")
}

func (c *SysParamsRestApi) DeleteBatch() {
	var sysParamsVOList []vo.SysParamsVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysParamsVOList)
	if err != nil {
		panic(err)
	}
	var sysParamsUidList []string
	for _, item := range sysParamsVOList {
		sysParamsUidList = append(sysParamsUidList, item.Uid)
	}
	if len(sysParamsUidList) > 0 {
		var sysParamsList []models.SysParams
		common.DB.Find(&sysParamsList, sysParamsUidList)
		var redisKeys []string
		for _, item := range sysParamsList {
			redisKeys = append(redisKeys, "SYSTEM_PARAMS:"+item.ParamsKey)
			if item.ParamsType == "1" {
				c.ThrowError("0", "内置参数不能删除")
				return
			}
		}
		if len(sysParamsUidList) > 0 {
			common.DB.Model(&sysParamsList).Select("status").Update("status", 0)
			common.RedisUtil.MultiDelete(redisKeys)
			c.SuccessWithMessage("删除成功")
			return
		}
	} else {
		c.ErrorWithMessage("删除失败")
	}
}

func (c *SysParamsRestApi) Add() {
	var sysParamsVO vo.SysParamsVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysParamsVO)
	if err != nil {
		panic(err)
	}
	var temp models.SysParams
	common.DB.Where("params_key=? and status=?", sysParamsVO.ParamsKey, 1).Last(&temp)
	if temp != (models.SysParams{}) {
		c.ErrorWithMessage("数据已存在")
		return
	}
	paramsType := common.InterfaceToString(sysParamsVO.ParamsType)
	sysParams := models.SysParams{
		Uid:         xid.New().String(),
		ParamsName:  sysParamsVO.ParamsName,
		ParamsKey:   sysParamsVO.ParamsKey,
		ParamsType:  paramsType,
		ParamsValue: sysParamsVO.ParamsValue,
		Remark:      sysParamsVO.Remark,
		Sort:        sysParamsVO.Sort,
	}
	common.DB.Create(&sysParams)
	c.SuccessWithMessage("新增成功")
}
