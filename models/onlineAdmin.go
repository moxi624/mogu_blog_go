package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/17 1:48 下午
 * @version 1.0
 */
type OnlineAdmin struct {
	/**
	 * 会话编号
	 */
	TokenId string `json:"tokenId"`

	/**
	 * 用户Token
	 */
	Token string `json:"token"`

	/**
	 * 管理员的UID
	 */
	AdminUid string `json:"adminUid"`

	/**
	 * 用户名称
	 */
	UserName string `json:"userName"`

	/**
	 * 登录IP地址
	 */
	Ipaddr string `json:"ipaddr"`

	/**
	 * 登录地址
	 */
	LoginLocation string `json:"loginLocation"`

	/**
	 * 浏览器类型
	 */
	Browser string `json:"browser"`

	/**
	 * 操作系统
	 */
	Os string `json:"os"`

	/**
	 * 角色名称
	 */
	RoleName string `json:"roleName"`

	/**
	 * 登录时间
	 */
	LoginTime string `json:"loginTime"`

	/**
	 * 过期时间
	 */
	ExpireTime string `json:"expireTime"`
}
