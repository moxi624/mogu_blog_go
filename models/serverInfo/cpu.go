package serverInfo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/12 1:20 下午
 * @version 1.0
 */

type Cpu struct {
	CpuNum int     `json:"cpuNum"`
	Total  float64 `json:"total"`
	Sys    float64 `json:"sys"`
	Used   float64 `json:"used"`
	Wait   float64 `json:"wait"`
	Free   float64 `json:"free"`
}
