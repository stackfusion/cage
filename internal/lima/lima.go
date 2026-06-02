package lima

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Instance represents a Lima VM as returned by limactl list --json.
type Instance struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	SSHLocalPort int    `json:"sshLocalPort"`
}

// List returns all Lima instances.
func List() ([]Instance, error) {
	out, err := run("limactl", "list", "--json")

	if err != nil {
		return nil, err
	}

	var instances []Instance

	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}

		var inst Instance

		if err := json.Unmarshal([]byte(line), &inst); err != nil {
			return nil, fmt.Errorf("parsing limactl output: %w", err)
		}

		instances = append(instances, inst)
	}

	return instances, nil
}

// Get returns a single instance by name, and whether it was found.
func Get(name string) (Instance, bool, error) {
	instances, err := List()

	if err != nil {
		return Instance{}, false, err
	}

	for _, inst := range instances {
		if inst.Name == name {
			return inst, true, nil
		}
	}

	return Instance{}, false, nil
}

// Exists reports whether a VM with the given name exists.
func Exists(name string) (bool, error) {
	_, found, err := Get(name)

	return found, err
}

// IsRunning reports whether the named VM is in Running status.
func IsRunning(name string) (bool, error) {
	inst, found, err := Get(name)

	if err != nil || !found {
		return false, err
	}

	return inst.Status == "Running", nil
}

// SSHPort returns the host SSH port for the named VM.
func SSHPort(name string) (int, error) {
	inst, _, err := Get(name)

	return inst.SSHLocalPort, err
}

// Create creates a new VM from a config file path.
func Create(name, configPath string) error {
	return runInteractive("limactl", "create", "--name="+name, configPath)
}

// Start starts a stopped VM.
func Start(name string) error {
	return runInteractive("limactl", "start", name)
}

// Stop stops a running VM.
func Stop(name string) error {
	return runInteractive("limactl", "stop", name)
}

// Delete deletes a VM.
func Delete(name string) error {
	return runInteractive("limactl", "delete", name)
}

// Shell opens an interactive shell (or runs a command) inside the VM.
// workdir is the path inside the VM to cd into.
// args are passed after "--"; if empty, an interactive shell is opened.
func Shell(name, workdir string, args ...string) error {
	argv := []string{"shell", "--workdir", workdir, name}

	if len(args) > 0 {
		argv = append(argv, "--")
		argv = append(argv, args...)
	}

	return runInteractive("limactl", argv...)
}

// GuestUser returns the Lima guest username for the current host user.
// Lima creates guests as <host-user>.guest.
func GuestUser() string {
	return os.Getenv("USER") + ".guest"
}

// MountPath maps a host path to the expected in-VM path.
// e.g. /Users/cr0t/Workspace/foo → /home/cr0t.guest/Workspace/foo
func MountPath(hostPath string) string {
	home := os.Getenv("HOME")
	rel := strings.TrimPrefix(hostPath, home)

	return "/home/" + GuestUser() + rel
}

// ---------------------------------------------------------------------------
// internal helpers
// ---------------------------------------------------------------------------

// run executes a command and returns its stdout. Stderr is discarded.
func run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}

	return stdout.String(), nil
}

// ExitError wraps a subprocess exit code so callers can distinguish a
// clean non-zero exit (e.g. Ctrl+C → 130) from an unexpected error.
type ExitError struct{ Code int }

func (e *ExitError) Error() string { return fmt.Sprintf("exit status %d", e.Code) }
func (e *ExitError) ExitCode() int { return e.Code }

// runInteractive runs a command with stdin/stdout/stderr attached to the
// terminal — used for limactl commands that need user interaction or live
// output (create, start, stop, shell).
func runInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return &ExitError{Code: exitErr.ExitCode()}
		}

		return err
	}

	return nil
}
