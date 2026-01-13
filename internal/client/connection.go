package client

import (
	"fmt"
	"net"
	"time"
)

type Connection struct {
	conn    net.Conn
	timeout time.Duration
}

// Dial opens a TCP connection to the Modbus target.
func Dial(address string, timeout time.Duration) (*Connection, error) {
	dialer := net.Dialer{
		Timeout:   timeout,
		KeepAlive: 30 * time.Second,
	}

	c, err := dialer.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	// Disable Nagle for lower latency
	if tcp, ok := c.(*net.TCPConn); ok {
		_ = tcp.SetNoDelay(true)
		_ = tcp.SetKeepAlive(true)
		_ = tcp.SetKeepAlivePeriod(30 * time.Second)
	}

	return &Connection{
		conn:    c,
		timeout: timeout,
	}, nil
}

// Close closes the underlying TCP connection.
func (c *Connection) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Write sends raw bytes to the socket.
func (c *Connection) Write(b []byte) error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
		return err
	}
	_, err := c.conn.Write(b)
	return err
}

// ReadFull reads exactly len(b) bytes.
func (c *Connection) ReadFull(b []byte) error {
	if err := c.conn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
		return err
	}

	n := 0
	for n < len(b) {
		r, err := c.conn.Read(b[n:])
		if err != nil {
			return err
		}
		n += r
	}
	return nil
}
