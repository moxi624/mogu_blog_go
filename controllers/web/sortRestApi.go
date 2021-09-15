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
	"encoding/json"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/service"
	"reflect"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/18 10:00 上午
 * @version 1.0
 */

type SortRestApi struct {
	base.BaseController
}

func (c *SortRestApi) GetSortList() {
	base.L.Print("获取归档日期")
	monthResult := common.RedisUtil.Get("MONTH_SET")
	if monthResult != "" {
		var list []interface{}
		err := json.Unmarshal([]byte(monthResult), &list)
		if err != nil {
			panic(err)
		}
		c.SuccessWithData(list)
		return
	}
	var list []models.BlogNoContent
	common.DB.Where("status=?,and is_publish=?", 1, "1").Order("create_time desc").Find(&list)
	list = service.BlogService.SetTagAndSortAndPictureByBlogListNoContent(list)
	m := map[string][]models.BlogNoContent{}
	var monthSet []string
	for _, blog := range list {
		createTime := time.Time(blog.CreatedAt)
		month := createTime.Format("2006年01月")
		monthSet = append(monthSet, month)
		if reflect.DeepEqual(m[month], models.BlogNoContent{}) {
			var blogList []models.BlogNoContent
			blogList = append(blogList, blog)
			m[month] = blogList
		} else {
			blogList := m[month]
			blogList = append(blogList, blog)
			m[month] = blogList
		}
	}
	for key, value := range m {
		b, _ := json.Marshal(value)
		common.RedisUtil.Set("BLOG_SORT_BY_MONTH:"+key, string(b))
	}
	monthSet = common.RemoveRepByMap(monthSet)
	b, _ := json.Marshal(monthSet)
	common.RedisUtil.Set("MONTH_SET", string(b))
	c.SuccessWithData(monthSet)
}

func (c *SortRestApi) GetArticleByMonth() {
	base.L.Print("通过月份获取文章列表")
	monthDate := c.GetString("monthDate")
	if monthDate == "" {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	contentResult := common.RedisUtil.Get("BLOG_SORT_BY_MONTH:" + monthDate)
	if contentResult != "" {
		var list []interface{}
		err := json.Unmarshal([]byte(contentResult), &list)
		if err != nil {
			panic(err)
		}
		c.SuccessWithData(list)
		return
	}
	var list []models.BlogNoContent
	common.DB.Where("status=? and is_publish=?", 1, "1").Order("create_time desc").Find(&list)
	list = service.BlogService.SetTagAndSortAndPictureByBlogListNoContent(list)
	m := map[string][]models.BlogNoContent{}
	var monthSet []string
	for _, blog := range list {
		createTime := time.Time(blog.CreatedAt)
		month := createTime.Format("2006年01月")
		monthSet = append(monthSet, month)
		if reflect.DeepEqual(m[month], models.BlogNoContent{}) {
			var blogList []models.BlogNoContent
			blogList = append(blogList, blog)
			m[month] = blogList
		} else {
			blogList := m[month]
			blogList = append(blogList, blog)
			m[month] = blogList
		}
	}
	for key, value := range m {
		b, _ := json.Marshal(value)
		common.RedisUtil.Set("BLOG_SORT_BY_MONTH:"+key, string(b))
	}
	monthSet = common.RemoveRepByMap(monthSet)
	b, _ := json.Marshal(monthSet)
	common.RedisUtil.Set("MONTH_SET", string(b))
	c.SuccessWithData(monthDate)
}
