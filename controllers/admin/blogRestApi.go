package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"mime/multipart"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/maps"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"reflect"
	"strconv"
	"strings"
)

const sortDesc = "sort desc"

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/11 9:57 上午
 * @version 1.0
 */
type BlogRestApi struct {
	base.BaseController
}

func (c *BlogRestApi) GetList() {
	var blogVO vo.BlogVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogVO)
	if err != nil {
		panic(err)
	}
	where := "status=?"
	if strings.TrimSpace(blogVO.Keyword) != "" {
		where += " and title like '%" + strings.TrimSpace(blogVO.Keyword) + "%'"
	}
	if blogVO.TagUid != "" {
		where += " and tag_uid like '%" + blogVO.TagUid + "%'"
	}
	if blogVO.BlogSortUid != "" {
		where += " and blog_sort_uid like '%" + blogVO.BlogSortUid + "%'"
	}
	levelKeyword := common.InterfaceToString(blogVO.LevelKeyword)
	if blogVO.LevelKeyword != "" {
		where += " and level='" + levelKeyword + "'"
	}
	if blogVO.IsPublish != "" {
		where += " and is_publish='" + blogVO.IsPublish + "'"
	}
	if blogVO.IsOriginal != "" {
		where += " and is_original='" + blogVO.IsOriginal + "'"
	}
	t := common.InterfaceToInt(blogVO.Type)
	if t != 0 {
		where += " and type=" + strconv.Itoa(t)
	}
	order := sortDesc
	if blogVO.OrderByAscColumn != "" {
		order = common.Camel2Case(blogVO.OrderByAscColumn) + " asc"
	} else if blogVO.OrderByDescColumn != "" {
		order = common.Camel2Case(blogVO.OrderByDescColumn) + " desc"
	} else if blogVO.UserSort == 0 {
		order = "create_time desc"
	}
	var pageList []models.Blog
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Blog{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((blogVO.CurrentPage - 1) * blogVO.PageSize).Limit(blogVO.PageSize).Order(order).Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	if len(pageList) == 0 {
		page := map[string]interface{}{}
		page["records"] = pageList
		c.SuccessWithData(page)
		return
	}
	var s []string
	var sortUids []string
	var tagUids []string
	for _, item := range pageList {
		if item.FileUid != "" {
			s = append(s, item.FileUid)
		}
		if item.BlogSortUid != "" {
			sortUids = append(sortUids, item.BlogSortUid)
		}
		if item.TagUid != "" {
			tagUidsTemp := strings.Split(item.TagUid, ",")
			for _, itemTagUid := range tagUidsTemp {
				tagUids = append(tagUids, itemTagUid)
			}
		}
	}
	pictureList := map[string]interface{}{}
	fileUids := strings.Join(s, ",")
	if fileUids != "" {
		pictureList = service.FileService.GetPicture(fileUids, ",")
	}
	picList := common.WebUtil.GetPictureMap(pictureList)
	var sortList []models.BlogSort
	var tagList []models.Tag
	c.Wg.Add(1)
	go func() {
		if len(sortUids) > 0 {
			common.DB.Find(&sortList, sortUids)
		}
		c.Wg.Done()
	}()
	c.Wg.Add(1)
	go func() {
		if len(tagUids) > 0 {
			common.DB.Find(&tagList, tagUids)
		}
		c.Wg.Done()
	}()
	c.Wg.Wait()
	sortMap := map[string]models.BlogSort{}
	tagMap := map[string]models.Tag{}
	pictureMap := map[string]string{}
	for _, item := range sortList {
		sortMap[item.Uid] = item
	}
	for _, item := range tagList {
		tagMap[item.Uid] = item
	}
	for _, item := range picList {
		pictureMap[item["uid"].(string)] = item["url"].(string)
	}
	for i, item := range pageList {
		if item.BlogSortUid != "" {
			pageList[i].BlogSort = sortMap[item.BlogSortUid]
		}
		if item.TagUid != "" {
			tagUidsTemp := strings.Split(item.TagUid, ",")
			var tagListTemp []models.Tag
			for _, tag := range tagUidsTemp {
				tagListTemp = append(tagListTemp, tagMap[tag])
			}
			pageList[i].TagList = tagListTemp
		}
		if item.FileUid != "" {
			pictureUidsTemp := strings.Split(item.FileUid, ",")
			var pictureListTemp []string
			for _, picture := range pictureUidsTemp {
				pictureListTemp = append(pictureListTemp, pictureMap[picture])
			}
			pageList[i].PhotoList = pictureListTemp
		}
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    blogVO.PageSize,
		Current: blogVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *BlogRestApi) Edit() {
	var blogVO vo.BlogVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogVO)
	if err != nil {
		panic(err)
	}
	var blog models.Blog
	common.DB.Where("uid=?", blogVO.Uid).Find(&blog)
	var count int64
	common.DB.Model(&models.Blog{}).Where("status=? and level= ?", 1, blogVO.Level).Count(&count)
	if !reflect.DeepEqual(blog, models.Blog{}) {
		if blog.Level != blogVO.Level {
			count++
		}
	}
	addVerdictResult := service.BlogService.AddVerdict(count, blogVO.Level)
	if addVerdictResult != "" {
		c.ErrorWithMessage(addVerdictResult)
		return
	}
	var admin models.Admin
	common.DB.Where("uid = ?", c.GetAdminUid()).Find(&admin)
	blog.AdminUid = admin.Uid
	if blogVO.IsOriginal == "1" {
		if admin.NickName != "" {
			blog.Author = admin.NickName
		} else {
			blog.Author = admin.UserName
		}
		projectName, _ := service.SysParamService.GetSysParamsValueByKey("PROJECT_NAME")
		blog.ArticlesPart = projectName
	} else {
		blog.Author = blogVO.Author
		blog.ArticlesPart = blogVO.ArticlesPart
	}
	t := common.InterfaceToInt(blogVO.Type)
	openComment := common.InterfaceToInt(blogVO.OpenComment)
	blog.Title = blogVO.Title
	blog.Summary = blogVO.Summary
	blog.Content = blogVO.Content
	blog.TagUid = blogVO.TagUid
	blog.BlogSortUid = blogVO.BlogSortUid
	blog.FileUid = blogVO.FileUid
	blog.Level = blogVO.Level
	blog.IsOriginal = blogVO.IsOriginal
	blog.IsPublish = blogVO.IsPublish
	blog.OpenComment = openComment
	blog.Type = t
	blog.OutsideLink = blogVO.OutsideLink
	err1 := common.DB.Save(&blog).Error
	isSave := true
	if err1 != nil {
		isSave = false
	}
	service.BlogService.UpdateSolrAndRedis(isSave, blog)
	common.RedisUtil.Delete("BLOG_LEVEL:" + strconv.Itoa(blog.Level))
	c.SuccessWithMessage("更新成功")
}

