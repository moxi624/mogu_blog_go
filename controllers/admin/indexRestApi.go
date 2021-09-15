package admin

import (
	"encoding/json"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/maps"
	"strings"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/22 4:14 下午
 * @version 1.0
 */
type IndexRestApi struct {
	base.BaseController
}

func (c *IndexRestApi) GetVisitByWeek() {
	weekVisitJson := common.RedisUtil.Get("DASHBOARD:WEEK_VISIT")
	if weekVisitJson != "" {
		weekVisitMap := map[string]interface{}{}
		json.Unmarshal([]byte(weekVisitJson), &weekVisitMap)
		c.Result("success", weekVisitMap)
		return
	}
	todayEndTime := common.DateUtils.GetToDayEndTime()
	sevenDays := common.DateUtils.GetDate(todayEndTime, -6)
	sevenDaysList := common.DateUtils.GetDaysByN(7, "2006-01-02")
	var pvMap []maps.PvMap
	common.DB.Table("t_web_visit").Distinct("DATE_FORMAT(create_time, '%Y-%m-%d') as date, COUNT(uid) as count").Where("create_time >= ? AND create_time <= ? ", sevenDays, todayEndTime).Group("DATE_FORMAT(create_time, '%Y-%m-%d')").Scan(&pvMap)
	var uvMap []maps.UvMap
	subQuery := common.DB.Table("t_web_visit").Distinct("DATE_FORMAT(create_time, '%Y-%m-%d') DATE, ip").Where("create_time >= ? AND create_time <= ?", sevenDays, todayEndTime).Group("DATE_FORMAT(create_time, '%Y-%m-%d'),ip")
	common.DB.Table("(?) as tmp", subQuery).Select("DATE, COUNT(ip) COUNT").Group("date").Scan(&uvMap)
	countPVMap := map[string]int{}
	countUVMap := map[string]int{}
	for _, item := range pvMap {
		countPVMap[item.Date] = item.Count
	}
	for _, item := range uvMap {
		countUVMap[item.Date] = item.Count
	}
	var pvList []int
	var uvList []int
	for _, day := range sevenDaysList {
		if countPVMap[day] != 0 {
			pvNumber := countPVMap[day]
			uvNumber := countUVMap[day]
			pvList = append(pvList, pvNumber)
			uvList = append(uvList, uvNumber)
		} else {
			pvList = append(pvList, 0)
			uvList = append(uvList, 0)
		}
	}
	resultMap := map[string]interface{}{}
	resultSevenDaysList := common.DateUtils.GetDaysByN(7, "01-02")
	resultMap["data"] = resultSevenDaysList
	resultMap["pv"] = pvList
	resultMap["uv"] = uvList
	b, _ := json.Marshal(resultMap)
	common.RedisUtil.SetEx("DASHBOARD:WEEK_VISIT", string(b), 10, time.Minute)
	c.Result("success", resultMap)
}

func (c *IndexRestApi) InitIndex() {
	m := map[string]interface{}{}
	c.Wg.Add(1)
	go func() {
		var blog models.Blog
		var blogCount int64
		common.DB.Model(&blog).Where("status=? and is_publish=?", 1, "1").Count(&blogCount)
		m["blogCount"] = blogCount
		c.Wg.Done()
	}()
	c.Wg.Add(1)
	go func() {
		var comment models.Comment
		var commentCount int64
		common.DB.Model(&comment).Where("status=?", 1).Count(&commentCount)
		m["commentCount"] = commentCount
		c.Wg.Done()
	}()
	c.Wg.Add(1)
	go func() {
		var user models.User
		var userCount int64
		common.DB.Model(&user).Where("status=?", 1).Count(&userCount)
		m["userCount"] = userCount
		c.Wg.Done()
	}()
	c.Wg.Add(1)
	go func() {
		startTime := common.DateUtils.GetToDayStartTime()
		endTime := common.DateUtils.GetToDayEndTime()
		var visitCount int64
		common.DB.Table("t_web_visit").Select("count(distinct(ip))").Where("create_time >= ? AND create_time <= ?", startTime, endTime).Count(&visitCount)
		m["visitCount"] = visitCount
		c.Wg.Done()
	}()
	c.Wg.Wait()
	c.Result("success", m)
}

func (c *IndexRestApi) GetBlogCountByTag() {
	jsonArrayList := common.RedisUtil.Get("DASHBOARD:BLOG_COUNT_BY_TAG")
	if jsonArrayList != "" {
		var jsonList []map[string]interface{}
		json.Unmarshal([]byte(jsonArrayList), &jsonList)
		c.Result("success", jsonList)
		return
	}
	var blogCountByTagMap []maps.BlogCountByTagMap
	common.DB.Table("t_blog").Select("tag_uid, COUNT(tag_uid) as count").Group("tag_uid").Scan(&blogCountByTagMap)
	tagMap := map[string]int{}
	for _, item := range blogCountByTagMap {
		tagUid := item.TagUid
		count := item.Count
		if len(tagUid) == 32 {
			if tagMap[tagUid] == 0 {
				tagMap[tagUid] = count
			} else {
				tempCount := tagMap[tagUid] + count
				tagMap[tagUid] = tempCount
			}
		} else {
			if tagUid != "" {
				strList := strings.Split(tagUid, ",")
				for _, strItem := range strList {
					if tagMap[strItem] != 0 {
						tagMap[strItem] = count
					} else {
						tempCount := tagMap[strItem] + count
						tagMap[strItem] = tempCount
					}
				}
			}

		}
	}
	tagUids := common.GetKeysInt(tagMap)
	tagUids = common.RemoveRepByMap(tagUids)
	var tagCollection []models.Tag
	if len(tagUids) > 0 {
		common.DB.Find(&tagCollection, tagUids)
	}
	tagEntityMap := map[string]string{}
	for _, tag := range tagCollection {
		if tag.Content != "" {
			tagEntityMap[tag.Uid] = tag.Content
		}
	}
	var resultList []map[string]interface{}
	for tagUid, count := range tagMap {
		if tagEntityMap[tagUid] != "" {
			tagName := tagEntityMap[tagUid]
			itemResultMap := map[string]interface{}{
				"tagUid": tagUid,
				"name":   tagName,
				"value":  count,
			}
			resultList = append(resultList, itemResultMap)
		}
	}
	if len(resultList) > 0 {
		b, _ := json.Marshal(resultList)
		common.RedisUtil.SetEx("DASHBOARD:BLOG_COUNT_BY_TAG", string(b), 2, time.Hour)
	}
	c.Result("success", resultList)
}

func (c *IndexRestApi) GetBlogCountByBlogSort() {
	jsonArrayList := common.RedisUtil.Get("DASHBOARD:BLOG_COUNT_BY_SORT")
	if jsonArrayList != "" {
		var jsonList []map[string]interface{}
		json.Unmarshal([]byte(jsonArrayList), &jsonList)
		c.Result("success", jsonList)
		return
	}
	var blogCountByBlogSortMap []maps.BlogCoutByBlogSortMap
	common.DB.Table("t_blog").Select("blog_sort_uid, COUNT(blog_sort_uid) AS count").Where("status = 1").Group("blog_sort_uid").Scan(&blogCountByBlogSortMap)
	blogSortMap := map[string]int{}
	for _, item := range blogCountByBlogSortMap {
		blogSortUid := item.BlogSortUid
		count := item.Count
		blogSortMap[blogSortUid] = count
	}
	blogSortUids := common.GetKeysInt(blogSortMap)
	blogSortUids = common.RemoveRepByMap(blogSortUids)
	var blogSortCollection []models.BlogSort
	if len(blogSortUids) > 0 {
		common.DB.Find(&blogSortCollection, blogSortUids)
	}
	blogSortEntityMap := map[string]string{}
	for _, blogSort := range blogSortCollection {
		if blogSort.SortName != "" {
			blogSortEntityMap[blogSort.Uid] = blogSort.SortName
		}
	}
	var resultList []map[string]interface{}
	for blogSortUid, count := range blogSortMap {
		if blogSortEntityMap[blogSortUid] != "" {
			blogSortName := blogSortEntityMap[blogSortUid]
			itemResultMap := map[string]interface{}{
				"blogSortUid": blogSortUid,
				"name":        blogSortName,
				"value":       count,
			}
			resultList = append(resultList, itemResultMap)
		}
	}
	if len(resultList) > 0 {
		b, _ := json.Marshal(resultList)
		common.RedisUtil.SetEx("DASHBOARD:BLOG_COUNT_BY_SORT", string(b), 2, time.Hour)
	}
	c.Result("success", resultList)
}

func (c *IndexRestApi) GetBlogContributeCount() {
	jsonMap := common.RedisUtil.Get("DASHBOARD:BLOG_CONTRIBUTE_COUNT")
	if jsonMap != "" {
		var tempMap map[string]interface{}
		json.Unmarshal([]byte(jsonMap), &tempMap)
		c.Result("success", tempMap)
		return
	}
	endTime := common.DateUtils.GetNowTime()
	startTime := common.DateUtils.GetDate(endTime, -365)
	var blogContributeMap []maps.BlogContributeMap
	common.DB.Table("t_blog").Distinct("DATE_FORMAT(create_time, '%Y-%m-%d') as date, COUNT(uid) as count").Where("1=1 && status = 1 && create_time >= ? && create_time < ?", startTime, endTime).Group("DATE_FORMAT(create_time, '%Y-%m-%d')").Scan(&blogContributeMap)
	dateList := common.DateUtils.GetDayBetweenDates(startTime, endTime)
	dataMap := map[string]int{}
	for _, itemMap := range blogContributeMap {
		dataMap[itemMap.Date] = itemMap.Count
	}
	var resultList [][]interface{}
	for _, item := range dateList {
		count := 0
		if dataMap[item] != 0 {
			count = dataMap[item]
		}
		var objectList []interface{}
		objectList = append(objectList, item)
		objectList = append(objectList, count)
		resultList = append(resultList, objectList)
	}
	resultMap := map[string]interface{}{}
	var countributeDateList []string
	countributeDateList = append(countributeDateList, startTime)
	countributeDateList = append(countributeDateList, endTime)
	resultMap["contributeDate"] = countributeDateList
	resultMap["blogContributeCount"] = resultList
	b, _ := json.Marshal(resultMap)
	common.RedisUtil.SetEx("DASHBOARD:BLOG_CONTRIBUTE_COUNT", string(b), 2, time.Hour)
	c.Result("success", resultMap)
}
