package global

import (
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gorm.io/gorm"
)

var (
	GlobalClient *client.Client
	GlobalDB     *gorm.DB
)

const (
	GlobalPrometheusYmlInit = "./scripts/init_prometheus_yml.sh"
	GlobalPrometheusYml     = "/etc/prometheus/prometheus.yml"
)
