package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler interface {
	Health(*gin.Context)
}

type healthHandler struct{}

func NewHealthHandler() HealthHandler {
	return &healthHandler{}
}

func (h *healthHandler) Health(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
