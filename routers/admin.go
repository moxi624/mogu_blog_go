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

	//登录管理
	login := web.NewNamespace(prefix + "/auth",
		web.NSRouter("/login", &admin.LoginRestAPI{}, "post:Login"),
		web.NSRouter("/info", &admin.LoginRestAPI{}, "get:Info"),
		web.NSRouter("/getMenu", &admin.LoginRestAPI{}, "get:GetMenu"),
		web.NSRouter("/logout", &admin.LoginRestAPI{}, "post:Logout"),
		web.NSRouter("/getWebSiteName", &admin.LoginRestAPI{}, "get:GetWebSiteName"),
	)

	// 首页相关
	index := web.NewNamespace(prefix + "/index",
		web.NSRouter("/getVisitByWeek", &admin.IndexRestApi{}, "get:GetVisitByWeek"),
		web.NSRouter("/init", &admin.IndexRestApi{}, "get:InitIndex"),
		web.NSRouter("/getBlogCountByTag", &admin.IndexRestApi{}, "get:GetBlogCountByTag"),
		web.NSRouter("/getBlogCountByBlogSort", &admin.IndexRestApi{}, "get:GetBlogCountByBlogSort"),
		web.NSRouter("/getBlogContributeCount", &admin.IndexRestApi{}, "get:GetBlogContributeCount"),
	)

	// 数据字典相关
	sysDictData := web.NewNamespace(prefix + "/sysDictData",
		web.NSRouter("/getListByDictType", &admin.SysDictDataRestApi{}, "post:GetListByDictType"),
		web.NSRouter("/getListByDictTypeList", &admin.SysDictDataRestApi{}, "post:GetListByDictTypeList"),
		web.NSRouter("/getList", &admin.SysDictDataRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.SysDictDataRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.SysDictDataRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.SysDictDataRestApi{}, "post:Add"),
	)

	// 系统相关
	system := web.NewNamespace(prefix + "/system",
		web.NSRouter("/getMe", &admin.SystemRestApi{}, "get:GetMe"),
		web.NSRouter("/editMe", &admin.SystemRestApi{}, "post:EditMe"),
		web.NSRouter("/changePwd", &admin.SystemRestApi{}, "post:ChangePwd"),
	)

	// 系统配置相关
	systemConfig := web.NewNamespace(prefix + "/systemConfig",
		web.NSRouter("/getSystemConfig", &admin.SystemConfigRestApi{}, "get:GetSystemConfig"),
		web.NSRouter("/editSystemConfig", &admin.SystemConfigRestApi{}, "post:EditSystemConfig"),
		web.NSRouter("/cleanRedisByKey", &admin.SystemConfigRestApi{}, "post:CleanRedisByKey"),
	)

	// 代办事项相关
	todo := web.NewNamespace(prefix + "/todo",
		web.NSRouter("/getList", &admin.TodoRestApi{}, "post:GetList"),
		web.NSRouter("/add", &admin.TodoRestApi{}, "post:Add"),
		web.NSRouter("/edit", &admin.TodoRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.TodoRestApi{}, "post:Delete"),
	)

	// 访问相关
	webVisit := web.NewNamespace(prefix + "/webVisit",
		web.NSRouter("/getList", &admin.WebVisitRestApi{}, "post:GetList"),
	)

	// 用户相关
	user := web.NewNamespace(prefix + "/user",
		web.NSRouter("/getList", &admin.UserRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.UserRestApi{}, "post:Edit"),
		web.NSRouter("/resetUserPassword", &admin.UserRestApi{}, "post:ResetUserPassword"),
		web.NSRouter("/delete", &admin.UserRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.UserRestApi{}, "post:Add"),
	)

	// 评论相关
	comment := web.NewNamespace(prefix + "/comment",
		web.NSRouter("/getList", &admin.CommentRestApi{}, "post:GetList"),
		web.NSRouter("/delete", &admin.CommentRestApi{}, "post:Delete"),
		web.NSRouter("/deleteBatch", &admin.CommentRestApi{}, "post:DeleteBatch"),
	)

	// 标签相关
	tag := web.NewNamespace(prefix + "/tag",
		web.NSRouter("/getList", &admin.TagRestApi{}, "post:GetList"),
		web.NSRouter("/stick", &admin.TagRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.TagRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.TagRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.TagRestApi{}, "post:Add"),
		web.NSRouter("/tagSortByClickCount", &admin.TagRestApi{}, "post:TagSortByClickCount"),
		web.NSRouter("/tagSortByCite", &admin.TagRestApi{}, "post:TagSortByCite"),
	)

	// 博客分类相关
	blogSort := web.NewNamespace(prefix + "/blogSort",
		web.NSRouter("/getList", &admin.BlogSortRestApi{}, "post:GetList"),
		web.NSRouter("/stick", &admin.BlogSortRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.BlogSortRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.BlogSortRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.BlogSortRestApi{}, "post:Add"),
		web.NSRouter("/blogSortByClickCount", &admin.BlogSortRestApi{}, "post:BlogSortByClickCount"),
		web.NSRouter("/blogSortByCite", &admin.BlogSortRestApi{}, "post:BlogSortByCite"),
	)

	// 博客相关
	blog := web.NewNamespace(prefix + "/blog",
		web.NSRouter("/getList", &admin.BlogRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.BlogRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.BlogRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.BlogRestApi{}, "post:Add"),
		web.NSRouter("/uploadLocalBlog", &admin.BlogRestApi{}, "post:UploadLocalBlog"),
		web.NSRouter("/deleteBatch", &admin.BlogRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/editBatch", &admin.BlogRestApi{}, "post:EditBatch"),
	)

	// 管理员相关
	adminApi := web.NewNamespace(prefix + "/admin",
		web.NSRouter("/getOnlineAdminList", &admin.AdminRestApi{}, "post:GetOnlineAdminList"),
		web.NSRouter("/forceLogout", &admin.AdminRestApi{}, "post:ForceLogout"),
		web.NSRouter("/getList", &admin.AdminRestApi{}, "post:GetList"),
		web.NSRouter("/restPwd", &admin.AdminRestApi{}, "post:RestPwd"),
		web.NSRouter("/edit", &admin.AdminRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.AdminRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.AdminRestApi{}, "post:Add"),
	)

	// 监控相关
	server := web.NewNamespace(prefix + "/monitor",
		web.NSRouter("/getServerInfo", &admin.ServerMonitorRestApi{}, "get:GetInfo"),
	)

	// 资源分类相关
	resourceSort := web.NewNamespace(prefix + "/resourceSort",
		web.NSRouter("/getList", &admin.ResourceSortRestApi{}, "post:GetList"),
		web.NSRouter("/stick", &admin.ResourceSortRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.ResourceSortRestApi{}, "post:Edit"),
		web.NSRouter("/add", &admin.ResourceSortRestApi{}, "post:Add"),
		web.NSRouter("/deleteBatch", &admin.ResourceSortRestApi{}, "post:DeleteBatch"),
	)

	// 图片分类相关
	pictureSort := web.NewNamespace(prefix + "/pictureSort",
		web.NSRouter("/getList", &admin.PictureSortRestApi{}, "post:GetList"),
		web.NSRouter("/stick", &admin.PictureSortRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.PictureSortRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.PictureSortRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.PictureSortRestApi{}, "post:Add"),
	)

	// 图片相关
	picture := web.NewNamespace(prefix + "/picture",
		web.NSRouter("/getList", &admin.PictureRestApi{}, "post:GetList"),
		web.NSRouter("/delete", &admin.PictureRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.PictureRestApi{}, "post:Add"),
		web.NSRouter("/setCover", &admin.PictureRestApi{}, "post:SetCover"),
	)

	// 学习视频相关
	studyVideo := web.NewNamespace(prefix + "/studyVideo",
		web.NSRouter("/getList", &admin.StudyVideoRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.StudyVideoRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.StudyVideoRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.StudyVideoRestApi{}, "post:Add"),
	)

	// 反馈相关
	feedback := web.NewNamespace(prefix + "/feedback",
		web.NSRouter("/getList", &admin.FeedbackRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.FeedbackRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.FeedbackRestApi{}, "post:DeleteBatch"),
	)

	// 角色相关
	role := web.NewNamespace(prefix + "/role",
		web.NSRouter("/getList", &admin.RoleRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.RoleRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.RoleRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.RoleRestApi{}, "post:Add"),
	)

	// 菜单相关
	categoryMenu := web.NewNamespace(prefix + "/categoryMenu",
		web.NSRouter("/getAll", &admin.CategoryMenuRestApi{}, "get:GetAll"),
		web.NSRouter("/getButtonAll", &admin.CategoryMenuRestApi{}, "get:GetButtonAll"),
		web.NSRouter("/stick", &admin.CategoryMenuRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.CategoryMenuRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.CategoryMenuRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.CategoryMenuRestApi{}, "post:Add"),
	)

	// 网站配置相关
	webConfig := web.NewNamespace(prefix + "/webConfig",
		web.NSRouter("/getWebConfig", &admin.WebConfigRestApi{}, "get:GetWebConfig"),
		web.NSRouter("/editWebConfig", &admin.WebConfigRestApi{}, "post:EditWebConfig"),
	)

	// 字典类型相关
	sysDictType := web.NewNamespace(prefix + "/sysDictType",
		web.NSRouter("/getList", &admin.SysDictTypeRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.SysDictTypeRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.SysDictTypeRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.SysDictTypeRestApi{}, "post:Add"),
	)

	// 系统参数相关
	sysParams := web.NewNamespace(prefix + "/sysParams",
		web.NSRouter("/getList", &admin.SysParamsRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.SysParamsRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.SysParamsRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.SysParamsRestApi{}, "post:Add"),
	)

	// 链接相关
	link := web.NewNamespace(prefix + "/link",
		web.NSRouter("/getList", &admin.LinkRestApi{}, "post:GetList"),
		web.NSRouter("/stick", &admin.LinkRestApi{}, "post:Stick"),
		web.NSRouter("/edit", &admin.LinkRestApi{}, "post:Edit"),
		web.NSRouter("/delete", &admin.LinkRestApi{}, "post:Delete"),
		web.NSRouter("/add", &admin.LinkRestApi{}, "post:Add"),
	)

	// 专题相关
	subject := web.NewNamespace(prefix + "/subject",
		web.NSRouter("/getList", &admin.SubjectRestApi{}, "post:GetList"),
		web.NSRouter("/edit", &admin.SubjectRestApi{}, "post:Edit"),
		web.NSRouter("/deleteBatch", &admin.SubjectRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/add", &admin.SubjectRestApi{}, "post:Add"),
		web.NSRouter("/getItemList", &admin.SubjectItemRestApi{}, "post:GetList"),
	)

	// 专题列表
	subjectItem := web.NewNamespace(prefix + "/subjectItem",
		web.NSRouter("/add", &admin.SubjectItemRestApi{}, "post:Add"),
		web.NSRouter("/getList", &admin.SubjectItemRestApi{}, "post:GetList"),
		web.NSRouter("/deleteBatch", &admin.SubjectItemRestApi{}, "post:DeleteBatch"),
		web.NSRouter("/edit", &admin.SubjectItemRestApi{}, "post:Edit"),
		web.NSRouter("/sortByCreateTime", &admin.SubjectItemRestApi{}, "post:SortByCreateTime"),
	)

	// 日志相关
	log := web.NewNamespace(prefix + "/log",
		web.NSRouter("/getExceptionList", &admin.LogRestApi{}, "post:GetExceptionList"),
		web.NSRouter("/getLogList", &admin.LogRestApi{}, "post:GetLogList"),
	)

	// 导航栏相关
	webNavBar := web.NewNamespace(prefix + "/webNavbar",
		web.NSRouter("/getAllList", &admin.WebNavbarRestApi{}, "get:GetAllList"),
		web.NSRouter("/edit", &admin.WebNavbarRestApi{}, "post:Edit"),
		web.NSRouter("/add", &admin.WebNavbarRestApi{}, "post:Add"),
		web.NSRouter("/delete", &admin.WebNavbarRestApi{}, "post:Delete"),
	)

	// 注册全部路由
	web.AddNamespace(login, index, sysDictData, system, systemConfig, todo, webVisit, user, comment, tag, blogSort, blog,
		adminApi, server, resourceSort, pictureSort, picture, studyVideo, feedback, role, categoryMenu, webConfig, sysDictType,
		sysParams, link, subject, subjectItem, log, webNavBar)
}
