package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/22 11:16 上午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type CategoryMenu struct {
	Uid                string         `gorm:"primaryKey" json:"uid"`
	Name               string         `json:"name"`
	MenuLevel          int            `json:"menuLevel"`
	Summary            string         `json:"summary"`
	ParentUid          string         `json:"parentUid"`
	Url                string         `json:"url"`
	Icon               string         `json:"icon"`
	Sort               int            `json:"sort"`
	Status             int8           `gorm:"default:1" json:"status"`
	CreatedAt          time.Time      `gorm:"column:create_time" json:"createTime"`
	UpdatedAt          time.Time      `gorm:"column:update_time" json:"updateTime"`
	IsShow             int            `json:"isShow"`
	MenuType           int            `json:"menuType"`
	IsJumpExternalUrl  int            `json:"isJumpExternalUrl"`
	ParentCategoryMenu []CategoryMenu `gorm:"-" json:"parentCategoryMenu"`
	ChildCategoryMenu  []CategoryMenu `gorm:"-" json:"childCategoryMenu"`
}

func (CategoryMenu) TableName() string {
	return "t_category_menu"
}
