/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Fri Oct 20 11:41:45 2023 +0800
 */
package service

import (
	"errors"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"openeuler.org/PilotGo/prometheus-plugin/server/dao"
	"openeuler.org/PilotGo/prometheus-plugin/server/model"
)

var ResultOptMsg = []string{"安装成功", "卸载成功"}

const (
	CommandInstall_Type = "install"
	CommandRemove_Type  = "uninstall"

	CommandInstall_Cmd = "yum install -y golang-github-prometheus-node_exporter && (echo '安装成功'; systemctl start node_exporter) || echo '安装失败'"
	CommandRemove_Cmd  = "yum remove -y golang-github-prometheus-node_exporter && echo '卸载成功' || echo '卸载失败'"
)

func ProcessResult(res *common.CmdResult, command_type string) error {
	result := &model.PrometheusTarget{
		UUID:     res.MachineUUID,
		TargetIP: res.MachineIP,
		Port:     "9100",
	}
	logger.Info("node-exporter安装状态:\n%v", res)
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
