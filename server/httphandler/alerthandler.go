package httphandler

import (
	"encoding/json"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/prometheus-plugin/dao"
	"openeuler.org/PilotGo/prometheus-plugin/model"
	"openeuler.org/PilotGo/prometheus-plugin/service"
	"openeuler.org/PilotGo/prometheus-plugin/utils"
)

func QuerySearchAlerts(c *gin.Context) {
	query := &response.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	alertName := c.Query("alertName")
	ip := c.Query("ip")
	level := c.Query("level")
	handleState := c.Query("handleState")
	alertState := c.Query("state")
	alertStart := c.Query("alertStart")
	alertEnd := c.Query("alertEnd")
	search := c.Query("search")

	var as model.AlertTime
	if len(alertStart) != 0 {
		err = json.Unmarshal([]byte(alertStart), &as)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
	}

	var ae model.AlertTime
	if len(alertEnd) != 0 {
		err = json.Unmarshal([]byte(alertEnd), &ae)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
	}

	Ifpage := c.Query("paged")
	switch Ifpage {
	case "false":
		if search == "true" {
			data, total, err := service.SearchAlerts(alertName, ip, level, handleState, alertState, as, ae)
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			response.DataPagination(c, data, total, query)
		} else {
			data, total, err := dao.QueryAlerts()
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			response.DataPagination(c, data, total, query)
		}

	default:
		if search == "true" {
			data, total, err := service.SearchAlerts(alertName, ip, level, handleState, alertState, as, ae)
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			lists, err := utils.DataPaging(query, data, total)
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			response.DataPagination(c, lists, total, query)
		} else {
			data, total, err := dao.QueryAlerts()
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			lists, err := utils.DataPaging(query, data, total)
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			response.DataPagination(c, lists, total, query)
		}
	}
}
func UpdateHandleState(c *gin.Context) {
	ids := &struct {
		Ids         []int  `json:"ids"`
		HandleState string `json:"state"`
	}{}
	if err := c.Bind(&ids); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	for _, id := range ids.Ids {
		if err := service.UpdateHandleState(id, ids.HandleState); err != nil {
			logger.Error("id=%v 更新处理状态失败：%v", id, err.Error())
		}
	}

	response.Success(c, nil, "已更新处理状态")
}
