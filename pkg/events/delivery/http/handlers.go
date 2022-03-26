package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sparrowganz/TestTask-events/pkg/message"
)

type Handler struct {
	broker message.Broker
}

func New(broker message.Broker) *Handler {
	return &Handler{
		broker: broker,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/api/events", h.SetEvents)
}

// SetEvents POST /api/events
func (h *Handler) SetEvents(ctx *gin.Context) {
	//todo implement
}
