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
 * @date  2021/2/1 11:01 上午
 * @version 1.0
 */

type LinkRestApi struct {
	base.BaseController
}

func (c *LinkRestApi) GetList() {
	base.L.Print("获取友情链接列表")
	var linkVO vo.LinkVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &linkVO)
	if err != nil {
		panic(err)
	}
	where := "status != ?"
	if strings.TrimSpace(linkVO.Keyword) != "" {
		where += " and title like \"%" + strings.TrimSpace(linkVO.Keyword) + "%\""
	}
	if linkStatus := common.InterfaceToString(linkVO.LinkStatus); linkStatus != "" {
		where += " and link_status=" + linkStatus
	}
	var pageList []models.Link
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Link{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 0).Offset((linkVO.CurrentPage - 1) * linkVO.PageSize).Limit(linkVO.PageSize).Order("sort desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var s []string
	for _, item := range pageList {
		if item.FileUid != "" {
			s = append(s, item.FileUid)
		}
	}
	fileUids := strings.Join(s, ",")
	pictureMap := map[string]string{}
	pictureResult := map[string]interface{}{}
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
		Size:    linkVO.PageSize,
		Current: linkVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *LinkRestApi) Stick() {
	var linkVO vo.LinkVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &linkVO)
	if err != nil {
		panic(err)
	}
	var link models.Link
	common.DB.Where("uid=?", linkVO.Uid).Find(&link)
	var pageList []models.Link
	common.DB.Offset(0).Limit(1).Order("sort desc").Find(&pageList)
	maxSort := pageList[0]
	switch {
	case maxSort.Uid == "":
		c.ErrorWithMessage("传入参数有误")
	case maxSort.Uid == link.Uid:
		c.ErrorWithMessage("该链接已经在顶端")
	default:
		sortCount := maxSort.Sort + 1
		link.Sort = sortCount
		common.DB.Save(&link)
		c.SuccessWithMessage("操作成功")
	}
}

func (c *LinkRestApi) Edit() {
	var linkVO vo.LinkVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &linkVO)
	if err != nil {
		panic(err)
	}
	var link models.Link
	common.DB.Where("uid = ?", linkVO.Uid).Find(&link)
	linkStatus := common.InterfaceToInt(linkVO.LinkStatus)
	link.Title = linkVO.Title
	link.Summary = linkVO.Summary
	link.LinkStatus = linkStatus
	link.Url = linkVO.Url
	link.Sort = linkVO.Sort
	link.Email = linkVO.Email
	link.FileUid = linkVO.FileUid
	common.DB.Save(&link)
	c.SuccessWithMessage("更新成功")
}

func (c *LinkRestApi) Delete() {
	var linkVO vo.LinkVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &linkVO)
	if err != nil {
		panic(err)
	}
	var link models.Link
	common.DB.Where("uid=?", linkVO.Uid).Find(&link)
	common.DB.Model(&link).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}

func (c *LinkRestApi) Add() {
	var linkVO vo.LinkVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &linkVO)
	if err != nil {
		panic(err)
	}
	linkStatus := common.InterfaceToInt(linkVO.LinkStatus)
	link := models.Link{
		Uid:        xid.New().String(),
		Summary:    linkVO.Summary,
		Title:      linkVO.Title,
		Url:        linkVO.Url,
		LinkStatus: linkStatus,
		Sort:       linkVO.Sort,
		Email:      linkVO.Email,
		FileUid:    linkVO.FileUid,
	}
	common.DB.Create(&link)
	c.SuccessWithMessage("新增成功")
}
