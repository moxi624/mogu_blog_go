package base

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/rs/xid"
	"mime/multipart"
	"mogu-go-v2/common"
	"mogu-go-v2/models"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/12 5:55 下午
 * @version 1.0
 */

type BaseController struct {
	web.Controller
	Wg sync.WaitGroup
}

var L = logs.GetLogger()

func (c *BaseController) Result(code string, data interface{}) {
	c.Data["json"] = map[string]interface{}{
		"code": code,
		"data": data,
	}
	err := c.ServeJSON()
	if err != nil {
		panic(err)
	}
}

func (c *BaseController) GetOsAndBrowserInfo() map[string]string {
	userAgent := c.Ctx.Request.Header.Get("User-Agent")
	user := strings.ToLower(userAgent)
	var os1 string
	var browser = ""
	switch {
	case strings.Contains(user, "windows"):
		os1 = "Windows"
	case strings.Contains(user, "mac"):
		os1 = "Mac"
	case strings.Contains(user, "x11"):
		os1 = "Unix"
	case strings.Contains(user, "android"):
		os1 = "Android"
	case strings.Contains(user, "iphone"):
		os1 = "Iphone"
	default:
		os1 = "UnKnown, More-Info: " + userAgent
	}
	switch {
	case strings.Contains(user, "edge"):
		sub := string([]byte(userAgent)[strings.Index(userAgent, "Edge"):])
		split := strings.Split(sub, " ")
		browser = strings.ReplaceAll(split[0], "/", "-")
	case strings.Contains(user, "msie"):
		sub := string([]byte(userAgent)[strings.Index(userAgent, "MSIE"):])
		split := strings.Split(sub, ";")
		split1 := strings.Split(split[0], " ")
		browser = strings.ReplaceAll(split1[0], "MSIE", "IE") + "-" + split1[1]
	case strings.Contains(user, "safari") && strings.Contains(user, "version"):
		sub := string([]byte(userAgent)[strings.Index(userAgent, "Safari"):])
		split := strings.Split(sub, " ")
		split1 := strings.Split(split[0], "/")
		sub1 := string([]byte(userAgent)[strings.Index(userAgent, "Version"):])
		split2 := strings.Split(sub1, " ")
		split3 := strings.Split(split2[0], "/")
		browser = split1[0] + "-" + split3[1]
	case strings.Contains(user, "opr") || strings.Contains(user, "opera"):
		if strings.Contains(user, "opera") {
			sub := string([]byte(userAgent)[strings.Index(userAgent, "Opera"):])
			split := strings.Split(sub, " ")
			split1 := strings.Split(split[0], "/")
			sub1 := string([]byte(userAgent)[strings.Index(userAgent, "Version"):])
			split2 := strings.Split(sub1, " ")
			split3 := strings.Split(split2[0], "/")
			browser = split1[0] + "-" + split3[1]
		} else if strings.Contains(user, "opr") {
			sub := string([]byte(userAgent)[strings.Index(userAgent, "OPR"):])
			split := strings.Split(sub, " ")
			rep := strings.ReplaceAll(split[0], "/", "-")
			browser = strings.ReplaceAll(rep, "OPR", "Opera")
		}
	case strings.Contains(user, "chrome"):
		sub := string([]byte(userAgent)[strings.Index(userAgent, "Chrome"):])
		split := strings.Split(sub, " ")
		browser = strings.ReplaceAll(split[0], "/", "-")
	case strings.Contains(user, "mozilla") || strings.Contains(user, "netscape"):
		browser = "Netscape-?"
	case strings.Contains(user, "firefox"):
		sub := string([]byte(userAgent)[strings.Index(userAgent, "Firefox"):])
		split := strings.Split(sub, " ")
		browser = strings.ReplaceAll(split[0], "/", "-")
	case strings.Contains(user, "rv"):
		sub := string([]byte(userAgent)[strings.Index(userAgent, "rv"):])
		split := strings.Split(sub, " ")
		IEVersion := strings.ReplaceAll(split[0], "rv", "-")
		browser = "IE" + string([]byte(IEVersion)[:len(IEVersion)])
	default:
		browser = "UnKnown"
	}
	result := map[string]string{}
	result["OS"] = os1
	result["BROWSER"] = browser
	return result

}
func (c *BaseController) AddOnlineAdmin(admin models.Admin, expirationSecond int64) {
	ip := c.GetIP()
	m := c.GetOsAndBrowserInfo()
	os1 := m["OS"]
	browser := m["BROWSER"]
	onlineAdmin := models.OnlineAdmin{
		AdminUid:   admin.Uid,
		TokenId:    admin.TokenUid,
		Token:      admin.ValidCode,
		Os:         os1,
		Browser:    browser,
		Ipaddr:     ip,
		LoginTime:  common.DateUtils.GetNowTime(),
		RoleName:   admin.Role.RoleName,
		UserName:   admin.UserName,
		ExpireTime: common.DateUtils.GetDateStr(time.Now(), expirationSecond),
	}
	jsonResult := common.RedisUtil.Get("IP_SOURCE:" + ip)
	if jsonResult == "" {
		address := common.IpUtils.GetAddresses(ip)
		if address != "" {
			onlineAdmin.LoginLocation = address
			common.RedisUtil.SetEx("IP_SOURCE:"+ip, address, 24, time.Hour)
		}
	} else {
		onlineAdmin.LoginLocation = jsonResult
	}
	b, _ := json.Marshal(onlineAdmin)
	common.RedisUtil.SetEx("LOGIN_TOKEN_KEY:"+admin.ValidCode, string(b), expirationSecond, time.Second)
	common.RedisUtil.SetEx("LOGIN_UUID_KEY:"+admin.TokenUid, admin.ValidCode, expirationSecond, time.Second)
}

