package admin

import (
	"encoding/json"
	"github.com/rs/xid"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"reflect"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/28 2:01 下午
 * @version 1.0
 */

type WebConfigRestApi struct {
	base.BaseController
}

func (c *WebConfigRestApi) GetWebConfig() {
	var webConfig models.WebConfig
	common.DB.Order("create_time desc").First(&webConfig)
	if !reflect.DeepEqual(webConfig, models.WebConfig{}) && webConfig.Logo != "" {
		pictureList := service.FileService.GetPicture(webConfig.Logo, ",")
		webConfig.PhotoList = common.WebUtil.GetPicture(pictureList)
	}
	if !reflect.DeepEqual(webConfig, models.WebConfig{}) && webConfig.AliPay != "" {
		pictureList := service.FileService.GetPicture(webConfig.AliPay, ",")
		if len(pictureList) > 0 {
			webConfig.AliPayPhoto = common.WebUtil.GetPicture(pictureList)[0]
		}
	}
	if !reflect.DeepEqual(webConfig, models.WebConfig{}) && webConfig.WeixinPay != "" {
		pictureList := service.FileService.GetPicture(webConfig.WeixinPay, ",")
		if len(pictureList) > 0 {
			webConfig.WeixinPayPhoto = common.WebUtil.GetPicture(pictureList)[0]
		}
	}
	c.SuccessWithData(webConfig)
}

func (c *WebConfigRestApi) EditWebConfig() {
	var webConfigVO vo.WebConfigVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &webConfigVO)
	if err != nil {
		panic(err)
	}
	if webConfigVO.Uid == "" {
		webConfig := models.WebConfig{
			Uid:                  xid.New().String(),
			Logo:                 webConfigVO.Logo,
			Name:                 webConfigVO.Name,
			Summary:              webConfigVO.Summary,
			Keyword:              webConfigVO.Keyword,
			Author:               webConfigVO.Author,
			RecordNum:            webConfigVO.RecordNum,
			Title:                webConfigVO.Title,
			AliPay:               webConfigVO.AliPay,
			WeixinPay:            webConfigVO.WeixinPay,
			OpenComment:          webConfigVO.OpenComment,
			OpenMobileComment:    webConfigVO.OpenMobileComment,
			OpenAdmiration:       webConfigVO.OpenAdmiration,
			OpenMobileAdmiration: webConfigVO.OpenMobileAdmiration,
			Github:               webConfigVO.Github,
			Gitee:                webConfigVO.Gitee,
			QqNumber:             webConfigVO.QqNumber,
			QqGroup:              webConfigVO.QqGroup,
			WeChat:               webConfigVO.WeChat,
			Email:                webConfigVO.Email,
			ShowList:             webConfigVO.ShowList,
			LoginTypeList:        webConfigVO.LoginTypeList,
			PhotoList:            webConfigVO.PhotoList,
			AliPayPhoto:          webConfigVO.AliPayPhoto,
			WeixinPayPhoto:       webConfigVO.WeixinPayPhoto,
		}
		common.DB.Create(&webConfig)
	} else {
		var webConfig models.WebConfig
		common.DB.Where("uid = ?", webConfigVO.Uid).Find(&webConfig)
		webConfig.Logo = webConfigVO.Logo
		webConfig.Name = webConfigVO.Name
		webConfig.Summary = webConfigVO.Summary
		webConfig.Keyword = webConfigVO.Keyword
		webConfig.Author = webConfigVO.Author
		webConfig.RecordNum = webConfigVO.RecordNum
		webConfig.Title = webConfigVO.Title
		webConfig.AliPay = webConfigVO.AliPay
		webConfig.WeixinPay = webConfigVO.WeixinPay
		webConfig.OpenComment = webConfigVO.OpenComment
		webConfig.OpenMobileComment = webConfigVO.OpenMobileComment
		webConfig.OpenAdmiration = webConfigVO.OpenAdmiration
		webConfig.OpenMobileAdmiration = webConfigVO.OpenMobileAdmiration
		webConfig.Github = webConfigVO.Github
		webConfig.Gitee = webConfigVO.Gitee
		webConfig.QqNumber = webConfigVO.QqNumber
		webConfig.QqGroup = webConfigVO.QqGroup
		webConfig.WeChat = webConfigVO.WeChat
		webConfig.Email = webConfigVO.Email
		webConfig.ShowList = webConfigVO.ShowList
		webConfig.LoginTypeList = webConfigVO.LoginTypeList
		webConfig.PhotoList = webConfigVO.PhotoList
		webConfig.AliPayPhoto = webConfigVO.AliPayPhoto
		webConfig.WeixinPayPhoto = webConfigVO.WeixinPayPhoto
		common.DB.Save(&webConfig)
	}
	common.RedisUtil.Delete("WEB_CONFIG")
	keySet := common.RedisUtil.Keys("LOGIN_TYPE*")
	common.RedisUtil.MultiDelete(keySet)
	c.SuccessWithMessage("更新成功")
}
