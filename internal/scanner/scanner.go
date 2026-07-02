// Package scanner provides port scanning utilities.
package scanner

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"github.com/Vinay-Madarkhandi/portview/internal/types"
)

func ScanPorts() ([]types.PortInfo, error) {
	var (
		ports []types.PortInfo
		err   error
	)

	switch runtime.GOOS {
	case "darwin":
		ports, err = scanDarwin()
	default:
		ports, err = scanLinux()
	}
	if err != nil {
		return nil, err
	}

	sortPorts(ports)
	return ports, nil
}

func scanLinux() ([]types.PortInfo, error) {
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
	return ports, nil
}

func scanDarwin() ([]types.PortInfo, error) {
	path, err := exec.LookPath("lsof")
	if err != nil {
		return nil, fmt.Errorf("lsof command not found: %w", err)
	}

	cmd := exec.Command(path, "-nP", "-iTCP", "-sTCP:LISTEN", "-iUDP", "-FpcPn")
	output, err := cmd.CombinedOutput()
	if err != nil {
		outStr := strings.TrimSpace(string(output))
		if outStr == "" {
			return nil, nil
		}
		if strings.Contains(outStr, "permission") || strings.Contains(outStr, "Operation not permitted") {
			return nil, fmt.Errorf("insufficient permissions: try running with sudo")
		}
		return nil, fmt.Errorf("lsof command failed: %w (%s)", err, outStr)
	}

	return ParseLsofOutput(string(output)), nil
}

func sortPorts(ports []types.PortInfo) {
	sort.Slice(ports, func(i, j int) bool {
		if ports[i].Port == ports[j].Port {
			return ports[i].Protocol < ports[j].Protocol
		}
		return ports[i].Port < ports[j].Port
	})
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