func (c *BaseController) SuccessWithData(data interface{}) {
	c.Data["json"] = map[string]interface{}{
		"code": "success",
		"data": data,
	}
	err := c.ServeJSON()
	if err != nil {
		panic(err)
	}
}

func (c *BaseController) ErrorWithMessage(message string) {
	c.Data["json"] = map[string]interface{}{
		"code":    "error",
		"message": message,
	}
	err := c.ServeJSON()
	if err != nil {
		panic(err)
	}
}

func (c *BaseController) SuccessWithMessage(message string) {
	c.Data["json"] = map[string]interface{}{
		"code":    "success",
		"message": message,
	}
	err := c.ServeJSON()
	if err != nil {
		panic(err)
	}
}

func (c *BaseController) GetIP() string {
	header := c.Ctx.Request.Header
	ipAddress := header.Get("x-forwarded-for")
	/*	ip := exnet.ClientPublicIP(c.Ctx.Request)*/
	/*	if ip == ""{*/
	/*		ip = exnet.ClientIP(c.Ctx.Request)*/
	/*	}*/
	/*	L.Print("这个是IP:"+ip)*/
	if ipAddress == "" || strings.EqualFold(ipAddress, "unknown") {
		ipAddress = header.Get("Proxy-Client-IP")
	}
	if ipAddress == "" || strings.EqualFold(ipAddress, "unknown") {
		ipAddress = header.Get("WL-Proxy-Client-IP")
	}
	if ipAddress == "" || strings.EqualFold(ipAddress, "unknown") {
		ipAddress = c.Ctx.Request.RemoteAddr
		if strings.Contains(ipAddress, "[::1]") {
			conn, err := net.Dial("udp", "baidu.com:80")
			if err != nil {
				fmt.Println(err.Error())
			}
			defer conn.Close()
			ipAddress = conn.LocalAddr().String()
			s := strings.Split(ipAddress, ":")
			ipAddress = s[0]
		}
	}
	if ipAddress != "" && len(ipAddress) > 15 {
		if strings.Contains(ipAddress, ",") {
			s := strings.Split(ipAddress, ",")
			ipAddress = s[0]
		}
	}
	if ipAddress != "" && strings.Contains(ipAddress, ":") {
		s := strings.Split(ipAddress, ":")
		ipAddress = s[0]
	}
	return ipAddress
}

