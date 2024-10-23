package httphandler

import (
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/prometheus-plugin/model"
	"openeuler.org/PilotGo/prometheus-plugin/service"
)

func AddRuleHandler(c *gin.Context) {
	var alert model.Rule
	if err := c.Bind(&alert); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	if len(alert.AlertTargets) == 0 {
		response.Fail(c, nil, "请选择监控机器")
		return
	}
	err := service.AddRule(&alert)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "添加告警配置成功")
}
