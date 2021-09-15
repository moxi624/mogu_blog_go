package admin

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/service"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

/**
 * 登录管理
 * @author  镜湖老杨
 * @date  2020/12/10 8:46 上午
 * @version 1.0
 */
type LoginRestAPI struct {
	base.BaseController
}

func (c *LoginRestAPI) Login() {
	username := strings.Trim(c.GetString("username"), " ")
	password := common.SHA256(strings.TrimSpace(c.GetString("password")))
	base.L.Print("密码：" + password)
	isRememberMe, _ := c.GetBool("isRemember")
	if username == "" || password == "" {
		c.Result("error", "账号或密码不能为空")
		return
	}
	ip := c.GetIP()
	limitCount := common.RedisUtil.Get("LOGIN_LIMIT:" + ip)
	var tmplimitCount int
	var maxLoginError int
	if limitCount != "" {
		tmplimitCount, _ = strconv.Atoi(limitCount)
		maxLoginError, _ = beego.AppConfig.Int("maxLoginError")
		if tmplimitCount >= maxLoginError {
			c.Result("error", "密码输错次数过多,已被锁定30分钟")
			return
		}
	}
	isEmail := common.CheckEmail(username)
	isMobile := common.CheckMobile(username)
	var admin models.Admin
	switch {
	case isEmail:
		common.DB.Where("email = ? and status=?", username, 1).First(&admin)
	case isMobile:
		common.DB.Where("mobile = ? and status=?", username, 1).First(&admin)
	default:
		common.DB.Where("user_name = ? and status=?", username, 1).First(&admin)
	}
	if reflect.DeepEqual(admin, models.Admin{}) {
		t := c.setLoginCommit()
		c.Result("error", "用户名或密码错误或者账号不存在，错误"+strconv.Itoa(t)+"次后，账户将被锁定30分钟")
		return
	}
	passwordDB := admin.PassWord
	if password != passwordDB {
		t := c.setLoginCommit()
		c.Result("error", "用户名或密码错误，错误"+strconv.Itoa(t)+"次后，账户将被锁定30分钟")
		return
	}
	var roleUids []string
	roleUids = append(roleUids, admin.RoleUid)
	var roles []models.Role
	common.DB.Where("uid = ?", roleUids).Find(&roles)
	if len(roles) == 0 {
		c.Result("error", "没有分配角色权限")
		return
	}
	roleName := roles[0].RoleName
	var expiration int64
	if isRememberMe {
		expiration, _ = beego.AppConfig.Int64("isRememberMeExpiresSecond")
	} else {
		expiration, _ = beego.AppConfig.Int64("audience_expiresSecond")
	}
	audience, _ := beego.AppConfig.String("audience_clientId")
	isuser, _ := beego.AppConfig.String("audience_name")
	base64Secret, _ := beego.AppConfig.String("audience_base64Secret")
	jwtToken := common.Jwx.CreateJWT(
		admin.UserName,
		admin.Uid,
		roleName,
		audience,
		isuser,
		expiration*1000,
		base64Secret,
	)
	s, _ := beego.AppConfig.String("tokenHead")
	token := s + jwtToken
	base.L.Print(token)
	result := map[string]string{}
	result["token"] = token
	count := admin.LoginCount + 1
	admin.LoginCount = count
	admin.LastLoginIp = ip
	admin.LastLoginTime = time.Now()
	if admin.Birthday.Format("2006-01-02") == "0001-01-01" {
		admin.Birthday = t
	}
	common.DB.Save(&admin)
	admin.ValidCode = token
	admin.TokenUid = common.ULID()
	admin.Role = roles[0]
	c.AddOnlineAdmin(admin, expiration)
	c.Result("success", result)
}

func (c *LoginRestAPI) Info() {
	tokenData := c.GetString("token")
	m := map[string]interface{}{}
	if tokenData == "" {
		c.Result("error", "token用户过期")
		return
	}
	s, _ := beego.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, s)
	token := common.Jwx.ParseToken(tokenString)
	var admin models.Admin
	a, _ := token.Get("adminUid")
	common.DB.Where("uid=?", a.(string)).Find(&admin)
	m["token"] = token
	if admin.Avatar != "" {
		pictureList := service.FileService.GetPicture(admin.Avatar, ",")
		list := common.WebUtil.GetPicture(pictureList)
		if len(list) > 0 {
			m["avatar"] = list[0]
		} else {
			m["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
		}
	}
	var roleUid []string
	roleUid = append(roleUid, admin.RoleUid)
	var roleList []models.Role
	common.DB.Find(&roleList, roleUid)
	m["roles"] = roleList
	c.Result("success", m)
}

