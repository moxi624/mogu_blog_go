package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/7 4:25 下午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type WebVisit struct {
	Uid             string    `gorm:primaryKey" json:"uid"`
	UserUid         string    `json:"userUid"`
	Ip              string    `json:"ip"`
	Behavior        string    `json:"behavior"`
	ModuleUid       string    `json:"moduleUid"`
	OtherData       string    `json:"otherData"`
	Status          int8      `gorm:"default:1" json:"status"`
	CreatedAt       time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt       time.Time `gorm:"column:update_time" json:"updateTime"`
	Os              string    `json:"os"`
	Browser         string    `json:"browser"`
	IpSource        string    `json:"ipSource"`
	Content         string    `gorm:"-" json:"content"`
	BehaviorContent string    `gorm:"-" json:"behaviorContent"`
}

func (WebVisit) TableName() string {
	return "t_web_visit"
}
