package controller

import (
	"github.com/vrutkovs/ci-chart/pkg/event"
	"github.com/vrutkovs/ci-chart/pkg/mustgather"
)

type Controller struct {
	mustgather mustgather.Parser
	eventStore event.Store
}

func New(path string, eventStore event.Store) *Controller {
	parser := mustgather.NewParser(path)
	controller := &Controller{
		mustgather: parser,
		eventStore: eventStore,
	}

	return controller
}

// ParseMustGather runs must gather parsing
func (c *Controller) ParseMustGather() error {
	return c.mustgather.ParseMustGather()
}

// FindPodTransitions will parse must gather and fill in eventStore with pod state transitions
func (c *Controller) FindPodTransitions() {
	for _, ns := range c.mustgather.Namespaces() {
		for _, i := range *c.mustgather.PodEvents(ns) {
			c.eventStore.Add(i)
		}
	}
}

// FindOperatorTransitions will parse must gather and fill in eventStore with pod state transitions
func (c *Controller) FindOperatorTransitions() {
	for _, ns := range c.mustgather.Namespaces() {
		for _, i := range *c.mustgather.OperatorEvents(ns) {
			c.eventStore.Add(i)
		}
	}
}
