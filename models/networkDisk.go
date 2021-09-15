package models

import (
	_ "gorm.io/gorm"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/13 3:22 下午
 * @version 1.0
 */

type NetworkDisk struct {
	Uid         string    `gorm:"primaryKey" json:"uid"`
	CreatedAt   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt   time.Time `gorm:"column:update_time" json:"updateTime"`
	Status      int8      `gorm:"default:1" json:"status"`
	AdminUid    string    `json:"adminUid"`
	ExtendName  string    `gorm:"default:null" json:"extendName"`
	FileName    string    `json:"fileName"`
	FilePath    string    `json:"filePath"`
	FileSize    int64     `json:"fileSize"`
	IsDir       int       `json:"isDir"`
	LocalUrl    string    `json:"localUrl"`
	QiNiuUrl    string    `json:"qiNiuUrl"`
	FileOldName string    `json:"fileOldName"`
	MinioUrl    string    `json:"minioUrl"`
	OldFilePath string    `gorm:"-" json:"oldFilePath"`
	NewFilePath string    `gorm:"-" json:"newFilePath"`
	Files       string    `gorm:"-" json:"files"`
	FileType    int       `gorm:"-" json:"fileType"`
	FileUrl     string    `gorm:"-" json:"fileUrl"`
}

func (NetworkDisk) TableName() string {
	return "t_network_disk"
}
