package picture

import (
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/14 1:20 下午
 * @version 1.0
 */

type StorageRestApi struct {
	base.BaseController
}

func (c *StorageRestApi) GetStorage() {
	c.CheckLogin()
	storage := c.GetStorageByAdmin()
	c.SuccessWithData(storage)
}

func (c *StorageRestApi) UploadFile() {
	adminUid := c.CheckLogin()
	filepath := c.GetString("filePath")
	fileDatas, _ := c.GetFiles("filedatas")
	systemConfig := c.GetSystemConfig()
	if systemConfig == (models.SystemConfig{}) {
		c.ThrowError("00101", "系统配置有误")
	}
	var newStorageSize int64
	var storageSize int64
	for _, fileData := range fileDatas {
		newStorageSize += fileData.Size
	}
	storage := c.GetStorageByAdmin()
	if storage != (models.Storage{}) {
		storageSize = storage.StorageSize + newStorageSize
		if storage.MaxStorageSize < storageSize {
			c.ThrowError("00300", "上传失败，您可用的空间已经不足！")
			return
		}
		storage.StorageSize = storageSize
	} else {
		c.ThrowError("00300", "上传失败，您没有分配可用的上传空间！")
		return
	}
	result := c.BatchUploadFile(fileDatas, systemConfig, "filedatas")
	var fileList []models.File
	if result["code"] == "error" {
		c.ErrorWithMessage(result["data"].(string))
		return
	}
	fileList = result["data"].([]models.File)
	var networkDiskList []models.NetworkDisk
	for _, file := range fileList {
		var saveNetworkDist models.NetworkDisk
		saveNetworkDist.AdminUid = adminUid
		saveNetworkDist.FilePath = filepath
		saveNetworkDist.QiNiuUrl = file.QiNiuUrl
		saveNetworkDist.LocalUrl = file.PicURL
		saveNetworkDist.MinioUrl = file.MinioUrl
		saveNetworkDist.FileSize = file.FileSize
		saveNetworkDist.FileName = file.PicName
		saveNetworkDist.ExtendName = file.PicExpandedName
		saveNetworkDist.FileOldName = file.FileOldName
		saveNetworkDist.Uid = file.Uid
		networkDiskList = append(networkDiskList, saveNetworkDist)
	}
	common.DB.Create(&networkDiskList)
	common.DB.Save(&storage)
	c.SuccessWithMessage("操作成功")
}
