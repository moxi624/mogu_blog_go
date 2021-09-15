package admin

import (
	"encoding/json"
	"github.com/rs/xid"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/11 9:42 上午
 * @version 1.0
 */
type BlogSortRestApi struct {
	base.BaseController
}

func (c *BlogSortRestApi) GetList() {
	var blogSortVO vo.BlogSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogSortVO)
	if err != nil {
		panic(err)
	}
	base.L.Print("获取博客分类列表")
	where := "status=?"
	if strings.TrimSpace(blogSortVO.Keyword) != "" {
		where += " and sort_name like '%" + strings.TrimSpace(blogSortVO.Keyword) + "%'"
	}
	var pageList []models.BlogSort
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.BlogSort{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((blogSortVO.CurrentPage - 1) * blogSortVO.PageSize).Limit(blogSortVO.PageSize).Order("sort desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    blogSortVO.PageSize,
		Current: blogSortVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *BlogSortRestApi) Stick() {
	base.L.Print("置顶分类")
	var blogSortVO vo.BlogSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogSortVO)
	if err != nil {
		panic(err)
	}
	var blogSort models.BlogSort
	common.DB.Where("uid=?", blogSortVO.Uid).Find(&blogSort)
	var pageList []models.BlogSort
	common.DB.Offset(0).Limit(1).Order("sort desc").Find(&pageList)
	maxSort := pageList[0]
	switch {
	case maxSort.Uid == "":
		c.ErrorWithMessage("传入参数有误")
	case maxSort.Uid == blogSort.Uid:
		c.ErrorWithMessage("该分类已经在顶端")
	default:
		blogSort.Sort = maxSort.Sort + 1
		common.DB.Save(&blogSort)
		c.SuccessWithMessage("操作成功")
	}
}

func (c *BlogSortRestApi) Edit() {
	base.L.Print("编辑博客分类")
	var blogSortVO vo.BlogSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogSortVO)
	if err != nil {
		panic(err)
	}
	var blogSort models.BlogSort
	common.DB.Where("uid=?", blogSortVO.Uid).Find(&blogSort)
	var tempSort models.BlogSort
	if blogSort.SortName != blogSortVO.SortName {
		common.DB.Where("sort_name=? and status=?", blogSortVO.SortName, 1).First(&tempSort)
		if tempSort != (models.BlogSort{}) {
			c.ErrorWithMessage("数据已存在")
			return
		}
	}
	blogSort.Content = blogSortVO.Content
	blogSort.SortName = blogSortVO.SortName
	blogSort.Sort = blogSortVO.Sort
	common.DB.Save(&blogSort)
	service.BlogService.DeleteRedisByBlogSort()
	c.SuccessWithMessage("更新成功")
}

func (c *BlogSortRestApi) DeleteBatch() {
	base.L.Print("批量删除博客分类")
	var blogSortVOList []vo.BlogSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogSortVOList)
	if err != nil {
		panic(err)
	}
	if len(blogSortVOList) == 0 {
		c.ErrorWithMessage("参数错误")
		return
	}
	var uids []string
	for _, item := range blogSortVOList {
		uids = append(uids, item.Uid)
	}
	var blogCount int64
	common.DB.Model(&models.Blog{}).Where("status=? and blog_sort_uid in ?", 1, uids).Count(&blogCount)
	if blogCount > 0 {
		c.ErrorWithMessage("该分类下还有博客")
		return
	}
	var blogSortList []models.BlogSort
	common.DB.Find(&blogSortList, uids)
	save := common.DB.Model(&blogSortList).Select("status").Update("status", 0).Error
	if save == nil {
		service.BlogService.DeleteRedisByBlogSort()
		c.SuccessWithMessage("删除成功")
	} else {
		c.ErrorWithMessage("删除失败")
	}
}

func (c *BlogSortRestApi) Add() {
	base.L.Print("增加博客分类")
	var blogSortVO vo.BlogSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogSortVO)
	if err != nil {
		panic(err)
	}
	var tempSort models.BlogSort
	common.DB.Where("sort_name=? and status=?", blogSortVO.SortName, 1).First(&tempSort)
	if tempSort != (models.BlogSort{}) {
		c.ErrorWithMessage("数据已存在")
		return
	}
	blogSort := models.BlogSort{
		Uid:      xid.New().String(),
		Content:  blogSortVO.Content,
		SortName: blogSortVO.SortName,
		Sort:     blogSortVO.Sort,
	}
	common.DB.Create(&blogSort)
	c.SuccessWithMessage("新增成功")
}

func (c *BlogSortRestApi) BlogSortByClickCount() {
	base.L.Print("通过点击量排序博客分类")
	var blogSortList []models.BlogSort
	common.DB.Where("status=?", 1).Order("click_count desc").Find(&blogSortList)
	for i, item := range blogSortList {
		blogSortList[i].Sort = item.ClickCount
	}
	common.DB.Save(&blogSortList)
	c.SuccessWithMessage("排序成功")
}

func (c *BlogSortRestApi) BlogSortByCite() {
	base.L.Print("通过引用量排序博客分类")
	m := map[string]int{}
	var blogSortList []models.BlogSort
	common.DB.Where("status=?", 1).Find(&blogSortList)
	for _, item := range blogSortList {
		m[item.Uid] = 0
	}
	var blogList []models.Blog
	common.DB.Where("status=? and is_publish=?", 1, "1").Select("blog_sort_uid").Find(&blogList)
	for _, item := range blogList {
		blogSortUid := item.BlogSortUid
		if _, ok := m[blogSortUid]; ok {
			count := m[blogSortUid] + 1
			m[blogSortUid] = count
		} else {
			m[blogSortUid] = 0
		}
	}
	for i, item := range blogSortList {
		blogSortList[i].Sort = m[item.Uid]
	}
	common.DB.Save(&blogSortList)
	c.SuccessWithMessage("排序成功")
}
