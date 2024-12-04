package hub

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan Message
	ID       int64  `json:"id"`
	RoomID   int64  `json:"roomid"`
	Username string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   int64  `json:"roomid"`
	Username string `json:"username"`
}
