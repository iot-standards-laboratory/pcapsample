package transport

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/quic-go/quic-go"
)



// Start a server that echos all data on the first stream opened by the client
func echoQUICServer() error {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	conn, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}
		go func() {
			io.Copy(loggingWriter{stream}, stream)
		}()
	}
	// Echo through the loggingWriter
	return err
}

func quicClientMain() error {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	conn, err := quic.DialAddr(context.Background(), addr, tlsConf, nil)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			return err
		}

		fmt.Printf("Client: Sending '%s'\n", message)
		_, err = stream.Write([]byte(message))
		if err != nil {
			return err
		}

		buf := make([]byte, len(message))
		_, err = io.ReadFull(stream, buf)
		if err != nil {
			return err
		}
		fmt.Printf("Client: Got '%s'\n", buf)
	}

	return nil
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}

func GenQUIC() {
	go func() { log.Fatal(echoQUICServer()) }()

	err := quicClientMain()
	if err != nil {
		panic(err)
	}
}
