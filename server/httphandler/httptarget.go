package httphandler

import (
	"net/http"

	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/prometheus-plugin/dao"
	"openeuler.org/PilotGo/prometheus-plugin/model"
)

func DBTargets(c *gin.Context) {
	targets, err := dao.GetPrometheusTarget()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	objs := []model.PrometheusObject{
		{
			Targets: targets,
		},
	}
	c.JSON(http.StatusOK, objs)
}