func (c *BlogRestApi) Delete() {
	var blogVO vo.BlogVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogVO)
	if err != nil {
		panic(err)
	}
	var blog models.Blog
	common.DB.Where("uid=?", blogVO.Uid).Find(&blog)
	save := common.DB.Model(&blog).Select("status").Update("status", 0).Error
	if save == nil {
		blogUidList := []string{blogVO.Uid}
		service.SubjectItemService.DeleteBatchSubjectItemByBlogUid(blogUidList)
	}
	common.RedisUtil.Delete("BLOG_LEVEL:" + strconv.Itoa(blog.Level))
	c.SuccessWithMessage("删除成功")
}

func (c *BlogRestApi) Add() {
	var blogVO vo.BlogVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogVO)
	if err != nil {
		panic(err)
	}
	var count int64
	common.DB.Model(&models.Blog{}).Where("status=? and level= ?", 1, blogVO.Level).Count(&count)
	addVerdictResult := service.BlogService.AddVerdict(count, blogVO.Level)
	if addVerdictResult != "" {
		c.ErrorWithMessage(addVerdictResult)
		return
	}
	var blog models.Blog
	projectName, _ := service.SysParamService.GetSysParamsValueByKey("PROJECT_NAME")
	if blogVO.IsOriginal == "1" {
		var admin models.Admin
		common.DB.Where("uid = ?", c.GetAdminUid()).Find(&admin)
		if !reflect.DeepEqual(admin, models.Admin{}) {
			if admin.NickName != "" {
				blog.Author = admin.NickName
			} else {
				blog.Author = admin.UserName
			}
			blog.AdminUid = admin.Uid
		}
		blog.ArticlesPart = projectName
	} else {
		blog.Author = blogVO.Author
		blog.ArticlesPart = blogVO.ArticlesPart
	}
	t := common.InterfaceToInt(blogVO.Type)
	openComment := common.InterfaceToInt(blogVO.OpenComment)
	blog.Uid = xid.New().String()
	blog.Title = blogVO.Title
	blog.Summary = blogVO.Summary
	blog.Content = blogVO.Content
	blog.TagUid = blogVO.TagUid
	blog.BlogSortUid = blogVO.BlogSortUid
	blog.FileUid = blogVO.FileUid
	blog.Level = blogVO.Level
	blog.IsOriginal = blogVO.IsOriginal
	blog.IsPublish = blogVO.IsPublish
	blog.OpenComment = openComment
	blog.Type = t
	blog.OutsideLink = blogVO.OutsideLink
	common.DB.Create(&blog)
	common.RedisUtil.Delete("BLOG_LEVEL:" + strconv.Itoa(blog.Level))
	c.SuccessWithMessage("新增成功")
}

