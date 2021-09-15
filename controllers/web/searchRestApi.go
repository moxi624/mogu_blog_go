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
	"math"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/service"
	"strconv"
	"strings"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/18 3:57 下午
 * @version 1.0
 */
type SearchRestApi struct {
	base.BaseController
}

func (c *SearchRestApi) SqlSearchBlog() {
	keywords := c.GetString("keywords")
	currentPage, _ := c.GetInt("currentPage", 1)
	pageSize, _ := c.GetInt("pageSize", 10)
	keyword := strings.TrimSpace(keywords)
	if keyword == "" {
		c.ErrorWithMessage("关键字不能为空")
	} else {
		var blogList []models.BlogNoContent
		where := "status=? and is_publish=? and (title like '%" + keyword + "%'" + "or summary like '%" + keyword + "%')"
		common.DB.Where(where, 1, "1").Offset((currentPage - 1) * pageSize).Limit(pageSize).Order("click_count desc").Find(&blogList)
		var blogSortUidList []string
		pictureMap := map[string]string{}
		var fileUids strings.Builder
		for i, item := range blogList {
			blogSortUidList = append(blogSortUidList, item.BlogSortUid)
			if item.FileUid != "" {
				fileUids.WriteString(item.FileUid + ",")
			}
			blogList[i].Title = service.BlogService.GetHitCode(item.Title, keyword)
			blogList[i].Summary = service.BlogService.GetHitCode(item.Summary, keyword)
		}
		var pictureList map[string]interface{}
		if fileUids.Len() > 0 {
			pictureList = service.FileService.GetPicture(fileUids.String(), ",")
		}
		picList := common.WebUtil.GetPictureMap(pictureList)
		for _, item := range picList {
			pictureMap[item["uid"].(string)] = item["url"].(string)
		}
		var blogSortList []models.BlogSort
		if len(blogSortUidList) > 0 {
			common.DB.Find(&blogSortList, blogSortUidList)
		}
		blogSortMap := map[string]string{}
		for _, item := range blogSortList {
			blogSortMap[item.Uid] = item.SortName
		}
		for i, item := range blogList {
			if blogSortMap[item.BlogSortUid] != "" {
				blogList[i].BlogSortName = blogSortMap[item.BlogSortUid]
			}
			if item.FileUid != "" {
				pictureUidsTemp := strings.Split(item.FileUid, ",")
				var pictureListTemp []string
				for _, picture := range pictureUidsTemp {
					pictureListTemp = append(pictureListTemp, pictureMap[picture])
				}
				if len(pictureListTemp) > 0 {
					blogList[i].PhotoUrl = pictureListTemp[0]
				} else {
					blogList[i].PhotoUrl = ""
				}
			}
		}
		m := map[string]interface{}{}
		m["total"] = len(blogList)
		f := float64(len(blogList)) / float64(pageSize)
		m["totalPages"] = math.Ceil(f)
		m["pageSize"] = pageSize
		m["currentPage"] = currentPage
		m["blogList"] = blogList
		c.SuccessWithData(m)
	}
}

func (c *SearchRestApi) SearchBlogByTag() {
	tagUid := c.GetString("tagUid")
	currentPage, _ := c.GetInt("currentPage", 1)
	pageSize, _ := c.GetInt("pageSize", 10)
	var tag models.Tag
	common.DB.Where("uid = ?", tagUid).Find(&tag)
	if tag != (models.Tag{}) {
		ip := c.GetIP()
		jsonResult := common.RedisUtil.Get("TAG_CLICK:" + ip + "#" + tagUid)
		if jsonResult == "" {
			clickCount := tag.ClickCount
			tag.ClickCount = clickCount + 1
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
	pageList = service.BlogService.SetTagAndSortAndPictureByBlogListNoContent(pageList)
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    pageSize,
		Current: currentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *SearchRestApi) SearchBlogByAuthor() {
	author := c.GetString("author")
	currentPage, _ := c.GetInt("currentPage", 1)
	pageSize, _ := c.GetInt("pageSize", 10)
	var pageList []models.BlogNoContent
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.BlogNoContent{}).Where("status=? and is_publish=? and author=?", 1, "1", author).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where("status=? and is_publish=? and author=?", 1, "1", author).Offset((currentPage - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	pageList = service.BlogService.SetTagAndSortAndPictureByBlogListNoContent(pageList)
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    pageSize,
		Current: currentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *SearchRestApi) SearchBlogBySort() {
	blogSortUid := c.GetString("blogSortUid")
	currentPage, _ := c.GetInt("currentPage", 1)
	pageSize, _ := c.GetInt("pageSize", 10)
	var blogSort models.BlogSort
	common.DB.Where("uid = ?", blogSortUid).Find(&blogSort)
	if blogSort != (models.BlogSort{}) {
		ip := c.GetIP()
		jsonResult := common.RedisUtil.Get("TAG_CLICK:" + ip + "#" + blogSortUid)
		if jsonResult == "" {
			clickCount := blogSort.ClickCount
			blogSort.ClickCount = clickCount + 1
			common.DB.Save(&blogSort)
			common.RedisUtil.SetEx("TAG_CLICK:"+ip+"#"+blogSortUid, strconv.Itoa(clickCount), 24, time.Hour)
		}
	}
	var pageList []models.BlogNoContent
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.BlogNoContent{}).Where("status=? and is_publish=? and blog_sort_uid=?", 1, "1", blogSortUid).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where("status=? and is_publish=? and blog_sort_uid=?", 1, "1", blogSortUid).Offset((currentPage - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	pageList = service.BlogService.SetTagAndSortAndPictureByBlogListNoContent(pageList)
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    pageSize,
		Current: currentPage,
	}
	c.SuccessWithData(iPage)
}
