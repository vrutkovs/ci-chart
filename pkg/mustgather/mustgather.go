package mustgather

import (
	"github.com/vrutkovs/ci-chart/pkg/event"
)

type parser struct {
	path        string
	tmplocation string
}

type Parser interface {
	ParseMustGather() error
	Namespaces() []string
	PodEvents(ns string) []event.Input
	OperatorEvents(ns string) []event.Input
}

func NewParser(path string) Parser {
	// Unpack tar.gz to tmploc
	parser := &parser{path: path}
	return parser
}

func (s *parser) ParseMustGather() error {
	err := s.unpackMustGather()
	if err != nil {
		return err
	}
	return nil
}

func (s *parser) unpackMustGather() error {
	return nil
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
