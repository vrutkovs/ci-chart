package mustgather

import "github.com/vrutkovs/ci-chart/pkg/event"

type parser struct {
	path string
}

type Parser interface {
	Namespaces() []string
	PodEvents(ns string) []event.Input
	OperatorEvents(ns string) []event.Input
}

func NewParser(path string) Parser {
	return &parser{path: path}
}

func (s *parser) Namespaces() []string {
	return []string{"foo"}
}

func (s *parser) PodEvents(ns string) []event.Input {
	return []event.Input{}
}

func (s *parser) OperatorEvents(ns string) []event.Input {
	return []event.Input{}
}
