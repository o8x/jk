package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/o8x/jk/v2/logger"
	"github.com/o8x/jk/v2/uniqid"
)

var DefaultClient = &http.Client{}

var DefaultHeader = map[string]string{
	"User-Agent": "jk-http-client/1.7.1",
}

func Get(url string, headers ...map[string]string) ([]byte, error) {
	var h map[string]string
	if headers != nil {
		h = headers[0]
	}

	return Raw(http.MethodGet, url, nil, h)
}

func Post(url string, data any, headers ...map[string]string) ([]byte, error) {
	var h map[string]string
	if headers != nil {
		h = headers[0]
	}

	marshal, _ := json.Marshal(data)
	return Raw(http.MethodPost, url, marshal, h)
}

func Raw(method, url string, data []byte, headers map[string]string) ([]byte, error) {
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
	logger.WithError(err).
		WithField("request id", rid).
		WithField("data", string(all)).
		Debug("http response")
	return all, err
}
