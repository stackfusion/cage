package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/lima"
	"github.com/stackfusion/cage/internal/ui"
)

var hookCmd = &cobra.Command{
	Use:   "hook [bash|zsh|fish]",
	Short: "Print shell integration snippet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runHook(args[0])
	},
}

var chpwdCmd = &cobra.Command{
	Use:    "chpwd",
	Short:  "Check caged status on directory change (called by shell hook)",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Run swallows any errors and exit code is always 0.
		// chpwd must never interrupt the user's workflow.
		_ = runChpwd()
	},
}

var promptCmd = &cobra.Command{
	Use:    "prompt",
	Short:  "Check if in a caged directory or if VM is running (called by shell prompt)",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Run swallows any errors and exit code is always 0.
		// prompt must never interrupt the user's workflow.
		_ = runPrompt()
	},
}

func init() {
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(chpwdCmd)
	rootCmd.AddCommand(promptCmd)
}

const zshHook = `
# cage shell integration (zsh)
_cage_chpwd() { cage chpwd 2>/dev/null || true; }
autoload -Uz add-zsh-hook
add-zsh-hook chpwd _cage_chpwd
`

const bashHook = `
# cage shell integration (bash)
_cage_prompt_command() { cage chpwd 2>/dev/null || true; }
PROMPT_COMMAND="_cage_prompt_command${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
`

const fishHook = `
# cage shell integration (fish)
function _cage_chpwd --on-variable PWD
    cage chpwd
end
`

func runHook(shell string) error {
	switch shell {
	case "zsh":
		fmt.Print(zshHook)
	case "bash":
		fmt.Print(bashHook)
	case "fish":
		fmt.Print(fishHook)
	default:
		return fmt.Errorf("unknown shell %q — use bash, zsh, or fish", shell)
	}

	return nil
}

func runChpwd() error {
	cwd, err := os.Getwd()

	if err != nil {
		return nil
	}

	cageDir := config.FindCageDir(cwd)

	if cageDir == "" {
		return nil
	}

	name := config.VMName(cageDir)
	running, _ := lima.IsRunning(name)

	if config.IsAcknowledged(cageDir) {
		if running {
			ui.Subtle("%s", ui.Dim("vm is running · `cage` to enter it"))
		} else {
			ui.Subtle("%s", ui.Dim("vm is stopped · `cage` to start & enter it"))
		}
	} else {
		if running {
			ui.Info("caged directory, vm is %s", ui.Green("running"))
			ui.Info("run %s to enter the VM, or %s to mute this banner", ui.Bold("`cage`"), ui.Bold("`cage ack`"))
		} else {
			ui.Warn("caged directory, vm is %s", ui.Red("stopped"))
			ui.Warn("run %s to start and enter the VM, or %s to mute this banner", ui.Bold("`cage`"), ui.Bold("`cage ack`"))
		}
	}

	return nil
}

func runPrompt() error {
	cwd, err := os.Getwd()

	if err != nil {
		return nil // stay silent
	}

	cageDir := config.FindCageDir(cwd)

	if cageDir == "" {
		fmt.Print("uncaged")

		return nil
	}

	name := config.VMName(cageDir)
	running, _ := lima.IsRunning(name)

	if running {
		fmt.Print("running")
	} else {
		fmt.Print("stopped")
	}

	return nil
}
