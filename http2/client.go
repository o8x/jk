package http2

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/net/http2"

	"github.com/o8x/jk/v2/cert"
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
