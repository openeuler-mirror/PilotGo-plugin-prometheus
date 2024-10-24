package main

import (
	"fmt"
	"os"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"openeuler.org/PilotGo/prometheus-plugin/config"
	"openeuler.org/PilotGo/prometheus-plugin/db"
	"openeuler.org/PilotGo/prometheus-plugin/plugin"
	"openeuler.org/PilotGo/prometheus-plugin/router"
	"openeuler.org/PilotGo/prometheus-plugin/service"
	prometheus "openeuler.org/PilotGo/prometheus-plugin/service/prometheus"
)

func main() {
	fmt.Println("hello prometheus")

	config.Init()

	if err := logger.Init(config.Config().Logopts); err != nil {
		fmt.Printf("logger init failed, please check the config file: %s", err)
		os.Exit(-1)
	}

	if err := prometheus.Init(config.Config().HttpServer.Addr); err != nil {
		logger.Error("check prometheus error: %s", err)
		os.Exit(-1)
	}

	if err := db.MysqldbInit(config.Config().Mysql); err != nil {
		logger.Error("mysql db init failed, please check again: %s", err)
		os.Exit(-1)
	}

	server := router.InitRouter()

	plugin.Client = client.DefaultClient(plugin.Init(config.Config().PluginPrometheus, config.Config().PrometheusServer))
	router.RegisterAPIs(server)
	router.StaticRouter(server)
	service.GetTags()       // pilotgo机器列表tag标签
	service.AddExtentions() // 添加扩展点

	if err := server.Run(config.Config().HttpServer.Addr); err != nil {
		logger.Fatal("failed to run server: %v", err)
	}
}
