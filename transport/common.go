package transport

import (
	"fmt"
	"io"
)

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

const addr = "localhost:4242"
const message = "Computer Networks Packet Capture"

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}
