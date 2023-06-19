package application

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GenHTTP() {
	mux := newRouter()
	srv := &http.Server{
		Addr:    ":4242",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	time.Sleep(time.Second)

	for i := 0; i < 5; i++ {
		resp, err := http.Get("http://" + addr)
		if err != nil {
			panic(err)
		}

		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("buf:", string(buf))
	}

	srv.Close()
}
