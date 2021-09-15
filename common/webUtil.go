package common

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"mogu-go-v2/models"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/21 4:01 下午
 * @version 1.0
 */
type webUtil struct{}

var l = logs.GetLogger()

func (webUtil) GetPicture(result map[string]interface{}) []string {
	picturePriority := ""
	localPictureBaseUrl := ""
	qiNiuPictureBaseUrl := ""
	minioPictureBaseUrl := ""
	systemConfigJson := RedisUtil.Get("SYSTEM_CONFIG")
	var systemConfig models.SystemConfig
	if systemConfigJson == "" {
		DB.Where("status = ?", 1).First(&systemConfig)
		if systemConfig == (models.SystemConfig{}) {
			l.Println("系统配置不存在，请检查t_system_config表是否有数据，并重新导入数据库")
		} else {
			j, _ := json.Marshal(systemConfig)
			RedisUtil.SetEx("SYSTEM_CONFIG", string(j), 24, time.Hour)
		}
		picturePriority = systemConfig.PicturePriority
		localPictureBaseUrl = systemConfig.LocalPictureBaseUrl
		qiNiuPictureBaseUrl = systemConfig.QiNiuPictureBaseUrl
		minioPictureBaseUrl = systemConfig.MinioPictureBaseUrl
	} else {
		err := json.Unmarshal([]byte(systemConfigJson), &systemConfig)
		if err != nil {
			l.Println("系统配置转换错误，请检查系统配置，或者清空Redis后重试！")
		}
		picturePriority = systemConfig.PicturePriority
		localPictureBaseUrl = systemConfig.LocalPictureBaseUrl
		qiNiuPictureBaseUrl = systemConfig.QiNiuPictureBaseUrl
		minioPictureBaseUrl = systemConfig.MinioPictureBaseUrl
	}
	var picUrls []string
	if result["code"] == "success" {
		picData := result["data"].([]map[string]interface{})
		if len(picData) > 0 {
			for i := 0; i < len(picData); i++ {
				if picturePriority == "1" {
					picUrls = append(picUrls, qiNiuPictureBaseUrl+picData[i]["qiNiuUrl"].(string))
				} else if picturePriority == "2" {
					picUrls = append(picUrls, minioPictureBaseUrl+picData[i]["minioUrl"].(string))
				} else {
					picUrls = append(picUrls, localPictureBaseUrl+picData[i]["url"].(string))
				}
			}
		}
	}
	return picUrls
}

func (webUtil) GetPictureMap(result map[string]interface{}) []map[string]interface{} {
	picturePriority := ""
	localPictureBaseUrl := ""
	qiNiuPictureBaseUrl := ""
	minioPictureBaseUrl := ""
	systemConfigJson := RedisUtil.Get("SYSTEM_CONFIG")
	var systemConfig models.SystemConfig
	if systemConfigJson == "" {
		DB.Where("status = ?", 1).First(&systemConfig)
		if systemConfig == (models.SystemConfig{}) {
			l.Println("系统配置不存在，请检查t_system_config表是否有数据，并重新导入数据库")
		} else {
			j, _ := json.Marshal(systemConfig)
			RedisUtil.SetEx("SYSTEM_CONFIG", string(j), 24, time.Hour)
		}
		picturePriority = systemConfig.PicturePriority
		localPictureBaseUrl = systemConfig.LocalPictureBaseUrl
		qiNiuPictureBaseUrl = systemConfig.QiNiuPictureBaseUrl
		minioPictureBaseUrl = systemConfig.MinioPictureBaseUrl
	} else {
		err := json.Unmarshal([]byte(systemConfigJson), &systemConfig)
		if err != nil {
			l.Println("系统配置转换错误，请检查系统配置，或者清空Redis后重试！")
		}
		picturePriority = systemConfig.PicturePriority
		localPictureBaseUrl = systemConfig.LocalPictureBaseUrl
		qiNiuPictureBaseUrl = systemConfig.QiNiuPictureBaseUrl
		minioPictureBaseUrl = systemConfig.MinioPictureBaseUrl
	}
	var resultList []map[string]interface{}
	if result["code"] == "success" {
		picData := result["data"].([]map[string]interface{})
		if len(picData) > 0 {
			for i := 0; i < len(picData); i++ {
				m := map[string]interface{}{}
				if picData[i]["uid"] == "" {
					continue
				}
				if picturePriority == "1" {
					m["url"] = qiNiuPictureBaseUrl + picData[i]["qiNiuUrl"].(string)
				} else if picturePriority == "2" {
					m["url"] = minioPictureBaseUrl + picData[i]["minioUrl"].(string)
				} else {
					m["url"] = localPictureBaseUrl + picData[i]["url"].(string)
				}
				m["uid"] = picData[i]["uid"]
				resultList = append(resultList, m)
			}
		}
	}
	return resultList
}

var WebUtil = &webUtil{}
