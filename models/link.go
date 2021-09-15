package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/8 8:50 上午
 * @version 1.0
 */

import (
	_ "gorm.io/gorm"
	"time"
)

type Link struct {
	Uid        string    `gorm:"primaryKey" json:"uid"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	Url        string    `json:"url"`
	ClickCount int       `json:"clickCount"`
	AdminUid   string    `json:"adminUid"`
	UserUid    string    `json:"userUid"`
	Status     int8      `gorm:"default:1" json:"status"`
	CreatedAt  time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt  time.Time `gorm:"column:update_time" json:"updateTime"`
	Sort       int       `json:"sort"`
	LinkStatus int       `json:"linkStatus"`
	Email      string    `json:"email"`
	FileUid    string    `json:"fileUid"`
	PhotoList  []string  `gorm:"-" json:"photoList"`
}

func (Link) TableName() string {
	return "t_link"
}
