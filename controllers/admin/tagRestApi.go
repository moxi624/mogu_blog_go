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
 * @date  2021/1/11 9:13 上午
 * @version 1.0
 */
type TagRestApi struct {
	base.BaseController
}

func (c *TagRestApi) GetList() {
	var tagVO vo.TagVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &tagVO)
	if err != nil {
		panic(err)
	}
	base.L.Print("获取标签列表")
	where := "status=?"
	order := "sort desc"
	if strings.TrimSpace(tagVO.Keyword) != "" {
		where += " and content like '%" + strings.TrimSpace(tagVO.Keyword) + "%'"
	}
	if tagVO.OrderByAscColumn != "" {
		order = common.Camel2Case(tagVO.OrderByAscColumn) + " asc"
	} else if tagVO.OrderByDescColumn != "" {
		order = common.Camel2Case(tagVO.OrderByDescColumn) + " desc"
	}
	var pageList []models.Tag
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Tag{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((tagVO.CurrentPage - 1) * tagVO.PageSize).Limit(tagVO.PageSize).Order(order).Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    tagVO.PageSize,
		Current: tagVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *TagRestApi) Stick() {
	base.L.Print("置顶标签")
	var tagVO vo.TagVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &tagVO)
	if err != nil {
		panic(err)
	}
	var tag models.Tag
	common.DB.Where("uid=?", tagVO.Uid).Find(&tag)
	var pageList []models.Tag
	common.DB.Offset(0).Limit(1).Order("sort desc").Find(&pageList)
	maxSort := pageList[0]
	switch maxSort.Uid {
	case "":
		c.ErrorWithMessage("传入参数有误")
	case tag.Uid:
		c.ErrorWithMessage("该分类已经在顶端")
	default:
		tag.Sort = maxSort.Sort + 1
		common.DB.Save(&tag)
		c.SuccessWithMessage("置顶成功")
	}
}

func (c *TagRestApi) Edit() {
	base.L.Print("编辑标签")
	var tagVO vo.TagVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &tagVO)
	if err != nil {
		panic(err)
	}
	var tag models.Tag
	common.DB.Where("uid=?", tagVO.Uid).Find(&tag)
	var tempTag models.Tag
	if tag != (models.Tag{}) && tag.Content == tagVO.Content {
		common.DB.Where("content=? and status=?", tagVO.Content, 1).First(&tempTag)
		if tempTag != (models.Tag{}) {
			c.ErrorWithMessage("记录已存在")
			return
		}
	}
	tag.Content = tagVO.Content
	tag.Sort = tagVO.Sort
	common.DB.Save(&tag)
	service.BlogService.DeleteRedisByBlogTag()
	c.SuccessWithMessage("更新成功")
}

func (c *TagRestApi) DeleteBatch() {
	base.L.Print("批量删除标签")
	var tagVOList []vo.TagVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &tagVOList)
	if err != nil {
		panic(err)
	}
	if len(tagVOList) == 0 {
		c.ErrorWithMessage("参数错误")
		return
	}
	var uids []string
	for _, item := range tagVOList {
		uids = append(uids, item.Uid)
	}
	var blogCount int64
	common.DB.Model(&models.Blog{}).Where("status=? and blog_sort_uid in ?", 1, uids).Count(&blogCount)
	if blogCount > 0 {
		c.ErrorWithMessage("该分类下还有博客")
		return
	}
	var tagList []models.Tag
	common.DB.Find(&tagList, uids)
	save := common.DB.Model(&tagList).Select("status").Update("status", 0).Error
	if save == nil {
		service.BlogService.DeleteRedisByBlogTag()
		c.SuccessWithMessage("删除成功")
	} else {
		c.ErrorWithMessage("删除失败")
	}
}

func (c *TagRestApi) Add() {
	base.L.Print("编辑标签")
	var tagVO vo.TagVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &tagVO)
	if err != nil {
		panic(err)
	}
	var tempTag models.Tag
	common.DB.Where("content=? and status=?", tagVO.Content, 1).First(&tempTag)
	if tempTag != (models.Tag{}) {
		c.ErrorWithMessage("数据已存在")
		return
	}
	tag := models.Tag{
		Uid:     xid.New().String(),
		Content: tagVO.Content,
		Sort:    tagVO.Sort,
	}
	common.DB.Create(&tag)
	c.SuccessWithMessage("新增成功")
}

func (c *TagRestApi) TagSortByClickCount() {
	base.L.Print("通过点击量排序标签")
	var tagList []models.Tag
	common.DB.Where("status=?", 1).Order("click_count desc").Find(&tagList)
	for i, item := range tagList {
		tagList[i].Sort = item.ClickCount
	}
	common.DB.Save(&tagList)
	c.SuccessWithMessage("排序成功")
}

func (c *TagRestApi) TagSortByCite() {
	base.L.Print("通过引用量排序标签")
	m := map[string]int{}
	var tagList []models.Tag
	common.DB.Where("status=?", 1).Find(&tagList)
	for _, item := range tagList {
		m[item.Uid] = 0
	}
	var blogList []models.Blog
	common.DB.Where("status=? and is_publish=?", 1, "1").Select("tag_uid").Find(&blogList)
	for _, item := range blogList {
		tagUids := item.TagUid
		tagUidList := strings.Split(tagUids, ",")
		for _, tagUid := range tagUidList {
			if _, ok := m[tagUid]; ok {
				count := m[tagUid] + 1
				m[tagUid] = count
			} else {
				m[tagUid] = 0
			}
		}
	}
	for i, item := range tagList {
		tagList[i].Sort = m[item.Uid]
	}
	common.DB.Save(&tagList)
	c.SuccessWithMessage("排序成功")
}
