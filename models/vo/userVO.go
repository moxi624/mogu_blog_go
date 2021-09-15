package vo

import (
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/9 9:52 上午
 * @version 1.0
 */
type UserVO struct {
	Keyword                string      `json:"keyword"`
	CurrentPage            int         `json:"currentPage"`
	PageSize               int         `json:"pageSize"`
	Uid                    string      `json:"uid"`
	Status                 int         `json:"status"`
	UserName               string      `json:"userName"`
	PassWord               string      `json:"passWord"`
	NickName               string      `json:"nickName"`
	Gender                 interface{} `json:"gender"`
	Avatar                 string      `json:"avatar"`
	Email                  string      `json:"email"`
	Birthday               time.Time   `gorm:"type:date" json:"birthday"`
	Mobile                 string      `json:"mobile"`
	QqNumber               string      `json:"qqNumber"`
	WeChat                 string      `json:"weChat"`
	Occupation             string      `json:"occupation"`
	Summary                string      `json:"summary"`
	Source                 string      `json:"source"`
	Uuid                   string      `json:"uuid"`
	CommentStatus          interface{} `json:"commentStatus"`
	StartEmailNotification int         `json:"startEmailNotification"`
	UserTag                int         `json:"userTag"`
	PhotoUrl               string      `json:"photoUrl"`
}
