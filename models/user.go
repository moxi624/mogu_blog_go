package models

import (
	_ "gorm.io/gorm"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/25 11:11 上午
 * @version 1.0
 */
type User struct {
	Uid                    string    `gorm:"primaryKey" json:"uid"`
	UserName               string    `json:"userName"`
	PassWord               string    `json:"passWord"`
	Gender                 int       `json:"gender"`
	Avatar                 string    `json:"avatar"`
	Email                  string    `json:"email"`
	Birthday               time.Time `gorm:"type:date" json:"birthday"`
	Mobile                 string    `json:"mobile"`
	ValidCode              string    `json:"validCode"`
	Summary                string    `json:"summary"`
	LoginCount             int       `json:"loginCount"`
	LastLoginTime          time.Time `json:"lastLoginTime"`
	LastLoginIp            string    `json:"lastLoginIp"`
	Status                 int8      `gorm:"default:1" json:"status"`
	CreatedAt              time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt              time.Time `gorm:"column:update_time" json:"updateTime"`
	NickName               string    `json:"nickName"`
	Source                 string    `json:"source"`
	Uuid                   string    `json:"uuid"`
	QqNumber               string    `json:"qqNumber"`
	WeChat                 string    `json:"weChat"`
	Occupation             string    `json:"occupation"`
	CommentStatus          int       `json:"commentStatus"`
	IpSource               string    `json:"ipSource"`
	Browser                string    `json:"browser"`
	Os                     string    `json:"os"`
	StartEmailNotification int       `json:"startEmailNotification"`
	UserTag                int       `json:"userTag"`
	PhotoUrl               string    `gorm:"-" json:"photoUrl"`
}

func (User) TableName() string {
	return "t_user"
}
