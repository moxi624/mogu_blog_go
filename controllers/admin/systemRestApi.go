package admin

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/vo"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/30 8:03 上午
 * @version 1.0
 */
type SystemRestApi struct {
	base.BaseController
}

func (c *SystemRestApi) GetMe() {
	c.Result("success", c.GetMeService())
}

func (c *SystemRestApi) EditMe() {
	var adminVO vo.AdminVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &adminVO)
	if err != nil {
		panic(err)
	}
	//tokenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	if tokenData == "" {
		c.ErrorWithMessage("传入参数有误")
		return
	}
	tokenHead, _ := beego.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	t := common.Jwx.ParseToken(tokenString)
	a, _ := t.Get("adminUid")
	admin := models.Admin{
		Uid:            a.(string),
		UserName:       adminVO.UserName,
		PassWord:       adminVO.PassWord,
		Gender:         adminVO.Gender,
		Avatar:         adminVO.Avatar,
		Email:          adminVO.Email,
		Birthday:       adminVO.Birthday,
		Mobile:         adminVO.Mobile,
		NickName:       adminVO.NickName,
		QqNumber:       adminVO.QqNumber,
		WeChat:         adminVO.WeChat,
		Occupation:     adminVO.Occupation,
		Summary:        adminVO.Summary,
		Github:         adminVO.Github,
		Gitee:          adminVO.Gitee,
		RoleUid:        adminVO.RoleUid,
		PersonResume:   adminVO.PersonResume,
		StorageSize:    adminVO.StorageSize,
		MaxStorageSize: adminVO.MaxStorageSize,
	}
	common.DB.Updates(&admin)
	c.SuccessWithMessage("操作成功")
}

func (c *SystemRestApi) ChangePwd() {
	oldPwd := c.GetString("oldPwd")
	newPwd := c.GetString("newPwd")
	if oldPwd == "" || newPwd == "" {
		c.ErrorWithMessage("传入参数心得")
		return
	}
	//tokenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	tokenHead, _ := beego.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	t := common.Jwx.ParseToken(tokenString)
	a, _ := t.Get("adminUid")
	var admin models.Admin
	common.DB.Where("uid=?", a.(string)).Find(&admin)
	p := common.SHA256(oldPwd)
	isPassword := p == admin.PassWord
	if isPassword {
		admin.PassWord = common.SHA256(newPwd)
		common.DB.Save(&admin)
		c.SuccessWithMessage("更新成功")
	} else {
		c.ErrorWithMessage("密码错误")
	}
}
