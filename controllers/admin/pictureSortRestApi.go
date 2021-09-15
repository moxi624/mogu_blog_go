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
 * @date  2021/1/19 3:15 下午
 * @version 1.0
 */

type PictureSortRestApi struct {
	base.BaseController
}

func (c *PictureSortRestApi) GetList() {
	var pictureSortVO vo.PictureSortVO
	base.L.Print("获取图片分类列表")
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pictureSortVO)
	if err != nil {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	where := "status=?"
	if strings.TrimSpace(pictureSortVO.Keyword) != "" {
		where += " and name like '%" + strings.TrimSpace(pictureSortVO.Keyword) + "%'"
	}
	var pageList []models.PictureSort
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.PictureSort{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((pictureSortVO.CurrentPage - 1) * pictureSortVO.PageSize).Limit(pictureSortVO.PageSize).Order("sort desc").Find(&pageList)
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
		pictureMap[item["uid"].(string)] = item["url"].(string)
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
		Size:    pictureSortVO.PageSize,
		Current: pictureSortVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *PictureSortRestApi) Stick() {
	base.L.Print("置顶图片分类")
	var pictureSortVO vo.PictureSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pictureSortVO)
	if err != nil {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	var pictureSort models.PictureSort
	common.DB.Where("uid=?", pictureSortVO.Uid).Find(&pictureSort)
	var pageList []models.PictureSort
	common.DB.Offset(0).Limit(1).Order("sort desc").Find(&pageList)
	maxSort := pageList[0]
	switch {
	case maxSort.Uid == "":
		c.ErrorWithMessage("传入参数有误")
	case maxSort.Uid == pictureSort.Uid:
		c.ErrorWithMessage("该分类已经在顶端")
	default:
		sortCount := maxSort.Sort + 1
		pictureSort.Sort = sortCount
		common.DB.Save(&pictureSort)
		c.SuccessWithMessage("操作成功")
	}
}

func (c *PictureSortRestApi) Edit() {
	base.L.Print("编辑图片分类")
	var pictureSortVO vo.PictureSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pictureSortVO)
	if err != nil {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	var pictureSort models.PictureSort
	common.DB.Where("uid=?", pictureSortVO.Uid).Find(&pictureSort)
	pictureSort.Name = pictureSortVO.Name
	pictureSort.ParentUid = pictureSortVO.ParentUid
	pictureSort.Sort = pictureSortVO.Sort
	pictureSort.FileUid = pictureSortVO.FileUid
	pictureSort.IsShow = pictureSortVO.IsShow
	common.DB.Save(&pictureSort)
	c.SuccessWithMessage("更新成功")
}

func (c *PictureSortRestApi) Delete() {
	base.L.Print("删除图片分类")
	var pictureSortVO vo.PictureSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pictureSortVO)
	if err != nil {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	var pictureCount int64
	var picture models.Picture
	common.DB.Model(&picture).Where("status=? and picture_sort_uid in ?", 1, pictureSortVO.Uid).Count(&pictureCount)
	if pictureCount > 0 {
		c.ErrorWithMessage("该分类下还有图片")
		return
	}
	var pictureSort models.PictureSort
	common.DB.Where("uid=?", pictureSortVO.Uid).Find(&pictureSort)
	common.DB.Model(&pictureSort).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}

func (c *PictureSortRestApi) Add() {
	base.L.Print("增加图片分类")
	var pictureSortVO vo.PictureSortVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pictureSortVO)
	if err != nil {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	var pictureSort models.PictureSort
	pictureSort.Name = pictureSortVO.Name
	pictureSort.ParentUid = pictureSortVO.ParentUid
	pictureSort.Sort = pictureSortVO.Sort
	pictureSort.FileUid = pictureSortVO.FileUid
	pictureSort.IsShow = pictureSortVO.IsShow
	pictureSort.Uid = xid.New().String()
	common.DB.Create(&pictureSort)
	c.SuccessWithMessage("新增成功")
}
