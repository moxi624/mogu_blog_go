package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/21 2:43 下午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type File struct {
	Uid             string    `gorm:"primaryKey" json:"uid"`
	FileOldName     string    `json:"fileOldName"`
	PicName         string    `json:"picName"`
	PicURL          string    `json:"picURL"`
	PicExpandedName string    `json:"picExpandedName"`
	FileSize        int64     `json:"fileSize"`
	FileSortUid     string    `json:"fileSortUid"`
	AdminUid        string    `json:"adminUid"`
	UserUid         string    `json:"userUid"`
	Status          int8      `gorm:"default:1" json:"status"`
	CreatedAt       time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt       time.Time `gorm:"column:update_time" json:"updateTime"`
	QiNiuUrl        string    `json:"qiNiuYun"`
	MinioUrl        string    `json:"minioUrl"`
}

func (File) TableName() string {
	return "t_file"
}
