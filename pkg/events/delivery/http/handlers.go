package http

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sparrowganz/TestTask-events/pkg/events"
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
	router.POST("/events", h.SetEvents)
}

// SetEvents POST /api/events
func (h *Handler) SetEvents(c *gin.Context) {
	defer c.Request.Body.Close()

	var err error

	scanner := bufio.NewScanner(c.Request.Body)
	for scanner.Scan() {
		//todo продумать что будет если отправить сериализованные данные
		err = h.service.SendEvent(scanner.Bytes())
		if err != nil {
			h.logger.Println("(ERROR) ", errors.Wrap(err, "failed send event"))
		}
	}
	if err := scanner.Err(); err != nil {
		h.logger.Println("(ERROR) ", errors.Wrap(err, "failed scan"))
		c.Status(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
}
