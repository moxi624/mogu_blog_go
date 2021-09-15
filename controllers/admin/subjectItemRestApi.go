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
	"sort"
	"strconv"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/5 10:09 上午
 * @version 1.0
 */

type SubjectItemRestApi struct {
	base.BaseController
}

func (c *SubjectItemRestApi) Add() {
	var subjectItemVOList []vo.SubjectItemVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subjectItemVOList)
	if err != nil {
		panic(err)
	}
	var blogUidList []string
	var subjectUid string
	for _, subjectItemVO := range subjectItemVOList {
		blogUidList = append(blogUidList, subjectItemVO.BlogUid)
		if subjectUid == "" && subjectItemVO.SubjectUid != "" {
			subjectUid = subjectItemVO.SubjectUid
		}
	}
	var reapeatSubjectItemList []models.SubjectItem
	common.DB.Where(`subject_uid=? and blog_uid in ?`, subjectUid, blogUidList).Find(&reapeatSubjectItemList)
	var repeatBlogList []string
	for _, item := range reapeatSubjectItemList {
		repeatBlogList = append(repeatBlogList, item.BlogUid)
	}
	var subjectItemList []models.SubjectItem
	for _, subjectItemVO := range subjectItemVOList {
		if subjectItemVO.SubjectUid == "" || subjectItemVO.BlogUid == "" {
			c.ErrorWithMessage("参数错误")
			return
		}
		if common.SliceFind(repeatBlogList, subjectItemVO.BlogUid) {
			continue
		} else {
			subjectItem := models.SubjectItem{
				Uid:        xid.New().String(),
				SubjectUid: subjectItemVO.SubjectUid,
				BlogUid:    subjectItemVO.BlogUid,
			}
			subjectItemList = append(subjectItemList, subjectItem)
		}
	}
	if len(subjectItemList) == 0 {
		if len(repeatBlogList) == 0 {
			c.ErrorWithMessage("操作失败")
			return
		}
		c.ErrorWithMessage("操作失败，已跳过" + strconv.Itoa(len(repeatBlogList)) + "个重复数据")
		return
	}
	common.DB.Create(&subjectItemList)
	if len(repeatBlogList) == 0 {
		c.SuccessWithMessage("操作成功")
	} else {
		c.SuccessWithMessage("操作成功，已跳过" + strconv.Itoa(len(repeatBlogList)) + "个重复数据，成功新增" + strconv.Itoa(len(subjectItemList)-len(repeatBlogList)) + `条数据`)
	}
}

func (c *SubjectItemRestApi) GetList() {
	var subjectItemVO vo.SubjectItemVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subjectItemVO)
	if err != nil {
		panic(err)
	}
	where := "status=?"
	if subjectItemVO.SubjectUid != "" {
		where += " and subject_uid ='" + subjectItemVO.SubjectUid + "'"
	}
	var pageList []models.SubjectItem
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.SubjectItem{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((subjectItemVO.CurrentPage - 1) * subjectItemVO.PageSize).Limit(subjectItemVO.PageSize).Order("sort desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var blogUidList []string
	for _, item := range pageList {
		blogUidList = append(blogUidList, item.BlogUid)
	}
	var blogCollection []models.Blog
	if len(blogUidList) > 0 {
		common.DB.Find(&blogCollection, blogUidList)
		if len(blogCollection) > 0 {
			blogList := service.BlogService.SetTagAndSortAndPictureByBlogList(blogCollection)
			blogMap := map[string]models.Blog{}
			for _, item := range blogList {
				blogMap[item.Uid] = item
			}
			for i, item := range pageList {
				pageList[i].Blog = blogMap[item.BlogUid]
			}
		}
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    subjectItemVO.PageSize,
		Current: subjectItemVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c SubjectItemRestApi) DeleteBatch() {
	var subjectItemVOList []vo.SubjectItemVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subjectItemVOList)
	if err != nil {
		panic(err)
	}
	if len(subjectItemVOList) == 0 {
		c.ErrorWithMessage("参数错误")
		return
	}
	var uids []string
	for _, item := range subjectItemVOList {
		uids = append(uids, item.Uid)
	}
	common.DB.Model(&models.SubjectItem{}).Where("uid in ?", uids).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}

func (c *SubjectItemRestApi) Edit() {
	var subjectItemVOList []vo.SubjectItemVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subjectItemVOList)
	if err != nil {
		panic(err)
	}
	var subjectItemUidList []string
	for _, item := range subjectItemVOList {
		subjectItemUidList = append(subjectItemUidList, item.Uid)
	}
	var subjectItemCollection []models.SubjectItem
	if len(subjectItemUidList) > 0 {
		common.DB.Find(&subjectItemCollection, subjectItemUidList)
		if len(subjectItemCollection) > 0 {
			subjectItemVOHashMap := map[string]vo.SubjectItemVO{}
			for _, item := range subjectItemVOList {
				subjectItemVOHashMap[item.Uid] = item
			}
			var subjectItemList []models.SubjectItem
			for _, item := range subjectItemCollection {
				subjectItemVO := subjectItemVOHashMap[item.Uid]
				item.SubjectUid = subjectItemVO.SubjectUid
				item.BlogUid = subjectItemVO.BlogUid
				item.Sort = subjectItemVO.Sort
				subjectItemList = append(subjectItemList, item)
			}
			common.DB.Save(&subjectItemList)
		}
	}
	c.SuccessWithData("更新成功")
}

func (c *SubjectItemRestApi) SortByCreateTime() {
	base.L.Print("通过点击量排序博客分类")
	subjectUid := c.GetString("subjectUid")
	isDesc, _ := c.GetBool("isDesc")
	var subjectItemList []models.SubjectItem
	common.DB.Where("status=? and subject_uid=?", 1, subjectUid).Find(&subjectItemList)
	var blogUidList []string
	for _, item := range subjectItemList {
		blogUidList = append(blogUidList, item.BlogUid)
	}
	if len(blogUidList) == 0 {
		c.SuccessWithData("排序失败")
		return
	}
	var blogList []models.Blog
	common.DB.Find(&blogList, blogUidList)
	if isDesc {
		sort.SliceStable(blogList, func(i, j int) bool {
			return common.IsTimeAAfterTimeB(blogList[i].CreateTime, blogList[j].CreateTime)
		})
	} else {
		sort.SliceStable(blogList, func(i, j int) bool {
			return common.IsTimeABeforeTimeB(blogList[i].CreateTime, blogList[j].CreateTime)
		})
	}
	maxSort := len(blogList)
	subjectItemSortMap := map[string]int{}
	for _, item := range blogList {
		maxSort--
		subjectItemSortMap[item.Uid] = maxSort
	}
	for i, item := range subjectItemList {
		subjectItemList[i].Sort = subjectItemSortMap[item.BlogUid]
	}
	common.DB.Save(&subjectItemList)
	c.SuccessWithMessage("排序成功")
}
