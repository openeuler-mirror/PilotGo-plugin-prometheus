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

func ProcessResult(res *common.RunResult, command_type string) error {
	if res.Error != nil {
		return errors.New(res.Error.(string))
	}
	result := &model.PrometheusTarget{
		UUID:     res.CmdResult.MachineUUID,
		TargetIP: res.CmdResult.MachineIP,
		Port:     "9100",
	}

	ok, err := dao.IsExistTargetUUID(res.CmdResult.MachineUUID)
	if err != nil {
		return err
	}

	if !ok && command_type == CommandInstall_Type && ResultOptStdout(res) {
		if Err := dao.AddPrometheusTarget(result); Err != nil {
			return errors.New("保存结果失败：" + Err.Error())
		}
	}
	if ok && command_type == CommandRemove_Type && ResultOptStdout(res) {
		if Err := dao.DeletePrometheusTarget(result); Err != nil {
			return errors.New("删除prometheus target失败: " + Err.Error())
		}
	}
	return nil
}

func ResultOptStdout(res *common.RunResult) bool {
	stdout := res.CmdResult.Stdout
	for _, msg := range ResultOptMsg {
		if strings.Contains(stdout, msg) {
			return true
		}
	}
	return false
}
