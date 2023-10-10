package main

import (
	"fmt"
	"os"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"openeuler.org/PilotGo/prometheus-plugin/config"
	"openeuler.org/PilotGo/prometheus-plugin/db"
	"openeuler.org/PilotGo/prometheus-plugin/global"
	"openeuler.org/PilotGo/prometheus-plugin/httphandler/service"
	"openeuler.org/PilotGo/prometheus-plugin/plugin"
	"openeuler.org/PilotGo/prometheus-plugin/router"
)

func main() {
	fmt.Println("hello prometheus")

	config.Init()

	if err := logger.Init(config.Config().Logopts); err != nil {
		fmt.Printf("logger init failed, please check the config file: %s", err)
		os.Exit(-1)
	}

	if err := service.InitPrometheus(config.Config().HttpServer.Addr); err != nil {
		logger.Error("check prometheus error: %s", err)
		os.Exit(-1)
	}

	if err := db.MysqldbInit(config.Config().Mysql); err != nil {
		logger.Error("mysql db init failed, please check again: %s", err)
		os.Exit(-1)
	}

	server := router.InitRouter()

	global.GlobalClient = client.DefaultClient(plugin.Init(config.Config().PluginPrometheus, config.Config().PrometheusServer))
	router.RegisterAPIs(server)
	router.StaticRouter(server)
	global.GlobalClient.Server = config.Config().HttpServer.Addr

	if err := server.Run(config.Config().HttpServer.Addr); err != nil {
		logger.Fatal("failed to run server")
	}
}
