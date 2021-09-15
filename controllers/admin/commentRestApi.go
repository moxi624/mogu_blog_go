package admin

import (
	"encoding/json"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"strconv"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/9 2:43 下午
 * @version 1.0
 */
type CommentRestApi struct {
	base.BaseController
}

func (c *CommentRestApi) GetList() {
	var commentVO vo.CommentVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &commentVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		panic(err)
	}
	base.L.Print("获取评论列表")
	where := "status=?"
	if strings.TrimSpace(commentVO.Keyword) != "" {
		where += " and content like \"%" + strings.TrimSpace(commentVO.Keyword) + "%\""
	}
	if t := common.InterfaceToInt(commentVO.Type); t != 0 {
		where += " and type=" + strconv.Itoa(t)
	}
	if commentVO.Source != "" && commentVO.Source != "all" {
		where += " and source=\"" + commentVO.Source + "\""
	}
	if commentVO.UserName != "" {
		userName := commentVO.UserName
		var list []models.User
		common.DB.Where("nick_name like ? and status=?", userName, 1).Find(&list)
		if len(list) > 0 {
			var userUid []string
			for _, item := range list {
				userUid = append(userUid, item.Uid)
			}
			where += " and user_uid in ('" + strings.Join(userUid, "','") + "')"
		} else {
			where += " and user_uid = 'uid00000000000000000000000000000'"
		}
	}
	var pageList []models.Comment
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.Comment{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((commentVO.CurrentPage - 1) * commentVO.PageSize).Limit(commentVO.PageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var userUidSet []string
	var blogUidSet []string
	for _, item := range pageList {
		if item.UserUid != "" {
			userUidSet = append(userUidSet, item.UserUid)
		}
		if item.ToUserUid != "" {
			userUidSet = append(userUidSet, item.ToUserUid)
		}
		if item.BlogUid != "" {
			blogUidSet = append(blogUidSet, item.BlogUid)
		}
	}
	var userCollection []models.User
	blogMap := map[string]models.Blog{}
	c.Wg.Add(1)
	go func() {
		var blogList []models.Blog
		if len(blogUidSet) > 0 {
			common.DB.Find(&blogList, common.RemoveRepByMap(blogUidSet))
		}
		for _, item := range blogList {
			item.Content = ""
			blogMap[item.Uid] = item
		}
		c.Wg.Done()
	}()
	c.Wg.Add(1)
	go func() {
		if len(userUidSet) > 0 {
			common.DB.Find(&userCollection, common.RemoveRepByMap(userUidSet))
		}
		c.Wg.Done()
	}()
	c.Wg.Wait()
	var s []string
	for _, item := range userCollection {
		if item.Avatar != "" {
			s = append(s, item.Avatar)
		}
	}
	fileUids := strings.Join(s, ",")
	pictureList := map[string]interface{}{}
	if fileUids != "" {
		pictureList = service.FileService.GetPicture(fileUids, ",")
	}
	pictureMap := map[string]string{}
	picList := common.WebUtil.GetPictureMap(pictureList)
	for _, item := range picList {
		pictureMap[item["uid"].(string)] = item["url"].(string)
	}
	userMap := map[string]models.User{}
	for _, item := range userCollection {
		if pictureMap[item.Avatar] != "" {
			item.PhotoUrl = pictureMap[item.Avatar]
		}
		userMap[item.Uid] = item
	}
	eCommentSource := common.Emu.CommentSourceEmu()
	for k, item := range pageList {
		commentSource := eCommentSource[item.Source]
		item.SourceName = commentSource["name"]
		if item.UserUid != "" {
			item.User = userMap[item.ToUserUid]
		}
		if item.ToUserUid != "" {
			item.ToUser = userMap[item.ToUserUid]
		}
		if item.BlogUid != "" {
			item.Blog = blogMap[item.BlogUid]
		}
		pageList[k] = item
	}
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    commentVO.PageSize,
		Current: commentVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *CommentRestApi) Delete() {
	base.L.Print("获取评论列表")
	var commentVO vo.CommentVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &commentVO)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		return
	}
	var comment models.Comment
	common.DB.Where("uid=?", commentVO.Uid).Find(&comment)
	common.DB.Model(&comment).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}

func (c *CommentRestApi) DeleteBatch() {
	var commentVOList []vo.CommentVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &commentVOList)
	if err != nil {
		c.ThrowError("0", "传入参数有误！")
		panic(err)
	}
	if len(commentVOList) <= 0 {
		c.ErrorWithMessage("传入参数有误！")
		return
	}
	var uids []string
	for _, item := range commentVOList {
		uids = append(uids, item.Uid)
	}
	var commentList []models.Comment
	common.DB.Find(&commentList, uids)
	common.DB.Model(&commentList).Select("status").Update("status", 0)
	c.SuccessWithMessage("删除成功")
}
