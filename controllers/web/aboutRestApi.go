//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package web

import (
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/service"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/10 2:00 下午
 * @version 1.0
 */

type AboutRestApi struct {
	base.BaseController
}

func (c *AboutRestApi) GetMe() {
	base.L.Print("获取关于我的信息")
	var admin models.Admin
	common.DB.Where("user_name=?", "admin").Last(&admin)
	admin.PassWord = ""
	if admin.Avatar != "" {
		pictureList := service.FileService.GetPicture(admin.Avatar, ",")
		admin.PhotoList = common.WebUtil.GetPicture(pictureList)
	}
	result := models.Admin{
		NickName:     admin.NickName,
		Occupation:   admin.Occupation,
		Summary:      admin.Summary,
		Avatar:       admin.Avatar,
		PhotoList:    admin.PhotoList,
		PersonResume: admin.PersonResume,
	}
	c.SuccessWithData(result)
}
