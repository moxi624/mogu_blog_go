package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/7 2:42 下午
 * @version 1.0
 */

type WebVisitVO struct {
	Keyword         string `json:"keyword"`
	CurrentPage     int    `json:"currentPage"`
	PageSize        int    `json:"pageSize"`
	Uid             string `json:"uid"`
	Status          int    `json:"status"`
	UserUid         string `json:"userUid"`
	Ip              string `json:"ip"`
	Os              string `json:"os"`
	Browser         string `json:"browser"`
	Behavior        string `json:"behavior"`
	ModuleUid       string `json:"moduleUid"`
	OtherData       string `json:"otherData"`
	StartTime       string `json:"startTime"`
	Content         string `gorm:"-" json:"content"`
	BehaviorContent string `gorm:"-" json:"behaviorContent"`
}
