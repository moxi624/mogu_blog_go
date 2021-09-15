package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/19 2:30 下午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type ResourceSort struct {
	Uid        string    `gorm:"primaryKey" json:"uid"`
	Status     int8      `gorm:"default:1" json:"status"`
	CreatedAt  time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt  time.Time `gorm:"column:update_time" json:"updateTime"`
	FileUid    string    `json:"fileUid"`
	SortName   string    `json:"sortName"`
	Content    string    `json:"content"`
	ClickCount string    `json:"clickCount"`
	ParentUid  string    `json:"parentUid"`
	Sort       int       `json:"sort"`
	PhotoList  []string  `gorm:"-" json:"photoList"`
}

func (ResourceSort) TableName() string {
	return "t_resource_sort"
}
