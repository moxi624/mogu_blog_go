package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/12 8:24 下午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type Role struct {
	Uid              string    `gorm:"primaryKey" json:"uid"`
	RoleName         string    `json:"roleName"`
	Summary          string    `json:"summary"`
	Status           int8      `gorm:"default:1" json:"status"`
	CreatedAt        time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt        time.Time `gorm:"column:update_time" json:"updateTime"`
	CategoryMenuUids string    `gorm:"type:text" json:"categoryMenuUids"`
}

func (Role) TableName() string {
	return "t_role"
}
