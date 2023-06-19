package transport

import (
	"fmt"
	"log"
	"net"
)

func echoUDPServer() error {
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		p := make([]byte, 1024)
		n, addr, err := listener.ReadFromUDP(p)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Printf("Server: Got '%s'\n", string(p[:n]))
		// Echo the data back to the client.
		listener.WriteToUDP(p[:n], addr)
	}

	return nil
}

func udpClientMain() error {
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := 0; i < 5; i++ {
		fmt.Printf("Client: Sending '%s'\n", message)
		_, err = conn.Write([]byte(message))
		if err != nil {
			return err
		}

		// Read data from the server.
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			return err
		}
		fmt.Printf("Client: Got '%s'\n", buf)
	}

	return nil
}

func GenUDP() {
	go func() {
		log.Fatal(echoUDPServer())
	}()

	err := udpClientMain()
	if err != nil {
		panic(err)
	}
}
