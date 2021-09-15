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
 * @date  2021/1/19 4:10 下午
 * @version 1.0
 */

type PictureRestApi struct {
	base.BaseController
}

func (c *PictureRestApi) GetList() {
	var pictureVO vo.PictureVO
	base.L.Print("获取图片列表")
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pictureVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
	}
	where := "status=? and picture_sort_uid=?"
	if strings.TrimSpace(pictureVO.Keyword) != "" {
		where += " and sort_name like '%" + strings.TrimSpace(pictureVO.Keyword) + "%'"
	}
	var pageList []models.Picture
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Picture{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1, pictureVO.PictureSortUid).Offset((pictureVO.CurrentPage - 1) * pictureVO.PageSize).Limit(pictureVO.PageSize).Order("create_time desc").Find(&pageList)
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
			pageList[i].PictureUrl = pictureMap[item.FileUid]
		}
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    pictureVO.PageSize,
		Current: pictureVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *PictureRestApi) Delete() {
	base.L.Print("删除图片")
	var pictureVO vo.PictureVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pictureVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		panic(err)
	}
	uidStr := pictureVO.Uid
	if uidStr == "" {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	uids := strings.Split(pictureVO.Uid, ",")
	var pictureList []models.Picture
	common.DB.Find(&pictureList, uids)
	common.DB.Model(&pictureList).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}

func (c *PictureRestApi) Add() {
	base.L.Print("增加图片")
	var pictureVOList []vo.PictureVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pictureVOList)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		return
	}
	var pictureList []models.Picture
	if len(pictureVOList) > 0 {
		for _, pictureVO := range pictureVOList {
			var picture models.Picture
			picture.FileUid = pictureVO.FileUid
			picture.PictureSortUid = pictureVO.PictureSortUid
			picture.PicName = pictureVO.PicName
			picture.Status = 1
			picture.Uid = xid.New().String()
			pictureList = append(pictureList, picture)
		}
		common.DB.Create(&pictureList)
		c.SuccessWithMessage("新增成功")
	} else {
		c.ErrorWithMessage("新增失败")
	}
}

func (c *PictureRestApi) SetCover() {
	base.L.Print("设置图片分类封面")
	var pictureVO vo.PictureVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pictureVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		return
	}
	var pictureSort models.PictureSort
	common.DB.Where("uid=?", pictureVO.PictureSortUid).Find(&pictureSort)
	var picture models.Picture
	if !reflect.DeepEqual(pictureSort, models.PictureSort{}) {
		common.DB.Where("uid=?", pictureVO.Uid).Find(&picture)
		if picture != (models.Picture{}) {
			pictureSort.FileUid = picture.FileUid
			common.DB.Save(&pictureSort)
		} else {
			c.ErrorWithMessage("图片不存在")
			return
		}
	}
	if reflect.DeepEqual(pictureSort, models.PictureSort{}) {
		c.ErrorWithMessage("图片分类不存在")
	} else {
		c.SuccessWithMessage("更新成功")
	}
}
