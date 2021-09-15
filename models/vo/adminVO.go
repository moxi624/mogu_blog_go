package vo

import "time"

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/30 9:26 上午
 * @version 1.0
 */
type AdminVO struct {
	Keyword        string    `json:"keyword"`
	CurrentPage    int       `json:"currentPage"`
	PageSize       int       `json:"pageSize"`
	Uid            string    `json:"uid"`
	Status         int       `json:"status"`
	UserName       string    `json:"userName"`
	PassWord       string    `json:"passWord"`
	Gender         string    `json:"gender"`
	Avatar         string    `json:"avatar"`
	Email          string    `json:"email"`
	Birthday       time.Time `gorm:"type:date" json:"birthday"`
	Mobile         string    `json:"mobile"`
	NickName       string    `json:"nickName"`
	QqNumber       string    `json:"qqNumber"`
	WeChat         string    `json:"weChat"`
	Occupation     string    `json:"occupation"`
	Summary        string    `json:"summary"`
	Github         string    `json:"github"`
	Gitee          string    `json:"gitee"`
	RoleUid        string    `json:"roleUid"`
	PersonResume   string    `gorm:"type:text" json:"personResume"`
	StorageSize    int64     `json:"storageSize"`
	MaxStorageSize int64     `json:"maxStorageSize"`
}
