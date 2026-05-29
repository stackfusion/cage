package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/lima"
	"github.com/stackfusion/cage/internal/ui"
)

var shellCmd = &cobra.Command{
	Use:   "shell [-- command]",
	Short: "Open a shell inside the VM, or run a command",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		return runShellFrom(cwd, args)
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

// runShell is a convenience wrapper for use within the same directory.
func runShell(args []string) error {
	cwd, _ := os.Getwd()
	return runShellFrom(cwd, args)
}

// runShellFrom opens a shell inside the VM associated with cageDir,
// landing at workDir inside the VM. workDir may differ from cageDir
// when called from a subdirectory.
func runShellFrom(workDir string, args []string) error {
	if err := requireLima(); err != nil {
		ui.Die("%s", err)
	}

	cageDir := config.FindCageDir(workDir)

	if cageDir == "" {
		ui.Die("no %s file in current directory or any parent — run %s first", ui.Bold(".cage"), ui.Bold("cage init"))
	}

	name := config.VMName(cageDir)

	if err := requireRunning(name); err != nil {
		ui.Die("%s", err)
	}

	vmPath := lima.MountPath(workDir)

	if len(args) == 0 {
		port, _ := lima.SSHPort(name)

		ui.Info("connecting to VM %s (SSH :%d)...", ui.Bold(name), port)
	}

	return lima.Shell(name, vmPath, args...)
}
