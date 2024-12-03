/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Dec 2 09:19:50 2024 +0800
 */
package global

import (
	"encoding/binary"
	"net"
)

func IntToIP(_ipNum uint32) string {
	ip := make(net.IP, 4)
	binary.NativeEndian.PutUint32(ip, _ipNum)
	return ip.String()
}

func IsAddrValid(_addr string) bool {
	conn, err := net.Dial("tcp", _addr)
	if err != nil {
		return false	
	}
	defer conn.Close()
	return true
}