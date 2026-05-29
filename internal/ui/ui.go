package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	labelInfo    = color.New(color.FgCyan, color.Bold).Sprint("cage:")
	labelSuccess = color.New(color.FgGreen, color.Bold).Sprint("cage:")
	labelWarn    = color.New(color.FgYellow, color.Bold).Sprint("cage:")
	labelError   = color.New(color.FgRed, color.Bold).Sprint("cage: error:")

	bold = color.New(color.Bold).SprintFunc()
	dim  = color.New(color.Faint).SprintFunc()
)

func Info(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", labelInfo, fmt.Sprintf(format, a...))
}

func Success(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", labelSuccess, fmt.Sprintf(format, a...))
}

func Warn(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", labelWarn, fmt.Sprintf(format, a...))
}

func Die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", labelError, fmt.Sprintf(format, a...))

	os.Exit(1)
}

func Bold(s string) string { return bold(s) }
func Dim(s string) string  { return dim(s) }

// Ask prints a styled prompt and reads a line from stdin.
// Returns def if the user presses Enter without typing anything.
func Ask(prompt, def string) string {
	hint := ""

	if def != "" {
		hint = " " + dim("["+def+"]")
	}

	fmt.Fprintf(os.Stderr, "%s %s%s: ", color.YellowString("?"), bold(prompt), hint)

	var answer string

	fmt.Fscanln(os.Stdin, &answer)

	if answer == "" {
		return def
	}

	return answer
}

// Confirm prints a yes/no prompt and returns true if the user answers y/Y.
// def should be "y" or "n" and is used when the user presses Enter.
func Confirm(prompt, def string) bool {
	hint := "y/N"

	if def == "y" {
		hint = "Y/n"
	}

	fmt.Fprintf(os.Stderr, "%s %s %s: ", color.YellowString("?"), bold(prompt), dim("["+hint+"]"))

	var answer string

	fmt.Fscanln(os.Stdin, &answer)

	if answer == "" {
		answer = def
	}

	return answer == "y" || answer == "Y"
}
