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
	"fmt"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/service"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/9 1:01 下午
 * @version 1.0
 */

type IndexRestApi struct {
	base.BaseController
}

func (c *IndexRestApi) GetWebConfig() {
	base.L.Print("获取网站配置")
	webConfigResult := common.RedisUtil.Get("WEB_CONFIG")
	if webConfigResult != "" {
		var webConfig models.WebConfig
		err := json.Unmarshal([]byte(webConfigResult), &webConfig)
		if err != nil {
			panic(err)
		}
		c.SuccessWithData(webConfig)
		return
	}
	var webConfig models.WebConfig
	common.DB.Order("create_time desc").First(&webConfig)
	if reflect.DeepEqual(webConfig, models.WebConfig{}) {
		c.ThrowError("00101", "系统配置不存在")
		return
	}
	var stringBuilder strings.Builder
	pictureResult := map[string]interface{}{}
	fmt.Println(webConfig)
	if webConfig.Logo != "" {
		stringBuilder.WriteString(webConfig.Logo + ",")
	}
	if webConfig.AliPay != "" {
		stringBuilder.WriteString(webConfig.AliPay + ",")
	}
	if webConfig.WeixinPay != "" {
		stringBuilder.WriteString(webConfig.WeixinPay + ",")
	}
	if !reflect.DeepEqual(stringBuilder, strings.Builder{}) {
		pictureResult = service.FileService.GetPicture(stringBuilder.String(), ",")
	}
	pictureList := common.WebUtil.GetPictureMap(pictureResult)
	pictureMap := map[string]string{}
	for _, item := range pictureList {
		pictureMap[item["uid"].(string)] = item["url"].(string)
	}
	if webConfig.Logo != "" && pictureMap[webConfig.Logo] != "" {
		webConfig.LogoPhoto = pictureMap[webConfig.Logo]
	}
	if webConfig.Logo != "" && pictureMap[webConfig.AliPay] != "" {
		webConfig.AliPayPhoto = pictureMap[webConfig.AliPay]
	}
	if webConfig.Logo != "" && pictureMap[webConfig.WeixinPay] != "" {
		webConfig.WeixinPayPhoto = pictureMap[webConfig.WeixinPay]
	}
	showListJson := webConfig.ShowList
	email := webConfig.Email
	qqNumber := webConfig.QqNumber
	qqGroup := webConfig.QqGroup
	github := webConfig.Github
	gitee := webConfig.Gitee
	webChat := webConfig.WeChat
	webConfig.Email, webConfig.QqGroup, webConfig.QqNumber, webConfig.Github, webConfig.Gitee, webConfig.WeChat = "", "", "", "", "", ""
	var showList []string
	err := json.Unmarshal([]byte(showListJson), &showList)
	if err != nil {
		panic(err)
	}
	for _, item := range showList {
		switch item {
		case "1":
			webConfig.Email = email
		case "2":
			webConfig.Email = qqNumber
		case "3":
			webConfig.Email = qqGroup
		case "4":
			webConfig.Email = github
		case "5":
			webConfig.Email = gitee
		case "6":
			webConfig.Email = webChat
		}
	}
	b, _ := json.Marshal(webConfig)
	common.RedisUtil.SetEx("WEB_CONFIG", string(b), 24, time.Hour)
	c.SuccessWithData(webConfig)
}

