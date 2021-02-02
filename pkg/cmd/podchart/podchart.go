package podchart

import (
	"flag"

	"github.com/spf13/cobra"

	"github.com/vrutkovs/ci-chart/pkg/client"
	"github.com/vrutkovs/ci-chart/pkg/cmd"
	"github.com/vrutkovs/ci-chart/pkg/controller"
	"github.com/vrutkovs/ci-chart/pkg/event"
	"github.com/vrutkovs/ci-chart/pkg/ui"
)

func NewCommand(name string) *cobra.Command {
	f := client.NewFactory(name)

	c := &cobra.Command{
		Use:   name,
		Short: "Monitor ClusterOperator phase transitions over time in a OpenShift cluster.",
		Run: func(c *cobra.Command, args []string) {
			cmd.CheckError(run(c, f))
		},
	}

	f.BindFlags(c.PersistentFlags())
	c.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return c
}

func run(c *cobra.Command, f client.Factory) error {
	eventStore := event.NewStore()
	controller := controller.New(f.MustGather(), eventStore)

	go ui.Run(eventStore, f.Port())
	if err := controller.Run(); err != nil {
		return err
	}

	return nil
}
