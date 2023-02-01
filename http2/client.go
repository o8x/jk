package http2

import (
	"crypto/tls"
	"net/http"

	"github.com/o8x/jk/cert"
	"golang.org/x/net/http2"
)

func NewClient(f *cert.Folder) (*http.Client, error) {
	pool, err := f.GetCertPool()
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: pool,
			},
		},
	}, nil
}
