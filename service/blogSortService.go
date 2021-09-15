package service

import (
	"mogu-go-v2/common"
	"mogu-go-v2/models"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/4 10:57 上午b
 * @version 1.0
 */

type blogSortService struct{}

func (blogSortService) GetTopOne() models.BlogSort {
	var blogSort models.BlogSort
	common.DB.Where("status=?", 1).Order("sort desc").Last(&blogSort)
	return blogSort
}

var BlogSortService = &blogSortService{}