func (c *BaseController) GetSystemConfig() models.SystemConfig {
	//token := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	paramsToken := c.GetString("token")
	platform := c.GetString("platform")
	systemConfigMap := map[string]interface{}{}
	if platform == "web" || (paramsToken != "" && len(paramsToken) == 32) {
		systemConfigMap = c.getSystemConfigByWebToken(paramsToken)
	} else {
		if token != "" {
			systemConfigMap = c.GetSystemConfigMap(token)
		} else {
			systemConfigMap = c.GetSystemConfigMap(paramsToken)
		}
	}
	if len(systemConfigMap) == 0 {
		c.ThrowError("00107", "请先配置七牛云，或者重新登录")
		return models.SystemConfig{}
	}
	uploadQiNiu := systemConfigMap["uploadQiNiu"].(string)
	uploadLocal := systemConfigMap["uploadLocal"].(string)
	localPictureBaseUrl := systemConfigMap["localPictureBaseUrl"].(string)
	qiNiuPictureBaseUrl := systemConfigMap["qiNiuPictureBaseUrl"].(string)
	qiNiuAccessKey := systemConfigMap["qiNiuAccessKey"].(string)
	qiNiuSecretKey := systemConfigMap["qiNiuSecretKey"].(string)
	qiNiuBucket := systemConfigMap["qiNiuBucket"].(string)
	qiNiuArea := systemConfigMap["qiNiuArea"].(string)

	minioEndPoint := systemConfigMap["minioEndPoint"].(string)
	minioAccessKey := systemConfigMap["minioAccessKey"].(string)
	minioSecretKey := systemConfigMap["minioSecretKey"].(string)
	minioBucket := systemConfigMap["minioBucket"].(string)
	uploadMinio := systemConfigMap["uploadMinio"].(string)
	minioPictureBaseUrl := systemConfigMap["minioPictureBaseUrl"].(string)
	if uploadQiNiu == "1" && (qiNiuPictureBaseUrl == "" || qiNiuAccessKey == "" || qiNiuSecretKey == "" || qiNiuBucket == "" || qiNiuArea == "") {
		L.Print("请先配置七牛云")
		return models.SystemConfig{}
	}
	if uploadLocal == "1" && localPictureBaseUrl == "" {
		L.Print("请先配置本地存储")
		return models.SystemConfig{}
	}
	if uploadMinio == "1" && (minioEndPoint == "" || minioPictureBaseUrl == "" || minioAccessKey == "" || minioSecretKey == "" || minioBucket == "") {
		L.Print("请先配置Minio上传服务")
		return models.SystemConfig{}
	}
	systemConfig := models.SystemConfig{
		QiNiuAccessKey:      qiNiuAccessKey,
		QiNiuSecretKey:      qiNiuSecretKey,
		QiNiuBucket:         qiNiuBucket,
		QiNiuArea:           qiNiuArea,
		UploadQiNiu:         uploadQiNiu,
		UploadLocal:         uploadLocal,
		PicturePriority:     systemConfigMap["picturePriority"].(string),
		LocalPictureBaseUrl: systemConfigMap["localPictureBaseUrl"].(string),
		QiNiuPictureBaseUrl: systemConfigMap["qiNiuPictureBaseUrl"].(string),
		MinioEndPoint:       minioEndPoint,
		MinioAccessKey:      minioAccessKey,
		MinioSecretKey:      minioSecretKey,
		MinioBucket:         minioBucket,
		MinioPictureBaseUrl: minioPictureBaseUrl,
		UploadMinio:         uploadMinio,
	}
	return systemConfig
}

func (c *BaseController) GetAdminUid() string {
	//okenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	if tokenData == "" {
		tokenData = c.GetString("token")
	}
	s, _ := web.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, s)
	token := common.Jwx.ParseToken(tokenString)
	a, _ := token.Get("adminUid")
	return a.(string)
}

func (c *BaseController) CheckLogin() string {
	if c.GetAdminUid() == "" {
		L.Print("用户未登录")
		c.ThrowError("00008", "Token令牌未被识别,请重新登录")
		return ""
	} else {
		return c.GetAdminUid()
	}
}

func (c *BaseController) ThrowError(code string, message string) {
	c.Data["json"] = map[string]interface{}{
		"data":    code,
		"message": message,
	}
	err := c.ServeJSON()
	if err != nil {
		panic(err)
	}
}

func (c *BaseController) getSystemConfigByWebToken(token string) map[string]interface{} {
	webUserJsonResult := common.RedisUtil.Get("USER_TOKEN:" + token)
	if webUserJsonResult == "" {
		c.ThrowError("00008", "token令牌未被识别，请重新登录")
		return map[string]interface{}{}
	}
	resultMap := map[string]interface{}{}
	jsonResult := common.RedisUtil.Get("SYSTEM_CONFIG")
	if jsonResult != "" {
		err := json.Unmarshal([]byte(jsonResult), &resultMap)
		if err != nil {
			panic(err)
		}
	} else {
		resultTempMap := service.AuthService.GetSystemConfig(token)
		if resultTempMap["code"] == "success" {
			resultData := resultTempMap["data"]
			b, _ := json.Marshal(resultData)
			common.RedisUtil.SetEx("SYSTEM_CONFIG", string(b), 30, time.Minute)
			err := json.Unmarshal(b, &resultMap)
			if err != nil {
				panic(err)
			}
		}
	}
	return resultMap
}

