package admin

import (
	"encoding/json"
	"github.com/rs/xid"
	"math"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"reflect"
	"strconv"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/11 4:22 下午
 * @version 1.0
 */
type AdminRestApi struct {
	base.BaseController
}

func (c *AdminRestApi) GetOnlineAdminList() {
	var adminVO vo.AdminVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &adminVO)
	if err != nil {
		panic(err)
	}
	keys := common.RedisUtil.Keys("LOGIN_TOKEN_KEY*")
	onlineAdminJsonList := common.RedisUtil.MultiGet(keys)
	pageSize := adminVO.PageSize
	currentPage := adminVO.CurrentPage
	total := len(onlineAdminJsonList)
	startIndex := math.Max(float64((currentPage-1)*pageSize), 0)
	endIndex := math.Min(float64(currentPage+pageSize), float64(total))
	onlineAdminSubList := onlineAdminJsonList[int(startIndex):int(endIndex)]
	var onlineAdminList []models.OnlineAdmin
	for _, item := range onlineAdminSubList {
		var onlineAdmin models.OnlineAdmin
		err := json.Unmarshal([]byte(item.(string)), &onlineAdmin)
		if err != nil {
			panic(err)
		}
		onlineAdmin.Token = ""
		onlineAdminList = append(onlineAdminList, onlineAdmin)
	}
	page := map[string]interface{}{}
	page["records"] = onlineAdminList
	c.SuccessWithData(page)
}

func (c *AdminRestApi) ForceLogout() {
	var tokenUidList []string
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &tokenUidList)
	if err != nil {
		panic(err)
	}
	if len(tokenUidList) == 0 {
		c.ErrorWithMessage("传入参数有误")
		return
	}
	var tokenList []string
	for _, item := range tokenUidList {
		token := common.RedisUtil.Get("LOGIN_UUID_KEY:" + item)
		if token != "" {
			tokenList = append(tokenList, token)
		}
	}
	var keyList []string
	keyPrefix := "LOGIN_TOKEN_KEY:"
	for _, token := range tokenList {
		keyList = append(keyList, keyPrefix+token)
	}
	common.RedisUtil.MultiDelete(keyList)
	c.SuccessWithMessage("操作成功")
}

func (c *AdminRestApi) GetList() {
	var adminVO vo.AdminVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &adminVO)
	if err != nil {
		panic(err)
	}
	where := "status != ?"
	if strings.TrimSpace(adminVO.Keyword) != "" {
		where += " and (user_name like \"%" + strings.TrimSpace(adminVO.Keyword) + "%\" or nick_name like \"%" + strings.TrimSpace(adminVO.Keyword) + "%\")"
	}
	var pageList []models.Admin
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Admin{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 0).Offset((adminVO.CurrentPage - 1) * adminVO.PageSize).Limit(adminVO.PageSize).Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var s []string
	var adminUidList []string
	for _, item := range pageList {
		if item.Avatar != "" {
			s = append(s, item.Avatar)
		}
		adminUidList = append(adminUidList, item.Uid)
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
	storageList := service.StorageService.GetStorageByAdminUid(adminUidList)
	storageMap := map[string]models.Storage{}
	for _, item := range storageList {
		storageMap[item.AdminUid] = item
	}
	for i, item := range pageList {
		var role models.Role
		common.DB.Where("uid=?", item.RoleUid).Find(&role)
		pageList[i].Role = role
		if item.Avatar != "" {
			pictureUidsTemp := strings.Split(item.Avatar, ",")
			var pictureListTemp []string
			for _, picture := range pictureUidsTemp {
				if pictureMap[picture] != "" {
					pictureListTemp = append(pictureListTemp, pictureMap[picture])
				}
			}
			pageList[i].PhotoList = pictureListTemp
		}
		storage := storageMap[item.Uid]
		if storage != (models.Storage{}) {
			pageList[i].StorageSize = storage.StorageSize
			pageList[i].MaxStorageSize = storage.MaxStorageSize
		} else {
			pageList[i].StorageSize = 0
			pageList[i].MaxStorageSize = 0
		}
		pageList[i].PassWord = ""
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    adminVO.PageSize,
		Current: adminVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *AdminRestApi) RestPwd() {
	var adminVO vo.AdminVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &adminVO)
	if err != nil {
		panic(err)
	}
	defaultPassword, b := service.SysParamService.GetSysParamsValueByKey("SYS_DEFAULT_PASSWORD")
	if !b {
		c.BaseController.ThrowError("00106", defaultPassword)
		return
	}
	adminUid := c.GetAdminUid()
	var admin models.Admin
	common.DB.Where("uid = ?", adminVO.Uid).Find(&admin)
	if admin.UserName == "admin" && admin.Uid != adminUid {
		c.ErrorWithMessage("更新管理员密码失败")
	} else {
		admin.PassWord = common.SHA256(defaultPassword)
		if adminVO.Birthday.Format("2006-01-02") == "0001-01-01" {
			admin.Birthday = t
		} else {
			admin.Birthday = adminVO.Birthday
		}
		if admin.LastLoginTime.Format("2006-01-02") == "0001-01-01" {
			admin.LastLoginTime = t
		}
		common.DB.Save(&admin)
		c.SuccessWithMessage("操作成功")
	}
}

