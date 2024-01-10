package service

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"openeuler.org/PilotGo/prometheus-plugin/plugin"
)

func AddExtentions() {
	var ex []common.Extention
	me1 := &common.MachineExtention{
		Type: common.ExtentionMachine,
		Name: "安装exporter",
		URL:  "/plugin/prometheus/run?type=" + CommandInstall_Type,
	}
	me2 := &common.MachineExtention{
		Type: common.ExtentionMachine,
		Name: "卸载exporter",
		URL:  "/plugin/prometheus/run?type=" + CommandRemove_Type,
	}
	pe := &common.PageExtention{
		Type:    common.ExtentionPage,
		Name:    "plugin-prometheus",
		IsIndex: true,
	}
	ex = append(ex, me1, me2, pe)
	plugin.Client.RegisterExtention(ex)
}
