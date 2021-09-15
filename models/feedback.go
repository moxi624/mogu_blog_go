package models

import "time"

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/23 8:10 上午
 * @version 1.0
 */

type Feedback struct {
	Uid            string    `gorm:"primaryKey" json:"uid"`
	UserUid        string    `json:"userUid"`
	Content        string    `gorm:"type:text" json:"content"`
	Status         int8      `gorm:"default:1" json:"status"`
	CreatedAt      time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt      time.Time `gorm:"column:update_time" json:"updateTime"`
	Title          string    `json:"title"`
	FeedbackStatus int       `json:"feedbackStatus"`
	Reply          string    `json:"reply"`
	AdminUid       string    `json:"adminUid"`
	User           User      `gorm:"foreignKey:UserUid;references:Uid" json:"user"`
}

func (Feedback) TableName() string {
	return "t_feedback"
}
