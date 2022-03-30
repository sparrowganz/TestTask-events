package service

import (
	"TestTask-events/pkg/events"
)

type Service struct {
	repo events.Repository
}

func NewSender(repo events.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s Service) SendEvent(event *events.Event) error {
	event.IP = "8.8.8.8"
	event.ServerTime = "2020-12-01 23:53:00"
	s.repo.Collect(event)
	return nil
}