func (c *BaseController) GetSystemConfigMap(token string) map[string]interface{} {
	adminJsonResult := common.RedisUtil.Get("LOGIN_TOKEN_KEY:" + token)
	if adminJsonResult == "" {
		c.ThrowError("00008", "token令牌未被识别，请重新登录")
		return map[string]interface{}{}
	}
	resultMap := map[string]interface{}{}
	jsonResult := common.RedisUtil.Get("SYSTEM_CONFIG")
	if jsonResult != "" {
		err := json.Unmarshal([]byte(jsonResult), &resultMap)
		if err != nil {
			panic(err)
		}
	} else {
		resultTempMap := service.SystemConfigService.GetConfig()
		if resultTempMap["code"] == "success" {
			resultData := resultTempMap["data"]
			b, _ := json.Marshal(resultData)
			common.RedisUtil.SetEx("SYSTEM_CONFIG", string(b), 30, time.Minute)
			err := json.Unmarshal(b, &resultMap)
			if err != nil {
				panic(err)
			}
		}
	}
	return resultMap
}

func (c *BaseController) GetStorageByAdmin() models.Storage {
	//tokenData := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	if tokenData == "" {
		tokenData = c.GetString("token")
	}
	tokenHead, _ := web.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	token := common.Jwx.ParseToken(tokenString)
	a, _ := token.Get("adminUid")
	var reStorage models.Storage
	common.DB.Where("status = ? and admin_uid=?", 1, a.(string)).Limit(1).First(&reStorage)
	return reStorage
}

func (c *BaseController) BatchUploadFile(filedatas []*multipart.FileHeader, systemConfig models.SystemConfig, filename string) map[string]interface{} {
	uploadQiNiu := systemConfig.UploadQiNiu
	uploadLocal := systemConfig.UploadLocal
	uploadMinio := systemConfig.UploadMinio
	source := c.GetString("source")
	userUid := ""
	adminUid := ""
	projectName := ""
	sortName := ""
	if source == "picture" {
		userUid = c.GetString("userUid")
		adminUid = c.GetString("adminUid")
		projectName = c.GetString("projectName")
		sortName = c.GetString("sortName")
	} else if source == "admin" {
		/*userUid = c.GetString("userUid")
		adminUid = c.GetString("adminUid")
		projectName = c.GetString("projectName")
		sortName = c.GetString("sortName")*/
	} else {

	}
	if projectName == "" {
		projectName = "base"
	}
	if userUid == "" && adminUid == "" {
		m := map[string]interface{}{
			"code": "error",
			"data": "请先注册",
		}
		return m
	}
	var fileSorts []models.FileSort
	common.DB.Where("sort_name=? and project_name=? and status=?", sortName, projectName, 1).Find(&fileSorts)
	var fileSort models.FileSort
	if len(fileSorts) < 1 {
		m := map[string]interface{}{
			"code": "error",
			"data": "文件不被允许上传",
		}
		return m
	}
	fileSort = fileSorts[0]
	sortUrl := fileSort.Url
	if sortUrl == "" {
		sortUrl = "base/common/"
	} else {
		sortUrl = fileSort.Url
	}
	var lists []models.File
	if len(filedatas) > 0 {
		for _, filedata := range filedatas {
			oldName := filedata.Filename
			size := filedata.Size
			picExpandedName := path.Ext(oldName)
			var newFilename string
			if filename == "file" {
				newFilename = strconv.FormatInt(time.Now().Unix(), 10) + ".png"
			}
			newFilename = strconv.FormatInt(time.Now().Unix(), 10) + picExpandedName
			localUrl := ""
			qiNiuUrl := ""
			minioUrl := ""
			tempFileData := filedata
			if uploadQiNiu == "1" {
				qiNiuUrl = c.uploadSingleFile(tempFileData, filename)
			}
			if uploadMinio == "1" {

			}
			if uploadLocal == "1" {

			}
			var file models.File
			if filename == "file" {
				file = models.File{
					Uid:             xid.New().String(),
					FileSortUid:     fileSort.Uid,
					FileOldName:     "",
					FileSize:        size,
					PicExpandedName: "png",
					PicName:         newFilename,
					PicURL:          localUrl,
					Status:          1,
					UserUid:         userUid,
					AdminUid:        adminUid,
					QiNiuUrl:        qiNiuUrl,
					MinioUrl:        minioUrl,
				}
			}
			file = models.File{
				Uid:             xid.New().String(),
				FileSortUid:     fileSort.Uid,
				FileOldName:     oldName,
				FileSize:        size,
				PicExpandedName: picExpandedName,
				PicName:         newFilename,
				PicURL:          localUrl,
				Status:          1,
				UserUid:         userUid,
				AdminUid:        adminUid,
				QiNiuUrl:        qiNiuUrl,
				MinioUrl:        minioUrl,
			}
			common.DB.Create(&file)
			lists = append(lists, file)
		}
		m := map[string]interface{}{
			"code": "success",
			"data": lists,
		}
		return m
	}
	m := map[string]interface{}{
		"code": "error",
		"data": "请上传图片",
	}
	return m
}

