package service

import (
	"TestTask-events/pkg/app"
	"TestTask-events/pkg/events"
	"TestTask-events/pkg/workers"
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
	if len(event) == 0 {
		return nil
	}
	s.analyticsPipe.Send(event)
	return nil
}
