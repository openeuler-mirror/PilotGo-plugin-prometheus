/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Oct 23 16:52:16 2024 +0800
 */
package initprometheus

import (
	"errors"
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"openeuler.org/PilotGo/prometheus-plugin/server/utils"
)

const (
	PrometheusYmlInit = "scripts/init_prometheus_yml.sh"
	PrometheusYml     = "/etc/prometheus/prometheus.yml"
)

func Init(httpaddr string) error {
	checkPrometheusVersion := `/usr/bin/prometheus --version | grep -oP '(2+\.\d+)\.\d+' | awk -F '.' '{print $2}'`
	_, stdout, _, _ := utils.RunCommand(checkPrometheusVersion)
	version, _ := strconv.Atoi(strings.Trim(stdout, "\n"))
	if version <= 28 {
		return errors.New(`please install prometheus greater than 2.28`)
	}

	ok, err := checkPrometheusYmlConsistency(httpaddr)
	if ok && err == nil {
		err = initPrometheusYml(httpaddr)
		if err != nil {
			return err
		}
	}
	return nil
}

func initPrometheusYml(httpaddr string) error {
	if err := backup(); err != nil {
		return err
	}

	if err := initYml(httpaddr); err != nil {
		return err
	}

	logger.Debug("prometheus init success")
	return nil
}

func backup() error {
	cmd := "cp " + PrometheusYml + " " + PrometheusYml + ".bak"
	exitcode, _, stderr, err := utils.RunCommand(cmd)
	if exitcode == 0 && stderr == "" && err == nil {
		return nil
	}
	return err
}

func initYml(httaddr string) error {
	cmd := "sh " + PrometheusYmlInit + " " + httaddr + " " + PrometheusYml
	exitcode, _, stderr, err := utils.RunCommand(cmd)
	if exitcode == 0 && stderr == "" && err == nil {
		return nil
	}
	return errors.New(stderr)
}
