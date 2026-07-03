package clipboard

import (
	"os/exec"
	"testing"
)

func TestClipboardCommand(t *testing.T) {
	originalLookPath := lookPath
	t.Cleanup(func() { lookPath = originalLookPath })

	available := map[string]bool{
		"pbcopy":  true,
		"cmd":     true,
		"wl-copy": true,
	}
	lookPath = func(file string) (string, error) {
		if available[file] {
			return "/usr/bin/" + file, nil
		}
		return "", exec.ErrNotFound
	}

	tests := []struct {
		goos string
		name string
		args []string
	}{
		{goos: "darwin", name: "pbcopy"},
		{goos: "windows", name: "cmd", args: []string{"/c", "clip"}},
		{goos: "linux", name: "wl-copy"},
	}

	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			cmd, err := clipboardCommand(tt.goos)
			if err != nil {
				t.Fatalf("clipboardCommand() error = %v", err)
			}
			if cmd.name != tt.name {
				t.Fatalf("command name = %q, want %q", cmd.name, tt.name)
			}
			if len(cmd.args) != len(tt.args) {
				t.Fatalf("args = %v, want %v", cmd.args, tt.args)
			}
			for i := range cmd.args {
				if cmd.args[i] != tt.args[i] {
					t.Fatalf("args = %v, want %v", cmd.args, tt.args)
				}
			}
		})
	}
}

func TestClipboardCommandLinuxFallback(t *testing.T) {
	originalLookPath := lookPath
	t.Cleanup(func() { lookPath = originalLookPath })

	lookPath = func(file string) (string, error) {
		if file == "xclip" {
			return "/usr/bin/xclip", nil
		}
		return "", exec.ErrNotFound
	}

	cmd, err := clipboardCommand("linux")
	if err != nil {
		t.Fatalf("clipboardCommand() error = %v", err)
	}
	if cmd.name != "xclip" {
		t.Fatalf("command name = %q, want xclip", cmd.name)
	}
}

func TestClipboardCommandErrors(t *testing.T) {
	originalLookPath := lookPath
	t.Cleanup(func() { lookPath = originalLookPath })

	lookPath = func(file string) (string, error) {
		return "", exec.ErrNotFound
	}

	if _, err := clipboardCommand("linux"); err == nil {
		t.Fatal("expected Linux error when no clipboard commands are installed")
	}
	if _, err := clipboardCommand("plan9"); err == nil {
		t.Fatal("expected unsupported OS error")
	}
	if err := Write("   "); err == nil {
		t.Fatal("expected empty clipboard write error")
	} else if err.Error() == "" {
		t.Fatal("expected non-empty error")
	}
}
