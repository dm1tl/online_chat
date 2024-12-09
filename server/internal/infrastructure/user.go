package infrastructure

import (
	"context"
	"net/http"
	appmodels "server/internal/app_models"
	"server/internal/utils/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) signUp(c *gin.Context) {
	var input appmodels.CreateUserReq
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := c.BindJSON(&input); err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	err := h.service.AuthManager.Create(ctx, input)
	if err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusInternalServerError, "couldn't create an account, try again")
		return
	}
	c.JSON(http.StatusOK, response.NewStatusResponse("you succesfully signed up!"))
}

func (h *Handler) signIn(c *gin.Context) {
	var input appmodels.LoginReq
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	token, err := h.service.AuthManager.Login(ctx, input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}
