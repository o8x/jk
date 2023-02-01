package http2

import (
	"io"
	"net/http"
)

type FlushWriter struct {
	Writer io.Writer
}

func (fw FlushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.Writer.Write(p)
	if f, ok := fw.Writer.(http.Flusher); ok {
		f.Flush()
	}
	return
}
