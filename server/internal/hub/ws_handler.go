package hub

import (
	"net/http"
	"server/internal/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

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

func (w *WSHandler) CreateRoom(c *gin.Context) {
	var input createRoomReq
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	w.hub.Rooms[input.ID] = &Room{
		ID:      input.ID,
		Name:    input.Name,
		Clients: map[int64]*Client{},
	}
	c.JSON(http.StatusOK, response.NewStatusResponse("you succesfully created a room"))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (w *WSHandler) JoinRoom(c *gin.Context) {
	roomID, err := strconv.ParseInt(c.Param("roomID"), 10, 64)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "incorrect room id")
		return
	}

	clientID, err := strconv.ParseInt(c.Query("userID"), 10, 64)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "incorrect client id")
		return
	}

	username := c.Query("username")
	if username == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "username is required")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Error(err)
		conn.WriteJSON(response.ErrorResponse{
			Message: "couldn't estalish connection",
		})
		return
	}
	defer conn.Close()

	cl := &Client{
		Conn:     conn,
		Message:  make(chan Message),
		RoomID:   roomID,
		ID:       clientID,
		Username: username,
	}

	msg := &Message{
		Content:  "user " + username + " has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	w.hub.Register <- cl
	w.hub.Broadcast <- msg
	go cl.writeMessages()
	cl.readMessages(w.hub)
}
