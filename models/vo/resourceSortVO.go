package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/19 2:16 下午
 * @version 1.0
 */

type ResourceSortVO struct {
	Uid         string `json:"uid"`
	Status      int    `json:"status"`
	Keyword     string `json:"keyword"`
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	SortName    string `json:"sortName"`
	Content     string `json:"content"`
	FileUid     string `json:"fileUid"`
	Sort        int    `json:"sort"`
}
