package common

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/7 2:09 下午
 * @version 1.0
 */

type emu struct{}

/*func (emu) Ebehavior() [][]string {
	BLOG_TAG := []string{"点击标签", "blog_tag"}
	BLOG_SORT := []string{"点击博客排序", "blog_sort"}
	BLOG_CONTENT := []string{"点击博客", "blog_content"}
	FRIENDSHIP_LINK := []string{"点击友情链接", "friendship_link"}
	BLOG_SEARCH := []string{"点击搜索", "blog_search"}
	STUDY_VIDEO := []string{"点击学习视频", "study_video"}
	VISIT_PAGE := []string{"访问页面", "visit_page"}
	VISIT_CLASSIFY := []string{"点击博客分类", "visit_classify"}
	VISIT_SORT := []string{"点击归档", "visit_sort"}
	BLOG_AUTHOR := []string{"点击作者", "blog_author"}
	PUBLISH_COMMENT := []string{"发表评论", "publish_comment"}
	DELETE_COMMENT := []string{"删除评论", "delete_comment"}
	REPORT_COMMENT := []string{"举报评论", "report_comment"}
	VISIT_TAG := []string{"点击博客标签页面", "visit_tag"}
	return [][]string{BLOG_TAG,BLOG_SORT,BLOG_CONTENT,FRIENDSHIP_LINK,BLOG_SEARCH,STUDY_VIDEO,
		VISIT_PAGE,VISIT_CLASSIFY,VISIT_SORT,BLOG_AUTHOR,PUBLISH_COMMENT,DELETE_COMMENT,REPORT_COMMENT,VISIT_TAG}
}*/

func (emu) BehaviorEmu() map[string]map[string]string {
	eBehavior := map[string]map[string]string{
		"BLOG_TAG":        {"content": "点击标签", "behavior": "blog_tag"},
		"BLOG_SORT":       {"content": "点击博客排序", "behavior": "blog_sort"},
		"BLOG_CONTENT":    {"content": "点击博客", "behavior": "blog_content"},
		"FRIENDSHIP_LINK": {"content": "点击友情链接", "behavior": "friendship_link"},
		"BLOG_SEARCH":     {"content": "点击搜索", "behavior": "blog_search"},
		"STUDY_VIDEO":     {"content": "点击学习视频", "behavior": "study_video"},
		"VISIT_PAGE":      {"content": "访问页面", "behavior": "visit_page"},
		"VISIT_CLASSIFY":  {"content": "点击博客分类", "behavior": "visit_classify"},
		"VISIT_SORT":      {"content": "点击归档", "behavior": "visit_sort"},
		"BLOG_AUTHOR":     {"content": "点击作者", "behavior": "blog_author"},
		"PUBLISH_COMMENT": {"content": "发表评论", "behavior": "publish_comment"},
		"DELETE_COMMENT":  {"content": "删除评论", "behavior": "delete_comment"},
		"REPORT_COMMENT":  {"content": "举报评论", "behavior": "report_comment"},
		"VISIT_TAG":       {"content": "点击博客标签页面", "behavior": "visit_tag"},
	}
	return eBehavior
}

func (emu) CommentSourceEmu() map[string]map[string]string {
	eCommentSource := map[string]map[string]string{
		"ABOUT":         {"code": "ABOUT", "name": "关于我"},
		"BLOG_INFO":     {"code": "BLOG_INFO", "name": "博客详情"},
		"MESSAGE_BOARD": {"code": "MESSAGE_BOARD", "name": "留言板"},
	}
	return eCommentSource
}

var Emu = &emu{}
