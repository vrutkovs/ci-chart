package mustgather

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/vrutkovs/ci-chart/pkg/event"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

type parser struct {
	path        string
	tmplocation string
}

type Parser interface {
	ParseMustGather() error
	Namespaces() []string
	PodEvents(ns string) *[]event.Input
	OperatorEvents(ns string) *[]event.Input
}

func NewParser(path string) Parser {
	// Unpack tar.gz to tmploc
	parser := &parser{path: path}
	return parser
}

func (s *parser) ParseMustGather() error {
	if err := s.unpackMustGather(); err != nil {
		return err
	}
	return nil
}

func (s *parser) unpackMustGather() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzf)

	parentDir := os.TempDir()
	s.tmplocation, err = ioutil.TempDir(parentDir, "mustgather")
	if err != nil {
		return err
	}

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(s.tmplocation, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tarReader); err != nil {
				return err
			}
			f.Close()
		}
	}
}

func (s *parser) Namespaces() []string {
	namespaces := []string{}
	paths, err := filepath.Glob(fmt.Sprintf("%s/*/namespaces/*", s.tmplocation))
	if err != nil {
		return namespaces
	}
	for _, p := range paths {
		namespaces = append(namespaces, filepath.Base(p))
	}
	return namespaces
}

func parseEventsYaml(inputs []corev1.Event, evfile string) []corev1.Event {
	bytes, err := ioutil.ReadFile(evfile)
	if err != nil {
		panic(err.Error())
	}
	// Open events file, deserialize it and return all inputs
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(bytes), nil, nil)
	if err != nil {
		fmt.Print(err)
		return inputs
	}

	for _, ev := range obj.(*corev1.EventList).Items {
		inputs = append(inputs, ev)
	}
	return inputs
}

func (s *parser) parseEventFiles(ns string) []corev1.Event {
	inputs := []corev1.Event{}
	eventfiles, err := filepath.Glob(fmt.Sprintf("%s/*/namespaces/%s/core/events.yaml", s.tmplocation, ns))
	if err != nil {
		return inputs
	}
	for _, evfile := range eventfiles {
		inputs = parseEventsYaml(inputs, evfile)
	}
	return inputs

}

func (s *parser) PodEvents(ns string) *[]event.Input {
	result := []event.Input{}

	inputs := s.parseEventFiles(ns)
	for _, ev := range inputs {
		if ev.InvolvedObject.Kind != "Pod" || ev.InvolvedObject.FieldPath != "" {
			continue
		}
		newInput := event.PodEventToInput(ev)
		result = append(result, newInput)
	}
	return &result
}

func (s *parser) OperatorEvents(ns string) *[]event.Input {
	result := []event.Input{}

	inputs := s.parseEventFiles(ns)
	for _, ev := range inputs {
		if ev.Reason != "OperatorStatusChanged" {
			continue
		}
		newInput := event.ClusterOperatorEventToInput(ev)
		result = append(result, newInput)
	}
	return &result
}