func (c *BaseController) uploadSingleFile(mutipartFile *multipart.FileHeader, filename string) string {
	url := ""
	systemConfig := c.GetSystemConfig()
	oldName := mutipartFile.Filename
	picExpandedName := path.Ext(oldName)
	var newFilename string
	if filename == "file" {
		newFilename = strconv.FormatInt(time.Now().Unix(), 10) + ".png"
	}
	newFilename = strconv.FormatInt(time.Now().Unix(), 10) + picExpandedName
	s, _ := web.AppConfig.String("fileUploadPath")
	dir := s + "/temp/"
	_, err := os.Stat(dir)
	if !os.IsExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	tempFiles := dir + newFilename
	err1 := c.SaveToFile(filename, tempFiles)
	if err1 != nil {
		panic(err1)
	}
	url = common.QiniuUtil.UploadQiNiu(tempFiles, systemConfig)
	return url
}

func (c *BaseController) DeleteFileService(networkDiskVO vo.NetworkDiskVO, qiNiuConfig map[string]interface{}) {
	uid := networkDiskVO.Uid
	if uid == "" {
		L.Print("删除的文件不能为空")
	}
	var networkDisk models.NetworkDisk
	common.DB.Where("uid=?", uid).Find(&networkDisk)
	uploadQiNiu := qiNiuConfig["uploadQiNiu"]
	uploadLocal := qiNiuConfig["uploadLocal"]
	uploadMinio := qiNiuConfig["uploadMinio"]
	networkDisk.Status = 0
	common.DB.Save(&networkDisk)
	if networkDisk.IsDir == 1 {
		path := networkDisk.FilePath + networkDisk.FileName
		var list []models.NetworkDisk
		common.DB.Where("status=1 and file_path like '" + path + "%'").Find(&list)
		if len(list) > 0 {
			for i := range list {
				list[i].Status = 0
			}
			if err := common.DB.Save(&list).Error; err == nil {
				if uploadLocal == "1" {

				}
				if uploadQiNiu == "1" {
					var fileList []string
					for _, item := range list {
						fileList = append(fileList, item.QiNiuUrl)
					}
					common.QiniuUtil.DeleteFileList(fileList, qiNiuConfig)
				}
				if uploadMinio == "1" {

				}
			}
		}
	} else {
		if uploadLocal == "1" {

		}
		if uploadQiNiu == "1" {
			qiNiuUrl := networkDisk.QiNiuUrl
			common.QiniuUtil.DeleteFile(qiNiuUrl, qiNiuConfig)
		}
		if uploadMinio == "1" {

		}
		storage := c.GetStorageByAdmin()
		storageSize := storage.StorageSize - networkDisk.FileSize
		if storageSize > 0 {
			storage.StorageSize = storageSize
		}
		storage.StorageSize = 0
		common.DB.Save(&storage)
	}
}

func (c *BaseController) GetMeService() models.Admin {
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	tokenHead, _ := web.AppConfig.String("tokenHead")
	tokenString := strings.TrimPrefix(tokenData, tokenHead)
	t := common.Jwx.ParseToken(tokenString)
	adminUid, _ := t.Get("adminUid")
	if adminUid == "" {
		return models.Admin{}
	}
	var admin models.Admin
	common.DB.Where("uid=?", adminUid.(string)).Find(&admin)
	admin.PassWord = ""
	if admin.Avatar != "" {
		pictureList := service.FileService.GetPicture(admin.Avatar, ",")
		admin.PhotoList = common.WebUtil.GetPicture(pictureList)
	}
	return admin
}
