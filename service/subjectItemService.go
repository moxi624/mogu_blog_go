package service

import (
	"mogu-go-v2/common"
	"mogu-go-v2/models"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/3 8:45 上午
 * @version 1.0
 */

type subjectItemService struct{}

func (subjectItemService) DeleteBatchSubjectItemByBlogUid(blogUid []string) string {
	checkSuccess := common.CheckUidList(blogUid)
	if !checkSuccess {
		return "Uid不合法"
	}
	var subjectItem []models.SubjectItem
	common.DB.Where("blog_uid in ?", blogUid).Find(&subjectItem)
	if len(subjectItem) > 0 {
		common.DB.Model(&subjectItem).Select("status").Update("status", 0)
		return "删除成功"
	}
	return "记录不存在"
}

var SubjectItemService = &subjectItemService{}
