package service

import (
	"mogu-go-v2/common"
	"mogu-go-v2/models"
	"reflect"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/23 8:40 上午
 * @version 1.0
 */

type userService struct{}

func (userService) GetUserListByIds(ids []string) []models.User {
	var userList []models.User
	if reflect.DeepEqual(ids, []string{}) {
		return userList
	}
	var userCollection []models.User
	common.DB.Find(&userCollection, ids)
	for _, item := range userCollection {
		userList = append(userList, item)
	}
	return userList
}

var UserService = &userService{}
