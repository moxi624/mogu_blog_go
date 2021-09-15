//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package service

import (
	"mogu-go-v2/common"
	"mogu-go-v2/models"
	"sort"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/3/11 9:07 上午
 * @version 1.0
 */

type webNavBarService struct{}

func (*webNavBarService) GetAllList() []models.WebNavbar {
	var list []models.WebNavbar
	common.DB.Where("navbar_level=? and status=?", 1, 1).Order("sort desc").Find(&list)
	var ids []string
	for _, item := range list {
		if item.Uid != "" {
			ids = append(ids, item.Uid)
		}
	}
	var childList []models.WebNavbar
	common.DB.Where("status=? and parent_uid in ?", 1, ids).Find(&childList)
	var tempList []models.WebNavbar
	for i, parentItem := range list {
		for _, item := range childList {
			if item.ParentUid == parentItem.Uid {
				tempList = append(tempList, item)
			}
		}
		sort.SliceStable(tempList, func(i, j int) bool {
			return tempList[i].Sort <= tempList[j].Sort
		})
		list[i].ChildWebNavBar = tempList
		tempList = []models.WebNavbar{}
	}
	return list
}

var WebNavBarService = &webNavBarService{}
