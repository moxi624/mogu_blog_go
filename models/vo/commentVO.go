package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/9 2:49 下午
 * @version 1.0
 */
type CommentVO struct {
	Keyword     string      `json:"keyword"`
	CurrentPage int         `json:"currentPage"`
	PageSize    int         `json:"pageSize"`
	Uid         string      `json:"uid"`
	Status      int         `json:"status"`
	UserUid     string      `json:"userUid"`
	ToUid       string      `json:"toUid"`
	ToUserUid   string      `json:"toUserUid"`
	UserName    string      `json:"userName"`
	Type        interface{} `json:"type"`
	Content     string      `json:"content"`
	BlogUid     string      `json:"blogUid"`
	Source      string      `json:"source"`
}
