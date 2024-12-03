package handler

import (
	"context"
	"net/http"
	appmodels "server/internal/app_models"
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
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	err := h.service.UserManager.Create(ctx, input)
	if err != nil {
		logrus.Error(err)
		NewErrorResponse(c, http.StatusInternalServerError, "couldn't create an account, try again")
		return
	}
	c.JSON(http.StatusOK, NewStatusResponse("you succesfully signed up!"))
}

func (h *Handler) signIn(c *gin.Context) {
	var input appmodels.LoginReq
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	token, err := h.service.UserManager.Login(ctx, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}