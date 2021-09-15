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
	"reflect"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/19 2:07 下午
 * @version 1.0
 */

type ResourceSortRestApi struct {
	base.BaseController
}

func (c *ResourceSortRestApi) GetList() {
	base.L.Print("获取分类资源列表")
	var resourceSortVO vo.ResourceSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &resourceSortVO)
	if err != nil {
		panic(err)
	}
	where := "status=?"
	if strings.TrimSpace(resourceSortVO.Keyword) != "" {
		where += " and sort_name like '%" + strings.TrimSpace(resourceSortVO.Keyword) + "%'"
	}
	var pageList []models.ResourceSort
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.ResourceSort{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((resourceSortVO.CurrentPage - 1) * resourceSortVO.PageSize).Limit(resourceSortVO.PageSize).Order("sort desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var s []string
	for _, item := range pageList {
		if item.FileUid != "" {
			s = append(s, item.FileUid+",")
		}
	}
	pictureResult := map[string]interface{}{}
	pictureMap := map[string]string{}
	fileUids := strings.Join(s, ",")
	if fileUids != "" {
		pictureResult = service.FileService.GetPicture(fileUids, ",")
	}
	picList := common.WebUtil.GetPictureMap(pictureResult)
	for _, item := range picList {
		pictureMap["uid"] = item["url"].(string)
	}
	for i, item := range pageList {
		if item.FileUid != "" {
			pictureUidsTemp := strings.Split(item.FileUid, ",")
			var pictureListTemp []string
			for _, picture := range pictureUidsTemp {
				pictureListTemp = append(pictureListTemp, pictureMap[picture])
			}
			pageList[i].PhotoList = pictureListTemp
		}
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    resourceSortVO.PageSize,
		Current: resourceSortVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *ResourceSortRestApi) Stick() {
	base.L.Print("置顶分类")
	var resourceSortVO vo.ResourceSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &resourceSortVO)
	if err != nil {
		panic(err)
	}
	var resourceSort models.ResourceSort
	common.DB.Where("uid=?", resourceSortVO.Uid).Find(&resourceSort)
	var pageList []models.ResourceSort
	common.DB.Offset(0).Limit(1).Order("sort desc").Find(&pageList)
	maxSort := pageList[0]
	switch {
	case maxSort.Uid == "":
		c.ErrorWithMessage("传入参数有误")
	case maxSort.Uid == resourceSort.Uid:
		c.ErrorWithMessage("该分类已经在顶端")
	default:
		sortCount := maxSort.Sort + 1
		resourceSort.Sort = sortCount
		common.DB.Save(&resourceSort)
		c.SuccessWithMessage("操作成功")
	}
}

func (c *ResourceSortRestApi) Edit() {
	base.L.Print("编辑资源分类")
	var resourceSortVO vo.ResourceSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &resourceSortVO)
	if err != nil {
		panic(err)
	}
	var resourceSort models.ResourceSort
	common.DB.Where("uid=?", resourceSortVO.Uid).Find(&resourceSort)
	var tempSort models.ResourceSort
	if resourceSort.SortName != resourceSortVO.SortName {
		common.DB.Where("sort_name=? and status=?", resourceSortVO.SortName, 1).First(&tempSort)
		if !reflect.DeepEqual(tempSort, models.ResourceSort{}) {
			c.ErrorWithMessage("该实体已存在")
			return
		}
	}
	if resourceSort.SortName == resourceSortVO.SortName || (resourceSort.SortName != resourceSortVO.SortName && reflect.DeepEqual(tempSort, models.ResourceSort{})) {
		resourceSort.SortName = resourceSortVO.SortName
		resourceSort.Content = resourceSortVO.Content
		resourceSort.FileUid = resourceSortVO.FileUid
		resourceSort.Sort = resourceSortVO.Sort
		common.DB.Save(&resourceSort)
		c.SuccessWithMessage("更新成功")
	}
}

func (c *ResourceSortRestApi) Add() {
	base.L.Print("增加资源分类")
	var resourceSortVO vo.ResourceSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &resourceSortVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		return
	}
	var tempSort models.ResourceSort
	common.DB.Where("sort_name=? and status=?", resourceSortVO.SortName, 1).First(&tempSort)
	if !reflect.DeepEqual(tempSort, models.ResourceSort{}) {
		c.ErrorWithMessage("该实体已存在")
		return
	}
	var resourceSort models.ResourceSort
	resourceSort.SortName = resourceSortVO.SortName
	resourceSort.Content = resourceSortVO.Content
	resourceSort.FileUid = resourceSortVO.FileUid
	resourceSort.Sort = resourceSortVO.Sort
	resourceSort.Uid = xid.New().String()
	common.DB.Create(&resourceSort)
	c.SuccessWithMessage("新增成功")
}

func (c *ResourceSortRestApi) DeleteBatch() {
	base.L.Print("批量删除资源分类")
	var resourceSortVOList []vo.ResourceSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &resourceSortVOList)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		return
	}
	if len(resourceSortVOList) == 0 {
		c.ErrorWithMessage("传入参数有误")
		return
	}
	var uids []string
	for _, item := range resourceSortVOList {
		uids = append(uids, item.Uid)
	}
	var studyVideo models.StudyVideo
	var count int64
	common.DB.Model(&studyVideo).Where("status=? and resource_sort_uid in ?", 1, uids).Count(&count)
	if count > 0 {
		c.ErrorWithMessage("该分类下还有资源")
	}
	var resourceSortList []models.ResourceSort
	common.DB.Find(&resourceSortList, uids)
	save := common.DB.Model(&resourceSortList).Select("status").Update("status", 0).Error
	if save == nil {
		c.SuccessWithMessage("删除成功")
	} else {
		c.ErrorWithMessage("删除失败")
	}
}
