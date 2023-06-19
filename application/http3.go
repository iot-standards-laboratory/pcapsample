package application

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
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

	var err error
	certs := make([]tls.Certificate, 1)
	certs[0], err = tls.LoadX509KeyPair("./assets/cert.pem", "./assets/priv.key")
	if err != nil {
		panic(err)
	}

	srv := &http3.Server{
		Addr:    ":4242",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServeTLS("./assets/cert.pem", "./assets/priv.key"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	time.Sleep(time.Second)

	http3ClientMain()

}

func http3ClientMain() {
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
			// KeyLogWriter:       keyLog,
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
			panic(err)
		}

		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("buf:", string(buf))
	}
}
