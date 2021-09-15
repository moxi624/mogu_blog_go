package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/26 9:20 上午
 * @version 1.0
 */
type RoleVO struct {
	Keyword          string `json:"keyword"`
	CurrentPage      int    `json:"currentPage"`
	PageSize         int    `json:"pageSize"`
	Uid              string `json:"uid"`
	Status           int    `json:"status"`
	RoleName         string `json:"roleName"`
	Summary          string `json:"summary"`
	CategoryMenuUids string `json:"categoryMenuUids"`
}
