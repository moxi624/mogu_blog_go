package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/11 9:46 上午
 * @version 1.0
 */
type BlogSortVO struct {
	Keyword     string `json:"keyword"`
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	Uid         string `json:"uid"`
	Status      int    `json:"status"`
	SortName    string `json:"sortName"`
	Content     string `json:"content"`
	Sort        int    `json:"sort"`
}
