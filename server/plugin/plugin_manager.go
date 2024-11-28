package plugin

import (
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"openeuler.org/PilotGo/prometheus-plugin/server/config"
)

var Client *client.Client

func Init(plugin *config.PluginPrometheus, prometheus *config.PrometheusServer) *client.PluginInfo {
	PluginInfo := client.PluginInfo{
		Name:        "prometheus",
		Version:     "1.0.1",
		Description: "Prometheus开源系统监视和警报工具包",
		Author:      "zhanghan",
		Email:       "zhanghan@kylinos.cn",
		Url:         plugin.URL,
		PluginType:  "iframe",
		ReverseDest: "http://" + prometheus.Addr,
	}

	return &PluginInfo
}
