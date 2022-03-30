package http

import (
	"TestTask-events/pkg/events"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	logger  *log.Logger
	service events.Service
}

func New(logger *log.Logger, service events.Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/events", h.setEventsHandler)
}

// SetEvents POST /api/events
func (h *Handler) setEventsHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	decoder := json.NewDecoder(c.Request.Body)
	var (
		err error
	)

	for decoder.More() {
		in := &events.Event{}
		err = decoder.Decode(in)
		if err != nil {
			log.Printf("[ERROR] failed parse request: %s\n", err.Error())
			continue
		}
		err = h.service.SendEvent(in)
		if err != nil {
			log.Printf("[ERROR] failed send event: %s\n", err.Error())
			continue
		}
	}

	c.Status(http.StatusOK)
}
