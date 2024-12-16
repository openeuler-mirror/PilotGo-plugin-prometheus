/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 26 16:42:38 2023 +0800
 */
package router

import (
	"net/http"
	"strings"

	"gitee.com/openeuler/PilotGo-plugins/event/sdk"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/prometheus-plugin/server/httphandler"
	"openeuler.org/PilotGo/prometheus-plugin/server/plugin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(logger.RequestLogger([]string{
		"/plugin/prometheus/target",
	}))
	router.Use(gin.Recovery())

	return router
}

func RegisterAPIs(router *gin.Engine) {
	logger.Debug("router register")
	plugin.Client.RegisterHandlers(router)
	sdk.RegisterEventHandlers(router, plugin.Client)
	sdk.UnPluginListenEventHandler()

	// prometheus api代理
	prometheus := router.Group("/plugin/" + plugin.Client.PluginInfo.Name + "/api/v1")
	{
		prometheus.GET("/query", func(c *gin.Context) {
			c.Set("query", plugin.Client.PluginInfo.ReverseDest)
			httphandler.Query(c)
		})
		prometheus.GET("/query_range", func(c *gin.Context) {
			c.Set("query_range", plugin.Client.PluginInfo.ReverseDest)
			httphandler.QueryRange(c)
		})
		prometheus.GET("/targets", func(c *gin.Context) {
			c.Set("targets", plugin.Client.PluginInfo.ReverseDest)
			httphandler.Targets(c)
		})
		prometheus.GET("/alerts", func(c *gin.Context) {
			c.Set("alerts", plugin.Client.PluginInfo.ReverseDest)
			httphandler.Alerts(c)
		})

	}

	// prometheus配置文件http方式获取监控target
	DBTarget := router.Group("/plugin/" + plugin.Client.PluginInfo.Name)
	{
		DBTarget.GET("target", httphandler.DBTargets)
	}

	//prometheus target crud
	targetManager := router.Group("/plugin/" + plugin.Client.PluginInfo.Name)
	{
		targetManager.POST("run", httphandler.RunCommand)
		targetManager.GET("monitorlist", httphandler.MonitorTargets)
	}

	//prometheus alert rule manager
	ruleManager := router.Group("/plugin/" + plugin.Client.PluginInfo.Name)
	{
		ruleManager.POST("ruleAdd", httphandler.AddRuleHandler)
		ruleManager.GET("ruleQuery", httphandler.QueryRules)
		ruleManager.GET("ruleDelete", httphandler.DeleteRuleList)
		ruleManager.POST("ruleUpdate", httphandler.UpdateRule)
		ruleManager.GET("ruleMetrics", httphandler.GetMonitorMetricsAndAlertLevel)
	}

	//prometheus alert  manager
	alertManager := router.Group("/plugin/" + plugin.Client.PluginInfo.Name)
	{
		alertManager.GET("alertQuery", httphandler.QuerySearchAlerts)
		alertManager.POST("alertUpdateState", httphandler.UpdateHandleState)
	}
}

func StaticRouter(router *gin.Engine) {
	router.Static("/plugin/prometheus/static", "../web/dist/static")
	router.StaticFile("/plugin/prometheus", "../web/dist/index.html")

	// 解决页面刷新404的问题
	router.NoRoute(func(c *gin.Context) {
		logger.Error("process noroute: %s", c.Request.URL.RawPath)
		if !strings.HasPrefix(c.Request.RequestURI, "/plugin/prometheus") {
			c.File("./web/dist/index.html")
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})
}
