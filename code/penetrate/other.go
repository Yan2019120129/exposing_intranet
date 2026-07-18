package penetrate

import (
	"errors"
	"io"
	transport "my-base/code/transport"
	"net"
	"time"

	"github.com/google/uuid"
)

// getKey 获取UUID键
func getKey() string {
	return uuid.NewString()
}

// Communication 通信
func Communication(conn, target *transport.Conn) error {
	return CommunicationWithTimeout(conn, target, 0)
}

func CommunicationWithTimeout(conn, target *transport.Conn, timeout time.Duration) error {
	conn.SetStatusActive()
	target.SetStatusActive()
	defer conn.Close()
	defer target.Close()

	results := make(chan error, 2)
	go func() {
		results <- copyHalfClose(target.GetConn(), conn.GetConn(), timeout)
	}()
	go func() {
		results <- copyHalfClose(conn.GetConn(), target.GetConn(), timeout)
	}()

	first := <-results
	if first != nil && !isNormalCopyEnd(first) {
		_ = conn.Close()
		_ = target.Close()
	}
	second := <-results
	if !isNormalCopyEnd(first) {
		return first
	}
	if !isNormalCopyEnd(second) {
		return second
	}
	return nil
}

func copyHalfClose(dst, src net.Conn, timeout time.Duration) error {
	var reader io.Reader = src
	var writer io.Writer = dst
	if timeout > 0 {
		reader = deadlineReader{conn: src, timeout: timeout}
		writer = deadlineWriter{conn: dst, timeout: timeout}
	}
	_, err := io.Copy(writer, reader)
	if tcpConn, ok := dst.(*net.TCPConn); ok {
		_ = tcpConn.CloseWrite()
	}
	return err
}

type deadlineReader struct {
	conn    net.Conn
	timeout time.Duration
}

func (r deadlineReader) Read(p []byte) (int, error) {
	if err := r.conn.SetReadDeadline(time.Now().Add(r.timeout)); err != nil {
		return 0, err
	}
	return r.conn.Read(p)
}

type deadlineWriter struct {
	conn    net.Conn
	timeout time.Duration
}

func (w deadlineWriter) Write(p []byte) (int, error) {
	if err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout)); err != nil {
		return 0, err
	}
	return w.conn.Write(p)
}

func isNormalCopyEnd(err error) bool {
	if err == nil || errors.Is(err, io.EOF) {
		return true
	}
	return errors.Is(err, net.ErrClosed)
}
