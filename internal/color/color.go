package color

import (
	"fmt"
	"os"
)

// ANSI escape codes.
const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
)

// enabled tracks whether color output is active.
var enabled = isTerminal()

func isTerminal() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func wrap(code, s string) string {
	if !enabled {
		return s
	}
	return code + s + reset
}

// Bold returns s in bold.
func Bold(s string) string { return wrap(bold, s) }

// Dim returns s in dim/faint.
func Dim(s string) string { return wrap(dim, s) }

// Red returns s in red.
func Red(s string) string { return wrap(red, s) }

// Green returns s in green.
func Green(s string) string { return wrap(green, s) }

// Yellow returns s in yellow.
func Yellow(s string) string { return wrap(yellow, s) }

// Cyan returns s in cyan.
func Cyan(s string) string { return wrap(cyan, s) }

// BoldGreen returns s in bold green.
func BoldGreen(s string) string { return wrap(bold+green, s) }

// Success prints a green check mark followed by the message.
func Success(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if enabled {
		fmt.Printf("%s%s %s%s\n", green, "✓", msg, reset)
	} else {
		fmt.Println(msg)
	}
}

// Warn prints a yellow warning message.
func Warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if enabled {
		fmt.Printf("%s%s %s%s\n", yellow, "!", msg, reset)
	} else {
		fmt.Printf("Warning: %s\n", msg)
	}
}

// Prompt prints a bold cyan prompt (no newline).
func Prompt(s string) {
	if enabled {
		fmt.Printf("%s%s%s ", bold+cyan, s, reset)
	} else {
		fmt.Printf("%s ", s)
	}
}