func (c *AdminRestApi) Edit() {
	var adminVO vo.AdminVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &adminVO)
	if err != nil {
		panic(err)
	}
	var admin models.Admin
	common.DB.Where("uid=?", adminVO.Uid).Find(&admin)
	var adminList []models.Admin
	if !reflect.DeepEqual(admin, models.Admin{}) {
		if admin.UserName == "admin" && adminVO.UserName != "admin" {
			c.ErrorWithMessage("超级管理员用户名必须为admin")
			return
		}
		common.DB.Where("status=? and user_name=?", 1, adminVO.UserName).Find(&adminList)
		if len(adminList) > 0 {
			for _, item := range adminList {
				if item.Uid == adminVO.Uid {
					continue
				} else {
					c.ErrorWithMessage("修改失败，用户名存在")
					return
				}
			}

		}
	}
	if adminVO.RoleUid != "" && admin.RoleUid != adminVO.RoleUid {
		common.RedisUtil.Delete("ADMIN_VISIT_MENU:" + admin.Uid)
	}
	admin.UserName = adminVO.UserName
	admin.Avatar = adminVO.Avatar
	admin.NickName = adminVO.NickName
	admin.Gender = adminVO.Gender
	admin.Email = adminVO.Email
	admin.QqNumber = adminVO.QqNumber
	admin.Github = adminVO.Github
	admin.Gitee = adminVO.Gitee
	admin.Occupation = adminVO.Occupation
	admin.Mobile = adminVO.Mobile
	admin.RoleUid = adminVO.RoleUid
	if adminVO.Birthday.Format("2006-01-02") == "0001-01-01" {
		admin.Birthday = t
	} else {
		admin.Birthday = adminVO.Birthday
	}
	common.DB.Save(&admin)
	result, b := service.StorageService.EditStorageSizes(admin.Uid, adminVO.MaxStorageSize*1024*1024)
	if b {
		c.SuccessWithMessage(result)
	} else {
		c.ErrorWithMessage(result)
	}
}

func (c *AdminRestApi) Delete() {
	adminUids := c.GetStrings("adminUids")
	var adminList []models.Admin
	common.DB.Find(&adminList, adminUids)
	common.DB.Model(&adminList).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}

func (c *AdminRestApi) Add() {
	var adminVO vo.AdminVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &adminVO)
	if err != nil {
		panic(err)
	}
	mobile := adminVO.Mobile
	userName := adminVO.UserName
	email := adminVO.Email
	if userName == "" {
		c.ErrorWithMessage("传入参数有误")
	} else {
		if email == "" && mobile == "" {
			c.ErrorWithMessage("邮箱和手机号至少一项不能为空")
			return
		}
		defaultPassword, b := service.SysParamService.GetSysParamsValueByKey("SYS_DEFAULT_PASSWORD")
		if !b {
			c.BaseController.ThrowError("00106", defaultPassword)
			return
		}
		var temp models.Admin
		common.DB.Where("user_name=?", userName).First(&temp)
		if reflect.DeepEqual(temp, models.Admin{}) {
			admin := models.Admin{
				Avatar:        adminVO.Avatar,
				Email:         adminVO.Email,
				Gender:        adminVO.Gender,
				UserName:      adminVO.UserName,
				NickName:      adminVO.NickName,
				RoleUid:       adminVO.RoleUid,
				PassWord:      common.SHA256(defaultPassword),
				Uid:           xid.New().String(),
				Birthday:      t,
				LastLoginTime: t,
			}
			common.DB.Create(&admin)
			maxStorageSize, b := service.SysParamService.GetSysParamsValueByKey("MAX_STORAGE_SIZE")
			if !b {
				c.BaseController.ThrowError("00106", defaultPassword)
				return
			}
			i, _ := strconv.Atoi(maxStorageSize)
			service.InitStorageSize(admin.Uid, int64(i)*1024*1024)
			c.SuccessWithMessage("新增成功")
		} else {
			c.ErrorWithMessage("该管理员已经存在")
		}
	}
}
