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
 * @date  2020/12/22 3:21 下午
 * @version 1.0
 */

type SystemConfigRestApi struct {
	base.BaseController
}

func (c *SystemConfigRestApi) GetSystemConfig() {
	result := service.SystemConfigService.GetConfig()
	c.Data["json"] = result
	err := c.ServeJSON()
	if err != nil {
		panic(err)
	}
}

func (c *SystemConfigRestApi) EditSystemConfig() {
	var systemConfigVO vo.SystemConfigVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &systemConfigVO)
	if err != nil {
		panic(err)
	}
	switch {
	case systemConfigVO.UploadLocal == "1" || systemConfigVO.UploadQiNiu == "0" || systemConfigVO.UploadMinio == "1":
		c.ErrorWithMessage("图片必须上传到一个区域,且当前版本不支持本地和Minio存储！")
		return
	case systemConfigVO.PicturePriority != "1":
		c.ErrorWithMessage("当前版本不支持本地和Minio存储！")
		return
	case systemConfigVO.StartEmailNotification == "1" && systemConfigVO.Email == "":
		c.ErrorWithMessage("开启邮件通知必须设置邮箱地址")
		return
	case systemConfigVO.Uid == "":
		editorModel := common.InterfaceToString(systemConfigVO.EditorModel)
		openDashboardNotification := common.InterfaceToString(systemConfigVO.OpenDashboardNotification)
		contentPicturePriority := common.InterfaceToString(systemConfigVO.ContentPicturePriority)
		systemConfig := models.SystemConfig{
			Uid:                       xid.New().String(),
			QiNiuAccessKey:            systemConfigVO.QiNiuAccessKey,
			QiNiuSecretKey:            systemConfigVO.QiNiuSecretKey,
			Email:                     systemConfigVO.Email,
			EmailUserName:             systemConfigVO.EmailUserName,
			EmailPassword:             systemConfigVO.EmailPassword,
			SmtpAddress:               systemConfigVO.SmtpAddress,
			SmtpPort:                  systemConfigVO.SmtpPort,
			QiNiuBucket:               systemConfigVO.QiNiuBucket,
			QiNiuArea:                 systemConfigVO.QiNiuArea,
			UploadQiNiu:               systemConfigVO.UploadQiNiu,
			UploadLocal:               systemConfigVO.UploadLocal,
			PicturePriority:           systemConfigVO.PicturePriority,
			QiNiuPictureBaseUrl:       systemConfigVO.QiNiuPictureBaseUrl,
			LocalPictureBaseUrl:       systemConfigVO.LocalPictureBaseUrl,
			StartEmailNotification:    systemConfigVO.StartEmailNotification,
			EditorModel:               editorModel,
			ThemeColor:                systemConfigVO.ThemeColor,
			MinioEndPoint:             systemConfigVO.MinioEndPoint,
			MinioAccessKey:            systemConfigVO.MinioAccessKey,
			MinioSecretKey:            systemConfigVO.MinioSecretKey,
			MinioBucket:               systemConfigVO.MinioBucket,
			UploadMinio:               systemConfigVO.UploadMinio,
			MinioPictureBaseUrl:       systemConfigVO.MinioPictureBaseUrl,
			OpenDashboardNotification: openDashboardNotification,
			DashboardNotification:     systemConfigVO.DashboardNotification,
			ContentPicturePriority:    contentPicturePriority,
		}
		common.DB.Create(&systemConfig)
	default:
		var systemConfig models.SystemConfig
		common.DB.Where("uid = ?", systemConfigVO.Uid).Find(&systemConfig)
		if systemConfigVO.PicturePriority != systemConfig.PicturePriority {
			service.BlogService.DeleteRedisByBlog()
		}
		editorModel := common.InterfaceToString(systemConfigVO.EditorModel)
		openDashboardNotification := common.InterfaceToString(systemConfigVO.OpenDashboardNotification)
		contentPicturePriority := common.InterfaceToString(systemConfigVO.ContentPicturePriority)
		systemConfig.QiNiuAccessKey = systemConfigVO.QiNiuAccessKey
		systemConfig.QiNiuSecretKey = systemConfigVO.QiNiuSecretKey
		systemConfig.Email = systemConfigVO.Email
		systemConfig.EmailUserName = systemConfigVO.EmailUserName
		systemConfig.EmailPassword = systemConfigVO.EmailPassword
		systemConfig.SmtpAddress = systemConfigVO.SmtpAddress
		systemConfig.SmtpPort = systemConfigVO.SmtpPort
		systemConfig.QiNiuBucket = systemConfigVO.QiNiuBucket
		systemConfig.QiNiuArea = systemConfigVO.QiNiuArea
		systemConfig.UploadQiNiu = systemConfigVO.UploadQiNiu
		systemConfig.UploadLocal = systemConfigVO.UploadLocal
		systemConfig.PicturePriority = systemConfigVO.PicturePriority
		systemConfig.QiNiuPictureBaseUrl = systemConfigVO.QiNiuPictureBaseUrl
		systemConfig.LocalPictureBaseUrl = systemConfigVO.LocalPictureBaseUrl
		systemConfig.StartEmailNotification = systemConfigVO.StartEmailNotification
		systemConfig.EditorModel = editorModel
		systemConfig.ThemeColor = systemConfigVO.ThemeColor
		systemConfig.MinioEndPoint = systemConfigVO.MinioEndPoint
		systemConfig.MinioAccessKey = systemConfigVO.MinioAccessKey
		systemConfig.MinioSecretKey = systemConfigVO.MinioSecretKey
		systemConfig.MinioBucket = systemConfigVO.MinioBucket
		systemConfig.UploadMinio = systemConfigVO.UploadMinio
		systemConfig.MinioPictureBaseUrl = systemConfigVO.MinioPictureBaseUrl
		systemConfig.OpenDashboardNotification = openDashboardNotification
		systemConfig.DashboardNotification = systemConfigVO.DashboardNotification
		systemConfig.ContentPicturePriority = contentPicturePriority
		common.DB.Save(&systemConfig)
	}
	common.RedisUtil.Delete("SYSTEM_CONFIG")
	c.SuccessWithMessage("更新成功")
}

func (c *SystemConfigRestApi) CleanRedisByKey() {
	var key []string
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &key)
	if err != nil {
		panic(err)
	}
	if len(key) == 0 {
		c.ErrorWithMessage("操作失败")
		return
	}
	for _, item := range key {
		if item == "ALL" {
			keys := common.RedisUtil.Keys("*")
			common.RedisUtil.MultiDelete(keys)
		} else {
			keys := common.RedisUtil.Keys(item + "*")
			common.RedisUtil.MultiDelete(keys)
		}
	}
	c.SuccessWithMessage("删除成功")
}
