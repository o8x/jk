package jshook

import (
	_ "embed"
	"fmt"
	"strings"

	error2 "github.com/o8x/jk/v2/error"

	"github.com/o8x/jk/v2/http"
	"github.com/o8x/jk/v2/jshook/ws"
	"github.com/o8x/jk/v2/logger"
	"github.com/o8x/jk/v2/rand"
	"github.com/o8x/jk/v2/response"
)

//go:embed js/inject.js
var js string

var script = `(function() {
	let script = document.createElement('script')
	script.type = 'text/javascript'
	script.src = "{{js_url}}"
	document.getElementsByTagName('head')[0].appendChild(script)
})()`

type HookFunc func(*ws.Conn)

type Hook struct {
	Listen string
	Name   string
	Func   HookFunc
}

func (h *Hook) Init() {
	if h.Name == "" {
		h.Name = fmt.Sprintf("%x", rand.Intn(1000))
	}

	if h.Listen == "" {
		h.Listen = "localhost:8080"
	}
}

func (h *Hook) BuildJsURL() string {
	h.Init()
	return fmt.Sprintf("http://%s/js?n=%s", h.Listen, h.Name)
}

func (h *Hook) BuildWSURL() string {
	h.Init()
	return fmt.Sprintf("ws://%s/ws?n=%s", h.Listen, h.Name)
}

func (h *Hook) ListenAndServe() error {
	h.Init()

	m := http.NewMux()
	m.Get("/ws", func(r http.Request) *response.Response {
		if err := ws.Upgrade(h.Name, r.Writer, r.Request); err != nil {
			return response.Error(err)
		}

		h.Func(ws.GetConn(h.Name))
		return response.OK(nil)
	})

	m.Any("/js", func(r http.Request) *response.Response {
		_, err := r.Writer.Write([]byte(strings.ReplaceAll(js, "{{ws_url}}", h.BuildWSURL())))
		error2.Hide(err)

		return nil
	})

	m.Any("/inject", func(r http.Request) *response.Response {
		_, err := r.Writer.Write([]byte(strings.ReplaceAll(script, "{{js_url}}", h.BuildJsURL())))
		error2.Hide(err)
		return nil
	})

	logger.WithField("addr", h.Listen).Info("listen on")
	logger.WithField("path", fmt.Sprintf("http://%s/inject", h.Listen)).Info("get inject script")

	return http.ListenAndServe(h.Listen, m)
}
