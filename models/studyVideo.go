package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/20 11:22 上午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type StudyVideo struct {
	Uid             string       `gorm:"primaryKey" json:"uid"`
	Status          int8         `gorm:"default:1" json:"status"`
	CreatedAt       time.Time    `gorm:"column:create_time" json:"createTime"`
	UpdatedAt       time.Time    `gorm:"column:update_time" json:"updateTime"`
	FileUid         string       `json:"fileUid"`
	ResourceSortUid string       `json:"resourceSortUid"`
	Name            string       `json:"name"`
	Summary         string       `json:"summary"`
	Content         string       `json:"content"`
	BaiduPath       string       `json:"baiduPath"`
	ClickCount      int          `json:"clickCount"`
	ParentUid       string       `json:"parentUid"`
	PhotoList       []string     `gorm:"-" json:"photoList"`
	ResourceSort    ResourceSort `gorm:"foreignKey:ResourceSortUid;references:Uid" json:"resourceSort"`
}

func (StudyVideo) TableName() string {
	return "t_study_video"
}
