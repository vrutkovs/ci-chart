package controller

import (
	"github.com/vrutkovs/ci-chart/pkg/event"
)

type Controller struct {
	mustgather string
	eventStore event.Store
}

func New(mustgather string, eventStore event.Store) *Controller {
	controller := &Controller{
		mustgather: mustgather,
		eventStore: eventStore,
	}

	return controller
}

// Run will parse must gather and fill in eventStore
func (c *Controller) Run() error {
	return nil
}
