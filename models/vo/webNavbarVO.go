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
 * @date  2021/3/11 1:33 下午
 * @version 1.0
 */

type WebNavbarVO struct {
	Keyword           string `json:"keyword"`
	CurrentPage       int    `json:"currentPage"`
	PageSize          int    `json:"pageSize"`
	Uid               string `json:"uid"`
	Status            int    `json:"status"`
	Name              string `json:"name"`
	NavbarLevel       int    `json:"navbarLevel"`
	Summary           string `json:"summary"`
	Icon              string `json:"icon"`
	ParentUid         string `json:"parentUid"`
	Url               string `json:"url"`
	Sort              int    `json:"sort"`
	IsShow            int    `json:"isShow"`
	IsJumpExternalUrl int    `json:"isJumpExternalUrl"`
}
