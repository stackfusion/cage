package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/lima"
	"github.com/stackfusion/cage/internal/ui"
)

func requireLima() error {
	if _, err := exec.LookPath("limactl"); err != nil {
		return fmt.Errorf("limactl not found — install Lima: %s", ui.Bold("brew install lima"))
	}

	return nil
}

func requireCageFile() error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	if !config.Exists(cwd) {
		return fmt.Errorf("no .cage file in current directory — run %s first", ui.Bold("cage init"))
	}

	return nil
}

func requireTemplate() error {
	if _, err := os.Stat(config.LimaTemplatePath()); os.IsNotExist(err) {
		return fmt.Errorf("Lima template not found — run %s first", ui.Bold("cage install"))
	}

	return nil
}

func requireRunning(name string) error {
	exists, err := lima.Exists(name)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("VM %s does not exist — run %s first", ui.Bold(name), ui.Bold("cage start"))
	}

	running, err := lima.IsRunning(name)

	if err != nil {
		return err
	}

	if !running {
		return fmt.Errorf("VM %s is not running — run %s first", ui.Bold(name), ui.Bold("cage start"))
	}

	return nil
}
