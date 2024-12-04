package hub

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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

func (c *Client) writeMessages() {
	defer func() {
		err := c.Conn.Close()
		logrus.Error("writeMessage", err)
	}()
	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessages(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		err := c.Conn.Close()
		logrus.Error("readMessage", err)
	}()
	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Error(err)
			}
			logrus.Error(err)
			break
		}
		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}
		hub.Broadcast <- msg
	}
}
