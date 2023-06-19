package transport

import (
	"fmt"
	"io"
	"log"
	"net"
)

func echoTCPServer() error {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		io.Copy(loggingWriter{conn}, conn)
	}

	// Echo through the loggingWriter
	return err
}

func tcpClientMain() error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		// Send data to the server.
		fmt.Printf("Client: Sending '%s'\n", message)
		n, err := conn.Write([]byte(message))
		if err != nil {
			return err
		}

		// Read data from the server.
		buf := make([]byte, n)
		_, err = conn.Read(buf)
		if err != nil {
			return err
		}
		fmt.Printf("Client: Got '%s'\n", buf)
	}

	return nil
}

func GenTCP() {
	go func() {
		log.Fatal(echoTCPServer())
	}()

	err := tcpClientMain()
	if err != nil {
		panic(err)
	}
}
