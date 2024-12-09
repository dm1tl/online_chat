package infrastructure

import "github.com/gin-gonic/gin"

type Router struct {
	handler   *Handler
	wsHandler *WSHandler
}

func NewRouter(handler *Handler, wsHandler *WSHandler) *Router {
	return &Router{
		handler:   handler,
		wsHandler: wsHandler,
	}
}

func (r *Router) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", r.handler.signUp)
		auth.POST("/sign-in", r.handler.signIn)
	}
	ws := router.Group("/ws", r.handler.userIdentity)
	{
		ws.POST("/createRoom", r.wsHandler.CreateRoom)
		ws.GET("/joinRoom/:roomID", r.wsHandler.JoinRoom)
		ws.GET("/getRooms", r.wsHandler.GetRooms)
	}
	return router
}
