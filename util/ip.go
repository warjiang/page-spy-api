package util

import (
	"math/rand"
	"net"
)

// GetLocalIP 获取本机ip 获取失败返回 ""
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return ""
}

func GetLocalIPList() []string {
	addrs, err := net.InterfaceAddrs()

	arr := []string{}
	if err != nil {
		return arr
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				arr = append(arr, ipNet.IP.String())
			}
		}
	}

	return arr
}

func IsIP(ipOrDns string) bool {
	if ip := net.ParseIP(ipOrDns); ip != nil {
		return true
	}
	return false
}

func ResolveIP(ipOrDns string) string {
	if IsIP(ipOrDns) {
		return ipOrDns
	}
	localIP := GetLocalIP()
	ips, err := net.LookupIP(ipOrDns)
	if err != nil {
		return ipOrDns
	}

	for _, ip := range ips {
		if ip.String() == localIP {
			return localIP
		}
	}
	// not found any matched ip address in ips
	if len(ips) > 0 {
		index := rand.Intn(len(ips))
		return ips[index].String()
	}
	return ipOrDns
}
