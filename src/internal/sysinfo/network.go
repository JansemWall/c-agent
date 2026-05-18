package sysinfo

import (
	"net"
	"os"
)

func GetHostInfo() (string, string) {
	hostname, _ := os.Hostname()
	addrs, _ := net.InterfaceAddrs()
	var ip string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return hostname, ip
}