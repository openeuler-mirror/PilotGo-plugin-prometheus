package service

import (
	"errors"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"openeuler.org/PilotGo/prometheus-plugin/dao"
	"openeuler.org/PilotGo/prometheus-plugin/model"
)

var ResultOptMsg = []string{"安装成功", "卸载成功"}

const (
	CommandInstall_Type = "install"
	CommandRemove_Type  = "remove"

	CommandInstall_Cmd = "yum install -y golang-github-prometheus-node_exporter && (echo '安装成功'; systemctl start node_exporter) || echo '安装失败'"
	CommandRemove_Cmd  = "yum remove -y golang-github-prometheus-node_exporter && echo '卸载成功' || echo '卸载失败'"
)

func ProcessResult(res *common.CmdResult, command_type string) error {
	result := &model.PrometheusTarget{
		UUID:     res.MachineUUID,
		TargetIP: res.MachineIP,
		Port:     "9100",
	}

	ok, err := dao.IsExistTargetUUID(res.MachineUUID)
	if err != nil {
		return err
	}

	if !ok && command_type == CommandInstall_Type && ResultOptStdout(res) {
		if err := dao.AddPrometheusTarget(result); err != nil {
			return errors.New("保存结果失败：" + err.Error())
		}
	}
	if ok && command_type == CommandRemove_Type && ResultOptStdout(res) {
		if err := dao.DeletePrometheusTarget(result); err != nil {
			return errors.New("删除prometheus target失败: " + err.Error())
		}
	}
	return nil
}

func ResultOptStdout(res *common.CmdResult) bool {
	stdout := res.Stdout
	for _, msg := range ResultOptMsg {
		if strings.Contains(stdout, msg) {
			return true
		}
	}
	return false
}
