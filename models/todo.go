package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/26 11:08 上午
 * @version 1.0
 */

import (
	_ "gorm.io/gorm"
	"time"
)

type Todo struct {
	Uid       string    `gorm:"primaryKey" json:"uid"`
	AdminUid  string    `json:"adminUid"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (Todo) TableName() string {
	return "t_todo"
}
