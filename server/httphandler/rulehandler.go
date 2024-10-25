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

func DeleteRuleList(c *gin.Context) {
	id := c.Query("id")
	err := service.DeleteRuleList(id)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "删除成功")
}
func QueryRules(c *gin.Context) {

	query := &response.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	search := c.Query("search")
	data, total, err := service.SearchRules(search, query)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.DataPagination(c, data, total, query)
}
