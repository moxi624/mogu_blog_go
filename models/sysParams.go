package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/29 3:34 下午
 * @version 1.0
 */

import (
	_ "gorm.io/gorm"
	"time"
)

type SysParams struct {
	Uid         string    `gorm:"primaryKey" json:"uid"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt   time.Time `gorm:"column:update_time" json:"updateTime"`
	ParamsType  string    `json:"paramsType"`
	ParamsName  string    `json:"paramsName"`
	ParamsKey   string    `json:"paramsKey"`
	Remark      string    `json:"remark"`
	ParamsValue string    `json:"paramsValue"`
	Sort        int       `json:"sort"`
}

func (SysParams) TableName() string {
	return "t_sys_params"
}
