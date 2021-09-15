package service

import (
	"mogu-go-v2/common"
	"mogu-go-v2/models"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/4 11:23 上午
 * @version 1.0
 */

type tagService struct{}

func (tagService) GetTopTag() models.Tag {
	var tag models.Tag
	common.DB.Where("status=?", 1).Order("sort desc").Last(&tag)
	return tag
}

var TagService = &tagService{}
