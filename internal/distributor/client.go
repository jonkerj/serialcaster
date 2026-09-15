package distributor

import (
	"net"
	"time"
)

// client represents a single connected TCP consumer.
type client struct {
	conn net.Conn
	ch   chan []byte
}

// newClient creates a client for conn with a buffered outbound channel.
func newClient(conn net.Conn) *client {
	return &client{
		conn: conn,
		ch:   make(chan []byte, 64),
	}
}

// addr returns the remote address of the underlying connection.
func (c *client) addr() string { return c.conn.RemoteAddr().String() }

// send enqueues data for delivery. Returns false and drops the frame if the
// buffer is full (slow consumer).
func (c *client) send(data []byte) bool {
	select {
	case c.ch <- data:
		return true
	default:
		return false
	}
}

// serve registers with h, forwards outbound frames to conn, and unregisters
// on exit. Intended to run in its own goroutine.
func (c *client) serve(h *Hub) {
	if tc, ok := c.conn.(*net.TCPConn); ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(30 * time.Second)
	}
	h.register(c)
	defer func() {
		h.unregister(c)
		c.conn.Close()
	}()
	for data := range c.ch {
		if _, err := c.conn.Write(data); err != nil {
			return
		}
	}
}
