package service

import (
	"mogu-go-v2/common"
	"mogu-go-v2/models"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/4 3:18 下午
 * @version 1.0
 */
type authService struct{}

func (authService) GetSystemConfig(token string) map[string]interface{} {
	userInfo := common.RedisUtil.Get("USER_TOKEN:" + token)
	if userInfo == "" {
		r := map[string]interface{}{
			"code": "error",
			"data": "token令牌未被识别",
		}
		return r
	}
	var systemConfig models.SystemConfig
	common.DB.Where("status=?", 1).Order("create_time desc").Limit(1).First(&systemConfig)
	r := map[string]interface{}{
		"code": "success",
		"data": systemConfig,
	}
	return r
}

var AuthService = &authService{}
