/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Dec 2 09:19:50 2024 +0800
 */
package flow

import (
	"fmt"
	"os"
	"time"

	"openeuler.org/PilotGo/prometheus-plugin/agent/probes/network/global"
)

func TcpNetFlow(stop chan os.Signal) {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			global.Global_ProcTcpManager.ProcTcpMetricsMap.Range(func(key, value any) bool {
				pid := key.(uint32)
				proc_tcp_metrics := value.(*global.ProcTcpMetrics)

				proc_tcp_flow_any, _ := global.Global_ProcTcpManager.ProcTcpFlowMap.Load(pid)
				proc_tcp_flow := proc_tcp_flow_any.(*global.ProcTcpFlow)
				proc_tcp_flow.TxFlow = proc_tcp_metrics.TcpMetrics.Tx - proc_tcp_metrics.TxLatest
				proc_tcp_flow.RxFlow = proc_tcp_metrics.TcpMetrics.Rx - proc_tcp_metrics.RxLatest

				proc_tcp_metrics.TxLatest = proc_tcp_metrics.TcpMetrics.Tx
				proc_tcp_metrics.RxLatest = proc_tcp_metrics.TcpMetrics.Rx

				// ttcode
				if proc_tcp_metrics.TcpMetrics.Tx < proc_tcp_metrics.TxLatest {
					fmt.Printf(">>>tx: %d, txlatest: %d\n", proc_tcp_metrics.TcpMetrics.Tx, proc_tcp_metrics.TxLatest)
				}
				clients := []string{}
				proc_tcp_metrics.TcpMetrics.Clients_addr_map.Range(func(key, value any) bool {
					addr := key.(string)
					clients = append(clients, addr)
					return true
				})
				fmt.Printf("\033[33mtime\033[0m: %s \033[33mcomm\033[0m: %s \033[33mpid\033[0m: %d \033[33mrole\033[0m: %s \033[33mclient\033[0m(%v)->\033[33mserver\033[0m(%s:%d) \033[32mtx: %d Byte/s rx: %d Byte/s\033[0m\n",
					time.Now().Format("15:04:05"),
					proc_tcp_metrics.TcpMetrics.Comm,
					pid,
					proc_tcp_metrics.TcpMetrics.Role,
					clients,
					proc_tcp_metrics.TcpMetrics.S_ip, proc_tcp_metrics.TcpMetrics.S_port,
					proc_tcp_flow.TxFlow,
					proc_tcp_flow.RxFlow,
				)
				return true
			})
		}
	}
}
