package service

import (
	"github.com/thinkeridea/go-extend/exunicode/exutf8"
	"mogu-go-v2/common"
	"mogu-go-v2/models"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/28 1:11 下午
 * @version 1.0
 */

type blogService struct {
	wg sync.WaitGroup
}

func (blogService) DeleteRedisByBlog() {
	common.RedisUtil.Delete("NEW_BLOG")
	common.RedisUtil.Delete("HOT_BLOG")
	common.RedisUtil.Delete("BLOG_LEVEL:1")
	common.RedisUtil.Delete("BLOG_LEVEL:2")
	common.RedisUtil.Delete("BLOG_LEVEL:3")
	common.RedisUtil.Delete("BLOG_LEVEL:4")
}

func (blogService) AddVerdict(count int64, level int) string {
	switch level {
	case 1:
		blogCount, _ := SysParamService.GetSysParamsValueByKey("BLOG_FIRST_COUNT")
		i, _ := strconv.Atoi(blogCount)
		if count > int64(i) {
			return "一级推荐不能超过" + blogCount + "个"
		}
		break
	case 2:
		blogCount, _ := SysParamService.GetSysParamsValueByKey("BLOG_SECOND_COUNT")
		i, _ := strconv.Atoi(blogCount)
		if count > int64(i) {
			return "二级推荐不能超过" + blogCount + "个"
		}
		break
	case 3:
		blogCount, _ := SysParamService.GetSysParamsValueByKey("BLOG_THIRD_COUNT")
		i, _ := strconv.Atoi(blogCount)
		if count > int64(i) {
			return "三级推荐不能超过" + blogCount + "个"
		}
		break
	case 4:
		blogCount, _ := SysParamService.GetSysParamsValueByKey("BLOG_FOURTH_COUNT")
		i, _ := strconv.Atoi(blogCount)
		if count > int64(i) {
			return "四级推荐不能超过" + blogCount + "个"
		}
		break
	default:
		return ""
	}
	return ""
}

func (blogService) UpdateSolrAndRedis(isSave bool, blog models.Blog) {
	if isSave && blog.IsPublish == "1" {
		/*m:=map[string]interface{}{
			"command":"add",
			"blog_uid":blog.Uid,
		}*/
	}
}

func (c *blogService) DeleteRedisByBlogSort() {
	common.RedisUtil.Delete("DASHBOARD:BLOG_COUNT_BY_SORT")
	c.DeleteRedisByBlog()
}

func (c *blogService) DeleteRedisByBlogTag() {
	common.RedisUtil.Delete("DASHBOARD:BLOG_COUNT_BY_TAG")
	c.DeleteRedisByBlog()
}

func (s *blogService) SetTagAndSortAndPictureByBlogList(list []models.Blog) []models.Blog {
	var sortUids, tagUids, fileUidsSet []string
	for _, item := range list {
		if item.FileUid != "" {
			fileUidsSet = append(fileUidsSet, item.FileUid)
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
	var pictureList map[string]interface{}
	var str []string
	var picList []map[string]interface{}
	fileUidsSet = common.RemoveRepByMap(fileUidsSet)
	if len(fileUidsSet) > 0 {
		count := 1
		var fileUids string
		for _, fileUid := range fileUidsSet {
			str = append(str, fileUid+",")
			l.Print(count % 10)
			if count%10 == 0 {
				fileUids = strings.Join(str, ",")
				pictureList = FileService.GetPicture(fileUids, ",")
				tempPicList := common.WebUtil.GetPictureMap(pictureList)
				picList = append(picList, tempPicList...)
				str = []string{}
			}
			count++
		}
		if len(str) >= 32 {
			pictureList = FileService.GetPicture(fileUids, ",")
			tempPicList := common.WebUtil.GetPictureMap(pictureList)
			picList = append(picList, tempPicList...)
		}
	}
	var sortList []models.BlogSort
	var tagList []models.Tag
	s.wg.Add(1)
	go func() {
		if len(str) > 0 {
			common.DB.Find(&sortList, sortUids)
		}
		s.wg.Done()
	}()
	s.wg.Add(1)
	go func() {
		if len(tagUids) > 0 {
			common.DB.Find(&tagList, tagUids)
		}
		s.wg.Done()
	}()
	s.wg.Wait()
	sortMap := map[string]models.BlogSort{}
	tagMap := map[string]models.Tag{}
	pictureMap := map[string]string{}
	s.wg.Add(1)
	go func() {
		for _, item := range sortList {
			sortMap[item.Uid] = item
		}
		s.wg.Done()
	}()
	s.wg.Add(1)
	go func() {
		for _, item := range tagList {
			tagMap[item.Uid] = item
		}
		s.wg.Done()
	}()
	s.wg.Add(1)
	go func() {
		for _, item := range picList {
			pictureMap[item["uid"].(string)] = item["url"].(string)
		}
		s.wg.Done()
	}()
	s.wg.Wait()
	for i, item := range list {
		s.wg.Add(3)
		item := item
		i := i
		go func() {
			if item.BlogSortUid != "" {
				list[i].BlogSort = sortMap[item.BlogSortUid]
				if sortMap[item.BlogSortUid] != (models.BlogSort{}) {
					list[i].BlogSortName = sortMap[item.BlogSortUid].SortName
				}
			}
			s.wg.Done()
		}()
		go func() {
			if item.TagUid != "" {
				tagUidsTemp := strings.Split(item.TagUid, ",")
				var tagListTemp []models.Tag
				for _, tag := range tagUidsTemp {
					tagListTemp = append(tagListTemp, tagMap[tag])
				}
				list[i].TagList = tagListTemp
			}
			s.wg.Done()
		}()
		go func() {
			if item.FileUid != "" {
				pictureUidsTemp := strings.Split(item.FileUid, ",")
				var pictureListTemp []string
				for _, picture := range pictureUidsTemp {
					pictureListTemp = append(pictureListTemp, picture)
				}
				list[i].PhotoList = pictureListTemp
				if len(pictureListTemp) > 0 {
					list[i].PhotoUrl = pictureListTemp[0]
				} else {
					list[i].PhotoUrl = ""
				}
			}
			s.wg.Done()
		}()
	}
	s.wg.Wait()
	return list
}

func (s *blogService) GetBlogPageByLevel(currentPage int, pageSize int, level int, userSort int) ([]models.BlogNoContent, int64) {
	var pageList []models.BlogNoContent
	var order string
	if userSort == 0 {
		order = "create_time desc"
	} else {
		order = "sort desc"
	}
	var total int64
	s.wg.Add(2)
	go func() {
		common.DB.Model(&models.Blog{}).Where("level=? and status=? and is_publish=?", level, 1, "1", 1).Count(&total)
		s.wg.Done()
	}()
	go func() {
		common.DB.Where("level=? and status=? and is_publish=?", level, 1, "1").Offset((currentPage - 1) * pageSize).Limit(pageSize).Order(order).Find(&pageList)
		s.wg.Done()
	}()
	s.wg.Wait()
	return pageList, total
}

func (s *blogService) SetBlog(list []models.BlogNoContent) []models.BlogNoContent {
	var fileUids strings.Builder
	var sortUids []string
	var tagUids []string
	for _, item := range list {
		if item.FileUid != "" {
			fileUids.WriteString(item.FileUid + ",")
		}
		if item.BlogSortUid != "" {
			sortUids = append(sortUids, item.BlogSortUid)
		}
		if item.TagUid != "" {
			tagUids = append(tagUids, item.TagUid)
		}
	}
	pictureList := map[string]interface{}{}
	if !reflect.DeepEqual(fileUids, strings.Builder{}) {
		pictureList = FileService.GetPicture(fileUids.String(), ",")
	}
	picList := common.WebUtil.GetPictureMap(pictureList)
	var sortList []models.BlogSort
	var tagList []models.Tag
	s.wg.Add(1)
	go func() {
		if len(sortUids) > 0 {
			common.DB.Find(&sortList, sortUids)
		}
		s.wg.Done()
	}()
	s.wg.Add(1)
	go func() {
		if len(tagUids) > 0 {
			common.DB.Find(&tagList, tagUids)
		}
		s.wg.Done()
	}()
	s.wg.Wait()
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
	for i, item := range list {
		if item.BlogSortUid != "" {
			list[i].BlogSort = sortMap[item.BlogSortUid]
		}
		if item.TagUid != "" {
			tagUidsTemp := strings.Split(item.TagUid, ",")
			var tagListTemp []models.Tag
			for _, tag := range tagUidsTemp {
				if tagMap[tag] != (models.Tag{}) {
					tagListTemp = append(tagListTemp, tagMap[tag])
				}
			}
			list[i].TagList = tagListTemp
		}
		if item.FileUid != "" {
			pictureUidsTemp := strings.Split(item.FileUid, ",")
			var pictureListTemp []string
			for _, picture := range pictureUidsTemp {
				pictureListTemp = append(pictureListTemp, pictureMap[picture])
			}
			list[i].PhotoList = pictureListTemp
		}
	}
	return list
}

func (s *blogService) SetTagAndSortAndPictureByBlogListNoContent(list []models.BlogNoContent) []models.BlogNoContent {
	var sortUids, tagUids, fileUidsSet []string
	for _, item := range list {
		if item.FileUid != "" {
			fileUidsSet = append(fileUidsSet, item.FileUid)
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
	var pictureList map[string]interface{}
	var str []string
	var picList []map[string]interface{}
	fileUidsSet = common.RemoveRepByMap(fileUidsSet)
	if len(fileUidsSet) > 0 {
		count := 1
		var fileUids string
		for _, fileUid := range fileUidsSet {
			str = append(str, fileUid+",")
			l.Print(count % 10)
			if count%10 == 0 {
				fileUids = strings.Join(str, ",")
				pictureList = FileService.GetPicture(fileUids, ",")
				tempPicList := common.WebUtil.GetPictureMap(pictureList)
				picList = append(picList, tempPicList...)
				str = []string{}
			}
			count++
		}
		if len(str) >= 32 {
			pictureList = FileService.GetPicture(fileUids, ",")
			tempPicList := common.WebUtil.GetPictureMap(pictureList)
			picList = append(picList, tempPicList...)
		}
	}
	var sortList []models.BlogSort
	var tagList []models.Tag
	s.wg.Add(2)
	go func() {
		if len(str) > 0 {
			common.DB.Find(&sortList, sortUids)
		}
		s.wg.Done()
	}()
	go func() {
		if len(tagUids) > 0 {
			common.DB.Find(&tagList, tagUids)
		}
		s.wg.Done()
	}()
	s.wg.Wait()
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
	for i, item := range list {
		item := item
		i := i
		go func() {
			if item.BlogSortUid != "" {
				list[i].BlogSort = sortMap[item.BlogSortUid]
				if sortMap[item.BlogSortUid] != (models.BlogSort{}) {
					list[i].BlogSortName = sortMap[item.BlogSortUid].SortName
				}
			}
		}()
		go func() {
			if item.TagUid != "" {
				tagUidsTemp := strings.Split(item.TagUid, ",")
				var tagListTemp []models.Tag
				for _, tag := range tagUidsTemp {
					tagListTemp = append(tagListTemp, tagMap[tag])
				}
				list[i].TagList = tagListTemp
			}
		}()
		go func() {
			if item.FileUid != "" {
				pictureUidsTemp := strings.Split(item.FileUid, ",")
				var pictureListTemp []string
				for _, picture := range pictureUidsTemp {
					pictureListTemp = append(pictureListTemp, picture)
				}
				list[i].PhotoList = pictureListTemp
				if len(pictureListTemp) > 0 {
					list[i].PhotoUrl = pictureListTemp[0]
				} else {
					list[i].PhotoUrl = ""
				}
			}
		}()
	}
	return list
}

func (*blogService) GetHitCode(str, keyword string) string {
	if str == "" || keyword == "" {
		return str
	}
	startStr := "<span style = 'color:red'>"
	endStr := "</span>"
	if str == keyword {
		return startStr + str + endStr
	}
	lowerCaseStr := strings.ToLower(str)
	lowerKeyword := strings.ToLower(keyword)
	lowerCaseArray := strings.Split(lowerCaseStr, lowerKeyword)
	isEndWith := strings.HasSuffix(lowerCaseStr, lowerKeyword)
	var count int
	var list []map[string]int
	var keyList []map[string]int
	for a := 0; a < len(lowerCaseArray); a++ {
		m := map[string]int{}
		keyMap := make(map[string]int)
		m["startIndex"] = count
		l := utf8.RuneCountInString(lowerCaseArray[a])
		count += l
		m["endIndex"] = count
		list = append(list, m)
		if a < len(lowerCaseArray)-1 || isEndWith {
			keyMap["startIndex"] = count
			count += utf8.RuneCountInString(keyword)
			keyMap["endIndex"] = count
			keyList = append(keyList, keyMap)
		}
	}
	var arrayList []string
	for _, item := range list {
		start := item["startIndex"]
		end := item["endIndex"]
		itemStr := exutf8.RuneSubString(str, start, end-start)
		arrayList = append(arrayList, itemStr)
	}
	var keyArrayList []string
	for _, item := range keyList {
		start := item["startIndex"]
		end := item["endIndex"]
		itemStr := exutf8.RuneSubString(str, start, end-start)
		keyArrayList = append(keyArrayList, itemStr)
	}
	var sb strings.Builder
	for a := 0; a < len(arrayList); a++ {
		sb.WriteString(arrayList[a])
		if a < len(arrayList)-1 || isEndWith {
			sb.WriteString(startStr)
			sb.WriteString(keyArrayList[a])
			sb.WriteString(endStr)
		}
	}
	return sb.String()
}

func (*blogService) SetTagByBlog(blog *models.Blog) {
	tagUid := blog.TagUid
	if tagUid != "" {
		uids := strings.Split(tagUid, ",")
		var tagList []models.Tag
		for _, uid := range uids {
			var tag models.Tag
			common.DB.Where("uid=?", uid).Find(&tag)
			if tag != (models.Tag{}) && tag.Status != 0 {
				tagList = append(tagList, tag)
			}
		}
		blog.TagList = tagList
	}
}

func (*blogService) SetSortByBlog(blog *models.Blog) {
	if !reflect.DeepEqual(blog, models.Blog{}) && blog.BlogSortUid != "" {
		var blogSort models.BlogSort
		common.DB.Where("uid=?", blog.BlogSortUid).Find(&blogSort)
		blog.BlogSort = blogSort
	}
}

func (s *blogService) SetTagAndSortByBlogList(list []models.Blog) []models.Blog {
	var sortUids []string
	var tagUids []string
	for _, item := range list {
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
	var sortList []models.BlogSort
	var tagList []models.Tag
	s.wg.Add(2)
	go func() {
		if len(sortUids) > 0 {
			common.DB.Find(&sortList, sortUids)
		}
		s.wg.Done()
	}()
	go func() {
		if len(tagUids) > 0 {
			common.DB.Find(&tagList, tagUids)
		}
		s.wg.Done()
	}()
	s.wg.Wait()
	sortMap := map[string]models.BlogSort{}
	tagMap := map[string]models.Tag{}
	for _, item := range sortList {
		sortMap[item.Uid] = item
	}
	for _, item := range tagList {
		tagMap[item.Uid] = item
	}
	for i, item := range list {
		if item.BlogSortUid != "" {
			list[i].BlogSort = sortMap[item.BlogSortUid]
		}
		if item.TagUid != "" {
			tagUidsTemp := strings.Split(item.TagUid, ",")
			var tagListTemp []models.Tag
			for _, tag := range tagUidsTemp {
				tagListTemp = append(tagListTemp, tagMap[tag])
			}
			list[i].TagList = tagListTemp
		}
	}
	return list
}

var BlogService = &blogService{}
