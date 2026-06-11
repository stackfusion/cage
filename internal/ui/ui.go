package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	Cyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	Green  = color.New(color.FgGreen, color.Bold).SprintFunc()
	Yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	Red    = color.New(color.FgRed, color.Bold).SprintFunc()
	Bold   = color.New(color.Bold).SprintFunc()
	Dim    = color.New(color.Faint).SprintFunc()
)

var (
	labelSubtle = Dim("cage:")
	labelInfo   = Cyan("cage:")
	labelWarn   = Yellow("cage:")
	labelOK     = Green("cage:")
	labelError  = Red("cage: error:")
)

func stderr(label, msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", label, msg)
}

func Subtle(format string, a ...any) {
	stderr(labelSubtle, fmt.Sprintf(format, a...))
}

func Info(format string, a ...any) {
	stderr(labelInfo, fmt.Sprintf(format, a...))
}

func Success(format string, a ...any) {
	stderr(labelOK, fmt.Sprintf(format, a...))
}

func Warn(format string, a ...any) {
	stderr(labelWarn, fmt.Sprintf(format, a...))
}

func Die(format string, a ...any) {
	stderr(labelError, fmt.Sprintf(format, a...))

	os.Exit(1)
}

// Ask prints a styled prompt and reads a line from stdin.
// Returns def if the user presses Enter without typing anything.
func Ask(prompt, def string) string {
	hint := ""

	if def != "" {
		hint = " " + Dim("["+def+"]")
	}

	fmt.Fprintf(os.Stderr, "%s %s%s: ", color.YellowString("?"), Bold(prompt), hint)

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

	fmt.Fprintf(os.Stderr, "%s %s %s: ", color.YellowString("?"), Bold(prompt), Dim("["+hint+"]"))

	var answer string

	fmt.Fscanln(os.Stdin, &answer)

	if answer == "" {
		answer = def
	}

	return answer == "y" || answer == "Y"
}
