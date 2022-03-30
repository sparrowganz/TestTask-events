package cache

import (
	"TestTask-events/pkg/events"
	"sync"
)

type Cache struct {
	data []*events.Event
	mu   *sync.Mutex
}

func New() *Cache {
	return &Cache{
		data: make([]*events.Event, 0, 1000000),
		mu:   &sync.Mutex{},
	}
}

func (c *Cache) Set(event *events.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = append(c.data, event)
}

func (c *Cache) GetAll() []*events.Event {
	c.mu.Lock()
	defer c.mu.Unlock()
	data := c.data
	c.data = make([]*events.Event, 0, 1000000)
	return data
}
