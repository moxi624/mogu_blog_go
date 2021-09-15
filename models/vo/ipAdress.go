package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/14 4:42 下午
 * @version 1.0
 */

type IpAddress struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ip      string `json:"ip"`
	Result  result `json:"result"`
}

type result struct {
	En_short string  `json:"en_short"`
	En_name  string  `json:"en_name"`
	Nation   string  `json:"nation"`
	Province string  `json:"province"`
	City     string  `json:"city"`
	District string  `json:"district"`
	Adcode   int     `json:"adcode"`
	Lat      float32 `json:"lat"`
	Lng      float32 `json:"lng"`
}
