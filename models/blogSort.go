package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/26 8:24 上午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type BlogSort struct {
	Uid        string    `gorm:"primaryKey" json:"uid"`
	SortName   string    `json:"sortName"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt  time.Time `gorm:"column:update_time" json:"updateTime"`
	Status     int8      `gorm:"default:1" json:"status"`
	Sort       int       `json:"sort"`
	ClickCount int       `json:"clickCount"`
}

func (BlogSort) TableName() string {
	return "t_blog_sort"
}
