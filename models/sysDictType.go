package models

import (
	_ "gorm.io/gorm"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/29 3:10 下午
 * @version 1.0
 */
type SysDictType struct {
	Uid         string    `gorm:"primaryKey" json:"uid"`
	Oid         int       `json:"oid"`
	DictName    string    `json:"dictName"`
	DictType    string    `json:"dictType"`
	CreateByUid string    `json:"createByUid"`
	UpdateByUid string    `json:"updateByUid"`
	Remark      string    `json:"remark"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt   time.Time `gorm:"column:update_time" json:"updateTime"`
	IsPublish   string    `json:"isPublish"`
	Sort        int       `json:"sort"`
}

func (SysDictType) TableName() string {
	return "t_sys_dict_type"
}
