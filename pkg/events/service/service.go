package service

import (
	"github.com/sparrowganz/TestTask-events/pkg/app"
	"github.com/sparrowganz/TestTask-events/pkg/events"
	"github.com/sparrowganz/TestTask-events/pkg/workers"
)

type Service struct {
	core          app.Core
	analyticsPipe workers.Pipeline
	repo          events.Repository
}

func NewSender(core app.Core, pipe workers.Pipeline) *Service {
	return &Service{
		core:          core,
		analyticsPipe: pipe,
	}
}

func (s Service) SendEvent(event []byte) error {
	s.analyticsPipe.Send(event)
	return nil
}
