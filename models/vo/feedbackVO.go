package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/22 11:19 上午
 * @version 1.0
 */

type FeedbackVO struct {
	Keyword        string `json:"keyword"`
	CurrentPage    int    `json:"currentPage"`
	PageSize       int    `json:"pageSize"`
	Uid            string `json:"uid"`
	Status         int    `json:"status"`
	UserUid        string `json:"userUid"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Reply          string `json:"reply"`
	FeedbackStatus int    `json:"feedbackStatus"`
}
