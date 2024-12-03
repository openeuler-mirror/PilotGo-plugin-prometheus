/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Dec 2 09:19:50 2024 +0800
 */
package src

import (
	"bytes"
	"encoding/binary"
	"log"
	"openeuler.org/PilotGo/prometheus-plugin/agent/probes/network/global"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type tcp_metrics bpf tcp_netflow.bpf.c -- -g -O2 -D __TARGET_ARCH_x86 -I./include

// func ltoa(ip uint32) string {
// 	return fmt.Sprintf("%v.%v.%v.%v", (ip>>24)&0xFF, (ip>>16)&0xFF, (ip>>8)&0xFF, (ip & 0xFF))
// }

// 更新 ProcTcpMetricsMap、ProcTcpFlowMap中的信息
// @return: debug
func AddProcTcp(_bpfMetrics bpfTcpMetrics) *global.TcpMetrics {
	client_addr := string(global.IntToIP(_bpfMetrics.C_ip)) + ":" + strconv.Itoa(int(_bpfMetrics.C_port))
	role := ""
	switch _bpfMetrics.Role {
	case 1:
		role = "client"
	case 0:
		role = "server"
	}

	tcpmetrics := &global.TcpMetrics{}
	tcpmetrics.Pid = _bpfMetrics.Pid
	tcpmetrics.S_ip = global.IntToIP(_bpfMetrics.S_ip)
	tcpmetrics.S_port = _bpfMetrics.S_port
	tcpmetrics.Family = _bpfMetrics.Family
	tcpmetrics.Role = role
	tcpmetrics.Comm = string(_bpfMetrics.Comm[:])
	tcpmetrics.Rx = tcpmetrics.Rx + _bpfMetrics.Rx
	tcpmetrics.Tx = tcpmetrics.Tx + _bpfMetrics.Tx
	tcpmetrics.Clients_addr_map.Store(client_addr, 1)

	if value, ok := global.Global_ProcTcpManager.ProcTcpMetricsMap.Load(_bpfMetrics.Pid); ok {
		proc_tcp_metrics := value.(*global.ProcTcpMetrics)
		proc_tcp_metrics.TcpMetrics.Rx = proc_tcp_metrics.TcpMetrics.Rx + _bpfMetrics.Rx
		proc_tcp_metrics.TcpMetrics.Tx = proc_tcp_metrics.TcpMetrics.Tx + _bpfMetrics.Tx
		proc_tcp_metrics.TcpMetrics.Role = tcpmetrics.Role
		proc_tcp_metrics.TcpMetrics.S_ip = tcpmetrics.S_ip
		proc_tcp_metrics.TcpMetrics.S_port = tcpmetrics.S_port

		if _, ok := proc_tcp_metrics.TcpMetrics.Clients_addr_map.Load(client_addr); !ok {
			// TODO: client port可能被其他进程占用
			proc_tcp_metrics.TcpMetrics.Clients_addr_map.Range(func(key, value any) bool {
				addr := key.(string)
				if !global.IsAddrValid(addr) {
					proc_tcp_metrics.TcpMetrics.Clients_addr_map.Delete(addr)
				}
				return true
			})
			proc_tcp_metrics.TcpMetrics.Clients_addr_map.Store(client_addr, 1)
		}

		proc_tcp_flow_any, _ := global.Global_ProcTcpManager.ProcTcpFlowMap.Load(_bpfMetrics.Pid)
		proc_tcp_flow := proc_tcp_flow_any.(*global.ProcTcpFlow)
		proc_tcp_flow.TcpMetrics.Role = tcpmetrics.Role
		proc_tcp_flow.TcpMetrics.S_ip = tcpmetrics.S_ip
		proc_tcp_flow.TcpMetrics.S_port = tcpmetrics.S_port
		if _, ok := proc_tcp_flow.TcpMetrics.Clients_addr_map.Load(client_addr); !ok {
			proc_tcp_flow.TcpMetrics.Clients_addr_map.Range(func(key, value any) bool {
				addr := key.(string)
				if !global.IsAddrValid(addr) {
					proc_tcp_flow.TcpMetrics.Clients_addr_map.Delete(addr)
				}
				return true
			})
			proc_tcp_flow.TcpMetrics.Clients_addr_map.Store(client_addr, 1)
		}
		return proc_tcp_metrics.TcpMetrics
	}

	proctcpmetrics := &global.ProcTcpMetrics{}
	proctcpmetrics.TcpMetrics = tcpmetrics
	global.Global_ProcTcpManager.ProcTcpMetricsMap.Store(tcpmetrics.Pid, proctcpmetrics)

	proctcpflow := &global.ProcTcpFlow{}
	proctcpflow.TcpMetrics = tcpmetrics
	global.Global_ProcTcpManager.ProcTcpFlowMap.Store(tcpmetrics.Pid, proctcpflow)
	return tcpmetrics
}

func Netflow(stop chan os.Signal) {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// kp, err := link.Kprobe("tcp_set_state", objs.TcpSetState, nil)
	// if err != nil {
	// 	log.Fatalf("opening kprobe: %s", err)
	// }
	// defer kp.Close()
	kp, err := link.Kprobe("tcp_v4_connect", objs.TcpV4Connect, nil)
	if err != nil {
		log.Fatalf("opening kprobe: %s", err)
	}
	defer kp.Close()
	kp, err = link.Kretprobe("inet_csk_accept", objs.InetCskAcceptExit, nil)
	if err != nil {
		log.Fatalf("opening kretprobe: %s", err)
	}
	defer kp.Close()
	kp, err = link.Kprobe("tcp_sendmsg", objs.TcpSendmsg, nil)
	if err != nil {
		log.Fatalf("opening kprobe: %s", err)
	}
	defer kp.Close()
	kp, err = link.Kretprobe("tcp_sendmsg", objs.TcpSendmsgExit, nil)
	if err != nil {
		log.Fatalf("opening kretprobe: %s", err)
	}
	defer kp.Close()
	kp, err = link.Kprobe("tcp_cleanup_rbuf", objs.TcpCleanupRbuf, nil)
	if err != nil {
		log.Fatalf("opening kprobe: %s", err)
	}
	defer kp.Close()
	kp, err = link.Kprobe("tcp_v4_destroy_sock", objs.TcpV4DestroySock, nil)
	if err != nil {
		log.Fatalf("opening kprobe: %s", err)
	}
	defer kp.Close()
	kp, err = link.Tracepoint("tcp", "tcp_destroy_sock", objs.BpfTraceTcpDestroySockFunc, nil)
	if err != nil {
		log.Fatalf("opening tracepoint: %s", err)
	}
	defer kp.Close()
	kp, err = link.AttachRawTracepoint(link.RawTracepointOptions{
		Name:    "tcp_destroy_sock",
		Program: objs.BpfRawTraceTcpDestroySock,
	})
	if err != nil {
		log.Fatalf("opening raw tracepoint: %s", err)
	}
	defer kp.Close()

	rd, err := ringbuf.NewReader(objs.TcpOutput)
	if err != nil {
		log.Fatalf("opening ringbuf reader: %s", err)
	}
	defer rd.Close()
	go func() {
		<-stop

		if err := rd.Close(); err != nil {
			log.Fatalf("closing ringbuf reader: %s", err)
		}
	}()

	var bpftcpmetrics bpfTcpMetrics
	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Println("Received signal, exiting..")
				return
			}
			log.Printf("reading from reader: %s", err)
			continue
		}

		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &bpftcpmetrics); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}

		AddProcTcp(bpftcpmetrics)

		// ttcode
		// clients := []string{}
		// tcpmetrics := AddProcTcp(bpftcpmetrics)
		// tcpmetrics.Clients_addr_map.Range(func(key, value any) bool {
		// 	addr := key.(string)
		// 	clients = append(clients, addr)
		// 	return true
		// })
		// fmt.Printf("pid: %d client(%v) server(%s:%d) family: %d socket role: %s comm: %v rx: %v tx:%v\n",
		// 	tcpmetrics.Pid,
		// 	clients,
		// 	tcpmetrics.S_ip, tcpmetrics.S_port,
		// 	tcpmetrics.Family,
		// 	tcpmetrics.Role,
		// 	tcpmetrics.Comm,
		// 	tcpmetrics.Rx, tcpmetrics.Tx,
		// )
	}
}
