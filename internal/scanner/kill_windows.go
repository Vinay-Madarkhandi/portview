//go:build windows

package scanner

import (
	"fmt"
	"os"
)

func KillProcess(pid int) error {
	if pid <= 4 {
		return fmt.Errorf("refusing to kill PID %d (system process)", pid)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find PID %d: %w", pid, err)
	}
	if err := process.Kill(); err != nil {
		return fmt.Errorf("failed to terminate PID %d: %w", pid, err)
	}
	return nil
}
