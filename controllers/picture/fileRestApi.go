package picture

import (
	"github.com/beego/beego/v2/core/logs"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/service"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/31 2:09 下午
 * @version 1.0
 */

type FileRestApi struct {
	base.BaseController
}

var l = logs.GetLogger()

func (c *FileRestApi) GetPicture() {
	fileIds := c.GetString("fileIds")
	code := c.GetString("code")
	l.Print("获取图片信息:", fileIds)
	result := service.FileService.GetPicture(fileIds, code)
	c.Data["json"] = result
	err := c.ServeJSON()
	if err != nil {
		panic(err)
	}
}

func (c *FileRestApi) CropperPicture() {
	file, _ := c.GetFiles("file")
	systemConfig := c.GetSystemConfig()
	qiNiuPictureBaseUrl := systemConfig.QiNiuPictureBaseUrl
	localPictureBaseUrl := systemConfig.LocalPictureBaseUrl
	minioPictureBaseUrl := systemConfig.MinioPictureBaseUrl
	picMap := c.BatchUploadFile(file, systemConfig, "file")
	var listMap []map[string]interface{}
	if picMap["code"] == "success" {
		picData := picMap["data"].([]models.File)
		if len(picData) > 0 {
			for i := 0; i < len(picData); i++ {
				item := map[string]interface{}{"uid": picData[i].Uid}
				if systemConfig.PicturePriority == "1" {
					item["url"] = qiNiuPictureBaseUrl + picData[i].QiNiuUrl
				} else if systemConfig.PicturePriority == "2" {
					item["url"] = minioPictureBaseUrl + picData[i].MinioUrl
				} else {
					item["url"] = localPictureBaseUrl + picData[i].PicURL
				}
				listMap = append(listMap, item)
			}
		}
	}
	c.Result("success", listMap)
}

func (c *FileRestApi) UploadPics() {
	fileDatas, _ := c.GetFiles("filedatas")
	systemConfig := c.GetSystemConfig()
	m := c.BatchUploadFile(fileDatas, systemConfig, "filedatas")
	c.Data["json"] = m
	err := c.ServeJSON()
	if err != nil {
		panic(err)
	}
}
