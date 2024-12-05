package handler

import (
	"context"
	"net/http"
	appmodels "server/internal/app_models"
	"server/internal/utils/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty header")
		return
	}
	logrus.Info(header)
	userId, err := h.service.Validate(ctx, appmodels.ValidateTokenReq{
		Token: header,
	})
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, "incorrect token value")
	}
	logrus.Info(userId)
	c.Set(userCtx, userId.ID)
	c.Next()
}
