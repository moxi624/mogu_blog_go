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
	beego "github.com/beego/beego/v2/server/web"
	"github.com/rs/xid"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/service"
	"reflect"
	"strconv"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/20 2:30 下午
 * @version 1.0
 */

type BlogContentRestApi struct {
	base.BaseController
}

func (c *BlogContentRestApi) GetBlogByUid() {
	uid := c.GetString("uid")
	oid, _ := c.GetInt("oid")
	ip := c.GetIP()
	if uid == "" && oid == 0 {
		c.ErrorWithMessage("传入参数有误")
		return
	}
	var blog models.Blog
	if uid != "" {
		common.DB.Where("uid=?", uid).Find(&blog)
	} else {
		common.DB.Where("oid=?", oid).Last(&blog)
	}
	if reflect.DeepEqual(blog, models.Blog{}) || blog.Status == 0 || blog.IsPublish == "0" {
		c.ErrorWithMessage("博客已被删除")
		return
	}
	setBlogCopyright(&blog)
	c.Wg.Add(3)
	go func() {
		service.BlogService.SetTagByBlog(&blog)
		c.Wg.Done()
	}()
	go func() {
		service.BlogService.SetSortByBlog(&blog)
		c.Wg.Done()
	}()
	go func() {
		jsonResult := common.RedisUtil.Get("BLOG_CLICK:" + ip + "#" + blog.Uid)
		if jsonResult == "" {
			blog.ClickCount++
			common.DB.Save(&blog)
			common.RedisUtil.SetEx("BLOG_CLICK:"+ip+"#"+blog.Uid, strconv.Itoa(blog.ClickCount), 24, time.Hour)
		}
		c.Wg.Done()
	}()
	c.Wg.Wait()
	setPhotoListByBlog(&blog)
	c.SuccessWithData(blog)
}

func setBlogCopyright(blog *models.Blog) {
	if blog.IsOriginal == "1" {
		originalTemplate, _ := beego.AppConfig.String("original_template")
		blog.Copyright = originalTemplate
	} else {
		reprintedTemplate, _ := beego.AppConfig.String("reprinted_template")
		variable := []string{
			blog.ArticlesPart,
			blog.ArticlesPart,
			blog.Author,
		}
		str := fmt.Sprintf(reprintedTemplate, variable)
		blog.Copyright = str
	}
}

func setPhotoListByBlog(blog *models.Blog) {
	if !reflect.DeepEqual(blog, models.Blog{}) && blog.FileUid != "" {
		result := service.FileService.GetPicture(blog.FileUid, ",")
		picList := common.WebUtil.GetPicture(result)
		if len(picList) > 0 {
			blog.PhotoList = picList
		}
	}
}

func (c *BlogContentRestApi) GetSameBlogByBlogUid() {
	blogUid := c.GetString("blogUid")
	if blogUid == "" {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	var blog models.Blog
	common.DB.Where("uid=?", blogUid).Find(&blog)
	var pageList []models.Blog
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Blog{}).Where("status=? and is_publish=? and blog_sort_uid=?", 1, "1", blog.BlogSortUid).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where("status=? and is_publish=? and blog_sort_uid=?", 1, "1", blog.BlogSortUid).Limit(10).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	pageList = service.BlogService.SetTagAndSortByBlogList(pageList)
	var newList []models.Blog
	for _, item := range pageList {
		if item.Uid == blogUid {
			continue
		}
		newList = append(newList, item)
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    10,
		Current: 1,
	}
	c.SuccessWithData(iPage)
}

func (c *BlogContentRestApi) PraiseBlogByUid() {
	uid := c.GetString("uid")
	if uid == "" {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	tokenJson := common.RedisUtil.Get("USER_TOKEN:" + token)
	var praise models.Comment
	if tokenJson == "" {
		c.ErrorWithMessage("请登录后点赞")
		return
	}
	var user models.User
	err := json.Unmarshal([]byte(tokenJson), &user)
	if err != nil {
		panic(err)
	}
	common.DB.Where("user_uid=? and blog_uid=? and type=?", user.Uid, uid, 1).Last(&praise)
	if !reflect.DeepEqual(praise, models.Comment{}) {
		c.ErrorWithMessage("你已经点过👍了")
		return
	}
	var blog models.Blog
	common.DB.Where("uid=?", uid).Find(&blog)
	praiseJsonResult := common.RedisUtil.Get("BLOG_PRAISE:" + uid)
	if praiseJsonResult == "" {
		common.RedisUtil.Set("BLOG_PRAISE:"+uid, "1")
		blog.CollectCount = 1
		common.DB.Save(&blog)
	} else {
		count := blog.CollectCount + 1
		common.RedisUtil.Set("BLOG_PRAISE:"+uid, strconv.Itoa(count))
		blog.CollectCount = count
		common.DB.Save(&blog)
	}
	comment := models.Comment{
		Uid:     xid.New().String(),
		UserUid: user.Uid,
		BlogUid: uid,
		Source:  "BLOG_INFO",
		Type:    1,
	}
	common.DB.Create(&comment)
	c.SuccessWithData(blog.CollectCount)
}
