//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"mogu-go-v2/controllers/admin"
	"mogu-go-v2/controllers/web"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/9 10:30 上午
 * @version 1.0
 */

func init() {
	// 请求前缀
	prefix := "/mogu-web"

	index := beego.NewNamespace(prefix + "/index",
		beego.NSRouter("/getWebConfig", &web.IndexRestApi{}, "get:GetWebConfig"),
		beego.NSRouter("/getBlogByLevel", &web.IndexRestApi{}, "get:GetBlogByLevel"),
		beego.NSRouter("/getNewBlog", &web.IndexRestApi{}, "get:GetNewBlog"),
		beego.NSRouter("/recorderVisitPage", &web.IndexRestApi{}, "get:RecorderVisitPage"),
		beego.NSRouter("/getHotBlog", &web.IndexRestApi{}, "get:GetHotBlog"),
		beego.NSRouter("/getHotTag", &web.IndexRestApi{}, "get:GetHotTag"),
		beego.NSRouter("/getLink", &web.IndexRestApi{}, "get:GetLink"),
		beego.NSRouter("/addLinkCount", &web.IndexRestApi{}, "get:AddLinkCount"),
		beego.NSRouter("/getWebNavbar", &web.IndexRestApi{}, "get:GetWebNavbar"),
	)
	about := beego.NewNamespace(prefix + "/about",
		beego.NSRouter("/getMe", &web.AboutRestApi{}, "get:GetMe"),
	)
	comment := beego.NewNamespace(prefix + "/web/comment",
		beego.NSRouter("/getList", &web.CommentRestApi{}, "post:GetList"),
		beego.NSRouter("/getListByUser", &web.CommentRestApi{}, "post:GetListByUser"),
		beego.NSRouter("/getPraiseListByUser", &web.CommentRestApi{}, "post:GetPraiseListByUser"),
		beego.NSRouter("/add", &web.CommentRestApi{}, "post:Add"),
		beego.NSRouter("/delete", &web.CommentRestApi{}, "post:DeleteBatch"),
		beego.NSRouter("/report", &web.CommentRestApi{}, "post:Report"),
		beego.NSRouter("/closeEmailNotification/:userUid", &web.CommentRestApi{}, "get:CloseEmailNotification"),
		beego.NSRouter("/getUserReceiveCommentCount", &web.CommentRestApi{}, "get:GetUserReceiveCommentCount"),
		beego.NSRouter("/readUserReceiveCommentCount", &web.CommentRestApi{}, "post:ReadUserReceiveCommentCount"),
	)
	sort := beego.NewNamespace(prefix + "/sort",
		beego.NSRouter("/getSortList", &web.SortRestApi{}, "get:GetSortList"),
		beego.NSRouter("/getArticleByMonth", &web.SortRestApi{}, "get:GetArticleByMonth"),
	)
	classify := beego.NewNamespace(prefix + "/classify",
		beego.NSRouter("/getBlogSortList", &web.ClassifyRestApi{}, "get:GetBlogSortList"),
		beego.NSRouter("/getArticleByBlogSortUid", &web.ClassifyRestApi{}, "get:GetArticleByBlogSortUid"),
	)
	tag := beego.NewNamespace(prefix + "/tag",
		beego.NSRouter("/getTagList", &web.TagRestApi{}, "get:GetTagList"),
		beego.NSRouter("/getArticleByTagUid", &web.TagRestApi{}, "get:GetArticleByTagUid"),
	)
	search := beego.NewNamespace(prefix + "/search",
		beego.NSRouter("/sqlSearchBlog", &web.SearchRestApi{}, "get:SqlSearchBlog"),
		beego.NSRouter("/searchBlogByTag", &web.SearchRestApi{}, "get:SearchBlogByTag"),
		beego.NSRouter("/searchBlogByAuthor", &web.SearchRestApi{}, "get:SearchBlogByAuthor"),
		beego.NSRouter("/searchBlogBySort", &web.SearchRestApi{}, "get:SearchBlogBySort"),
	)
	login := beego.NewNamespace(prefix + "/login",
		beego.NSRouter("/login", &web.LoginRestAPI{}, "post:Login"),
		beego.NSRouter("/register", &web.LoginRestAPI{}, "post:Register"),
		beego.NSRouter("/activeUser/:token", &web.LoginRestAPI{}, "get:ActiveUser"),
	)
	oauth := beego.NewNamespace(prefix + "/oauth",
		beego.NSRouter("/verify/:accessToken", &web.AuthRestApi{}, "get:Verify"),
		beego.NSRouter("/getFeedbackList", &web.AuthRestApi{}, "get:GetFeedbackList"),
		beego.NSRouter("/delete/:accessToken", &web.AuthRestApi{}, "post:Delete"),
		beego.NSRouter("/editUser", &web.AuthRestApi{}, "post:EditUser"),
		beego.NSRouter("/bindUserEmail/:token", &web.AuthRestApi{}, "get:BindUserEmail"),
		beego.NSRouter("/updateUserPwd", &web.AuthRestApi{}, "post:UpdateUserPwd"),
		beego.NSRouter("/addFeedback", &web.AuthRestApi{}, "post:AddFeedback"),
		beego.NSRouter("/replyBlogLink", &web.AuthRestApi{}, "post:ReplyBlogLink"),
		beego.NSRouter("/render", &web.AuthRestApi{}, "post:RenderOauth"),
	)
	content := beego.NewNamespace(prefix + "/content",
		beego.NSRouter("/getBlogByUid", &web.BlogContentRestApi{}, "get:GetBlogByUid"),
		beego.NSRouter("/getSameBlogByBlogUid", &web.BlogContentRestApi{}, "get:GetSameBlogByBlogUid"),
		beego.NSRouter("/praiseBlogByUid", &web.BlogContentRestApi{}, "get:PraiseBlogByUid"),
	)

	sysDictData := beego.NewNamespace(prefix + "/sysDictData",
		beego.NSRouter("/getListByDictTypeList", &admin.SysDictDataRestApi{}, "post:GetListByDictTypeList"),
	)

	subject := beego.NewNamespace(prefix + "/subject",
		beego.NSRouter("/getList", &admin.SubjectRestApi{}, "post:GetList"),
		beego.NSRouter("/getItemList", &admin.SubjectItemRestApi{}, "post:GetList"),
	)

	beego.AddNamespace(index, about, comment, sort, classify, tag, search, login, oauth, content, sysDictData, subject)
}