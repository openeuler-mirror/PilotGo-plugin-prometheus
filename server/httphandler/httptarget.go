package httphandler

import (
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/prometheus-plugin/dao"
	prometheus "openeuler.org/PilotGo/prometheus-plugin/service/prometheus"
)

func DBTargets(c *gin.Context) {
	targets, err := prometheus.PrometheusTargetsUpdate()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	c.JSON(200, targets)
}

func MonitorTargets(c *gin.Context) {
	targets, err := dao.QueryPrometheusTargets()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, targets, "获取到prometheus监控列表")
}
