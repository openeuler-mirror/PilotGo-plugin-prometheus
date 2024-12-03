/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Dec 2 09:19:50 2024 +0800
 */
package global

import "sync"

var Global_ProcTcpManager *ProcTcpManager

type TcpMetrics struct {
	Pid              uint32
	S_ip             string
	S_port           uint16
	Clients_addr_map sync.Map // key: "ip:port" value 1
	Family           uint16
	Role             string
	Comm             string
	OptFamily        uint16
	OptC_ip          uint32
	Rx               uint64
	Tx               uint64
}

type ProcTcpMetrics struct {
	TcpMetrics *TcpMetrics
	RxLatest   uint64 // bytes
	TxLatest   uint64 // bytes
}

type ProcTcpFlow struct {
	TcpMetrics *TcpMetrics
	RxFlow     uint64 // byte/s
	TxFlow     uint64 // byte/s
}

type ProcTcpManager struct {
	ProcTcpMetricsMap sync.Map
	ProcTcpFlowMap    sync.Map
}

func InitProcTcpManger() {
	Global_ProcTcpManager = &ProcTcpManager{
		ProcTcpMetricsMap: sync.Map{},
		ProcTcpFlowMap:    sync.Map{},
	}
}
