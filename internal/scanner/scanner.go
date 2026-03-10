// Package scanner provides port scanning utilities using the Linux ss command.
package scanner

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"syscall"

	"portview/internal/types"
)

func ScanPorts() ([]types.PortInfo, error) {
	path, err := exec.LookPath("ss")
	if err != nil {
		return nil, fmt.Errorf("ss command not found: %w (is iproute2 installed?)", err)
	}

	cmd := exec.Command(path, "-tulpnH")
	output, err := cmd.CombinedOutput()
	if err != nil {
		outStr := strings.TrimSpace(string(output))
		if strings.Contains(outStr, "permission") || strings.Contains(outStr, "Operation not permitted") {
			return nil, fmt.Errorf("insufficient permissions: try running with sudo")
		}
		return nil, fmt.Errorf("ss command failed: %w (%s)", err, outStr)
	}

	ports := ParseSSOutput(string(output))

	sort.Slice(ports, func(i, j int) bool {
		if ports[i].Port == ports[j].Port {
			return ports[i].Protocol < ports[j].Protocol
		}
		return ports[i].Port < ports[j].Port
	})

	return ports, nil
}

func KillProcess(pid int) error {
	if pid <= 1 {
		return fmt.Errorf("refusing to kill PID %d (system process)", pid)
	}
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to kill PID %d: %w", pid, err)
	}
	return nil
}
