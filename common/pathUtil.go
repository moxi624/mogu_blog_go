package common

import (
	"fmt"
	"net/url"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/13 3:44 下午
 * @version 1.0
 */

type pathUtil struct{}

func (pathUtil) UrlDecode(urlString string) string {
	decodeUrl, err := url.QueryUnescape(urlString)
	if err != nil {
		fmt.Println(err)
	}
	return decodeUrl
}

var PathUtil = &pathUtil{}
