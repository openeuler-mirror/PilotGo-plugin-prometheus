/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Fri Oct 20 11:41:45 2023 +0800
 */
package httphandler

import (
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/prometheus-plugin/server/dao"
	prometheus "openeuler.org/PilotGo/prometheus-plugin/server/service/prometheus"
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
