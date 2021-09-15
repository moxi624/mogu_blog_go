package serverInfo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/12 3:05 下午
 * @version 1.0
 */

type Mem struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
	Free  float64 `json:"free"`
	Usage float64 `json:"usage"`
}
