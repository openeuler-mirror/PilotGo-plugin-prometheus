/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 26 16:42:38 2023 +0800
 */
package main

import (
	"fmt"
	"os"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"openeuler.org/PilotGo/prometheus-plugin/server/config"
	"openeuler.org/PilotGo/prometheus-plugin/server/db"
	"openeuler.org/PilotGo/prometheus-plugin/server/plugin"
	"openeuler.org/PilotGo/prometheus-plugin/server/router"
	"openeuler.org/PilotGo/prometheus-plugin/server/service"
	prometheus "openeuler.org/PilotGo/prometheus-plugin/server/service/prometheus"
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

	if err := service.PullAlert(); err != nil {
		logger.Error("pull prometheus alerts failed, please check the code: %s", err)
		os.Exit(-1)
	}

	server := router.InitRouter()

	plugin.Client = client.DefaultClient(plugin.Init(config.Config().PluginPrometheus, config.Config().PrometheusServer))
	router.RegisterAPIs(server)
	router.StaticRouter(server)
	service.GetTags()        // pilotgo机器列表tag标签
	service.AddExtentions()  // 添加扩展点
	service.AddPermissions() // 添加权限

	if err := server.Run(config.Config().HttpServer.Addr); err != nil {
		logger.Fatal("failed to run server: %v", err)
	}
}
