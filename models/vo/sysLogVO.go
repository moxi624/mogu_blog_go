//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/25 10:04 上午
 * @version 1.0
 */

type SysLogVO struct {
	Keyword      string `json:"keyword"`
	CurrentPage  int    `json:"currentPage"`
	PageSize     int    `json:"pageSize"`
	Uid          string `json:"uid"`
	Status       int    `json:"status"`
	UserName     string `json:"userName"`
	adminUid     string `json:"adminUid"`
	Ip           string `json:"ip"`
	IpSource     string `json:"ipSource"`
	Url          string `json:"url"`
	Type         string `json:"type"`
	ClassPath    string `json:"classPath"`
	Method       string `json:"method"`
	Operation    string `json:"operation"`
	Params       string `json:"params"`
	StartTime    string `json:"startTime"`
	SpendTime    int64  `json:"spendTime"`
	SpendTimeStr string `json:"spendTimeStr"`
}
