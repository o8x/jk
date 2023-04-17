package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/o8x/jk/v2/logger"
	"github.com/o8x/jk/v2/uniqid"
)

var DefaultClient = &http.Client{}

type Response struct {
	*http.Response `json:"base"`
	StatusCode     int    `json:"status_code"`
	Body           []byte `json:"body"`
}

var DefaultHeader = map[string]string{
	"User-Agent": "jk-http-client/1.7.1",
}

func Get(url string, headers ...map[string]string) (*Response, error) {
	var h map[string]string
	if headers != nil {
		h = headers[0]
	}

	return Raw(http.MethodGet, url, nil, h)
}

func Output(url string, w io.Writer, headers ...map[string]string) (int, error) {
	var h map[string]string
	if headers != nil {
		h = headers[0]
	}

	raw, err := Raw(http.MethodGet, url, nil, h)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(raw.Body)
	return n, err
}

func Wget(url string, name string, headers ...map[string]string) (int, error) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return 0, err
	}

	return Output(url, f, headers...)
}

func Post(url string, data any, headers ...map[string]string) (*Response, error) {
	var h map[string]string
	if headers != nil {
		h = headers[0]
	}

	marshal, _ := json.Marshal(data)
	return Raw(http.MethodPost, url, marshal, h)
}

func Raw(method, url string, data []byte, headers map[string]string) (*Response, error) {
	rid := uniqid.String()
	logger.
		WithField("request id", rid).
		WithField("method", method).
		WithField("url", url).
		WithField("with data", string(data)).
		WithField("with header", headers).
		Debug("send http request")

	r, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		r.Header.Set(k, v)
	}

	// 设置默认 Header
	for k, v := range DefaultHeader {
		if _, ok := headers[k]; ok {
			continue
		}
		r.Header.Set(k, v)
	}

	resp, err := DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(resp.Body)

	var v string
	if len(all) < 1000 {
		v = string(all)
	}

	logger.WithError(err).
		WithField("request id", rid).
		WithField("data", v).
		WithField("length", len(all)).
		Debug("http response")

	return &Response{
		Response: resp,
		Body:     all,
	}, err
}
