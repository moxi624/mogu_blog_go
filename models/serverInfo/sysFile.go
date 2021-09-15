package serverInfo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/13 8:25 上午
 * @version 1.0
 */

type SysFile struct {
	DirName     string  `json:"dirName"`
	SysTypeName string  `json:"sysTypeName"`
	TypeName    string  `json:"typeName"`
	Total       string  `json:"total"`
	Free        string  `json:"free"`
	Used        string  `json:"used"`
	Usage       float64 `json:"usage"`
}
