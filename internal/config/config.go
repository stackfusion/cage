package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const CageFile = ".cage"

// Respecting XDG_CONFIG_HOME returns returns ~/.config/cage
func Dir() string {
	base := os.Getenv("XDG_CONFIG_HOME")

	if base == "" {
		base = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(base, "cage")
}

// LimaTemplatePath returns the path to the Lima YAML template.
func LimaTemplatePath() string {
	return filepath.Join(Dir(), "lima-template.yaml")
}

// AckPath returns the path to the acknowledgement registry file.
func AckPath() string {
	return filepath.Join(Dir(), "acknowledged")
}

// Exists reports whether a .cage file is present in dir.
func Exists(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, CageFile))

	return err == nil
}

// FindCageDir walks up from dir looking for a .cage file.
// Returns the directory containing .cage, or "" if not found.
func FindCageDir(dir string) string {
	for {
		if Exists(dir) {
			return dir
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			return "" // reached filesystem root
		}

		dir = parent
	}
}

// Field reads a simple "key: value" field from the .cage file in dir.
// Returns def if the file or key is not found.
func Field(dir, key, def string) string {
	data, err := os.ReadFile(filepath.Join(dir, CageFile))

	if err != nil {
		return def
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, key+":") {
			val := strings.TrimSpace(strings.TrimPrefix(line, key+":"))
			val = strings.Trim(val, `"`)

			if val != "" {
				return val
			}
		}
	}

	return def
}

// VMName returns the vm_name field from .cage, defaulting to "<dirname>-cage".
func VMName(dir string) string {
	return Field(dir, "vm_name", filepath.Base(dir)+"-cage")
}

// Write writes a minimal .cage file to dir.
func Write(dir, vmName string) error {
	path := filepath.Join(dir, CageFile)
	content := fmt.Sprintf("vm_name: %s\n", vmName)

	return os.WriteFile(path, []byte(content), 0644)
}

// Hash returns a short SHA-256 hash of the .cage file contents in dir.
func Hash(dir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(dir, CageFile))

	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(data)

	return fmt.Sprintf("%x", sum[:8]), nil
}

// IsAcknowledged reports whether the .cage file in dir has been acknowledged
// and hasn't changed since.
func IsAcknowledged(dir string) bool {
	hash, err := Hash(dir)

	if err != nil {
		return false
	}

	data, err := os.ReadFile(AckPath())

	if err != nil {
		return false
	}

	// Each line: "<absolute-dir> <hash>"
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Fields(line)

		if len(parts) == 2 && parts[0] == dir && parts[1] == hash {
			return true
		}
	}

	return false
}

// Acknowledge records the current .cage hash for dir.
func Acknowledge(dir string) error {
	hash, err := Hash(dir)

	if err != nil {
		return err
	}

	if err := os.MkdirAll(Dir(), 0755); err != nil {
		return err
	}

	// Read existing entries, replace or append this dir's entry.
	var lines []string
	data, err := os.ReadFile(AckPath())

	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			parts := strings.Fields(line)

			if len(parts) == 2 && parts[0] != dir {
				lines = append(lines, line)
			}
		}
	}

	lines = append(lines, fmt.Sprintf("%s %s", dir, hash))

	return os.WriteFile(AckPath(), []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
