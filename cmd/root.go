package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/lima"
	"github.com/stackfusion/cage/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:           "cage",
	Short:         "Safely work with untrusted projects via Lima VMs",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDefault()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		type exiter interface{ ExitCode() int }

		if exitErr, ok := err.(exiter); ok {
			// Clean non-zero exit from a subprocess (e.g. Ctrl+C) — no message needed
			os.Exit(exitErr.ExitCode())
		}

		// Real error — print it ourselves since cobra is silenced
		ui.Die("%s", err)
	}
}

func runDefault() error {
	if err := requireLima(); err != nil {
		ui.Die("%s", err)
	}

	// 1. Check install
	if _, err := os.Stat(config.LimaTemplatePath()); os.IsNotExist(err) {
		ui.Warn("Lima template not found at %s", ui.Bold(config.LimaTemplatePath()))

		if !ui.Confirm("Run 'cage install' now?", "y") {
			ui.Die("run %s to set up cage first", ui.Bold("cage install"))
		}

		if err := runInstall(""); err != nil {
			return err
		}
	}

	cwd, _ := os.Getwd()

	// 2. Find the nearest caged directory — cwd itself or a parent.
	//    If none found, offer to init the current directory.
	cageDir := config.FindCageDir(cwd)

	if cageDir == "" {
		ui.Warn("no %s file in current directory or any parent", ui.Bold(".cage"))

		if !ui.Confirm("Run 'cage init' now?", "y") {
			ui.Die("run %s to configure this project", ui.Bold("cage init"))
		}

		if err := runInit(); err != nil {
			return err
		}

		cageDir = cwd
	} else if cageDir != cwd {
		ui.Info("using caged parent directory: %s", ui.Bold(cageDir))
	}

	// 3. Start if needed, then shell — landing in cwd inside the VM,
	//    not necessarily cageDir.
	name := config.VMName(cageDir)
	running, err := lima.IsRunning(name)

	if err != nil {
		return err
	}

	if !running {
		// runStart reads .cage from the working directory, so chdir temporarily.
		orig, _ := os.Getwd()

		if err := os.Chdir(cageDir); err != nil {
			return err
		}

		if err := runStart(); err != nil {
			return err
		}

		_ = os.Chdir(orig)
	}

	// Pass cwd as workDir so the shell lands where the user actually is.
	return runShellFrom(cwd, nil)
}
