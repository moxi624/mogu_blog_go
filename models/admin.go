package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/11 4:47 下午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type Admin struct {
	Uid            string    `gorm:"primaryKey" json:"uid"`
	UserName       string    `json:"userName"`
	PassWord       string    `json:"passWord"`
	Gender         string    `json:"gender"`
	Avatar         string    `json:"avatar"`
	Email          string    `json:"email"`
	Birthday       time.Time `gorm:"type:date" json:"birthday"`
	Mobile         string    `json:"mobile"`
	ValidCode      string    `json:"validCode"`
	Summary        string    `json:"summary"`
	LoginCount     int       `json:"loginCount"`
	LastLoginTime  time.Time `json:"lastLoginTime"`
	LastLoginIp    string    `json:"lastLoginIp"`
	Status         int8      `gorm:"default:1" json:"status"`
	CreatedAt      time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt      time.Time `gorm:"column:update_time" json:"updateTime"`
	NickName       string    `json:"nickName"`
	QqNumber       string    `json:"qqNumber"`
	WeChat         string    `json:"weChat"`
	Occupation     string    `json:"occupation"`
	Github         string    `json:"github"`
	Gitee          string    `json:"gitee"`
	RoleUid        string    `json:"roleUid"`
	PersonResume   string    `gorm:"type:text" json:"personResume"`
	TokenUid       string    `gorm:"-" json:"tokenUid"`
	Role           Role      `gorm:"foreignKey:RoleUid;references:Uid" json:"role"`
	PhotoList      []string  `gorm:"-" json:"photoList"`
	StorageSize    int64     `gorm:"-" json:"storageSize"`
	MaxStorageSize int64     `gorm:"-" json:"maxStorageSize"`
}

func (Admin) TableName() string {
	return "t_admin"
}
