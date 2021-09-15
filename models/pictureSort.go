package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/19 3:52 下午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type PictureSort struct {
	Uid       string    `gorm:"primaryKey" json:"uid"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"updateTime"`
	FileUid   string    `json:"fileUid"`
	Name      string    `json:"name"`
	ParentUid string    `json:"parentUid"`
	Sort      int       `json:"sort"`
	IsShow    int8      `gorm:"default:1" json:"isShow"`
	PhotoList []string  `gorm:"-" json:"photoList"`
}

func (PictureSort) TableName() string {
	return "t_picture_sort"
}
