package main

import (
	"flag"
	"pcapture/application"
	"pcapture/icmp"
	"pcapture/transport"
)

func main() {
	protocol := flag.String("p", "icmp", "protocol used for generate traffic")
	flag.Parse()

	switch *protocol {
	case "icmp":
		icmp.Gen()
	case "tcp":
		transport.GenTCP()
	case "udp":
		transport.GenUDP()
	case "quic":
		transport.GenQUIC()
	case "http":
		application.GenHTTP()
	case "dns":
		application.GenDNS()
	case "coap":
		application.GenCoAP()
	case "http3":
		application.GenHTTP3()

	}
}
