package web

import (
	"context"
	"fmt"
	"github.com/sparrowganz/TestTask-events/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sparrowganz/TestTask-events/config"
	pkgEventsRouter "github.com/sparrowganz/TestTask-events/pkg/events/delivery/http"
	pkgEventService "github.com/sparrowganz/TestTask-events/pkg/events/service"
)

type Main interface {
	Start()
}

type mainData struct {
	server http.Server
	core   app.Core
}

func NewMain(config *config.Data, core app.Core) Main {
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

	//Init router
	router := gin.Default()
	api := router.Group("/api")

	//Init workers Pipeline
	eventService, err := pkgEventService.New(m.core)
	if err != nil {
		m.core.Logger().Fatal("Failed start eventService")
	}

	//Init events handler
	eventHandler := pkgEventsRouter.New(m.core.Logger(), eventService)
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
