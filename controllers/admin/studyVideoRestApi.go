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
 * @date  2021/1/20 4:03 下午
 * @version 1.0
 */

type StudyVideoRestApi struct {
	base.BaseController
}

func (c *StudyVideoRestApi) GetList() {
	base.L.Print("获取学习视频")
	var studyVideoVO vo.StudyVideoVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &studyVideoVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		return
	}
	where := "status=?"
	if strings.TrimSpace(studyVideoVO.Keyword) != "" {
		where += " and name like '%" + strings.TrimSpace(studyVideoVO.Keyword) + "%'"
	}
	var pageList []models.StudyVideo
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.StudyVideo{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((studyVideoVO.CurrentPage - 1) * studyVideoVO.PageSize).Limit(studyVideoVO.PageSize).Order("create_time desc").Find(&pageList)
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
		if item.ResourceSortUid != "" {
			var resourceSort models.ResourceSort
			common.DB.Where("uid=?", item.ResourceSortUid).Find(&resourceSort)
			pageList[i].ResourceSort = resourceSort
		}
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    studyVideoVO.PageSize,
		Current: studyVideoVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *StudyVideoRestApi) Edit() {
	base.L.Print("编辑学习视频")
	var studyVideeVO vo.StudyVideoVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &studyVideeVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		return
	}
	var studyVideo models.StudyVideo
	common.DB.Where("uid=?", studyVideeVO.Uid).Find(&studyVideo)
	studyVideo.Name = studyVideeVO.Name
	studyVideo.Summary = studyVideeVO.Summary
	studyVideo.Content = studyVideeVO.Content
	studyVideo.FileUid = studyVideeVO.FileUid
	studyVideo.BaiduPath = studyVideeVO.BaiduPath
	studyVideo.ResourceSortUid = studyVideeVO.ResouceSortUid
	common.DB.Save(&studyVideo)
	c.SuccessWithMessage("更新成功")
}

func (c *StudyVideoRestApi) DeleteBatch() {
	base.L.Print("删除学习视频")
	var studyVideeVOList []vo.StudyVideoVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &studyVideeVOList)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		return
	}
	if len(studyVideeVOList) == 0 {
		c.ErrorWithMessage("参数传入有错误")
		return
	}
	var uids []string
	for _, item := range studyVideeVOList {
		uids = append(uids, item.Uid)
	}
	var blogSortList []models.StudyVideo
	common.DB.Find(&blogSortList, uids)
	save := common.DB.Model(&blogSortList).Select("status").Update("status", 0).Error
	if save == nil {
		c.SuccessWithMessage("删除成功")
	} else {
		c.ErrorWithMessage("删除失败")
	}
}

func (c *StudyVideoRestApi) Add() {
	base.L.Print("增加学习视频")
	var studyVideeVO vo.StudyVideoVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &studyVideeVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		return
	}
	var studyVideo models.StudyVideo
	studyVideo.Name = studyVideeVO.Name
	studyVideo.Summary = studyVideeVO.Summary
	studyVideo.Content = studyVideeVO.Content
	studyVideo.FileUid = studyVideeVO.FileUid
	studyVideo.BaiduPath = studyVideeVO.BaiduPath
	studyVideo.ResourceSortUid = studyVideeVO.ResouceSortUid
	studyVideo.ClickCount = 0
	studyVideo.Uid = xid.New().String()
	common.DB.Create(&studyVideo)
	c.SuccessWithMessage("新增成功")
}
