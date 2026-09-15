package distributor

import "log/slog"

// Hub distributes data to all registered clients. It must be driven by a
// single call to Run, which serialises all mutations to the client set.
type Hub struct {
	clients  map[*client]struct{}
	joins    chan *client
	leaves   chan *client
	outbound chan []byte
}

// NewHub returns an initialised Hub ready to run.
func NewHub() *Hub {
	return &Hub{
		clients:  make(map[*client]struct{}),
		joins:    make(chan *client, 32),
		leaves:   make(chan *client, 32),
		outbound: make(chan []byte, 256),
	}
}

func (h *Hub) register(c *client)   { h.joins <- c }
func (h *Hub) unregister(c *client) { h.leaves <- c }

// Broadcast enqueues data for delivery to all currently registered clients.
func (h *Hub) Broadcast(data []byte) { h.outbound <- data }

// Run is the single goroutine that owns h.clients — no locks needed.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.joins:
			h.clients[c] = struct{}{}
			slog.Info("client connected", "addr", c.addr(), "total", len(h.clients))

		case c := <-h.leaves:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.ch)
				slog.Info("client disconnected", "addr", c.addr(), "total", len(h.clients))
			}

		case data := <-h.outbound:
			for c := range h.clients {
				c.send(data)
			}
		}
	}
}
