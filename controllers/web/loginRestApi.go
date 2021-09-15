//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package web

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/rs/xid"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/19 3:57 下午
 * @version 1.0
 */

type LoginRestAPI struct {
	base.BaseController
}

var t, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")

func (c *LoginRestAPI) Login() {
	var userVO vo.UserVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userVO)
	if err != nil {
		panic(err)
	}
	isOpenLoginType, message := service.WebConfigService.IsOpenLoginType("PASSWORD")
	if !isOpenLoginType {
		if message != "" {
			c.ThrowError("00101", message)
			return
		}
		c.Result("error", "后台未开启该登录方式!")
		return
	}
	userName := userVO.UserName
	var user models.User
	common.DB.Where("user_name = ? or email = ?", userName, userName).Last(&user)
	if user == (models.User{}) || user.Status == 0 {
		c.Result("error", "用户不存在")
		return
	}
	if user.Status == 2 {
		c.ErrorWithMessage("用户账号未激活")
		return
	}
	if user.PassWord != "" && user.PassWord == common.MD5(userVO.PassWord) {
		ip := c.GetIP()
		userMap := c.GetOsAndBrowserInfo()
		user.Browser = userMap["BROWSER"]
		user.Os = userMap["OS"]
		user.LastLoginIp = ip
		user.LastLoginTime = time.Now()
		common.DB.Save(&user)
		if user.Avatar != "" {
			avatarResult := service.FileService.GetPicture(user.Avatar, ",")
			picList := common.WebUtil.GetPicture(avatarResult)
			if len(picList) > 0 {
				user.PhotoUrl = picList[0]
			}
		}
		token := xid.New().String()
		user.PassWord = ""
		b, _ := json.Marshal(user)
		t, _ := beego.AppConfig.Int64("user_token_survival_time")
		common.RedisUtil.SetEx("USER_TOKEN:"+token, string(b), t, time.Hour)
		c.SuccessWithData(token)
		base.L.Print("登录成功，返回token: ", token)
	} else {
		c.Result("error", "账号或密码错误")
	}
}

func (c *LoginRestAPI) Register() {
	var userVO vo.UserVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &userVO)
	if err != nil {
		panic(err)
	}
	isOpenLoginType, message := service.WebConfigService.IsOpenLoginType("PASSWORD")
	if !isOpenLoginType {
		if message != "" {
			c.ThrowError("00101", message)
			return
		}
		c.ErrorWithMessage("后台未开启该注册方式!")
		return
	}
	if len(userVO.UserName) < 5 || len(userVO.UserName) >= 20 || len(userVO.PassWord) < 5 || len(userVO.PassWord) >= 20 {
		c.Result("error", "用户名和密码长度必须在5-19")
		return
	}
	ip := c.GetIP()
	m := c.GetOsAndBrowserInfo()
	var user models.User
	common.DB.Where("status = ? and (user_name=? or email=?)", 1, userVO.UserName, userVO.Email).Last(&user)
	if user != (models.User{}) {
		c.Result("error", "用户或者邮件已存在")
		return
	}
	user.Uid = xid.New().String()
	user.UserName = userVO.UserName
	user.NickName = userVO.NickName
	user.PassWord = common.MD5(userVO.PassWord)
	user.Email = userVO.Email
	user.Source = "MOGU"
	user.LastLoginIp = ip
	user.Browser = m["BROWSER"]
	user.Os = m["OS"]
	user.Status = 2
	user.Birthday = t
	user.LastLoginTime = t
	common.DB.Create(&user)
	token := xid.New().String()
	user.PassWord = ""
	b, _ := json.Marshal(user)
	common.RedisUtil.SetEx("ACTIVATE_USER:"+token, string(b), 1, time.Hour)
	common.Email.SendActiveEmail(user, token)
	c.Result("success", "注册成功，请登录邮箱进行账号激活")
}

func (c *LoginRestAPI) ActiveUser() {
	token := c.GetString(":token")
	userInfo := common.RedisUtil.Get("ACTIVATE_USER:" + token)
	if userInfo == "" {
		c.Data["title"] = "出错了"
		c.Data["subtitle"] = "token令牌未被识别，请重新登录"
		c.TplName = bindUserEmail
		err := c.Render()
		if err != nil {
			panic(err)
		}
		return
	}
	var user models.User
	err := json.Unmarshal([]byte(userInfo), &user)
	if err != nil {
		panic(err)
	}
	if user.Status != 2 {
		c.Data["title"] = "出错了"
		c.Data["subtitle"] = "用户账号已经被激活,请勿重复激活"
		c.TplName = bindUserEmail
		err := c.Render()
		if err != nil {
			panic(err)
		}
		return
	}
	user.Status = 1
	common.DB.Model(&user).Select("status").Update("status", 1)
	var userList []models.User
	common.DB.Where("user_name=? and uid!=? and status!=?", user.UserName, user.Uid, 1).Find(&userList)
	if len(userList) > 0 {
		var uidList []string
		for _, item := range userList {
			uidList = append(uidList, item.Uid)
		}
		common.DB.Model(&userList).Where("uid in ?", uidList).Select("status").Update("status", 0)
	}
	c.Data["title"] = "恭喜你"
	c.Data["subtitle"] = "邮箱激活成功"
	c.TplName = bindUserEmail
	err1 := c.Render()
	if err1 != nil {
		panic(err1)
	}
}
