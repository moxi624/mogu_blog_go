package models

import (
	_ "gorm.io/gorm"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/25 10:54 上午
 * @version 1.0
 */
type Comment struct {
	Uid             string    `gorm:"primaryKey" json:"uid"`
	UserUid         string    `json:"userUid"`
	ToUid           string    `json:"toUid"`
	ToUserUid       string    `json:"toUserUid"`
	Content         string    `gorm:"type:text" json:"content"`
	BlogUid         string    `json:"blogUid"`
	Status          int8      `gorm:"default:1" json:"status"`
	CreatedAt       time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt       time.Time `gorm:"column:update_time" json:"updateTime"`
	Source          string    `json:"source"`
	Type            int       `json:"type"`
	FirstCommentUid string    `json:"firstCommentUid"`
	User            User      `gorm:"foreignKey:UserUid;references:Uid" json:"user"`
	UserName        string    `gorm:"-" json:"userName"`
	ToUserName      string    `gorm:"-" json:"toUserName"`
	ToUser          User      `gorm:"foreignKey:ToUserUid;references:Uid" json:"toUser"`
	ReplyList       []Comment `gorm:"-" json:"replyList"`
	ToComment       *Comment  `gorm:"-" json:"toComment"`
	SourceName      string    `gorm:"-" json:"sourceName"`
	Blog            Blog      `gorm:"foreignKey:BlogUid;references:Uid" json:"blog"`
}

func (Comment) TableName() string {
	return "t_comment"
}
