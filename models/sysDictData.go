package models

import (
	_ "gorm.io/gorm"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/29 4:20 下午
 * @version 1.0
 */
type SysDictData struct {
	Uid         string      `gorm:"primaryKey" json:"uid"`
	Oid         int         `json:"oid"`
	DictTypeUid string      `json:"dictTypeUid"`
	DictLabel   string      `json:"dictLabel"`
	DictValue   string      `json:"dictValue"`
	CssClass    string      `json:"cssClass"`
	ListClass   string      `json:"listClass"`
	IsDefault   int         `json:"isDefault"`
	CreateByUid string      `json:"createByUid"`
	UpdateByUid string      `json:"updateByUid"`
	Remark      string      `json:"remark"`
	Status      int8        `gorm:"default:1" json:"status"`
	CreatedAt   time.Time   `gorm:"column:create_time" json:"createTime"`
	UpdatedAt   time.Time   `gorm:"column:update_time" json:"updateTime"`
	IsPublish   string      `json:"isPublish"`
	Sort        int         `json:"sort"`
	SysDictType SysDictType `gorm:"foreignKey:DictTypeUid;references:Uid" json:"sysDictType"`
}

func (SysDictData) TableName() string {
	return "t_sys_dict_data"
}
