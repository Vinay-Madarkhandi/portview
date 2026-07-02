// Package scanner provides port scanning utilities.
package scanner

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"

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
	case "windows":
		ports, err = scanWindows()
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

func scanWindows() ([]types.PortInfo, error) {
	path, err := findPowerShell()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(
		path,
		"-NoProfile",
		"-NonInteractive",
		"-ExecutionPolicy",
		"Bypass",
		"-Command",
		windowsScanScript,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outStr := strings.TrimSpace(string(output))
		if outStr == "" {
			return nil, nil
		}
		if isPermissionError(outStr) {
			return nil, fmt.Errorf("insufficient permissions: try running as Administrator")
		}
		return nil, fmt.Errorf("PowerShell port scan failed: %w (%s)", err, outStr)
	}

	return ParsePowerShellCSVOutput(string(output)), nil
}

func findPowerShell() (string, error) {
	for _, name := range []string{"powershell.exe", "powershell", "pwsh.exe", "pwsh"} {
		path, err := exec.LookPath(name)
		if err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("PowerShell command not found")
}

func isPermissionError(output string) bool {
	output = strings.ToLower(output)
	return strings.Contains(output, "access is denied") ||
		strings.Contains(output, "permission") ||
		strings.Contains(output, "operation not permitted")
}

func sortPorts(ports []types.PortInfo) {
	sort.Slice(ports, func(i, j int) bool {
		if ports[i].Port == ports[j].Port {
			return ports[i].Protocol < ports[j].Protocol
		}
		return ports[i].Port < ports[j].Port
	})
}

const windowsScanScript = `
$ErrorActionPreference = 'Stop'

function New-PortRow($protocol, $address, $port, $pid) {
    $processName = 'unknown'
    $process = Get-Process -Id $pid -ErrorAction SilentlyContinue
    if ($process) {
        $processName = $process.ProcessName
    }

    [PSCustomObject]@{
        Protocol = $protocol
        Address  = $address
        Port     = $port
        Process  = $processName
        PID      = $pid
    }
}

$rows = @()
$rows += Get-NetTCPConnection -State Listen -ErrorAction SilentlyContinue | ForEach-Object {
    New-PortRow 'tcp' $_.LocalAddress $_.LocalPort $_.OwningProcess
}
$rows += Get-NetUDPEndpoint -ErrorAction SilentlyContinue | ForEach-Object {
    New-PortRow 'udp' $_.LocalAddress $_.LocalPort $_.OwningProcess
}

$rows | ConvertTo-Csv -NoTypeInformation
`