func (c *LoginRestAPI) GetMenu() {
	m := map[string]interface{}{}
	var categoryMenuList []models.CategoryMenu
	//tokenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	if tokenData == "" {
		c.Result("error", "token用户过期")
		return
	}
	tokenHead, _ := beego.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	t := common.Jwx.ParseToken(tokenString)
	a, _ := t.Get("adminUid")
	var admin models.Admin
	common.DB.Where("uid=?", a.(string)).Find(&admin)
	roleUid := []string{admin.RoleUid}
	var roleList []models.Role
	common.DB.Find(&roleList, roleUid)
	var categoryMenuUids []string
	for _, item := range roleList {
		categoryMenuUid := item.CategoryMenuUids
		s1 := strings.Replace(categoryMenuUid, "[", "", -1)
		s2 := strings.Replace(s1, "]", "", -1)
		s3 := strings.Replace(s2, "\"", "", -1)
		categoryMenuUids = strings.Split(s3, ",")
	}
	common.DB.Find(&categoryMenuList, categoryMenuUids)
	var buttonList []models.CategoryMenu
	var secondMenuUidList []string
	for _, item := range categoryMenuList {
		if item.MenuType == 0 && item.MenuLevel == 2 {
			secondMenuUidList = append(secondMenuUidList, item.Uid)
			secondMenuUidList = common.RemoveRepByMap(secondMenuUidList)
		}
		if item.MenuType == 1 && item.ParentUid != "" {
			secondMenuUidList = append(secondMenuUidList, item.ParentUid)
			secondMenuUidList = common.RemoveRepByMap(secondMenuUidList)
			buttonList = append(buttonList, item)
		}
	}
	var childCategoryMenuList []models.CategoryMenu
	var parentCategoryMenuList []models.CategoryMenu
	var parentCategoryMenuUids []string
	if len(secondMenuUidList) > 0 {
		common.DB.Find(&childCategoryMenuList, secondMenuUidList)
	}
	for _, item := range childCategoryMenuList {
		if item.MenuLevel == 2 {
			if item.ParentUid != "" {
				parentCategoryMenuUids = append(parentCategoryMenuUids, item.ParentUid)
			}
		}
	}
	if len(parentCategoryMenuUids) > 0 {
		common.DB.Find(&parentCategoryMenuList, parentCategoryMenuUids)
	}
	sort.SliceStable(parentCategoryMenuList, func(i, j int) bool {
		return parentCategoryMenuList[i].Sort > parentCategoryMenuList[j].Sort
	})
	m["parentList"] = parentCategoryMenuList
	m["sonList"] = childCategoryMenuList
	m["buttonList"] = buttonList
	c.Result("success", m)
}

func (c *LoginRestAPI) Logout() {
	tokenData := c.Ctx.GetCookie("Admin-Token")
	if tokenData == "" {
		c.Result("error", "操作失败")
		return
	}
	adminJson := common.RedisUtil.Get("LOGIN_TOKEN_KEY:" + tokenData)
	if adminJson != "" {
		var onlineAdmin models.OnlineAdmin
		err := json.Unmarshal([]byte(adminJson), &onlineAdmin)
		if err != nil {
			panic(err)
		}
		tokenUid := onlineAdmin.TokenId
		common.RedisUtil.Delete("LOGIN_UUID_KEY:" + tokenUid)
	}
	common.RedisUtil.Delete("LOGIN_TOKEN_KEY:" + tokenData)
	c.Result("success", "操作成功")
}

func (c *LoginRestAPI) setLoginCommit() int {
	ip := c.GetIP()
	count := common.RedisUtil.Get("LOGIN_LIMIT:" + ip)
	surplusCount := 5
	if count != "" {
		countTemp, _ := strconv.Atoi(count)
		countTemp++
		surplusCount -= countTemp
		common.RedisUtil.SetEx("LOGIN_LIMIT:"+ip, strconv.Itoa(countTemp), 10, time.Minute)
	} else {
		surplusCount--
		common.RedisUtil.SetEx("LOGIN_LIMIT:"+ip, "1", 30, time.Minute)
	}
	return surplusCount
}

func (c *LoginRestAPI) GetWebSiteName() {
	var webConfig models.WebConfig
	common.DB.Last(&webConfig)
	if webConfig.Name != "" {
		c.SuccessWithData(webConfig.Name)
	} else {
		c.SuccessWithData("")
	}
}
