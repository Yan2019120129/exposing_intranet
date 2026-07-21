package penetrate

import (
	"errors"
	"io"
	transport "my-base/code/transport"
	"net"
	"time"

	"github.com/google/uuid"
)

// getKey 生成用于标识连接的 UUID。
func getKey() string {
	return uuid.NewString()
}

// Communication 在两个连接间进行无限期双向数据转发。
func Communication(conn, target *transport.Conn) error {
	return CommunicationWithTimeout(conn, target, 0)
}

// CommunicationWithTimeout 在两个连接间进行带空闲超时的双向数据转发。
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

// copyHalfClose 将源连接数据复制到目标连接，并在完成后关闭目标写端。
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

// deadlineReader 在每次读取前设置读超时。
type deadlineReader struct {
	// conn 是需要设置读截止时间的连接。
	conn net.Conn
	// timeout 是单次读取允许空闲的最长时间。
	timeout time.Duration
}

// Read 设置读截止时间后从底层连接读取数据。
func (r deadlineReader) Read(p []byte) (int, error) {
	if err := r.conn.SetReadDeadline(time.Now().Add(r.timeout)); err != nil {
		return 0, err
	}
	return r.conn.Read(p)
}

// deadlineWriter 在每次写入前设置写超时。
type deadlineWriter struct {
	// conn 是需要设置写截止时间的连接。
	conn net.Conn
	// timeout 是单次写入允许阻塞的最长时间。
	timeout time.Duration
}

// Write 设置写截止时间后向底层连接写入数据。
func (w deadlineWriter) Write(p []byte) (int, error) {
	if err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout)); err != nil {
		return 0, err
	}
	return w.conn.Write(p)
}

// isNormalCopyEnd 判断复制操作是否因正常 EOF 或连接关闭而结束。
func isNormalCopyEnd(err error) bool {
	if err == nil || errors.Is(err, io.EOF) {
		return true
	}
	return errors.Is(err, net.ErrClosed)
}
