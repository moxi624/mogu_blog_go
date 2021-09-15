package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/28 3:05 下午
 * @version 1.0
 */

type WebConfigVO struct {
	Keyword              string   `json:"keyword"`
	CurrentPage          int      `json:"currentPage"`
	PageSize             int      `json:"pageSize"`
	Uid                  string   `json:"uid"`
	Status               int      `json:"status"`
	Logo                 string   `json:"logo"`
	Name                 string   `json:"name"`
	Summary              string   `json:"summary"`
	Author               string   `json:"author"`
	RecordNum            string   `json:"recordNum"`
	Title                string   `json:"title"`
	AliPay               string   `json:"aliPay"`
	WeixinPay            string   `json:"weixinPay"`
	OpenComment          string   `json:"openComment"`
	OpenMobileComment    string   `json:"openMobileComment"`
	OpenAdmiration       string   `json:"openAdmiration"`
	OpenMobileAdmiration string   `json:"openMobileAdmiration"`
	Github               string   `json:"github"`
	Gitee                string   `json:"gitee"`
	QqNumber             string   `json:"qqNumber"`
	QqGroup              string   `json:"qqGroup"`
	WeChat               string   `json:"weChat"`
	Email                string   `json:"email"`
	ShowList             string   `json:"showList"`
	LoginTypeList        string   `json:"loginTypeList"`
	PhotoList            []string `json:"photoList"`
	AliPayPhoto          string   `json:"aliPayPhoto"`
	WeixinPayPhoto       string   `json:"weixinPayPhoto"`
}
