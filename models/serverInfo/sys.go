package serverInfo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/12 4:00 下午
 * @version 1.0
 */

type Sys struct {
	ComputerName string `json:"computerName"`
	ComputerIp   string `json:"computerIp"`
	UserDir      string `json:"userDir"`
	OsName       string `json:"osName"`
	OsArch       string `json:"osArch"`
}
