package common

import "regexp"

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/11 2:40 下午
 * @version 1.0
 */

func CheckEmail(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func CheckMobile(mobileNumber string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNumber)
}

func CheckPicture(picture string) bool {
	regular := "<img\\s+(?:[^>]*)src\\s*=\\s*([^>]+)>"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(picture)
}
