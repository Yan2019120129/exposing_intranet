package utils

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Ws struct {
	header http.Header
	dialer *websocket.Dialer
	conn   *websocket.Conn
	url    string
	Err    error
}

func NewWs(u string) *Ws {
	return &Ws{
		url:    u,
		dialer: websocket.DefaultDialer,
	}
}

func (i *Ws) Run() *Ws {
	var err error
	i.conn, _, err = i.dialer.Dial(i.url, i.header)
	if err != nil {
		i.Err = err
	}
	return i
}

// Set 设置请求头信息
func (i *Ws) Set(key, val string) *Ws {
	if i.header == nil {
		i.header = make(http.Header)
	}
	i.header.Set(key, val)
	return i
}

func (i *Ws) Send(s string) *Ws {
	err := i.conn.WriteMessage(websocket.TextMessage, []byte(s))
	if err != nil {
		return i
	}
	return i
}

func (i *Ws) Ping(t time.Duration, s string) *Ws {
	go func() {
		for {
			err := i.conn.WriteMessage(websocket.PingMessage, []byte(s))
			if err != nil {
				_ = i.conn.Close()
				return
			}
			time.Sleep(t * time.Second)
		}
	}()
	return i
}

func (i *Ws) Read(fu func([]byte)) {
	for {
		_, message, err := i.conn.ReadMessage()
		if err != nil {
			_ = i.conn.Close()
			i.Err = err
			return
		}
		fu(message)
	}
}
