package util

import (
	"fmt"
	"net"
)

func LocalIp() string {
	ip := "127.0.0.1"

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Printf("获取本地IP发生异常 : %v\n", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()

				break
			}
		}
	}

	return ip
}
