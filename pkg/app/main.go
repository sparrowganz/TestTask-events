package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sparrowganz/TestTask-events/config"
	pkgEventsRouter "github.com/sparrowganz/TestTask-events/pkg/events/delivery/http"
	pkgKafka "github.com/sparrowganz/TestTask-events/pkg/message/kafka"
)

type Main interface {
	Start()
}

type mainData struct {
	server http.Server
	core   Core
}

func NewMain(config *config.Data, core Core) Main {
	return &mainData{
		server: http.Server{
			Addr:        fmt.Sprintf(":%d", config.Server.Port),
			ReadTimeout: config.Server.ReadTimeOut,
			ErrorLog:    core.Logger(),
		},
		core: core,
	}
}

func (m *mainData) Start() {

	router := gin.Default()
	api := router.Group("/api")

	messageBroker := pkgKafka.New()

	//Init events handler
	eventHandler := pkgEventsRouter.New(messageBroker)
	eventHandler.RegisterRoutes(api)

	//Start listening server
	m.core.Group().Go(func() error {
		m.server.Handler = router
		return m.server.ListenAndServe()
	})
	//Stop listening server
	m.core.Group().Go(func() error {
		<-m.core.Context().Done()
		return m.server.Shutdown(context.Background())
	})
}
