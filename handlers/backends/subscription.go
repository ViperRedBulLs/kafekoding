package backends

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type subscription struct {
	conn *connection
	room string
	hub  *Hub
}

func (s *subscription) readPump() {
	c := s.conn
	defer func() {
		c.ws.Close()
		s.hub.unregister <- s
	}()
	for {
		message := NewMessage()
		err := c.ws.ReadJSON(message)
		if err != nil {
			return
		}
		message.SetRoom(s.room)
		message.SetTime(time.Now())
		s.hub.broadcast <- message
	}
}

func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(time.Second * 50)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.ws.WriteJSON(message)
			if err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, room string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &connection{
		ws:   ws,
		send: make(chan *Message, 256),
	}

	s := &subscription{
		room: room,
		hub:  hub,
		conn: c,
	}

	s.hub.register <- s

	go s.writePump()
	go s.readPump()
}
