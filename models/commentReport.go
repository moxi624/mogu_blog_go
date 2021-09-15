//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package models

import "time"

/**
 *
 * @author  镜湖老杨
 * @date  2021/3/4 8:23 上午
 * @version 1.0
 */

type CommentReport struct {
	Uid              string    `gorm:"primaryKey" json:"uid"`
	UserUid          string    `json:"userUid"`
	ReportCommentUid string    `json:"reportCommentUid"`
	ReportUserUid    string    `json:"reportUserUid"`
	Content          string    `gorm:"type:text" json:"content"`
	Progress         int8      `json:"progress"`
	Status           int8      `gorm:"default:1" json:"status"`
	CreatedAt        time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt        time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (CommentReport) TableName() string {
	return "t_comment_report"
}
