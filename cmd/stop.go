package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/lima"
	"github.com/stackfusion/cage/internal/ui"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Lima VM",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStop()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStop() error {
	if err := requireLima(); err != nil {
		ui.Die("%s", err)
	}

	if err := requireCageFile(); err != nil {
		ui.Die("%s", err)
	}

	cwd, _ := os.Getwd()
	name := config.VMName(cwd)
	exists, err := lima.Exists(name)

	if err != nil {
		return err
	}

	if !exists {
		ui.Die("VM %s does not exist", ui.Bold(name))
	}

	ui.Info("stopping VM %s...", ui.Bold(name))

	if err := lima.Stop(name); err != nil {
		return err
	}

	ui.Success("stopped")

	return nil
}
