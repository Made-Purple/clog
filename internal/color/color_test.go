package color

import (
	"io"
	"os"
	"strings"
	"testing"
)

// setEnabled overrides the package-level color toggle for the duration of a
// test, restoring it afterwards. Tests run sequentially, so mutating the global
// is safe as long as it is reset.
func setEnabled(t *testing.T, v bool) {
	t.Helper()
	orig := enabled
	enabled = v
	t.Cleanup(func() { enabled = orig })
}

// captureStdout runs fn with os.Stdout redirected to a pipe and returns
// everything written to it.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	done := make(chan string, 1)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()

	fn()

	os.Stdout = orig
	w.Close()
	out := <-done
	r.Close()
	return out
}

func TestWrap(t *testing.T) {
	setEnabled(t, false)
	if got := wrap(red, "hi"); got != "hi" {
		t.Errorf("wrap disabled = %q, want %q", got, "hi")
	}

	setEnabled(t, true)
	if got := wrap(red, "hi"); got != red+"hi"+reset {
		t.Errorf("wrap enabled = %q, want %q", got, red+"hi"+reset)
	}
}

func TestColorFuncsDisabled(t *testing.T) {
	setEnabled(t, false)
	fns := map[string]func(string) string{
		"Bold": Bold, "Dim": Dim, "Red": Red, "Green": Green,
		"Yellow": Yellow, "Cyan": Cyan, "BoldGreen": BoldGreen,
	}
	for name, fn := range fns {
		if got := fn("x"); got != "x" {
			t.Errorf("%s disabled = %q, want %q", name, got, "x")
		}
	}
}

func TestColorFuncsEnabled(t *testing.T) {
	setEnabled(t, true)
	cases := []struct {
		name string
		fn   func(string) string
		code string
	}{
		{"Bold", Bold, bold},
		{"Dim", Dim, dim},
		{"Red", Red, red},
		{"Green", Green, green},
		{"Yellow", Yellow, yellow},
		{"Cyan", Cyan, cyan},
		{"BoldGreen", BoldGreen, bold + green},
	}
	for _, c := range cases {
		want := c.code + "x" + reset
		if got := c.fn("x"); got != want {
			t.Errorf("%s = %q, want %q", c.name, got, want)
		}
	}
}

func TestSuccessWarnPromptPlain(t *testing.T) {
	setEnabled(t, false)

	if got := captureStdout(t, func() { Success("hello %s", "world") }); got != "hello world\n" {
		t.Errorf("Success plain = %q, want %q", got, "hello world\n")
	}
	if got := captureStdout(t, func() { Warn("be careful") }); got != "Warning: be careful\n" {
		t.Errorf("Warn plain = %q, want %q", got, "Warning: be careful\n")
	}
	if got := captureStdout(t, func() { Prompt("Continue?") }); got != "Continue? " {
		t.Errorf("Prompt plain = %q, want %q", got, "Continue? ")
	}
}

func TestSuccessWarnPromptColored(t *testing.T) {
	setEnabled(t, true)

	if got := captureStdout(t, func() { Success("done %d", 1) }); !strings.Contains(got, "✓") || !strings.Contains(got, "done 1") {
		t.Errorf("Success colored = %q, want a check mark and message", got)
	}
	if got := captureStdout(t, func() { Warn("watch out") }); !strings.Contains(got, "watch out") || !strings.Contains(got, yellow) {
		t.Errorf("Warn colored = %q, want yellow and message", got)
	}
	if got := captureStdout(t, func() { Prompt("Pick:") }); !strings.Contains(got, "Pick:") || !strings.Contains(got, reset) {
		t.Errorf("Prompt colored = %q, want message and reset", got)
	}
}

func TestIsTerminal(t *testing.T) {
	// NO_COLOR disables color regardless of the underlying stream.
	t.Setenv("NO_COLOR", "1")
	if isTerminal() {
		t.Errorf("isTerminal() = true with NO_COLOR set")
	}

	// A regular file is not a character device, so color is off.
	t.Setenv("NO_COLOR", "")
	f, err := os.CreateTemp(t.TempDir(), "stdout")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	orig := os.Stdout
	os.Stdout = f
	t.Cleanup(func() { os.Stdout = orig; f.Close() })
	if isTerminal() {
		t.Errorf("isTerminal() = true for a regular file")
	}
}

func TestIsTerminalCharDevice(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	// /dev/null is a character device, so isTerminal() should report true.
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open %s: %v", os.DevNull, err)
	}
	orig := os.Stdout
	os.Stdout = devnull
	t.Cleanup(func() { os.Stdout = orig; devnull.Close() })
	if !isTerminal() {
		t.Errorf("isTerminal() = false for %s (a character device)", os.DevNull)
	}
}
