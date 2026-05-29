package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/lima"
	"github.com/stackfusion/cage/internal/ui"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Create (if needed) and start the Lima VM",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStart()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart() error {
	if err := requireLima(); err != nil {
		ui.Die("%s", err)
	}

	if err := requireCageFile(); err != nil {
		ui.Die("%s", err)
	}

	if err := requireTemplate(); err != nil {
		ui.Die("%s", err)
	}

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	name := config.VMName(cwd)
	exists, err := lima.Exists(name)

	if err != nil {
		return err
	}

	if exists {
		running, err := lima.IsRunning(name)

		if err != nil {
			return err
		}

		if running {
			ui.Info("VM %s is already running", ui.Bold(name))

			return nil
		}

		ui.Info("starting VM %s...", ui.Bold(name))

		return lima.Start(name)
	}

	// New VM — render the template and create
	tmp, err := renderTemplate(cwd)

	if err != nil {
		return err
	}

	defer os.Remove(tmp)

	ui.Info("creating VM %s...", ui.Bold(name))

	if err := lima.Create(name, tmp); err != nil {
		return err
	}

	if err := lima.Start(name); err != nil {
		return err
	}

	port, err := lima.SSHPort(name)

	if err != nil {
		return err
	}

	ui.Success("VM %s running — SSH port %s", ui.Bold(name), ui.Bold(fmt.Sprintf("%d", port)))

	return nil
}

// renderTemplate substitutes mount paths into the Lima template and writes
// the result to a temp file. The caller must remove the file when done.
func renderTemplate(hostDir string) (string, error) {
	data, err := os.ReadFile(config.LimaTemplatePath())

	if err != nil {
		return "", err
	}

	vmPath := lima.MountPath(hostDir)
	rendered := strings.ReplaceAll(string(data), "CAGE_MOUNT_HOST", hostDir)
	rendered = strings.ReplaceAll(rendered, "CAGE_MOUNT_VM", vmPath)
	tmp, err := os.CreateTemp("", "cage-lima-*.yaml")

	if err != nil {
		return "", err
	}

	defer tmp.Close()

	if _, err := tmp.WriteString(rendered); err != nil {
		os.Remove(tmp.Name())

		return "", err
	}

	return tmp.Name(), nil
}
