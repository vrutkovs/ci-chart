package client

import (
	"fmt"

	"github.com/spf13/pflag"
)

const (
	defaultHTTPPort  = 3001
	defaultCacheSize = 100
)

// Factory knows how to create a Kubernetes client.
type Factory interface {
	// BindFlags binds common flags (--kubeconfig, --namespace) to the passed-in FlagSet.
	BindFlags(flags *pflag.FlagSet)
	// MustGather is a path to must-gather archive
	MustGather() string
	// Port returns the port to listen on
	Port() uint16
	// LogAllEvents returns whether we should log all events, including transitions between the same state
	LogAllEvents() bool
}

type factory struct {
	flags        *pflag.FlagSet
	mustgather   string
	baseName     string
	httpPort     uint16
	logAllEvents bool
}

// NewFactory returns a Factory.
func NewFactory(baseName string) Factory {
	f := &factory{
		flags:    pflag.NewFlagSet("", pflag.ContinueOnError),
		baseName: baseName,
	}

	f.flags.StringVar(&f.mustgather, "mustgather", "", "Path to the must-gather archive")
	f.flags.Uint16Var(&f.httpPort, "http-port", uint16(defaultHTTPPort), fmt.Sprintf("Port to serve charts on. Default: %d", defaultHTTPPort))
	f.flags.BoolVar(&f.logAllEvents, "log-all-events", false, fmt.Sprintf("Log all events, including same-state transitions"))

	return f
}

func (f *factory) BindFlags(flags *pflag.FlagSet) {
	flags.AddFlagSet(f.flags)
}

func (f *factory) Port() uint16 {
	return f.httpPort
}

func (f *factory) MustGather() string {
	return f.mustgather
}

func (f *factory) LogAllEvents() bool {
	return f.logAllEvents
}
