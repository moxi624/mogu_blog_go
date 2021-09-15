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
	"strconv"
	"strings"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/9 9:48 上午
 * @version 1.0
 */

type UserRestApi struct {
	base.BaseController
}

var t, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")

func (c *UserRestApi) GetList() {
	var userVO vo.UserVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userVO)
	if err != nil {
		panic(err)
	}
	base.L.Print("获取用户列表")
	where := "status != ?"
	if strings.TrimSpace(userVO.Keyword) != "" {
		where += " and (user_name like \"%" + strings.TrimSpace(userVO.Keyword) + "%\" or nick_name like \"%" + strings.TrimSpace(userVO.Keyword) + "%\")"
	}
	if strings.TrimSpace(userVO.Source) != "" {
		where += " and source=\"" + strings.TrimSpace(userVO.Source) + "\""
	}
	if commentStatus := common.InterfaceToInt(userVO.CommentStatus); commentStatus != 0 {
		where += " and comment_status=" + strconv.Itoa(commentStatus)
	}
	var pageList []models.User
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.User{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 0).Offset((userVO.CurrentPage - 1) * userVO.PageSize).Limit(userVO.PageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var s []string
	for _, item := range pageList {
		if item.Avatar != "" {
			s = append(s, item.Avatar)
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
		if item.Avatar != "" {
			pictureUidsTemp := strings.Split(item.Avatar, ",")
			var pictureListTemp []string
			for _, picture := range pictureUidsTemp {
				if pictureMap[picture] != "" {
					pictureListTemp = append(pictureListTemp, pictureMap[picture])
				}
			}
			if len(pictureListTemp) > 0 {
				pageList[i].PhotoUrl = pictureListTemp[0]
			}
		}
		pageList[i].PassWord = ""
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    userVO.PageSize,
		Current: userVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *UserRestApi) Edit() {
	base.L.Print("编辑用户")
	var userVO vo.UserVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userVO)
	if err != nil {
		panic(err)
	}
	var user models.User
	common.DB.Where("uid=?", userVO.Uid).Find(&user)
	g := common.InterfaceToInt(userVO.Gender)
	commentStatus := common.InterfaceToInt(userVO.CommentStatus)
	user.UserName = userVO.UserName
	user.Email = userVO.Email
	user.StartEmailNotification = userVO.StartEmailNotification
	user.Occupation = userVO.Occupation
	user.Gender = g
	user.QqNumber = userVO.QqNumber
	user.Summary = userVO.Summary
	user.Avatar = userVO.Avatar
	user.NickName = userVO.NickName
	user.UserTag = userVO.UserTag
	user.CommentStatus = commentStatus
	if userVO.Birthday.Format("2006-01-02") == "0001-01-01" {
		user.Birthday = t
	} else {
		user.Birthday = userVO.Birthday
	}
	common.DB.Save(&user)
	c.SuccessWithMessage("更新成功")
}

func (c *UserRestApi) ResetUserPassword() {
	base.L.Print("重置用户密码")
	var userVO vo.UserVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userVO)
	if err != nil {
		panic(err)
	}
	defaultPassword, b := service.SysParamService.GetSysParamsValueByKey("SYS_DEFAULT_PASSWORD")
	if !b {
		c.BaseController.ThrowError("00106", defaultPassword)
		return
	}
	var user models.User
	common.DB.Where("uid = ?", userVO.Uid).Find(&user)
	user.PassWord = common.SHA256(defaultPassword)
	if userVO.Birthday.Format("2006-01-02") == "0001-01-01" {
		user.Birthday = t
	} else {
		user.Birthday = userVO.Birthday
	}
	if user.LastLoginTime.Format("2006-01-02") == "0001-01-01" {
		user.LastLoginTime = t
	}
	common.DB.Save(&user)
	c.SuccessWithMessage("操作成功")
}

func (c UserRestApi) Delete() {
	base.L.Print("删除用户")
	var userVO vo.UserVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userVO)
	if err != nil {
		panic(err)
	}
	var user models.User
	common.DB.Where("uid = ?", userVO.Uid).Find(&user)
	common.DB.Model(&user).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}

func (c UserRestApi) Add() {
	base.L.Print("新增用户")
	var userVO vo.UserVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userVO)
	if err != nil {
		panic(err)
	}
	commentStatus := common.InterfaceToInt(userVO.CommentStatus)
	g := common.InterfaceToInt(userVO.Gender)
	user := models.User{
		Uid:                    xid.New().String(),
		UserName:               userVO.UserName,
		Email:                  userVO.Email,
		StartEmailNotification: userVO.StartEmailNotification,
		Occupation:             userVO.Occupation,
		Gender:                 g,
		QqNumber:               userVO.QqNumber,
		Summary:                userVO.Summary,
		Birthday:               t,
		Avatar:                 userVO.Avatar,
		NickName:               userVO.NickName,
		UserTag:                userVO.UserTag,
		CommentStatus:          commentStatus,
		LastLoginTime:          t,
	}
	defaultPassword, b := service.SysParamService.GetSysParamsValueByKey("SYS_DEFAULT_PASSWORD")
	if !b {
		c.BaseController.ThrowError("00106", defaultPassword)
		return
	}
	user.PassWord = defaultPassword
	user.Source = "MOGU"
	common.DB.Create(&user)
	c.SuccessWithMessage("新增成功")
}
