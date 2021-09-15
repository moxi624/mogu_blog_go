//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package service

import (
	"encoding/json"
	"mogu-go-v2/common"
	"mogu-go-v2/models"
	"reflect"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/19 4:08 下午
 * @version 1.0
 */

type webConfigService struct{}

func (*webConfigService) IsOpenLoginType(loginType string) (bool, string) {
	loginTypeJson := common.RedisUtil.Get("LOGIN_TYPE:" + loginType)
	if loginTypeJson != "" {
		return true, ""
	}
	var webConfig models.WebConfig
	common.DB.Order("create_time desc").First(&webConfig)
	if reflect.DeepEqual(webConfig, models.WebConfig{}) {
		return false, "系统配置不存在"
	}
	loginTypeListJson := webConfig.LoginTypeList
	var loginTypeList []string
	err := json.Unmarshal([]byte(loginTypeListJson), &loginTypeList)
	if err != nil {
		panic(err)
	}
	for _, item := range loginTypeList {
		switch item {
		case "1":
			common.RedisUtil.Set("LOGIN_TYPE:PASSWORD", "PASSWORD")
		case "2":
			common.RedisUtil.Set("LOGIN_TYPE:GITEE", "GITEE")
		case "3":
			common.RedisUtil.Set("LOGIN_TYPE:GITHUB", "GITHUB")
		case "4":
			common.RedisUtil.Set("LOGIN_TYPE:QQ", "QQ")
		case "5":
			common.RedisUtil.Set("LOGIN_TYPE:WECHAT", "WECHAT")
		}
	}
	loginTypeJson = common.RedisUtil.Get("LOGIN_TYPE:" + loginType)
	if loginTypeJson != "" {
		return true, ""
	} else {
		common.RedisUtil.Set("LOGIN_TYPE", "")
		return false, ""
	}
}

var WebConfigService = &webConfigService{}
