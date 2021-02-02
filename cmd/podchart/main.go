package main

import (
	"os"
	"path/filepath"

	"k8s.io/klog/v2"

	"github.com/vrutkovs/ci-chart/pkg/cmd"
	"github.com/vrutkovs/ci-chart/pkg/cmd/podchart"
)

func main() {
	defer klog.Flush()

	baseName := filepath.Base(os.Args[0])

	err := podchart.NewCommand(baseName).Execute()
	cmd.CheckError(err)
}
