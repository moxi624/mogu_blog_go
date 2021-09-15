package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/3 8:51 上午
 * @version 1.0
 */

import (
	_ "gorm.io/gorm"
	"time"
)

type SubjectItem struct {
	Uid        string    `gorm:"primaryKey" json:"uid"`
	CreatedAt  time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt  time.Time `gorm:"column:update_time" json:"updateTime"`
	Status     int8      `gorm:"default:1" json:"status"`
	SubjectUid string    `json:"subjectUid"`
	BlogUid    string    `json:"blogUid"`
	Sort       int       `json:"sort"`
	Blog       Blog      `gorm:"foreignKey:Uid;references:BlogUid" json:"blog"`
}

func (SubjectItem) TableName() string {
	return "t_subject_item"
}
