package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/27 3:51 下午
 * @version 1.0
 */

type CategoryMenuVO struct {
	Keyword           string `json:"keyword"`
	CurrentPage       int    `json:"currentPage"`
	PageSize          int    `json:"pageSize"`
	Uid               string `json:"uid"`
	Status            int    `json:"status"`
	Name              string `json:"Name"`
	MenuLevel         int    `json:"menuLevel"`
	MenuType          int    `json:"menuType"`
	Summary           string `json:"summary"`
	Icon              string `json:"icon"`
	ParentUid         string `json:"parentUid"`
	Url               string `json:"url"`
	Sort              int    `json:"sort"`
	IsShow            int    `json:"isShow"`
	IsJumpExternalUrl int    `json:"isJumpExternalUrl"`
}
