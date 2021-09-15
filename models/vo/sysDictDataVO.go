package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/29 8:51 上午
 * @version 1.0
 */

type SysDictDataVO struct {
	Keyword     string `json:"keyword"`
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	Uid         string `json:"uid"`
	Status      int    `json:"status"`
	Oid         int    `json:"oid"`
	DictLabel   string `json:"dictLabel"`
	DictValue   string `json:"dictValue"`
	DictTypeUid string `json:"dictTypeUid"`
	CssClass    string `json:"cssClass"`
	ListClass   string `json:"listClass"`
	IsDefault   int    `json:"isDefault"`
	IsPublish   string `json:"isPublish"`
	Remark      string `json:"remark"`
	Sort        int    `json:"sort"`
}
