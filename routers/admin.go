package routers

import (
	"github.com/beego/beego/v2/server/web"
	"mogu-go-v2/controllers/admin"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/31 2:33 下午
 * @version 1.0
 */
func init() {
	// 请求前缀
	prefix := "/mogu-admin"

	login := web.NewNamespace(prefix + "/auth",
		//登录管理
		web.NSRouter("/login", &admin.LoginRestAPI{}, "post:Login"),
		web.NSRouter("/info", &admin.LoginRestAPI{}, "get:Info"),
		web.NSRouter("/getMenu", &admin.LoginRestAPI{}, "get:GetMenu"),
		web.NSRouter("/logout", &admin.LoginRestAPI{}, "post:Logout"),
		web.NSRouter("/getWebSiteName", &admin.LoginRestAPI{}, "get:GetWebSiteName"),
	)
	index := web.NewNamespace(prefix + "/index",
		web.NSRouter("/getVisitByWeek", &admin.IndexRestApi{}, "get:GetVisitByWeek"),
		web.NSRouter("/init", &admin.IndexRestApi{}, "get:InitIndex"),
		web.NSRouter("/getBlogCountByTag", &admin.IndexRestApi{}, "get:GetBlogCountByTag"),
		web.NSRouter("/getBlogCountByBlogSort", &admin.IndexRestApi{}, "get:GetBlogCountByBlogSort"),
		web.NSRouter("/getBlogContributeCount", &admin.IndexRestApi{}, "get:GetBlogContributeCount"),
	)
	sysDictData := web.NewNamespace(prefix + "/sysDictData",
		web.NSRouter("/getListByDictType", &admin.SysDictDataRestApi{}, "post:GetListByDictType"),
		web.NSRouter("/getListByDictTypeList", &admin.SysDictDataRestApi{}, "post:GetListByDictTypeList"),
		web.NSRouter("/getList", &admin.SysDictDataRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.SysDictDataRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.SysDictDataRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.SysDictDataRestApi{}, "post:Add"),
	)
	system := web.NewNamespace(prefix + "/system",
		web.NSRouter("/getMe", &admin.SystemRestApi{}, "get:GetMe"),
		web.NSRouter("/editMe", &admin.SystemRestApi{}, "post:EditMe"),
		web.NSRouter("/changePwd", &admin.SystemRestApi{}, "post:ChangePwd"),
	)
	systemConfig := web.NewNamespace(prefix + "/systemConfig",
		web.NSRouter("/getSystemConfig", &admin.SystemConfigRestApi{}, "get:GetSystemConfig"),
		web.NSRouter("/editSystemConfig", &admin.SystemConfigRestApi{}, "post:EditSystemConfig"),
		web.NSRouter("/cleanRedisByKey", &admin.SystemConfigRestApi{}, "post:CleanRedisByKey"),
	)
	todo := web.NewNamespace(prefix + "/todo",
		web.NSRouter("/getList", &admin.TodoRestApi{}, "post:GetList"),
		web.NSRouter("/add", &admin.TodoRestApi{}, "post:Add"),
		web.NSRouter("/edit", &admin.TodoRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.TodoRestApi{}, "post:Delete"),
	)
	webVisit := web.NewNamespace(prefix + "/webVisit",
		web.NSRouter("/getList", &admin.WebVisitRestApi{}, "post:GetList"),
	)
	user := web.NewNamespace(prefix + "/user",
		web.NSRouter("/getList", &admin.UserRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.UserRestApi{}, "post:Edit"),
		web.NSRouter("/resetUserPassword", &admin.UserRestApi{}, "post:ResetUserPassword"),
		web.NSRouter("/delete", &admin.UserRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.UserRestApi{}, "post:Add"),
	)
	comment := web.NewNamespace(prefix + "/comment",
		web.NSRouter("/getList", &admin.CommentRestApi{}, "post:GetList"),
		web.NSRouter("/delete", &admin.CommentRestApi{}, "post:Delete"),
		web.NSRouter("/deleteBatch", &admin.CommentRestApi{}, "post:DeleteBatch"),
	)
	tag := web.NewNamespace(prefix + "/tag",
		web.NSRouter("/getList", &admin.TagRestApi{}, "post:GetList"),
		web.NSRouter("/stick", &admin.TagRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.TagRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.TagRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.TagRestApi{}, "post:Add"),
		web.NSRouter("/tagSortByClickCount", &admin.TagRestApi{}, "post:TagSortByClickCount"),
		web.NSRouter("/tagSortByCite", &admin.TagRestApi{}, "post:TagSortByCite"),
	)
	blogSort := web.NewNamespace(prefix + "/blogSort",
		web.NSRouter("/getList", &admin.BlogSortRestApi{}, "post:GetList"),
		web.NSRouter("/stick", &admin.BlogSortRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.BlogSortRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.BlogSortRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.BlogSortRestApi{}, "post:Add"),
		web.NSRouter("/blogSortByClickCount", &admin.BlogSortRestApi{}, "post:BlogSortByClickCount"),
		web.NSRouter("/blogSortByCite", &admin.BlogSortRestApi{}, "post:BlogSortByCite"),
	)
	blog := web.NewNamespace(prefix + "/blog",
		web.NSRouter("/getList", &admin.BlogRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.BlogRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.BlogRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.BlogRestApi{}, "post:Add"),
		web.NSRouter("/uploadLocalBlog", &admin.BlogRestApi{}, "post:UploadLocalBlog"),
		web.NSRouter("/deleteBatch", &admin.BlogRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/editBatch", &admin.BlogRestApi{}, "post:EditBatch"),
	)
	adminApi := web.NewNamespace(prefix + "/admin",
		web.NSRouter("/getOnlineAdminList", &admin.AdminRestApi{}, "post:GetOnlineAdminList"),
		web.NSRouter("/forceLogout", &admin.AdminRestApi{}, "post:ForceLogout"),
		web.NSRouter("/getList", &admin.AdminRestApi{}, "post:GetList"),
		web.NSRouter("/restPwd", &admin.AdminRestApi{}, "post:RestPwd"),
		web.NSRouter("/edit", &admin.AdminRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.AdminRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.AdminRestApi{}, "post:Add"),
	)
	server := web.NewNamespace(prefix + "/monitor",
		web.NSRouter("/getServerInfo", &admin.ServerMonitorRestApi{}, "get:GetInfo"),
	)
	resourceSort := web.NewNamespace(prefix + "/resourceSort",
		web.NSRouter("/getList", &admin.ResourceSortRestApi{}, "post:GetList"),
		web.NSRouter("/stick", &admin.ResourceSortRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.ResourceSortRestApi{}, "post:Edit"),
		web.NSRouter("/add", &admin.ResourceSortRestApi{}, "post:Add"),
		web.NSRouter("/deleteBatch", &admin.ResourceSortRestApi{}, "post:DeleteBatch"),
	)
	pictureSort := web.NewNamespace(prefix + "/pictureSort",
		web.NSRouter("/getList", &admin.PictureSortRestApi{}, "post:GetList"),
		web.NSRouter("/stick", &admin.PictureSortRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.PictureSortRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.PictureSortRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.PictureSortRestApi{}, "post:Add"),
	)
	picture := web.NewNamespace(prefix + "/picture",
		web.NSRouter("/getList", &admin.PictureRestApi{}, "post:GetList"),
		web.NSRouter("/delete", &admin.PictureRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.PictureRestApi{}, "post:Add"),
		web.NSRouter("/setCover", &admin.PictureRestApi{}, "post:SetCover"),
	)
	studyVideo := web.NewNamespace(prefix + "/studyVideo",
		web.NSRouter("/getList", &admin.StudyVideoRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.StudyVideoRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.StudyVideoRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.StudyVideoRestApi{}, "post:Add"),
	)
	feedback := web.NewNamespace(prefix + "/feedback",
		web.NSRouter("/getList", &admin.FeedbackRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.FeedbackRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.FeedbackRestApi{}, "post:DeleteBatch"),
	)
	role := web.NewNamespace(prefix + "/role",
		web.NSRouter("/getList", &admin.RoleRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.RoleRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.RoleRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.RoleRestApi{}, "post:Add"),
	)
	categoryMenu := web.NewNamespace(prefix + "/categoryMenu",
		web.NSRouter("/getAll", &admin.CategoryMenuRestApi{}, "get:GetAll"),
		web.NSRouter("/getButtonAll", &admin.CategoryMenuRestApi{}, "get:GetButtonAll"),
		web.NSRouter("/stick", &admin.CategoryMenuRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.CategoryMenuRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.CategoryMenuRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.CategoryMenuRestApi{}, "post:Add"),
	)
	webConfig := web.NewNamespace(prefix + "/webConfig",
		web.NSRouter("/getWebConfig", &admin.WebConfigRestApi{}, "get:GetWebConfig"),
		web.NSRouter("/editWebConfig", &admin.WebConfigRestApi{}, "post:EditWebConfig"),
	)
	sysDictType := web.NewNamespace(prefix + "/sysDictType",
		web.NSRouter("/getList", &admin.SysDictTypeRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.SysDictTypeRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.SysDictTypeRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.SysDictTypeRestApi{}, "post:Add"),
	)
	sysParams := web.NewNamespace(prefix + "/sysParams",
		web.NSRouter("/getList", &admin.SysParamsRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.SysParamsRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.SysParamsRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.SysParamsRestApi{}, "post:Add"),
	)
	link := web.NewNamespace(prefix + "/link",
		web.NSRouter("/getList", &admin.LinkRestApi{}, "post:GetList"),
		web.NSRouter("/stick", &admin.LinkRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.LinkRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.LinkRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.LinkRestApi{}, "post:Add"),
	)
	subject := web.NewNamespace(prefix + "/subject",
		web.NSRouter("/getList", &admin.SubjectRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.SubjectRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.SubjectRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.SubjectRestApi{}, "post:Add"),
		web.NSRouter("/getItemList", &admin.SubjectItemRestApi{}, "post:GetList"),
	)
	subjectItem := web.NewNamespace(prefix + "/subjectItem",
		web.NSRouter("/add", &admin.SubjectItemRestApi{}, "post:Add"),
		web.NSRouter("/getList", &admin.SubjectItemRestApi{}, "post:GetList"),
		web.NSRouter("/deleteBatch", &admin.SubjectItemRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/edit", &admin.SubjectItemRestApi{}, "post:Edit"),
		web.NSRouter("/sortByCreateTime", &admin.SubjectItemRestApi{}, "post:SortByCreateTime"),
	)
	log := web.NewNamespace(prefix + "/log",
		web.NSRouter("/getExceptionList", &admin.LogRestApi{}, "post:GetExceptionList"),
		web.NSRouter("/getLogList", &admin.LogRestApi{}, "post:GetLogList"),
	)
	webNavBar := web.NewNamespace(prefix + "/webNavbar",
		web.NSRouter("/getAllList", &admin.WebNavbarRestApi{}, "get:GetAllList"),
		web.NSRouter("/edit", &admin.WebNavbarRestApi{}, "post:Edit"),
		web.NSRouter("/add", &admin.WebNavbarRestApi{}, "post:Add"),
		web.NSRouter("/delete", &admin.WebNavbarRestApi{}, "post:Delete"),
	)
	web.AddNamespace(login, index, sysDictData, system, systemConfig, todo, webVisit, user, comment, tag, blogSort, blog,
		adminApi, server, resourceSort, pictureSort, picture, studyVideo, feedback, role, categoryMenu, webConfig, sysDictType,
		sysParams, link, subject, subjectItem, log, webNavBar)
}
