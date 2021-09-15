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
	"strconv"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/18 2:50 下午
 * @version 1.0
 */

type TagRestApi struct {
	base.BaseController
}

func (c *TagRestApi) GetTagList() {
	base.L.Print("获取标签信息")
	var tagList []models.Tag
	common.DB.Where("status=?", 1).Order("sort desc").Find(&tagList)
	c.SuccessWithData(tagList)
}

func (c *TagRestApi) GetArticleByTagUid() {
	tagUid := c.GetString("tagUid")
	currentPage, _ := c.GetInt("currentPage", 1)
	pageSize, _ := c.GetInt("pageSize", 10)
	if tagUid == "" {
		c.ErrorWithMessage("传入TagUid不能为空")
		return
	}
	base.L.Print("通过blogSortUid获取文章列表")
	var tag models.Tag
	common.DB.Where("uid=?", tagUid).Find(&tag)
	if tag != (models.Tag{}) {
		ip := c.GetIP()
		jsonResult := common.RedisUtil.Get("TAG_CLICK:" + ip + "#" + tagUid)
		if jsonResult == "" {
			clickCount := tag.ClickCount + 1
			tag.ClickCount = clickCount
			common.DB.Save(&tag)
			common.RedisUtil.SetEx("TAG_CLICK:"+ip+"#"+tagUid, strconv.Itoa(clickCount), 24, time.Hour)
		}
	}
	var pageList []models.BlogNoContent
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.BlogNoContent{}).Where("status=? and is_publish=? and tag_uid=?", 1, "1", tagUid).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where("status=? and is_publish=? and tag_uid=?", 1, "1", tagUid).Offset((currentPage - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&pageList)
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
