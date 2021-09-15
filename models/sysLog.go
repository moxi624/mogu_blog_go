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
 * @date  2021/2/25 2:03 下午
 * @version 1.0
 */

type SysLog struct {
	Uid       string    `gorm:"primaryKey" json:"uid"`
	UserName  string    `json:"userName"`
	AdminUid  string    `json:"adminUid"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"updateTime"`
	Ip        string    `json:"ip"`
	Url       string    `json:"url"`
	Type      string    `json:"type"`
	ClassPath string    `json:"classPath"`
	IpSource  string    `json:"ipSource"`
	Method    string    `json:"method"`
	Operation string    `json:"operation"`
	Params    string    `gorm:"type:text" json:"params"`
	SpendTime int       `json:"spendTime"`
}

func (*SysLog) TableName() string {
	return "t_sys_log"
}
