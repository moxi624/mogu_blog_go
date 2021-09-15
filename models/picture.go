package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/19 4:39 下午
 * @version 1.0
 */

import (
	_ "gorm.io/gorm"
	"time"
)

type Picture struct {
	Uid            string    `gorm:"primaryKey" json:"uid"`
	Status         int8      `gorm:"default:1" json:"status"`
	CreatedAt      time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt      time.Time `gorm:"column:update_time" json:"updateTime"`
	FileUid        string    `json:"fileUid"`
	PicName        string    `json:"picName"`
	PictureSortUid string    `json:"pictureSortUid"`
	PictureUrl     string    `gorm:"-" json:"pictureUrl"`
}

func (Picture) TableName() string {
	return "t_picture"
}
