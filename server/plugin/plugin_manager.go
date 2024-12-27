/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 26 16:42:38 2023 +0800
 */
package plugin

import (
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"openeuler.org/PilotGo/prometheus-plugin/server/config"
)

var Client *client.Client

func Init(plugin *config.PluginPrometheus, prometheus *config.PrometheusServer) *client.PluginInfo {
	PluginInfo := client.PluginInfo{
		MenuName:    "监控告警",
		Name:        "prometheus",
		Version:     "1.0.1",
		Description: "Prometheus开源系统监视和警报工具包",
		Author:      "zhanghan",
		Email:       "zhanghan@kylinos.cn",
		Url:         plugin.URL,
		Icon:        "Odometer",
		PluginType:  "micro-app",
		ReverseDest: "http://" + prometheus.Addr,
	}

	return &PluginInfo
}
