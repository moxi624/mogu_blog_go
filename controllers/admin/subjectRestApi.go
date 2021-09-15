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
 * @date  2021/2/5 9:15 上午
 * @version 1.0
 */

type SubjectRestApi struct {
	base.BaseController
}

func (c *SubjectRestApi) GetList() {
	var subjectVO vo.SubjectVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subjectVO)
	if err != nil {
		panic(err)
	}
	where := "status=?"
	if strings.TrimSpace(subjectVO.Keyword) != "" {
		where += " and subject_name like \"%" + strings.TrimSpace(subjectVO.Keyword) + "%\""
	}
	var pageList []models.Subject
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Subject{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((subjectVO.CurrentPage - 1) * subjectVO.PageSize).Limit(subjectVO.PageSize).Order("sort desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var s []string
	for _, item := range pageList {
		s = append(s, item.FileUid)
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
				if pictureMap[picture] != "" {
					pictureListTemp = append(pictureListTemp, pictureMap[picture])
				}
			}
			pageList[i].PhotoList = pictureListTemp
		}
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    subjectVO.PageSize,
		Current: subjectVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *SubjectRestApi) Edit() {
	var subjectVO vo.SubjectVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subjectVO)
	if err != nil {
		panic(err)
	}
	var subject models.Subject
	common.DB.Where("uid=?", subjectVO.Uid).Find(&subject)
	var tempSubject models.Subject
	if subject.SubjectName != subjectVO.SubjectName {
		common.DB.Where("subject_name=? and status=?", subjectVO.SubjectName, 1).First(&tempSubject)
		if !reflect.DeepEqual(tempSubject, models.Subject{}) {
			c.ErrorWithMessage("相同数据已存在")
			return
		}
	}
	clickCount := common.InterfaceToInt(subjectVO.ClickCount)
	collectCount := common.InterfaceToInt(subjectVO.CollectCount)
	subject.SubjectName = subjectVO.SubjectName
	subject.Summary = subjectVO.Summary
	subject.FileUid = subjectVO.FileUid
	subject.ClickCount = clickCount
	subject.CollectCount = collectCount
	subject.Sort = subjectVO.Sort
	common.DB.Save(&subject)
	c.SuccessWithMessage("更新成功")
}

func (c *SubjectRestApi) DeleteBatch() {
	var subjectVOList []vo.SubjectVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subjectVOList)
	if err != nil {
		panic(err)
	}
	if len(subjectVOList) == 0 {
		c.ErrorWithMessage("参数传入错误")
		return
	}
	var uids []string
	for _, item := range subjectVOList {
		uids = append(uids, item.Uid)
	}
	var count int64
	common.DB.Model(&models.SubjectItem{}).Where("status=? and subject_uid in ?", 1, uids).Count(&count)
	if count > 0 {
		c.ErrorWithMessage("该专题下还有内容")
		return
	}
	var subjectList []models.Subject
	common.DB.Find(&subjectList, uids)
	save := common.DB.Model(&subjectList).Select("status").Update("status", 0).Error
	if save == nil {
		c.SuccessWithMessage("删除成功")
	} else {
		c.ErrorWithMessage("删除失败")
	}
}

func (c *SubjectRestApi) Add() {
	var subjectVO vo.SubjectVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subjectVO)
	if err != nil {
		panic(err)
	}
	var tempSubject models.Subject
	common.DB.Where("subject_name=? and status=?", subjectVO.SubjectName, 1).Last(&tempSubject)
	if !reflect.DeepEqual(tempSubject, models.Subject{}) {
		c.ErrorWithMessage("数据已经存在")
		return
	}
	clickCount := common.InterfaceToInt(subjectVO.ClickCount)
	collectCount := common.InterfaceToInt(subjectVO.CollectCount)
	subject := models.Subject{
		Uid:          xid.New().String(),
		SubjectName:  subjectVO.SubjectName,
		Summary:      subjectVO.Summary,
		FileUid:      subjectVO.FileUid,
		ClickCount:   clickCount,
		CollectCount: collectCount,
		Sort:         subjectVO.Sort,
	}
	common.DB.Create(&subject)
	c.SuccessWithMessage("新增成功")
}
