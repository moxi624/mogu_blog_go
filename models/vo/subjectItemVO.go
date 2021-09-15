package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/5 10:22 上午
 * @version 1.0
 */

type SubjectItemVO struct {
	Uid         string `json:"uid"`
	Status      int    `json:"status"`
	Keyword     string `json:"keyword"`
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	SubjectUid  string `json:"subjectUid"`
	BlogUid     string `json:"blogUid"`
	Sort        int    `json:"sort"`
}
