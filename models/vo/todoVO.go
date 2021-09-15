package vo

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/26 9:05 上午
 * @version 1.0
 */

type TodoVO struct {
	Keyword     string `json:"keyword"`
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	Uid         string `json:"uid"`
	Status      int    `json:"status"`
	Text        string `json:"text"`
	Done        bool   `json:"done"`
}
