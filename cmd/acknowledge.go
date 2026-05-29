package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/ui"
)

var acknowledgeCmd = &cobra.Command{
	Use:     "acknowledge",
	Aliases: []string{"ack"},
	Short:   "Acknowledge the .cage file in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAcknowledge()
	},
}

func init() {
	rootCmd.AddCommand(acknowledgeCmd)
}

func runAcknowledge() error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	if !config.Exists(cwd) {
		ui.Die("no %s file in current directory", ui.Bold(".cage"))
	}

	if err := config.Acknowledge(cwd); err != nil {
		return err
	}

	ui.Success("acknowledged — banner suppressed until %s changes", ui.Bold(".cage"))

	return nil
}
