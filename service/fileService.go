package service

import (
	"github.com/beego/beego/v2/core/logs"
	"mogu-go-v2/common"
	"mogu-go-v2/models"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/21 2:24 下午
 * @version 1.0
 */
type fileService struct{}

var l = logs.GetLogger()

func (fileService) GetPicture(fileIds string, code string) map[string]interface{} {
	if code == "" {
		code = ","
	}
	if fileIds == "" {
		l.Println("图片UID不能为空")
		e := map[string]interface{}{
			"code": "error",
			"data": "图片UID不能为空"}
		return e
	}
	var list []map[string]interface{}
	changeStringToString := strings.Split(fileIds, code)
	var fileList []models.File
	common.DB.Find(&fileList, changeStringToString)
	if len(fileList) > 0 {
		for _, file := range fileList {
			if file != (models.File{}) {
				remap := map[string]interface{}{}
				remap["qiNiuUrl"] = file.QiNiuUrl
				remap["minioUrl"] = file.MinioUrl
				remap["url"] = file.PicURL
				remap["expandedName"] = file.PicExpandedName
				remap["name"] = file.PicName
				remap["uid"] = file.Uid
				remap["fileOldName"] = file.FileOldName
				list = append(list, remap)
			}
		}
	}
	e := map[string]interface{}{
		"code": "success",
		"data": list}
	return e
}

var FileService = &fileService{}
