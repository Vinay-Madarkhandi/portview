// Package clipboard provides small cross-platform clipboard helpers.
package clipboard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type commandSpec struct {
	name string
	args []string
}

var lookPath = exec.LookPath

func Write(text string) error {
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("nothing to copy")
	}

	cmd, err := clipboardCommand(runtime.GOOS)
	if err != nil {
		return err
	}

	c := exec.Command(cmd.name, cmd.args...)
	c.Stdin = strings.NewReader(text)
	if output, err := c.CombinedOutput(); err != nil {
		out := strings.TrimSpace(string(output))
		if out != "" {
			return fmt.Errorf("clipboard command failed: %w (%s)", err, out)
		}
		return fmt.Errorf("clipboard command failed: %w", err)
	}
	return nil
}

func clipboardCommand(goos string) (commandSpec, error) {
	switch goos {
	case "darwin":
		return requireCommand(commandSpec{name: "pbcopy"})
	case "windows":
		return requireCommand(commandSpec{name: "cmd", args: []string{"/c", "clip"}})
	case "linux":
		return firstAvailable(
			commandSpec{name: "wl-copy"},
			commandSpec{name: "xclip", args: []string{"-selection", "clipboard"}},
			commandSpec{name: "xsel", args: []string{"--clipboard", "--input"}},
		)
	default:
		return commandSpec{}, fmt.Errorf("clipboard unsupported on %s", goos)
	}
}

func firstAvailable(commands ...commandSpec) (commandSpec, error) {
	for _, cmd := range commands {
		if _, err := lookPath(cmd.name); err == nil {
			return cmd, nil
		}
	}
	return commandSpec{}, fmt.Errorf("no clipboard command found: install wl-clipboard, xclip, or xsel")
}

func requireCommand(cmd commandSpec) (commandSpec, error) {
	if _, err := lookPath(cmd.name); err != nil {
		return commandSpec{}, fmt.Errorf("%s command not found", cmd.name)
	}
	return cmd, nil
}
