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