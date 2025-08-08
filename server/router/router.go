/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 26 16:42:38 2023 +0800
 */
package router

import (
	"fmt"
	"net/http"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/prometheus-plugin/server/config"
	"openeuler.org/PilotGo/prometheus-plugin/server/httphandler"
	"openeuler.org/PilotGo/prometheus-plugin/server/plugin"
)

func StartHttpServer() {
	server := initRouter()

	plugin.Client.RegisterHandlers(server)

	plugin.GetTags() // pilotgo机器列表tag标签

	if err := server.Run(config.Config().HttpServer.Addr); err != nil {
		logger.Fatal("failed to run server: %v", err)
	}
}
func initRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(logger.RequestLogger([]string{
		"/plugin/prometheus/api/target",
	}))
	router.Use(gin.Recovery())

	registerStaticRouter(router)
	registerAPIs(router)
	return router
}

func registerAPIs(router *gin.Engine) {
	logger.Info("router register")
	// prometheus api代理
	api := router.Group("/plugin/prometheus/api")
	prometheus := api.Group("")
	{
		prometheus.GET("/query", func(c *gin.Context) {
			c.Set("query", fmt.Sprintf("http://%s", config.Config().PrometheusServer.Addr))
			httphandler.Query(c)
		})
		prometheus.GET("/query_range", func(c *gin.Context) {
			c.Set("query_range", fmt.Sprintf("http://%s", config.Config().PrometheusServer.Addr))
			httphandler.QueryRange(c)
		})
		prometheus.GET("/targets", func(c *gin.Context) {
			c.Set("targets", fmt.Sprintf("http://%s", config.Config().PrometheusServer.Addr))
			httphandler.Targets(c)
		})
		prometheus.GET("/alerts", func(c *gin.Context) {
			c.Set("alerts", fmt.Sprintf("http://%s", config.Config().PrometheusServer.Addr))
			httphandler.Alerts(c)
		})

	}

	// prometheus配置文件http方式获取监控target
	DBTarget := api.Group("")
	{
		DBTarget.GET("target", httphandler.DBTargets)
	}

	//prometheus target crud
	targetManager := api.Group("")
	{
		targetManager.POST("run", httphandler.RunCommand)
		targetManager.GET("monitorlist", httphandler.MonitorTargets)
	}

	//prometheus alert rule manager
	ruleManager := api.Group("")
	{
		ruleManager.POST("ruleAdd", httphandler.AddRuleHandler)
		ruleManager.GET("ruleQuery", httphandler.QueryRules)
		ruleManager.GET("ruleDelete", httphandler.DeleteRuleList)
		ruleManager.POST("ruleUpdate", httphandler.UpdateRule)
		ruleManager.GET("ruleMetrics", httphandler.GetMonitorMetricsAndAlertLevel)
	}

	//prometheus alert  manager
	alertManager := api.Group("")
	{
		alertManager.GET("alertQuery", httphandler.QuerySearchAlerts)
		alertManager.POST("alertUpdateState", httphandler.UpdateHandleState)
	}
}

func registerStaticRouter(router *gin.Engine) {
	router.Static("/plugin/prometheus/assets", "../web/dist/assets")
	router.StaticFile("/plugin/prometheus", "../web/dist/index.html")

	// 解决页面刷新404的问题
	router.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/plugin/prometheus/api") {
			c.File("../web/dist/index.html")
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})
}
