package service

import (
	"errors"
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo-plugins/sdk/logger"
	"gitee.com/openeuler/PilotGo-plugins/sdk/utils/command"
	"openeuler.org/PilotGo/prometheus-plugin/global"
)

func InitPrometheus(httpaddr string) error {
	checkNodeExporter := "rpm -qa |grep  golang-github-prometheus-node_exporter"
	_, stdout1, _, _ := command.RunCommand(checkNodeExporter)
	if len(strings.Trim(stdout1, "\n")) == 0 {
		return errors.New(`please use "yum install -y golang-github-prometheus-node_exporter" to install it`)
	}

	checkPrometheus := `/usr/bin/prometheus --version | grep -oP '(2+\.\d+)\.\d+' | awk -F '.' '{print $2}'`
	_, stdout2, _, _ := command.RunCommand(checkPrometheus)
	version, _ := strconv.Atoi(strings.Trim(stdout2, "\n"))
	if version <= 28 {
		return errors.New(`please install prometheus greater than 2.28`)
	}

	ok, err := CheckYMLHash()
	if ok && err == nil {
		err = initPrometheusYML(httpaddr)
		if err != nil {
			return err
		}
	}
	return nil
}

func initPrometheusYML(httpaddr string) error {
	if err := backup(); err != nil {
		return err
	}

	if err := initYML(httpaddr); err != nil {
		return err
	}

	logger.Debug("prometheus init success")
	return nil
}

func backup() error {
	cmd := "cp " + global.GlobalPrometheusYml + " " + global.GlobalPrometheusYml + ".bak"
	exitcode, _, stderr, err := command.RunCommand(cmd)
	if exitcode == 0 && stderr == "" && err == nil {
		return nil
	}
	return err
}

func initYML(httaddr string) error {
	cmd := "sh " + global.GlobalPrometheusYmlInit + " " + httaddr + " " + global.GlobalPrometheusYml
	exitcode, _, stderr, err := command.RunCommand(cmd)
	if exitcode == 0 && stderr == "" && err == nil {
		return nil
	}
	return err
}
