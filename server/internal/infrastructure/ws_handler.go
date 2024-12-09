package infrastructure

import (
	"context"
	"fmt"
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
		Clients: make(map[int64]*Client),
	}
	c.JSON(http.StatusOK, response.NewStatusResponse("you succesfully created a room"))
}

func (w *WSHandler) GetRooms(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	rooms, err := w.service.GetAllRooms(ctx)
	if err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "couldn't get all rooms, try again")
		return
	}
	for _, room := range rooms {
		w.hub.Rooms[room.ID] = &Room{
			ID:      room.ID,
			Name:    room.Name,
			Clients: make(map[int64]*Client),
		}
	}
	c.JSON(http.StatusOK, rooms)
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

	if err := w.service.ClientManager.AddClient(ctx, *input); err != nil {
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

	go func() {
		fmt.Println("Starting message processing")
		for msg := range w.hub.ProcessedQue {
			err := w.service.MessageManager.AddMessage(context.Background(), *msg)
			if err != nil {
				logrus.Errorf("Failed to save message: %v", err)
				sendErr := conn.WriteJSON(response.NewStatusResponse("failed to save message"))
				if sendErr != nil {
					logrus.Errorf("Failed to send error message to client: %v", sendErr)
				}
				continue
			}
			fmt.Println("Message processed and saved")
		}
	}()

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message),
		RoomID:   roomID,
		ID:       clientID,
		Username: username,
	}

	msg := &Message{
		Content:  "user " + username + " has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	//var messages appmodels.BackupMessages
	//for _, msg := range messages[roomID] {
	//	fmt.Print(msg)
	//	msgToClient := &Message{
	//		Content:  "testingBackup",
	//		RoomID:   int64(1),
	//		Username: "dima",
	//	}
	//	w.hub.Recover <- msgToClient
	//}
	//	for _, msg := range messages {
	//		if msg.RoomID == roomID {
	//			msgToClient := &Message{
	//				Content:  msg.Content,
	//				RoomID:   msg.RoomID,
	//				Username: msg.Username,
	//			}
	//			cl.Message <- *msgToClient
	//		}
	//	}

	w.hub.Register <- cl
	w.hub.Broadcast <- msg

	msgToClient := &Message{
		Content:  "testingBackup",
		RoomID:   int64(1),
		Username: "dima",
		UserID:   int64(1),
	}
	fmt.Println("rooms", w.hub.Rooms)
	fmt.Println("clients", w.hub.Rooms[roomID].Clients)
	w.hub.Recover <- msgToClient

	go cl.writeMessages()
	cl.readMessages(w.hub)
}

//TODO implement methods for saving messages in database, implement methods for state recovery
