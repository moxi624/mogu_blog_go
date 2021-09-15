package common

import (
"github.com/beego/beego/v2/core/logs"
)
/**
 *
 * @author  镜湖老杨
 * @date  2020/12/21 2:09 下午
 * @version 1.0
 */
func init() {
	logs.SetLogger("console")
}