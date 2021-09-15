package service

import (
	"mogu-go-v2/common"
	"mogu-go-v2/models"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/4 4:32 下午
 * @version 1.0
 */

type systemConfigfService struct{}

func (systemConfigfService) GetConfig() map[string]interface{} {
	var systemConfig models.SystemConfig
	common.DB.Where("status=?", 1).Order("create_time desc").Limit(1).First(&systemConfig)
	r := map[string]interface{}{
		"code": "success",
		"data": systemConfig,
	}
	return r
}

var SystemConfigService = &systemConfigfService{}
