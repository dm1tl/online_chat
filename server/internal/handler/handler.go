package handler

import (
	"server/internal/hub"
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service   *services.Service
	wsHandler *hub.WSHandler
}

func NewHandler(service *services.Service, wsHandler *hub.WSHandler) *Handler {
	return &Handler{
		service:   service,
		wsHandler: wsHandler,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	ws := router.Group("/ws")
	{
		ws.POST("/createRoom", h.wsHandler.CreateRoom)
		ws.GET("/joinRoom/:roomid", h.wsHandler.JoinRoom)
	}
	return router
}
