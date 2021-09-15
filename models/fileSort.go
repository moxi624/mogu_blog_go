package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/31 4:30 下午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type FileSort struct {
	Uid         string    `gorm:"primaryKey" json:"uid"`
	ProjectName string    `json:"projectName"`
	SortName    string    `json:"sortName"`
	Url         string    `json:"url"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt   time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (FileSort) TableName() string {
	return "t_file_sort"
}
