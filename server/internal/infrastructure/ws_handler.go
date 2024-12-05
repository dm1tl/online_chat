package infrastructure

import (
	"context"
	"net/http"
	appmodels "server/internal/app_models"
	"server/internal/services"
	"server/internal/utils/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WSHandler struct {
	service *services.Service
	hub     *Hub
}

func NewWSHandler(service *services.Service, hub *Hub) *WSHandler {
	return &WSHandler{
		service: service,
		hub:     hub,
	}
}

func (w *WSHandler) CreateRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var input appmodels.CreateRoomReq
	if err := c.BindJSON(&input); err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	id, err := w.service.RoomManager.CreateRoom(ctx, input)
	if err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "couldn't create room")
		return
	}

	w.hub.Rooms[id] = &Room{
		ID:      id,
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	roomID, err := strconv.ParseInt(c.Param("roomID"), 10, 64)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "incorrect room id")
		return
	}

	clientID, err := w.getUserId(c)
	if err != nil {
		logrus.Error("userId", err)
		response.NewErrorResponse(c, http.StatusBadRequest, "incorrect client id")
		return
	}
	//todo get username from cookie
	username := c.Query("username")
	if username == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "username is required")
		return
	}

	password := c.Query("password")

	input := &appmodels.AddClientReq{
		RoomID:   roomID,
		ClientID: clientID,
		Username: username,
		Password: password,
	}
	ok, err := w.service.RoomManager.GetRoom(ctx, *input)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"room_id": input.RoomID,
			"user_id": input.ClientID,
		}).Errorf("failed to join room: %v", err)

		response.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if !ok {
		logrus.WithFields(logrus.Fields{
			"room_id": input.RoomID,
			"user_id": input.ClientID,
		}).Warn("incorrect password provided")

		response.NewErrorResponse(c, http.StatusUnauthorized, "incorrect password")
		return
	}

	if err := w.service.RoomManager.AddClient(ctx, *input); err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "error while joining the room")
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
