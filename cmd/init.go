package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/ui"
)

var initForce bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Write a .cage config in the current directory",
	Long: `Write a .cage config in the current directory.

If a parent directory is already caged, cage will use that VM by default.
Use --force to cage the current directory explicitly anyway.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func init() {
	initCmd.Flags().BoolVar(&initForce, "force", false, "cage this directory even if a parent is already caged")
	rootCmd.AddCommand(initCmd)
}

func runInit() error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	if config.Exists(cwd) {
		ui.Die(".cage already exists (delete it first to re-init)")
	}

	// Warn if a parent is already caged, unless --force
	if !initForce {
		if parent := config.FindCageDir(filepath.Dir(cwd)); parent != "" {
			ui.Warn("parent directory is already caged: %s", ui.Bold(parent))
			ui.Warn("running %s here will start a separate VM for this subdirectory", ui.Bold("cage"))

			if !ui.Confirm("Continue anyway?", "n") {
				ui.Info("tip: run %s without init to use the parent VM", ui.Bold("cage"))

				return nil
			}
		}
	}

	defaultName := filepath.Base(cwd) + "-cage"
	vmName := ui.Ask("VM name", defaultName)

	if err := config.Write(cwd, vmName); err != nil {
		return err
	}

	ui.Success("wrote %s — run %s to provision the VM", ui.Bold(".cage"), ui.Bold("cage start"))

	return nil
}
