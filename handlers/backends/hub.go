package backends

type Hub struct {
	rooms      map[string]map[*connection]bool
	register   chan *subscription
	unregister chan *subscription
	broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*connection]bool),
		register:   make(chan *subscription),
		unregister: make(chan *subscription),
		broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}
			h.rooms[s.room][s.conn] = true
		case s := <-h.unregister:
			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		case msg := <-h.broadcast:
			connections := h.rooms[msg.Room]
			for c := range connections {
				select {
				case c.send <- msg:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, msg.Room)
					}
				}
			}
		}
	}
}
