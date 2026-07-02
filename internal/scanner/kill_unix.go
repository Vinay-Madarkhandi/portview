//go:build !windows

package scanner

import (
	"fmt"
	"syscall"
)

func KillProcess(pid int) error {
	if pid <= 1 {
		return fmt.Errorf("refusing to kill PID %d (system process)", pid)
	}
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to terminate PID %d: %w", pid, err)
	}
	return nil
}
