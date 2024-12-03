/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Fri Oct 20 11:41:45 2023 +0800
 */
package httphandler

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/prometheus-plugin/server/plugin"
	"openeuler.org/PilotGo/prometheus-plugin/server/service"
)

// 运行远程命令安装、卸载exporter
func RunCommand(c *gin.Context) {
	d := &struct {
		MachineUUIDs []string `json:"uuids"`
	}{}
	if err := c.ShouldBind(d); err != nil {
		logger.Debug("绑定批次参数失败：%s", err)
		response.Fail(c, nil, "parameter error")
		return
	}

	var command string
	command_type := c.Query("type")
	if command_type == service.CommandInstall_Type {
		command = service.CommandInstall_Cmd
	} else if command_type == service.CommandRemove_Type {
		command = service.CommandRemove_Cmd
	} else {
		response.Fail(c, nil, "请重新检查命令参数type")
		return
	}

	run_result := func(result []*common.CmdResult) {
		for _, res := range result {
			if err := service.ProcessResult(res, command_type); err != nil {
				logger.Error("处理结果失败：%v", err.Error())
			}
		}
	}
	dd := &common.Batch{
		MachineUUIDs: d.MachineUUIDs,
	}
	err := plugin.Client.RunCommandAsync(dd, command, run_result)
	if err != nil {
		logger.Error("远程调用失败：%v", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "指令下发完成")
}
