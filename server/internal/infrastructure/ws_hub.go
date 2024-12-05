package infrastructure

type Room struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Clients  map[int64]*Client
}

type Hub struct {
	Rooms      map[int64]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[int64]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if r, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			if r, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := r.Clients[cl.ID]; ok {
					msg := &Message{
						Content:  "user " + cl.Username + " has left the room",
						RoomID:   cl.RoomID,
						Username: cl.Username,
					}
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						h.Broadcast <- msg
					}
					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}
		case msg := <-h.Broadcast:
			if _, ok := h.Rooms[msg.RoomID]; ok {
				for _, cl := range h.Rooms[msg.RoomID].Clients {
					cl.Message <- *msg
				}
			}
		}
	}

}
