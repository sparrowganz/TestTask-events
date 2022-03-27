package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/roistat/go-clickhouse"
	"github.com/sparrowganz/TestTask-events/config"
	"github.com/sparrowganz/TestTask-events/pkg/app"
	pkgDBClickhouse "github.com/sparrowganz/TestTask-events/pkg/db/clickhouse"
	pkgEventsRouter "github.com/sparrowganz/TestTask-events/pkg/events/delivery/http"
	eventsPipeline "github.com/sparrowganz/TestTask-events/pkg/events/pipelines"
	pkgEventRepository "github.com/sparrowganz/TestTask-events/pkg/events/repository"
	pkgEventService "github.com/sparrowganz/TestTask-events/pkg/events/service"
	"net/http"
	"runtime/debug"
)

type mainData struct {
	server   http.Server
	database *clickhouse.Conn
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
	router := gin.Default()
	api := router.Group("/api")

	eventRepository := pkgEventRepository.New(m.database)
	err := eventRepository.InitMigrations()
	if err != nil {
		m.core.Logger().Fatal(err.Error())
	}

	analyticsPipe, err := eventsPipeline.CreateEventAnalyticsPipeline(m.core, eventRepository)
	if err != nil {
		m.core.Logger().Fatal(err.Error())
	}

	//Init workers Pipeline
	eventService := pkgEventService.NewSender(m.core, analyticsPipe)
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
