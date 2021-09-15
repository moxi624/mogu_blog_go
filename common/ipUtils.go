package common

import (
	"net"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/18 3:58 下午
 * @version 1.0
 */

type ipUtils struct{}

func (ipUtils) GetAddresses(ip string) string {
	//host := "https://ips.market.alicloudapi.com"
	//path := "/iplocaltion"
	//appcode, _ := beego.AppConfig.String("appCode")
	//urlSend := host + path + "?ip=" + ip
	//req := httplib.Get(urlSend)
	//req.Header("Authorization", "APPCODE "+appcode)
	//var ipAddress vo.IpAddress
	//err := req.ToJSON(&ipAddress)
	//if err != nil {
	//	panic(err)
	//}
	//if ipAddress.Code == 101 {
	//	return "内网访问"
	//}
	//return ipAddress.Result.Nation + ipAddress.Result.Province + ipAddress.Result.City + ipAddress.Result.District
	return "内网访问"
}

func (ipUtils) GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

var IpUtils = &ipUtils{}
