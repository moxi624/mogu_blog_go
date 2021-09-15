package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/21 8:09 上午
 * @version 1.0
 */

type StudyVideoVO struct {
	Uid            string `json:"uid"`
	Status         int    `json:"status"`
	Keyword        string `json:"keyword"`
	CurrentPage    int    `json:"currentPage"`
	PageSize       int    `json:"pageSize"`
	Name           string `json:"name"`
	Summary        string `json:"summary"`
	Content        string `json:"content"`
	BaiduPath      string `json:"baiduPath"`
	FileUid        string `json:"fileUid"`
	ResouceSortUid string `json:"resouceSortUid"`
}
