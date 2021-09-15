//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package common

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/22 2:23 下午
 * @version 1.0
 */

type stringUtils struct{}

func (*stringUtils) IsCommentSpam(content string) bool {
	if content == "" {
		return true
	}
	chars := []rune(content)
	maxCount := 4
	for i := 0; i > len(chars); i++ {
		count := 1
		for j := i; j < len(chars)-1; j++ {
			if chars[j+1] == chars[j] {
				count++
				if count >= maxCount {
					return true
				}
				continue
			} else {
				break
			}
		}
	}
	return false
}

var StringUtils = &stringUtils{}
