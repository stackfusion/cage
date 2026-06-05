package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackfusion/cage/internal/config"
	"github.com/stackfusion/cage/internal/ui"
)

// limaTemplate is embedded into the binary at build time.
// It is written to ~/.config/cage/lima-template.yaml by cage install.
const limaTemplate = `# cage Lima VM template
# CAGE_MOUNT_HOST and CAGE_MOUNT_VM are substituted at 'cage start' time.

vmType: vz
arch: default

cpus: 4
memory: 8GiB
disk: 40GiB

images:
  - location: https://cloud-images.ubuntu.com/releases/questing/release/ubuntu-25.10-server-cloudimg-amd64.img
    arch: x86_64
  - location: https://cloud-images.ubuntu.com/releases/questing/release/ubuntu-25.10-server-cloudimg-arm64.img
    arch: aarch64

mounts:
  - location: "CAGE_MOUNT_HOST"
    mountPoint: "CAGE_MOUNT_VM"
    writable: true

ssh:
  forwardAgent: true

provision:
  - mode: system
    script: |
      #!/bin/bash
      apt update -y
      apt install -y autoconf build-essential curl eza fish fop fzf git htop inotify-tools libncurses-dev libssl-dev libxml2-utils m4 parallel ripgrep tmux unixodbc-dev xsltproc
`

const rcMarker = "# cage shell integration"

type shellConfig struct {
	rc        string // path to rc file
	snippet   string // text to append
	sourceCmd string // how to reload after patching
}

var shells = map[string]func() shellConfig{
	"zsh": func() shellConfig {
		return shellConfig{
			rc: filepath.Join(os.Getenv("HOME"), ".zshrc"),
			snippet: `
# cage shell integration
if command -v cage &>/dev/null; then
  eval "$(cage hook zsh)"
fi
`,
			sourceCmd: "source ~/.zshrc",
		}
	},
	"bash": func() shellConfig {
		return shellConfig{
			rc: filepath.Join(os.Getenv("HOME"), ".bashrc"),
			snippet: `
# cage shell integration
if command -v cage &>/dev/null; then
  eval "$(cage hook bash)"
fi
`,
			sourceCmd: "source ~/.bashrc",
		}
	},
	"fish": func() shellConfig {
		return shellConfig{
			rc: filepath.Join(os.Getenv("HOME"), ".config", "fish", "config.fish"),
			snippet: `
# cage shell integration
if command -v cage &>/dev/null
  cage hook fish | source
end
`,
			sourceCmd: "source ~/.config/fish/config.fish",
		}
	},
}

var installShell string

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Write Lima template and patch shell rc",
	Long: `Write the Lima VM template to ~/.config/cage/ and patch your shell rc file.

Detects your current shell automatically, or pass --shell to be explicit:

  cage install --shell zsh
  cage install --shell bash
  cage install --shell fish`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstall(installShell)
	},
}

func init() {
	installCmd.Flags().StringVar(&installShell, "shell", "", "shell to patch (zsh, bash, fish); defaults to $SHELL")
	rootCmd.AddCommand(installCmd)
}

func runInstall(shell string) error {
	if err := os.MkdirAll(config.Dir(), 0755); err != nil {
		return err
	}

	// Write Lima template
	tmplPath := config.LimaTemplatePath()

	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		if err := os.WriteFile(tmplPath, []byte(limaTemplate), 0644); err != nil {
			return err
		}

		ui.Success("wrote Lima template → %s", ui.Bold(tmplPath))
	} else {
		ui.Warn("Lima template already exists at %s (skipping)", ui.Bold(tmplPath))
	}

	// Resolve shell
	if shell == "" {
		shell = detectShell()
	}

	cfgFn, ok := shells[shell]

	if !ok {
		return fmt.Errorf("unsupported shell %q — use zsh, bash, or fish", shell)
	}

	cfg := cfgFn()

	// Ensure parent dir exists (needed for fish)
	if err := os.MkdirAll(filepath.Dir(cfg.rc), 0755); err != nil {
		return err
	}

	if err := patchRC(cfg.rc, cfg.snippet); err != nil {
		return err
	}

	ui.Success("done — open a new shell or: %s", ui.Bold(cfg.sourceCmd))

	return nil
}

// detectShell returns the shell name inferred from the $SHELL env var.
func detectShell() string {
	shell := os.Getenv("SHELL")

	// $SHELL is typically a path like /bin/zsh or /usr/local/bin/fish
	return filepath.Base(shell)
}

func patchRC(rc, snippet string) error {
	data, err := os.ReadFile(rc)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if strings.Contains(string(data), rcMarker) {
		ui.Warn("%s already patched (skipping)", ui.Bold(rc))
		return nil
	}

	f, err := os.OpenFile(rc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.WriteString(snippet); err != nil {
		return err
	}

	ui.Success("patched %s", ui.Bold(rc))

	return nil
}
