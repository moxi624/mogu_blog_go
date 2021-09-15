package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/5 9:30 上午
 * @version 1.0
 */

import (
	_ "gorm.io/gorm"
	"time"
)

type Subject struct {
	Uid          string    `gorm:"primaryKey" json:"uid"`
	Status       int8      `gorm:"default:1" json:"status"`
	CreatedAt    time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt    time.Time `gorm:"column:update_time" json:"updateTime"`
	SubjectName  string    `json:"subjectName"`
	FileUid      string    `json:"fileUid"`
	Summary      string    `json:"summary"`
	ClickCount   int       `json:"clickCount"`
	CollectCount int       `json:"collectCount"`
	Sort         int       `json:"sort"`
	PhotoList    []string  `gorm:"-" json:"photoList"`
}

func (Subject) TableName() string {
	return "t_subject"
}
