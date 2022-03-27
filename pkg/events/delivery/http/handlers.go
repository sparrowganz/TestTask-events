package http

import (
	"TestTask-events/pkg/events"
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	//split body by end close tags
	//IMPORTANT: input data with structures become invalid
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		var startIndex int
		//Remove all splits and spaces
		for i := 0; i < len(data); i++ {
			if data[i] != ' ' && data[i] != '\n' {
				break
			}

			startIndex = i + 1
		}

		for i := startIndex; i < len(data)-1; i++ {
			if data[i] == '}' {
				return i + 2, data[startIndex : i+1], nil
			}
		}
		if !atEOF {
			return 0, nil, nil
		}

		return 0, data[startIndex:], bufio.ErrFinalToken
	})
	for scanner.Scan() {

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
