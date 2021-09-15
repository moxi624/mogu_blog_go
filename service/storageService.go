package service

import (
	"github.com/rs/xid"
	"mogu-go-v2/common"
	"mogu-go-v2/models"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/25 4:10 下午
 * @version 1.0
 */

type storageService struct{}

func (storageService) GetStorageByAdminUid(adminUidList []string) []models.Storage {
	var storageList []models.Storage
	common.DB.Where("status = ? and admin_uid in ?", 1, adminUidList).Find(&storageList)
	return storageList
}

func (storageService) EditStorageSizes(adminUid string, maxStorageSize int64) (string, bool) {
	var storage models.Storage
	common.DB.Where("admin_uid = ?", adminUid).Last(&storage)
	if storage == (models.Storage{}) {
		l.Print("未分配存储空间，重新初始化网盘空间")
		return InitStorageSize(adminUid, maxStorageSize), true
	}
	if maxStorageSize < storage.StorageSize {
		return "网盘容量不能小于当前已用空间", false
	}
	storage.MaxStorageSize = maxStorageSize
	common.DB.Save(&storage)
	return "操作成功", true
}

func InitStorageSize(adminUid string, maxStorageSize int64) string {
	saveStorage := models.Storage{
		AdminUid:       adminUid,
		StorageSize:    0,
		MaxStorageSize: maxStorageSize,
		Uid:            xid.New().String(),
	}
	common.DB.Create(&saveStorage)
	return "操作成功"
}

var StorageService = &storageService{}
