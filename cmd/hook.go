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
	// Use Run (not RunE) so errors are swallowed and exit code is always 0.
	// chpwd must never interrupt the user's workflow.
	Run: func(cmd *cobra.Command, args []string) {
		_ = runChpwd()
	},
}

func init() {
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(chpwdCmd)
}

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

// runChpwd is called by the shell on every cd. It must be fast and silent
// when the directory is not caged.
func runChpwd() error {
	cwd, err := os.Getwd()

	if err != nil {
		return nil // stay silent on error
	}

	cageDir := config.FindCageDir(cwd)

	if cageDir == "" {
		return nil // not a caged directory
	}

	name := config.VMName(cageDir)
	isChild := cageDir != cwd

	// Already acknowledged and unchanged → show quiet one-line indicator
	if config.IsAcknowledged(cageDir) {
		running, _ := lima.IsRunning(name)

		if running {
			fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("[✓ cage: "+name+"]"))
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("[⚠ cage: "+name+" — not running]"))
		}

		return nil
	}

	// Not yet acknowledged → loud banner
	running, _ := lima.IsRunning(name)
	fmt.Fprintln(os.Stderr, "")

	if isChild {
		ui.Warn("inside caged project %s (at %s)", ui.Bold(name), ui.Bold(cageDir))
	} else {
		ui.Warn("caged directory — VM %s", ui.Bold(name))
	}

	if running {
		ui.Info("run %s to enter the VM", ui.Bold("cage shell"))
	} else {
		ui.Info("VM is not running — run %s to start and enter it", ui.Bold("cage"))
		ui.Info("run %s to suppress this banner", ui.Bold("cage acknowledge"))
	}

	fmt.Fprintln(os.Stderr, "")

	return nil
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
