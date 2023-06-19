package application

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

var pool *x509.CertPool

func init() {
	var err error
	pool, err = x509.SystemCertPool()
	if err != nil {
		panic(err)
	}
}

func GenHTTP3() {
	mux := newRouter()

	srv := &http3.Server{
		Addr:      ":4242",
		Handler:   mux,
		TLSConfig: generateTLSConfig(),
	}

	go func() {

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	time.Sleep(time.Second)

	err := http3ClientMain()
	if err != nil {
		log.Fatalln(err)
	}

}

func http3ClientMain() error {
	keyLog, err := os.Create("http3_key.log")
	if err != nil {
		return err
	}
	defer keyLog.Close()

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
			KeyLogWriter:       keyLog,
		},
		QuicConfig: &quic.Config{},
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}

	for i := 0; i < 5; i++ {
		resp, err := hclient.Get("https://" + addr)
		if err != nil {
			return err
		}

		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Println("buf:", string(buf))
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
