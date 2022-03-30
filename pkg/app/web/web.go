package web

import (
	"TestTask-events/config"
	"TestTask-events/pkg/app"
	"TestTask-events/pkg/db/cache"
	pkgDBClickhouse "TestTask-events/pkg/db/clickhouse"
	pkgEventsRouter "TestTask-events/pkg/events/delivery/http"
	pkgEventRepository "TestTask-events/pkg/events/repository"
	pkgEventService "TestTask-events/pkg/events/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"runtime/debug"
)

type mainData struct {
	server   http.Server
	database *sqlx.DB
	core     app.Core
}

func NewMain(config *config.Data, core app.Core) (app.App, error) {

	db, err := pkgDBClickhouse.New(config.Database)
	if err != nil {
		return nil, err
	}

	return &mainData{
		server: http.Server{
			Addr:        fmt.Sprintf(":%d", config.Server.Port),
			ReadTimeout: config.Server.ReadTimeOut,
			ErrorLog:    core.Logger(),
		},
		database: db,
		core:     core,
	}, nil
}

func (m *mainData) Start() {

	//Init router
	//Use gin default if want write logs
	router := gin.Default()
	api := router.Group("/api")

	c := cache.New()

	eventRepository, err := pkgEventRepository.New(m.core, m.database, c)
	if err != nil {
		m.core.Logger().Fatal("Failed start eventRepository")
	}

	//Init workers Pipeline
	eventService := pkgEventService.NewSender(eventRepository)
	if err != nil {
		m.core.Logger().Fatal("Failed start eventService")
	}

	//Init events handler
	eventHandler := pkgEventsRouter.New(m.core.Logger(), eventService)
	eventHandler.RegisterRoutes(api)

	//Start listening server
	m.core.Group().Go(func() error {
		if r := recover(); r != nil {
			m.core.Logger().Printf("(PANIC) %v %v", r, string(debug.Stack()))
		}

		m.server.Handler = router
		return m.server.ListenAndServe()
	})
	//Stop listening server
	m.core.Group().Go(func() error {
		<-m.core.Context().Done()
		return m.server.Shutdown(context.Background())
	})
}
