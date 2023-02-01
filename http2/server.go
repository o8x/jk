package http2

import (
	"net/http"

	"github.com/o8x/jk/cert"
)

func ListenAndServeTLS(listen string, mux *http.ServeMux, folder *cert.Folder) error {
	srv := &http.Server{
		Addr:    listen,
		Handler: mux,
	}

	return srv.ListenAndServeTLS(folder.CrtFile, folder.KeyFile)
}
