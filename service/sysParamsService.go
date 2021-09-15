package service

import (
	"mogu-go-v2/common"
	"mogu-go-v2/models"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/29 3:59 下午
 * @version 1.0
 */

type sysParamsService struct{}

func (sysParamsService) GetSysParamsValueByKey(paramsKey string) (string, bool) {
	redisKey := "SYSTEM_PARAMS:" + paramsKey
	paramsValue := common.RedisUtil.Get(redisKey)
	if paramsValue == "" {
		sysParams := sysParamsByKey(paramsKey)
		if sysParams == (models.SysParams{}) || sysParams.ParamsValue == "" {
			return "请先配置博客相关参数", false
		}
		paramsValue = sysParams.ParamsValue
		common.RedisUtil.Set(redisKey, paramsValue)
	}
	return paramsValue, true
}

func sysParamsByKey(paramsKey string) models.SysParams {
	var sysParams models.SysParams
	common.DB.Where("params_key = ? and status=?", paramsKey, 1).Last(&sysParams)
	return sysParams
}

var SysParamService = &sysParamsService{}
