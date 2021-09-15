package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/29 8:36 上午
 * @version 1.0
 */

type SysDictTypeVO struct {
	Keyword     string `json:"keyword"`
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	Uid         string `json:"uid"`
	Status      int    `json:"status"`
	Oid         int    `json:"oid"`
	DictName    string `json:"dictName"`
	DictType    string `json:"dictType"`
	IsPublish   string `json:"isPublish"`
	Remark      string `json:"remark"`
	Sort        int    `json:"sort"`
}
