package main

import (
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

	src.Netflow(stopper)
}
