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

func (s *parser) PodEvents(ns string) []event.Input {
	return []event.Input{}
}

func (s *parser) OperatorEvents(ns string) []event.Input {
	return []event.Input{}
}
