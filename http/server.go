package http

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/o8x/jk/v2/response"
	"github.com/o8x/jk/v2/x"
)

type Handler func(Request) *response.Response

type Request struct {
	*http.Request
	Query url.Values
}

func (r *Request) Unmarshal(v any) error {
	all, err := io.ReadAll(r.Request.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(all, v)
}

func (r *Request) ReadBody() ([]byte, error) {
	return io.ReadAll(r.Request.Body)
}

func (r *Request) ReadBodyReckless() []byte {
	all, _ := io.ReadAll(r.Request.Body)
	return all
}

func (r *Request) RemoteHost() string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func (r *Request) RemoteHostIn(list ...string) bool {
	host := r.RemoteHost()

	for _, it := range list {
		if it == host {
			return true
		}
	}

	return false
}

func (r *Request) RemotePort() int {
	_, port, _ := net.SplitHostPort(r.RemoteAddr)
	return x.ParseInt(port, 0)
}

type Mux struct {
	mux *http.ServeMux
}

func NewMux() *Mux {
	return &Mux{
		mux: http.NewServeMux(),
	}
}

func (m *Mux) ListenAndServe(listen string) error {
	return ListenAndServe(listen, m)
}

func (m *Mux) Post(name string, fn Handler) {
	m.RegisterRoute(http.MethodPost, name, fn)
}

func (m *Mux) Get(name string, fn Handler) {
	m.RegisterRoute(http.MethodGet, name, fn)
}

func (m *Mux) Put(name string, fn Handler) {
	m.RegisterRoute(http.MethodPut, name, fn)
}

func (m *Mux) Trace(name string) {
	m.mux.HandleFunc(name, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodTrace {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		all, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, _ = w.Write(all)
	})
}

func (m *Mux) Delete(name string, fn Handler) {
	m.RegisterRoute(http.MethodDelete, name, fn)
}

func (m *Mux) Any(name string, fn Handler) {
	m.RegisterRoute("any", name, fn)
}

func (m *Mux) RegisterRoute(method, name string, fn Handler) {
	m.mux.HandleFunc(name, func(w http.ResponseWriter, r *http.Request) {
		if method != "any" && r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		resp := fn(Request{
			Request: r,
			Query:   r.URL.Query(),
		})
		if resp == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(resp.Dump())
	})
}

func ListenAndServe(listen string, mux *Mux) error {
	srv := &http.Server{
		Addr:    listen,
		Handler: mux.mux,
	}

	return srv.ListenAndServe()
}

func PostServer(listen, path string, fn Handler) error {
	mux := NewMux()
	mux.Post(path, fn)

	return ListenAndServe(listen, mux)
}

func GetServer(listen, path string, fn Handler) error {
	mux := NewMux()
	mux.Get(path, fn)

	return ListenAndServe(listen, mux)
}

func AnyServer(listen, path string, fn Handler) error {
	mux := NewMux()
	mux.Any(path, fn)

	return ListenAndServe(listen, mux)
}
