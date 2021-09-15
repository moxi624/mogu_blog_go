package service

import (
	"mogu-go-v2/common"
	"mogu-go-v2/models"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/4 1:40 下午
 * @version 1.0
 */

type pictureService struct{}

func (pictureService) GetTopOne() models.Picture {
	var picture models.Picture
	common.DB.Where("status=?", 1).Order("create_time desc").Last(&picture)
	return picture
}

var PictureService = &pictureService{}
