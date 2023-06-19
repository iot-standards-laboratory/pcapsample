package application

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/plgd-dev/go-coap/v3"
	coapMessage "github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	"github.com/plgd-dev/go-coap/v3/udp"
)

func coapServer() error {
	r := mux.NewRouter()

	r.DefaultHandleFunc(func(w mux.ResponseWriter, r *mux.Message) {
		err := w.SetResponse(codes.Content, coapMessage.TextPlain, bytes.NewReader([]byte(message)))
		if err != nil {
			log.Printf("cannot set response: %v", err)
		}
	})

	err := coap.ListenAndServe("udp", ":5683", r)

	return err
}

func coapClientMain() error {
	co, err := udp.Dial("localhost:5683")
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		msg, err := co.NewGetRequest(
			context.Background(),
			"/",
		)

		if err != nil {
			return err
		}

		resp, err := co.Do(msg)
		if err != nil {
			return err
		}

		b, err := io.ReadAll(resp.Body())
		if err != nil {
			return err
		}
		fmt.Printf("Client: Got '%s'\n", string(b))
	}

	return nil
}

func GenCoAP() {
	go func() {
		err := coapServer()
		if err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(time.Second)
	err := coapClientMain()
	if err != nil {
		log.Fatalln(err)
	}
}
