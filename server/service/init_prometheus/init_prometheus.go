package initprometheus

import (
	"errors"
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"openeuler.org/PilotGo/prometheus-plugin/utils"
)

const (
	GlobalPrometheusYmlInit = "scripts/init_prometheus_yml.sh"
	GlobalPrometheusYml     = "/etc/prometheus/prometheus.yml"
)

func InitPrometheus(httpaddr string) error {
	checkPrometheus := `/usr/bin/prometheus --version | grep -oP '(2+\.\d+)\.\d+' | awk -F '.' '{print $2}'`
	_, stdout, _, _ := utils.RunCommand(checkPrometheus)
	version, _ := strconv.Atoi(strings.Trim(stdout, "\n"))
	if version <= 28 {
		return errors.New(`please install prometheus greater than 2.28`)
	}

	ok, err := CheckYMLHash(httpaddr)
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
	cmd := "cp " + GlobalPrometheusYml + " " + GlobalPrometheusYml + ".bak"
	exitcode, _, stderr, err := utils.RunCommand(cmd)
	if exitcode == 0 && stderr == "" && err == nil {
		return nil
	}
	return err
}

func initYML(httaddr string) error {
	cmd := "sh " + GlobalPrometheusYmlInit + " " + httaddr + " " + GlobalPrometheusYml
	exitcode, _, stderr, err := utils.RunCommand(cmd)
	if exitcode == 0 && stderr == "" && err == nil {
		return nil
	}
	return errors.New(stderr)
}
