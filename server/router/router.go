package router

import (
	"net/http"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/prometheus-plugin/httphandler"
	"openeuler.org/PilotGo/prometheus-plugin/plugin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(logger.RequestLogger())
	router.Use(gin.Recovery())

	return router
}

func RegisterAPIs(router *gin.Engine) {
	logger.Debug("router register")
	plugin.Client.RegisterHandlers(router)

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
		targetManager.POST("addTarget", httphandler.AddPrometheusTarget)
		targetManager.DELETE("delTarget", httphandler.DeletePrometheusTarget)
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
