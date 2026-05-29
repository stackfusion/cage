package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/lima"
	"github.com/stackfusion/cage/internal/ui"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Stop and permanently delete the Lima VM",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDelete()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func runDelete() error {
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

	if !ui.Confirm("Delete VM '"+name+"' and all its data?", "n") {
		ui.Info("aborted")

		return nil
	}

	running, err := lima.IsRunning(name)

	if err != nil {
		return err
	}

	if running {
		if err := lima.Stop(name); err != nil {
			return err
		}
	}

	if err := lima.Delete(name); err != nil {
		return err
	}

	ui.Success("deleted VM %s", ui.Bold(name))

	return nil
}
