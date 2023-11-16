package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/o8x/jk/v2/logger"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var wss sync.Map

type Conn struct {
	*websocket.Conn
	Name string
}

func GetConn(name string) *Conn {
	value, ok := wss.Load(name)
	if ok {
		return value.(*Conn)
	}
	return nil
}

func Upgrade(name string, w http.ResponseWriter, r *http.Request) error {
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	wss.Store(name, &Conn{
		Conn: ws,
		Name: name,
	})
	return nil
}

type Message struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Result string `json:"result"`
}

func (m Message) Unmarshal(v any) error {
	return json.Unmarshal([]byte(m.Result), v)
}

type HandleFunc func(Message)

func (c *Conn) ReadJSON(fn HandleFunc) {
	for {
		var msg Message

		if err := c.Conn.ReadJSON(&msg); err != nil {
			if err.Error() == "websocket: close 1001 (going away)" {
				c.Conn.Close()
				logger.WithError(err).Error("websocket close")
				return
			}

			logger.WithError(err).Error("read json failed")
			continue
		}

		go fn(msg)
	}
}

func (c *Conn) WriteBytes(message []byte) error {
	return c.Conn.WriteMessage(websocket.TextMessage, message)
}

func (c *Conn) Write(data any) error {
	marshal, _ := json.Marshal(data)
	return c.WriteBytes(marshal)
}

func (c *Conn) ExecJS(code string) error {
	return c.Write(map[string]any{
		"name":   c.Name,
		"method": "exec_js",
		"code":   fmt.Sprintf(`(function (){%s ;})()`, code),
	})
}

func (c *Conn) ExecAsyncJS(code string) error {
	return c.Write(map[string]any{
		"name":   c.Name,
		"method": "exec_js",
		"code":   fmt.Sprintf(`(async function (){%s ;})()`, code),
	})
}
