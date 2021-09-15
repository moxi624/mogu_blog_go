package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/25 3:37 下午
 * @version 1.0
 */

import (
	_ "gorm.io/gorm"
	"time"
)

type Tag struct {
	Uid        string    `gorm:"primaryKey" json:"uid"`
	Content    string    `gorm:"type:text" json:"content"`
	Status     int8      `gorm:"default:1" json:"status"`
	ClickCount int       `json:"clickCount"`
	CreatedAt  time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt  time.Time `gorm:"column:update_time" json:"updateTime"`
	Sort       int       `json:"sort"`
}

func (Tag) TableName() string {
	return "t_tag"
}
