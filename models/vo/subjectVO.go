package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/5 9:19 上午
 * @version 1.0
 */

type SubjectVO struct {
	Keyword      string      `json:"keyword"`
	CurrentPage  int         `json:"currentPage"`
	PageSize     int         `json:"pageSize"`
	Uid          string      `json:"uid"`
	Status       int         `json:"status"`
	SubjectName  string      `json:"subjectName"`
	Summary      string      `json:"summary"`
	FileUid      string      `json:"fileUid"`
	Sort         int         `json:"sort"`
	ClickCount   interface{} `json:"clickCount"`
	CollectCount interface{} `json:"collectCount"`
}
