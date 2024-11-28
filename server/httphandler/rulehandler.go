package httphandler

import (
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/prometheus-plugin/server/model"
	"openeuler.org/PilotGo/prometheus-plugin/server/service"
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

func UpdateRule(c *gin.Context) {

	var a model.Rule
	if err := c.Bind(&a); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	if len(a.AlertTargets) == 0 {
		response.Fail(c, nil, "请选择监控机器")
		return
	}
	if err := service.UpdateRule(&a); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "已更新告警配置")
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
func GetMonitorMetricsAndAlertLevel(c *gin.Context) {
	data := &struct {
		Metrics    []string `json:"metrics"`
		RuleLevel  []string `json:"ruleLevel"`
		AlertState []string `json:"alertState"`
		AlertLevel []string `json:"alertLevel"`
	}{
		Metrics:    []string{"cpu使用率", "内存使用率", "网络流入", "网络流出", "磁盘容量", "服务器宕机", "TCP连接数"},
		RuleLevel:  service.GetRuleLevel(),
		AlertState: []string{"活跃", "待处理", "已处理"},
		AlertLevel: service.GetAlertLevel(),
	}
	response.Success(c, data, "获取到监控指标和告警级别")
}