func (c *IndexRestApi) GetBlogByLevel() {
	level, _ := c.GetInt("level")
	currentPage, _ := c.GetInt("currentPage")
	userSort, _ := c.GetInt("userSort")
	jsonResult := common.RedisUtil.Get("BLOG_LEVEL:" + strconv.Itoa(level))
	if jsonResult != "" {
		var jsonResult2List []interface{}
		err := json.Unmarshal([]byte(jsonResult), &jsonResult2List)
		if err != nil {
			panic(err)
		}
		page := map[string]interface{}{}
		page["records"] = jsonResult2List
		c.Result("success", page)
		return
	}
	var blogCount string
	switch level {
	case 0:
		blogCount, _ = service.SysParamService.GetSysParamsValueByKey("BLOG_NEW_COUNT")
	case 1:
		blogCount, _ = service.SysParamService.GetSysParamsValueByKey("BLOG_FIRST_COUNT")
	case 2:
		blogCount, _ = service.SysParamService.GetSysParamsValueByKey("BLOG_SECOND_COUNT")
	case 3:
		blogCount, _ = service.SysParamService.GetSysParamsValueByKey("BLOG_THIRD_COUNT")
	case 4:
		blogCount, _ = service.SysParamService.GetSysParamsValueByKey("BLOG_FOURTH_COUNT")
	}
	if blogCount == "" {
		c.ThrowError("000", "请配置系统参数")
		return
	}
	pageSize, _ := strconv.Atoi(blogCount)
	pageList, total := service.BlogService.GetBlogPageByLevel(currentPage, pageSize, level, userSort)
	if (level == 1 || level == 2) && len(pageList) == 0 {
		hotCurrentPage := 1
		blogHotCount, _ := service.SysParamService.GetSysParamsValueByKey("BLOG_HOT_COUNT")
		blogSecondCount, _ := service.SysParamService.GetSysParamsValueByKey("BLOG_SECOND_COUNT")
		if blogHotCount == "" || blogSecondCount == "" {
			c.ThrowError("000", "请先配置参数")
			return
		}
		hotPageSize, _ := strconv.Atoi(blogHotCount)
		var hotPageList []models.BlogNoContent
		var total int64
		c.Wg.Add(2)
		go func() {
			common.DB.Model(&models.BlogNoContent{}).Where("status=? and is_publish=?", 1, "1").Count(&total)
			c.Wg.Done()
		}()
		go func() {
			common.DB.Where("status=? and is_publish=?", 1, "1").Offset((hotCurrentPage - 1) * hotPageSize).Limit(hotPageSize).Order("click_count desc").Find(&hotPageList)
			c.Wg.Done()
		}()
		c.Wg.Wait()
		var secondBlogList []models.BlogNoContent
		var firstBlogList []models.BlogNoContent
		i, _ := strconv.Atoi(blogSecondCount)
		for a := 0; a < len(hotPageList); a++ {
			if len(hotPageList)-len(firstBlogList) > i {
				firstBlogList = append(firstBlogList, hotPageList[a])
			} else {
				secondBlogList = append(secondBlogList, hotPageList[a])
			}
		}
		firstBlogList = service.BlogService.SetBlog(firstBlogList)
		secondBlogList = service.BlogService.SetBlog(secondBlogList)
		if len(firstBlogList) > 0 {
			b, _ := json.Marshal(firstBlogList)
			common.RedisUtil.SetEx("BLOG_LEVEL:1", string(b), 1, time.Hour)
		}
		if len(secondBlogList) > 0 {
			b, _ := json.Marshal(secondBlogList)
			common.RedisUtil.SetEx("BLOG_LEVEL:2", string(b), 1, time.Hour)
		}
		switch level {
		case 1:
			pageList = firstBlogList
		case 2:
			pageList = secondBlogList
		}
		iPage := page.IPage{
			Records: pageList,
			Total:   total,
			Size:    hotPageSize,
			Current: hotCurrentPage,
		}
		c.SuccessWithData(iPage)
		return
	}
	pageList = service.BlogService.SetBlog(pageList)
	if len(pageList) > 0 {
		b, _ := json.Marshal(pageList)
		common.RedisUtil.SetEx("BLOG_LEVEL:"+strconv.Itoa(level), string(b), 1, time.Hour)
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    pageSize,
		Current: currentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *IndexRestApi) GetNewBlog() {
	base.L.Print("获取首页最新的博客")
	currentPage, _ := c.GetInt("currentPage")
	blogNewCount, _ := service.SysParamService.GetSysParamsValueByKey("BLOG_NEW_COUNT")
	if blogNewCount == "" {
		c.ThrowError("000", "请先配置系统参数")
		return
	}
	pageSize, _ := strconv.Atoi(blogNewCount)
	var pageList []models.BlogNoContent
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.BlogNoContent{}).Where("status=? and is_publish=?", 1, "1").Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where("status=? and is_publish=?", 1, "1").Offset((currentPage - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	if len(pageList) == 0 {
		page := map[string]interface{}{}
		page["records"] = pageList
		c.Result("success", page)
		return
	}
	pageList = service.BlogService.SetBlog(pageList)
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    pageSize,
		Current: currentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *IndexRestApi) RecorderVisitPage() {
	pageName := c.GetString("pageName")
	if pageName == "" {
		c.ErrorWithMessage("参数传入错误")
	} else {
		c.SuccessWithMessage("记录成功")
	}
}

func (c *IndexRestApi) GetHotBlog() {
	base.L.Print("获取首页排行博客")
	jsonResult := common.RedisUtil.Get("HOT_BLOG")
	if jsonResult != "" {
		var jsonResult2List []interface{}
		err := json.Unmarshal([]byte(jsonResult), &jsonResult2List)
		if err != nil {
			panic(err)
		}
		page := map[string]interface{}{}
		page["records"] = jsonResult2List
		c.Result("success", page)
		return
	}
	var currentPage int
	blogHotCount, _ := service.SysParamService.GetSysParamsValueByKey("BLOG_HOT_COUNT")
	if blogHotCount == "" {
		c.ThrowError("000", "请先配置参数")
		return
	}
	pageSize, _ := strconv.Atoi(blogHotCount)
	var pageList []models.BlogNoContent
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.BlogNoContent{}).Where("status=? and is_publish=?", 1, "1").Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where("status=? and is_publish=?", 1, "1").Offset((currentPage - 1) * pageSize).Limit(pageSize).Order("click_count desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	pageList = service.BlogService.SetBlog(pageList)
	if len(pageList) > 0 {
		b, _ := json.Marshal(pageList)
		common.RedisUtil.SetEx("HOT_BLOG", string(b), 1, time.Hour)
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    pageSize,
		Current: currentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *IndexRestApi) GetHotTag() {
	hotTagCount, _ := service.SysParamService.GetSysParamsValueByKey("HOT_TAG_COUNT")
	currentPage := 1
	pageSize, _ := strconv.Atoi(hotTagCount)
	var pageList []models.Tag
	common.DB.Where("status=?", 1).Offset((currentPage - 1) * pageSize).Limit(pageSize).Order("sort,click_count desc").Find(&pageList)
	c.SuccessWithData(pageList)
}

func (c *IndexRestApi) GetLink() {
	friendlyLinkCount, _ := service.SysParamService.GetSysParamsValueByKey("FRIENDLY_LINK_COUNT")
	currentPage := 1
	pageSize, _ := strconv.Atoi(friendlyLinkCount)
	var pageList []models.Link
	common.DB.Where("link_status=? and status=?", 1, 1).Offset((currentPage - 1) * pageSize).Limit(pageSize).Order("sort desc").Find(&pageList)
	c.Result("success", pageList)
}

func (c *IndexRestApi) AddLinkCount() {
	uid := c.GetString("uid")
	if uid == "" {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	var link models.Link
	common.DB.Where("uid=?", uid).Find(&link)
	if !reflect.DeepEqual(link, models.Link{}) {
		link.ClickCount++
		common.DB.Save(&link)
	} else {
		c.ErrorWithMessage("传入参数有误")
	}
	c.SuccessWithMessage("更新成功")
}

func (c *IndexRestApi) GetWebNavbar() {
	base.L.Print("获取首页导航")
	webNavbarList := service.WebNavBarService.GetAllList()
	c.SuccessWithData(webNavbarList)
}
