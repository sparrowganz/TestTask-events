package app

import (
	"TestTask-events/pkg/logger"
	"context"
	"golang.org/x/sync/errgroup"
	"log"
)

type App interface {
	Start()
}

type Core interface {
	Logger() *log.Logger
	Group() *errgroup.Group
	Context() context.Context
	Stop()
	Wait() error
}

type coreData struct {
	logger     *log.Logger
	group      *errgroup.Group
	context    context.Context
	cancelFunc context.CancelFunc
}

func NewCore(
	ctx context.Context,
	appName string,
) Core {

	ctx, cancel := context.WithCancel(ctx)
	group, ctx := errgroup.WithContext(ctx)

	return &coreData{
		logger:     logger.New(appName),
		group:      group,
		context:    ctx,
		cancelFunc: cancel,
	}
}

func (c *coreData) Group() *errgroup.Group {
	return c.group
}

func (c *coreData) Context() context.Context {
	return c.context
}

func (c *coreData) Stop() {
	c.cancelFunc()
}

func (c *coreData) Wait() error {
	return c.group.Wait()
}

func (c *coreData) Logger() *log.Logger {
	return c.logger
}
