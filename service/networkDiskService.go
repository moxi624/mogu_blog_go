package service

import (
	"mogu-go-v2/common"
	"mogu-go-v2/models"
	"mogu-go-v2/models/vo"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/19 9:26 上午
 * @version 1.0
 */

type networkDiskService struct{}

func (networkDiskService) UpdateFilepathByFilePath(networkDiskVO vo.NetworkDiskVO, isEdit bool) bool {
	oldFilePath := networkDiskVO.OldFilePath
	newFilePath := networkDiskVO.NewFilePath
	s := strings.TrimSuffix(newFilePath, networkDiskVO.FileName+"/")
	if oldFilePath == s && !isEdit {
		return false
	}
	fileName := networkDiskVO.FileName
	fileOldName := networkDiskVO.FileOldName
	extendName := networkDiskVO.ExtendName
	if extendName == "null" {
		extendName = ""
	}
	var networkDiskList []models.NetworkDisk
	where := "file_path=? and file_name=?"
	if extendName != "" {
		where += " and extend_name='" + extendName + "'"
	} else {
		where += " and extend_name is null"
	}
	common.DB.Where(where, oldFilePath, fileName).Find(&networkDiskList)
	for i := range networkDiskList {
		networkDiskList[i].FilePath = newFilePath
		networkDiskList[i].FileOldName = networkDiskVO.FileOldName
		if extendName == "" {
			networkDiskList[i].FileName = networkDiskVO.FileOldName
		}
	}
	if len(networkDiskList) > 0 {
		common.DB.Save(&networkDiskList)
	}
	oldFilePath += fileName + "/"
	newFilePath += fileOldName + "/"
	oldFilePath = strings.ReplaceAll(oldFilePath, "\\", "\\\\\\\\")
	oldFilePath = strings.ReplaceAll(oldFilePath, "'", "\\'")
	oldFilePath = strings.ReplaceAll(oldFilePath, "%", "\\%")
	oldFilePath = strings.ReplaceAll(oldFilePath, "_", "\\_")
	if extendName == "" {
		var childList []models.NetworkDisk
		where := "file_path like '" + oldFilePath + "%'"
		common.DB.Where(where).Find(&childList)
		for i, networkDisk := range childList {
			filePath := networkDisk.FilePath
			p := strings.ReplaceAll(filePath, oldFilePath, newFilePath)
			childList[i].FilePath = p
		}
		if len(childList) > 0 {
			common.DB.Save(&childList)
		}
	}
	return true
}

var NetworkDiskService = &networkDiskService{}