func (c *BlogRestApi) UploadLocalBlog() {
	fileDatas, _ := c.GetFiles("filedatas")
	m := service.SystemConfigService.GetConfig()
	systemConfig := m["data"].(models.SystemConfig)
	if systemConfig == (models.SystemConfig{}) {
		c.ErrorWithMessage("系统配置不存在,请重新登录")
		return
	}
	if systemConfig.UploadQiNiu == "1" && (systemConfig.QiNiuPictureBaseUrl == "" || systemConfig.QiNiuAccessKey == "") {
		c.ErrorWithMessage("请先设置七牛云")
		return
	}
	var fileList []*multipart.FileHeader
	var fileNameList []string
	for _, file := range fileDatas {
		fileOriginalName := file.Filename
		if !strings.HasSuffix(fileOriginalName, ".md") {
			c.ErrorWithMessage("目前仅支持Markdown文件")
			fileList = []*multipart.FileHeader{}
			break
		}
		fileList = append(fileList, file)
		fileNameList = append(fileNameList, common.FileUtil.GetFileName(fileOriginalName))
	}
	if len(fileList) == 0 {
		c.ErrorWithMessage("请选中需要上传的Markdown文件")
		return
	}
	var fileContentList []string
	for _, mutipartFile := range fileList {
		file, _ := mutipartFile.Open()
		fil := make([][]byte, 0)
		var b int64 = 0
		for {
			buffer := make([]byte, 1024)
			n, err := file.ReadAt(buffer, b)
			b += int64(n)
			fil = append(fil, buffer)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
		fileStream := bytes.Join(fil, []byte(""))
		blogContent := common.FileUtil.MarkdownToHTML(string(fileStream))
		fileContentList = append(fileContentList, blogContent)
		pictureList := c.GetString("pictureList")
		var list []maps.PictureListMap
		err := json.Unmarshal([]byte(pictureList), &list)
		if err != nil {
			panic(err)
		}
		pictureMap := map[string]string{}
		for _, item := range list {
			if systemConfig.PicturePriority == "1" {
				pictureMap[item.FileOldName] = item.QiNiuUrl
			} else {
				c.ErrorWithMessage("暂不支持其他存储")
				return
			}
		}
		matchUrlMap := map[string]string{}
		var matchList []string
		for _, blogContent := range fileContentList {
			if common.CheckPicture(blogContent) {
				matchList = append(matchList, blogContent)
			}
		}
		if len(matchList) > 0 {
			for _, matchStr := range matchList {
				splitList := strings.Split(matchStr, `\`)
				if len(splitList) >= 5 {
					var pictureUrl string
					if strings.Index(matchStr, "alt") > strings.Index(matchStr, "src") {
						pictureUrl = splitList[1]
					} else {
						pictureUrl = splitList[3]
					}
					if !strings.HasPrefix(pictureUrl, "http") {
						for key, value := range pictureMap {
							if strings.Contains(pictureUrl, key) {
								if systemConfig.PicturePriority == "1" {
									matchUrlMap[pictureUrl] = systemConfig.QiNiuPictureBaseUrl + value
								} else {
									c.ErrorWithMessage("暂不支持其他存储")
									return
								}
								break
							}
						}
					}
				}
			}
		}
		blogSort := service.BlogSortService.GetTopOne()
		tag := service.TagService.GetTopTag()
		picture := service.PictureService.GetTopOne()
		if blogSort == (models.BlogSort{}) || tag == (models.Tag{}) || picture == (models.Picture{}) {
			c.ErrorWithMessage("使用本地上传，请先确保博客分类，博客标签，博客图片中含有数据")
			return
		}
		admin := c.GetMeService()
		var blogList []models.Blog
		var count int
		projectName, _ := service.SysParamService.GetSysParamsValueByKey("PROJECT_NAME")
		var levels []string
		for _, content := range fileContentList {
			if len(matchUrlMap) > 0 {
				for key, value := range matchUrlMap {
					content = strings.ReplaceAll(content, key, value)
				}
			}
			blog := models.Blog{
				Uid:          xid.New().String(),
				BlogSortUid:  blogSort.Uid,
				TagUid:       tag.Uid,
				AdminUid:     admin.Uid,
				Author:       admin.NickName,
				ArticlesPart: projectName,
				Level:        0,
				Title:        fileNameList[count],
				Summary:      fileNameList[count],
				Content:      content,
				FileUid:      picture.FileUid,
				IsOriginal:   "1",
				IsPublish:    "0",
				OpenComment:  1,
				Type:         0,
			}
			blogList = append(blogList, blog)
			count++
			levels = append(levels, "BLOG_LEVEL:"+strconv.Itoa(blog.Level))
		}
		common.DB.Create(blogList)
		common.RedisUtil.MultiDelete(levels)
		c.SuccessWithMessage("新增成功")
	}
}

func (c *BlogRestApi) DeleteBatch() {
	var blogVOList []vo.BlogVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogVOList)
	if err != nil {
		panic(err)
	}
	if len(blogVOList) == 0 {
		c.ErrorWithMessage("参数错误")
	}
	var uidList []string
	var s []string
	for _, item := range blogVOList {
		uidList = append(uidList, item.Uid)
		s = append(s, item.Uid+",")
	}
	var blogList []models.Blog
	common.DB.Find(&blogList, uidList)
	save := common.DB.Model(&blogList).Select("status").Update("status", 0).Error
	if save == nil {
		service.SubjectItemService.DeleteBatchSubjectItemByBlogUid(uidList)
	}
	c.SuccessWithMessage("删除成功")
}

func (c *BlogRestApi) EditBatch() {
	var blogVOList []vo.BlogVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &blogVOList)
	if err != nil {
		panic(err)
	}
	if len(blogVOList) == 0 {
		c.ErrorWithMessage("参数错误")
	}
	var blogUidList []string
	blogVOMap := make(map[string]vo.BlogVO)
	for _, item := range blogVOList {
		blogUidList = append(blogUidList, item.Uid)
		blogVOMap[item.Uid] = item
	}
	var blogList []models.Blog
	common.DB.Find(&blogList, blogUidList)
	var levels []string
	for i, blog := range blogList {
		blogVO := blogVOMap[blog.Uid]
		if !reflect.DeepEqual(blogVO, vo.BlogVO{}) {
			_type := common.InterfaceToInt(blogVO.Type)
			blogList[i].Author = blogVO.Author
			blogList[i].ArticlesPart = blogVO.ArticlesPart
			blogList[i].Title = blogVO.Title
			blogList[i].Summary = blogVO.Summary
			blogList[i].Content = blogVO.Content
			blogList[i].TagUid = blogVO.TagUid
			blogList[i].BlogSortUid = blogVO.BlogSortUid
			blogList[i].FileUid = blogVO.FileUid
			blogList[i].Level = blogVO.Level
			blogList[i].IsOriginal = blogVO.IsOriginal
			blogList[i].IsPublish = blogVO.IsPublish
			blogList[i].Sort = blogVO.Sort
			blogList[i].Type = _type
			blogList[i].OutsideLink = blogVO.OutsideLink
		}
		levels = append(levels, "BLOG_LEVEL:"+strconv.Itoa(blog.Level))
	}
	common.DB.Save(&blogList)
	common.RedisUtil.MultiDelete(levels)
	c.SuccessWithMessage("更新成功")
}
