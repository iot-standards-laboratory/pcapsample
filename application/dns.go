package application

import (
	"fmt"
	"net"
)

func GenDNS() {
	domains := []string{"localhost", "naver.com", "google.com", "daum.net", "samsung.com"}
	for _, d := range domains {
		getIP(d)
	}
}

func getIP(domain string) {
	fmt.Println("Domain Name:", domain)
	ips, _ := net.LookupIP(domain)
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			fmt.Println("IPv4: ", ipv4)
		}
	}
}
