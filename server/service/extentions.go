/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jan 10 16:23:18 2024 +0800
 */
package service

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"openeuler.org/PilotGo/prometheus-plugin/server/plugin"
)

func AddExtentions() {
	var ex []common.Extention
	me1 := &common.MachineExtention{
		Type:       common.ExtentionMachine,
		Name:       "安装exporter",
		URL:        "/plugin/prometheus/run?type=" + CommandInstall_Type,
		Permission: "plugin.prometheus.agent/install",
	}
	me2 := &common.MachineExtention{
		Type:       common.ExtentionMachine,
		Name:       "卸载exporter",
		URL:        "/plugin/prometheus/run?type=" + CommandRemove_Type,
		Permission: "plugin.prometheus.agent/uninstall",
	}
	pe := &common.PageExtention{
		Type:       common.ExtentionPage,
		Name:       "plugin-prometheus",
		URL:        "/plugin/prometheus",
		Permission: "plugin.prometheus.page/menu",
	}
	ex = append(ex, me1, me2, pe)
	plugin.Client.RegisterExtention(ex)
}
