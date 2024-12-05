package hub

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	userCtx = "userId"
)

func (h *WSHandler) getUserId(c *gin.Context) (int64, error) {
	id := c.GetInt64(userCtx)
	logrus.Info(id)
	return id, nil
}
