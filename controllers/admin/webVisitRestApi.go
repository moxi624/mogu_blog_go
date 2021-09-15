package admin

import (
	"encoding/json"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
	"strconv"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/7 1:19 下午
 * @version 1.0
 */

type WebVisitRestApi struct {
	base.BaseController
}

func (c *WebVisitRestApi) GetList() {
	var webVisitVO vo.WebVisitVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &webVisitVO)
	if err != nil {
		panic(err)
	}
	where := "status=?"
	eBehavior := common.Emu.BehaviorEmu()
	if strings.TrimSpace(webVisitVO.Keyword) != "" {
		behavior := ""
		for _, v := range eBehavior {
			if v["content"] == strings.TrimSpace(webVisitVO.Keyword) {
				behavior = v["content"]
				break
			}
		}
		where += " and (ip like \"%" + strings.TrimSpace(webVisitVO.Keyword) + "%\" or behavior = \"" + behavior + "\" )"
	}
	if webVisitVO.StartTime != "" {
		time := strings.Split(webVisitVO.StartTime, ",")
		if len(time) == 2 {
			where += " and (create_time between \"" + common.DateUtils.Str2Date(time[0]) + "\" and \"" + common.DateUtils.Str2Date(time[1]) + "\")"
		}
	}
	var pageList []models.WebVisit
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.WebVisit{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((webVisitVO.CurrentPage - 1) * webVisitVO.PageSize).Limit(webVisitVO.PageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var blogUids []string
	var blogOids []string
	var tagUids []string
	var sortUids []string
	var linkUids []string
	for _, item := range pageList {
		if item.Behavior == eBehavior["BLOG_CONTENT"]["behavior"] {
			if item.ModuleUid != "" {
				blogUids = append(blogUids, item.ModuleUid)
			} else if item.OtherData != "" {
				blogOids = append(blogOids, item.OtherData)
			}
		} else if item.Behavior == eBehavior["BLOG_SORT"]["behavior"] || item.Behavior == eBehavior["VISIT_CLASSIFY"]["behavior"] {
			sortUids = append(sortUids, item.ModuleUid)
		} else if item.Behavior == eBehavior["BLOG_TAG"]["behavior"] || item.Behavior == eBehavior["VISIT_TAG"]["behavior"] {
			tagUids = append(tagUids, item.ModuleUid)
		} else if item.Behavior == eBehavior["FRIENDSHIP_LINK"]["behavior"] {
			linkUids = append(linkUids, item.ModuleUid)
		}
	}
	var blogList []models.Blog
	var blogListByOid []models.Blog
	var tagList []models.Tag
	var sortList []models.BlogSort
	var linkList []models.Link
	contentMap := map[string]string{}
	c.Wg.Add(1)
	go func() {
		if len(blogUids) > 0 {
			common.DB.Find(&blogList, blogUids)
		}
		for _, item := range blogList {
			contentMap[item.Uid] = item.Title
		}
		c.Wg.Done()
	}()
	c.Wg.Add(1)
	go func() {
		if len(blogOids) > 0 {
			common.DB.Where("oid in ?", blogOids).Find(&blogListByOid)
		}
		for _, item := range blogListByOid {
			contentMap[strconv.Itoa(item.Oid)] = item.Title
		}
		c.Wg.Done()
	}()
	c.Wg.Add(1)
	go func() {
		if len(tagUids) > 0 {
			common.DB.Find(&tagList, tagUids)
		}
		for _, item := range tagList {
			contentMap[item.Uid] = item.Content
		}
		c.Wg.Done()
	}()
	c.Wg.Add(1)
	go func() {
		if len(sortUids) > 0 {
			common.DB.Find(&sortList, sortUids)
		}
		for _, item := range sortList {
			contentMap[item.Uid] = item.SortName
		}
		c.Wg.Done()
	}()
	c.Wg.Add(1)
	go func() {
		if len(linkUids) > 0 {
			common.DB.Find(&linkList, linkUids)
		}
		for _, item := range linkList {
			contentMap[item.Uid] = item.Title
		}
		c.Wg.Done()
	}()
	c.Wg.Wait()
	for k, item := range pageList {
		for key := range eBehavior {
			if eBehavior[key]["behavior"] == item.Behavior {
				pageList[k].BehaviorContent = eBehavior[key]["content"]
				break
			}
		}
		if item.Behavior == eBehavior["BLOG_CONTENT"]["behavior"] || item.Behavior == eBehavior["BLOG_SORT"]["behavior"] ||
			item.Behavior == eBehavior["BLOG_TAG"]["behavior"] || item.Behavior == eBehavior["VISIT_TAG"]["behavior"] ||
			item.Behavior == eBehavior["VISIT_CLASSIFY"]["behavior"] || item.Behavior == eBehavior["FRIENDSHIP_LINK"]["behavior"] {
			if item.ModuleUid != "" {
				pageList[k].Content = contentMap[item.ModuleUid]
			} else {
				pageList[k].Content = contentMap[item.OtherData]
			}
		} else {
			pageList[k].Content = item.OtherData
		}
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    webVisitVO.PageSize,
		Current: webVisitVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}
