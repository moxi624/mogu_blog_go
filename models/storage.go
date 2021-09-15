package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/14 1:41 下午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type Storage struct {
	Uid            string    `gorm:"primaryKey" json:"uid"`
	CreatedAt      time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt      time.Time `gorm:"column:update_time" json:"updateTime"`
	Status         int8      `gorm:"default:1" json:"status"`
	AdminUid       string    `json:"adminUid"`
	StorageSize    int64     `json:"storageSize"`
	MaxStorageSize int64     `json:"maxStorageSize"`
}

func (Storage) TableName() string {
	return "t_storage"
}
