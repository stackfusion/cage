package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/lima"
	"github.com/stackfusion/cage/internal/ui"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove cage config/install; suggest Lima image cleanup",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPrune()
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}

func runPrune() error {
	ui.Warn("this will remove the cage installation and any .cage file here")

	if !ui.Confirm("Continue?", "n") {
		ui.Info("aborted")
		return nil
	}

	cwd, _ := os.Getwd()

	// Offer to delete the VM if .cage exists
	if config.Exists(cwd) {
		name := config.VMName(cwd)
		exists, err := lima.Exists(name)

		if err != nil {
			return err
		}

		if exists && ui.Confirm("Delete VM '"+name+"' too?", "n") {
			running, _ := lima.IsRunning(name)

			if running {
				if err := lima.Stop(name); err != nil {
					return err
				}
			}

			if err := lima.Delete(name); err != nil {
				return err
			}

			ui.Success("deleted VM %s", ui.Bold(name))
		}

		if err := os.Remove(config.CageFile); err != nil && !os.IsNotExist(err) {
			return err
		}

		ui.Success("removed %s", ui.Bold(".cage"))
	}

	// Offer to remove ~/.config/cage
	if _, err := os.Stat(config.Dir()); err == nil {
		if ui.Confirm("Remove "+config.Dir()+"?", "n") {
			if err := os.RemoveAll(config.Dir()); err != nil {
				return err
			}

			ui.Success("removed %s", ui.Bold(config.Dir()))
		}
	}

	ui.Warn("downloaded Lima images are kept in %s", ui.Bold("~/.lima"))
	ui.Info("to free disk space, run: %s", ui.Bold("limactl prune"))

	return nil
}
