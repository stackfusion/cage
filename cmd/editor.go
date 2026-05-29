package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/lima"
	"github.com/stackfusion/cage/internal/ui"
)

var zedCmd = &cobra.Command{
	Use:   "zed",
	Short: "Open the project in Zed via SSH remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEditor("zed")
	},
}

var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Open the project in VS Code via SSH remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEditor("code")
	},
}

func init() {
	rootCmd.AddCommand(zedCmd)
	rootCmd.AddCommand(codeCmd)
}

func runEditor(editor string) error {
	if err := requireLima(); err != nil {
		ui.Die("%s", err)
	}

	if err := requireCageFile(); err != nil {
		ui.Die("%s", err)
	}

	// Check the editor binary is available
	if _, err := exec.LookPath(editor); err != nil {
		ui.Die("%s not found — is it installed and on your PATH?", ui.Bold(editor))
	}

	cwd, _ := os.Getwd()
	name := config.VMName(cwd)

	if err := requireRunning(name); err != nil {
		ui.Die("%s", err)
	}

	vmPath := lima.MountPath(cwd)
	// For simplicity, we can add "Include ~/.lima/*/ssh.config" to the ~/.ssh/config,
	// so every new Lima VM will be accessible as lima-<name>
	sshHost := "lima-" + name

	var args []string
	switch editor {
	case "zed":
		// zed ssh://lima-<name>/path/to/project
		args = []string{fmt.Sprintf("ssh://%s%s", sshHost, vmPath)}
	case "code":
		// code --remote ssh-remote+lima-<name> /path/to/project
		args = []string{"--remote", "ssh-remote+" + sshHost, vmPath}
	default:
		return fmt.Errorf("unknown editor %q", editor)
	}

	ui.Info("opening %s → %s:%s", ui.Bold(editor), ui.Bold(sshHost), vmPath)

	cmd := exec.Command(editor, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
