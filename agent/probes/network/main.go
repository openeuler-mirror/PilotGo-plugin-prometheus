/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Dec 2 09:19:50 2024 +0800
 */
package main

import (
	"openeuler.org/PilotGo/prometheus-plugin/agent/probes/network/flow"
	"openeuler.org/PilotGo/prometheus-plugin/agent/probes/network/global"
	"openeuler.org/PilotGo/prometheus-plugin/agent/probes/network/src"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	global.InitProcTcpManger()
	go flow.TcpNetFlow(stopper)
	src.Netflow(stopper)
}
