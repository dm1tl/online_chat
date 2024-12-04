package hub

type WSHandler struct {
	hub *Hub
}

func NewWSHandler(hub *Hub) *WSHandler {
	return &WSHandler{
		hub: hub,
	}
}

type createRoomReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
