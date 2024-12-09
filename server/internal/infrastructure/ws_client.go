package infrastructure

import (
	"fmt"
	appmodels "server/internal/app_models"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       int64  `json:"id"`
	RoomID   int64  `json:"roomid"`
	Username string `json:"username"`
}

type Message struct {
	Content    string `json:"content"`
	RoomID     int64  `json:"roomid"`
	Username   string `json:"username"`
	FromUserID int64  `json:"-"`
	ToUserID   int64  `json:"-"`
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
			Content:    string(m),
			RoomID:     c.RoomID,
			Username:   c.Username,
			FromUserID: c.ID,
		}
		req := &appmodels.AddMessageReq{
			Content:  msg.Content,
			RoomID:   msg.RoomID,
			Username: msg.Username,
			UserID:   msg.FromUserID,
		}
		fmt.Println("Sending message to SaveQue:", req)
		hub.SaveQue <- req
		fmt.Println("Message sent to SaveQue")
		hub.Broadcast <- msg
		fmt.Println("Message broadcasted to clients")
	}
}
