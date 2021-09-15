package vo

import "mogu-go-v2/models"

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/11 9:59 上午
 * @version 1.0
 */
type BlogVO struct {
	Keyword           string          `json:"keyword"`
	CurrentPage       int             `json:"currentPage"`
	PageSize          int             `json:"pageSize"`
	Uid               string          `json:"uid"`
	Status            int             `json:"status"`
	Title             string          `json:"title"`
	Summary           string          `json:"summary"`
	TagUid            string          `json:"tagUid"`
	BlogSortUid       string          `json:"blogSortUid"`
	FileUid           string          `json:"fileUid"`
	AdminUid          string          `json:"adminUid"`
	IsPublish         string          `json:"isPublish"`
	IsOriginal        string          `json:"isOriginal"`
	Author            string          `json:"author"`
	ArticlesPart      string          `json:"articlesPart"`
	Level             int             `json:"level"`
	Type              interface{}     `json:"type"`
	OutsideLink       string          `json:"outsideLink"`
	Content           string          `json:"content"`
	TagList           []models.Tag    `json:"tagList"`
	PhotoList         []string        `json:"photoList"`
	BlogSort          models.BlogSort `json:"blogSort"`
	ParseCount        int             `json:"parseCount"`
	Copyright         string          `json:"copyright"`
	LevelKeyword      interface{}     `json:"levelKeyword"`
	UserSort          int             `json:"userSort"`
	Sort              int             `json:"sort"`
	OpenComment       interface{}     `json:"openComment"`
	OrderByDescColumn string          `json:"orderByDescColumn"`
	OrderByAscColumn  string          `json:"orderByAscColumn"`
}
