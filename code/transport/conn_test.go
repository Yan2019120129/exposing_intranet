package transport

import (
	"encoding/json"
	"my-base/code/message"
	"net"
	"sync"
	"testing"
)

func TestConnReusesDecoderAcrossMessages(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go func() {
		encoder := json.NewEncoder(server)
		_ = encoder.Encode(message.Message{Type: message.MsgTypePing, Symbol: "first"})
		_ = encoder.Encode(message.Message{Type: message.MsgTypePong, Symbol: "second"})
	}()

	conn := NewConn(client)
	var first, second message.Message
	if err := conn.ParseMsg(&first); err != nil {
		t.Fatalf("parse first message: %v", err)
	}
	if err := conn.ParseMsg(&second); err != nil {
		t.Fatalf("parse second message: %v", err)
	}
	if first.Symbol != "first" || second.Symbol != "second" {
		t.Fatalf("messages were not preserved: first=%+v second=%+v", first, second)
	}
}

func TestConnSerializesConcurrentSends(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	const count = 32
	readDone := make(chan error, 1)
	go func() {
		decoder := json.NewDecoder(server)
		for i := 0; i < count; i++ {
			var msg message.Message
			if err := decoder.Decode(&msg); err != nil {
				readDone <- err
				return
			}
		}
		readDone <- nil
	}()

	conn := NewConn(client)
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := conn.Send(message.Message{Type: message.MsgTypeDefault, Symbol: "message"}); err != nil {
				t.Errorf("send %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()
	if err := <-readDone; err != nil {
		t.Fatalf("decode concurrent sends: %v", err)
	}
}
