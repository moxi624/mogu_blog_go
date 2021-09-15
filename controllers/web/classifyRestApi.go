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
	"mogu-go-v2/models/page"
	"mogu-go-v2/service"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/18 2:06 下午
 * @version 1.0
 */

type ClassifyRestApi struct {
	base.BaseController
}

func (c *ClassifyRestApi) GetBlogSortList() {
	base.L.Print("获取分类信息")
	var blogSortList []models.BlogSort
	common.DB.Where("status=?", 1).Order("sort desc").Find(&blogSortList)
	c.SuccessWithData(blogSortList)
}

func (c *ClassifyRestApi) GetArticleByBlogSortUid() {
	blogSortUid := c.GetString("blogSortUid")
	currentPage, _ := c.GetInt("currentPage", 1)
	pageSize, _ := c.GetInt("pageSize", 10)
	if blogSortUid == "" {
		base.L.Print("点击分类,传入BlogSortUid不能为空")
		c.ErrorWithMessage("传入BlogSortUid不能为空")
		return
	}
	var pageList []models.BlogNoContent
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Blog{}).Where("status=? and is_publish=? and blog_sort_uid=?", 1, "1", blogSortUid).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where("status=? and is_publish=? and blog_sort_uid=?", 1, "1", blogSortUid).Offset((currentPage - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	list := service.BlogService.SetTagAndSortAndPictureByBlogListNoContent(pageList)
	iPage := page.IPage{
		Records: list,
		Total:   total,
		Size:    pageSize,
		Current: currentPage,
	}
	c.SuccessWithData(iPage)
}
