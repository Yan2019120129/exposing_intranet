package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Ws struct {
	header http.Header
	dialer *websocket.Dialer
	conn   *websocket.Conn
	url    string
	params url.Values
}

func NewWs(u string, arg ...map[string]any) *Ws {
	ws := &Ws{
		url: u,
		dialer: &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
		},
	}

	for _, a := range arg {
		if a == nil {
			continue
		}
		if v, ok := a["proxy"]; ok {
			u, _ = v.(string)
			proxyURL, _ := url.Parse(
				u,
			)
			ws.dialer.Proxy = http.ProxyURL(proxyURL)
		}

	}
	return ws
}

func (i *Ws) Run() *Ws {
	if i.params != nil {
		i.url += "?" + i.params.Encode()
	}

	conn, resp, err := i.dialer.Dial(i.url, i.header)
	if err != nil {
		if resp != nil {
			fmt.Println("status:", resp.Status)
		}
		panic(err)
	}

	i.conn = conn

	return i
}

// AddParam 添加请求参数
func (i *Ws) AddParam(key string, val string) *Ws {
	if i.params == nil {
		i.params = make(url.Values)
	}
	i.params.Set(key, val)
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
			return
		}
		fu(message)
	}
}
