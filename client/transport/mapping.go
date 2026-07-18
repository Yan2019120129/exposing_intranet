package transport

import (
	"log"
	"net"
	"time"

	"my-base/code/message"
	toolsPenetrate "my-base/code/transport"
)

// NewMapping creates a data connection to the requested local address.
func (c *Client) NewMapping(msg message.Message) *Client {
	local, ok := msg.Msg.(string)
	if !ok || local == "" {
		log.Println("invalid mapping address")
		return c
	}
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		log.Println("dial control mapping connection:", err)
		return c
	}
	handConn := toolsPenetrate.NewConn(conn)
	if err := handConn.Send(message.Message{
		Name:         msg.Name,
		Symbol:       msg.Symbol,
		TargetSymbol: msg.TargetSymbol,
		Type:         message.MsgTypeLink,
	}); err != nil {
		log.Println("send mapping handshake:", err)
		_ = conn.Close()
		return c
	}
	if err := conn.SetDeadline(time.Time{}); err != nil {
		log.Println("clear mapping deadline:", err)
		_ = conn.Close()
		return c
	}
	targetConn, err := net.Dial("tcp", local)
	if err != nil {
		log.Println("dial local address:", err)
		_ = conn.Close()
		return c
	}
	log.Println("new mapping success:", local)
	c.Copy(targetConn, conn)
	return c
}
