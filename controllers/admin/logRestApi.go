//Copyright (c) [2021] [YangLei]
//[mogu-go] is licensed under Mulan PSL v2.
//You can use this software according to the terms and conditions of the Mulan PSL v2.
//You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
//THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
//EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
//MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
//See the Mulan PSL v2 for more details.

package admin

import (
	"encoding/json"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/page"
	"mogu-go-v2/models/vo"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/2/25 8:49 上午
 * @version 1.0
 */

type LogRestApi struct {
	base.BaseController
}

func (c *LogRestApi) GetExceptionList() {
	var exceptionLogVO vo.ExceptionLogVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &exceptionLogVO)
	if err != nil {
		panic(err)
	}
	where := "status=?"
	if exceptionLogVO.Keyword != "" {
		where += " and content like '%" + exceptionLogVO.Keyword + "%'"
	}
	if exceptionLogVO.Operation != "" {
		time := strings.Split(exceptionLogVO.StartTime, ",")
		if len(time) == 2 {
			where += " and (create_time between \"" + common.DateUtils.Str2Date(time[0]) + "\" and \"" + common.DateUtils.Str2Date(time[1]) + "\")"
		}
	}
	var pageList []models.ExceptionLog
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.ExceptionLog{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((exceptionLogVO.CurrentPage - 1) * exceptionLogVO.PageSize).Limit(exceptionLogVO.PageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    exceptionLogVO.PageSize,
		Current: exceptionLogVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}

func (c *LogRestApi) GetLogList() {
	var sysLogVO vo.SysLogVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &sysLogVO)
	if err != nil {
		panic(err)
	}
	where := "status=?"
	if sysLogVO.UserName != "" {
		where += " and user_name='" + sysLogVO.UserName + "'"
	}
	if sysLogVO.Operation != "" {
		where += " and operator ='" + sysLogVO.Operation + "'"
	}
	if sysLogVO.Ip != "" {
		where += " and ip='" + sysLogVO.Ip + "'"
	}
	if sysLogVO.StartTime != "" {
		time := strings.Split(sysLogVO.StartTime, ",")
		if len(time) == 2 {
			where += " and (create_time between \"" + common.DateUtils.Str2Date(time[0]) + "\" and \"" + common.DateUtils.Str2Date(time[1]) + "\")"
		}
	}
	if sysLogVO.SpendTimeStr != "" {
		time := strings.Split(sysLogVO.SpendTimeStr, ",")
		if len(time) == 2 {
			where += " and (create_time between \"" + common.DateUtils.Str2Date(time[0]) + "\" and \"" + common.DateUtils.Str2Date(time[1]) + "\")"
		}
	}
	var pageList []models.SysLog
	var total int64
	c.Wg.Add(2)
	go func() {
		common.DB.Model(&models.SysLog{}).Where(where, 1).Count(&total)
		c.Wg.Done()
	}()
	go func() {
		common.DB.Where(where, 1).Offset((sysLogVO.CurrentPage - 1) * sysLogVO.PageSize).Limit(sysLogVO.PageSize).Order("create_time desc").Find(&pageList)
		c.Wg.Done()
	}()
	c.Wg.Wait()
	iPage := page.IPage{
		Records: pageList,
		Total:   total,
		Size:    sysLogVO.PageSize,
		Current: sysLogVO.CurrentPage,
	}
	c.SuccessWithData(iPage)
}
