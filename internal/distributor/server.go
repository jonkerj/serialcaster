package distributor

import (
	"log/slog"
	"net"
	"os"
)

// ServeTCP listens on addr, accepts incoming TCP connections, and spawns a
// goroutine per connection that registers with h and forwards broadcast frames.
// It exits the process if the listener cannot be established.
func ServeTCP(h *Hub, addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("listen failed", "addr", addr, "err", err)
		os.Exit(1)
	}
	slog.Info("TCP server listening", "addr", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("accept failed", "err", err)
			continue
		}
		go newClient(conn).serve(h)
	}
}
