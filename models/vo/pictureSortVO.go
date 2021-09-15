package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/19 3:19 下午
 * @version 1.0
 */
type PictureSortVO struct {
	Uid         string `json:"uid"`
	Status      int    `json:"status"`
	Keyword     string `json:"keyword"`
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	ParentUid   string `json:"parentUid"`
	Name        string `json:"name"`
	FileUid     string `json:"fileUid"`
	Sort        int    `json:"sort"`
	IsShow      int8   `json:"isShow"`
}
