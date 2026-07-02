package scanner

import (
	"strconv"
	"strings"

	"github.com/Vinay-Madarkhandi/portview/internal/types"
)

func ParseSSOutput(output string) []types.PortInfo {
	var ports []types.PortInfo

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		protocol := strings.ToLower(fields[0])
		state := strings.ToUpper(fields[1])

		if !isListeningState(protocol, state) {
			continue
		}

		localAddr := fields[4]
		address, port := parseAddress(localAddr)

		process, pid := parseProcessInfo(fields)

		ports = append(ports, types.PortInfo{
			Protocol: protocol,
			Port:     port,
			Address:  address,
			Process:  process,
			PID:      pid,
		})
	}

	return ports
}

func ParseLsofOutput(output string) []types.PortInfo {
	var ports []types.PortInfo

	currentPID := 0
	currentProcess := "unknown"
	currentProtocol := ""

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		prefix := line[0]
		value := ""
		if len(line) > 1 {
			value = line[1:]
		}

		switch prefix {
		case 'p':
			currentPID = parseInt(value)
			currentProcess = "unknown"
			currentProtocol = ""
		case 'c':
			if value != "" {
				currentProcess = value
			}
		case 'P':
			currentProtocol = strings.ToLower(value)
		case 'n':
			if currentProtocol != "tcp" && currentProtocol != "udp" {
				continue
			}

			address, port := parseLsofName(value)
			if port == 0 {
				continue
			}

			ports = append(ports, types.PortInfo{
				Protocol: currentProtocol,
				Port:     port,
				Address:  address,
				Process:  currentProcess,
				PID:      currentPID,
			})
		}
	}

	return ports
}

func isListeningState(protocol, state string) bool {
	switch protocol {
	case "tcp":
		return state == "LISTEN"
	case "udp":
		return state == "UNCONN"
	default:
		return false
	}
}

func parseAddress(localAddr string) (string, int) {
	// Find the last colon — port is always after the last ':'
	lastColon := strings.LastIndex(localAddr, ":")
	if lastColon < 0 {
		return localAddr, 0
	}

	address := localAddr[:lastColon]
	portStr := localAddr[lastColon+1:]

	// Clean up IPv6 bracket notation
	address = strings.TrimPrefix(address, "[")
	address = strings.TrimSuffix(address, "]")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return address, 0
	}

	return address, port
}

func parseLsofName(name string) (string, int) {
	name = strings.TrimSpace(name)
	if idx := strings.Index(name, " ("); idx >= 0 {
		name = strings.TrimSpace(name[:idx])
	}
	if idx := strings.Index(name, "->"); idx >= 0 {
		name = strings.TrimSpace(name[:idx])
	}

	fields := strings.Fields(name)
	if len(fields) > 1 {
		first := strings.ToLower(fields[0])
		if first == "tcp" || first == "udp" {
			name = fields[1]
		}
	}

	return parseAddress(name)
}

func parseProcessInfo(fields []string) (string, int) {
	process := "unknown"
	pid := 0

	if len(fields) < 7 {
		return process, pid
	}

	usersField := fields[len(fields)-1]
	if !strings.Contains(usersField, "pid=") {
		return process, pid
	}

	pid = extractPID(usersField)
	process = extractProcessName(usersField)

	return process, pid
}

func parseInt(value string) int {
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return n
}

func extractPID(usersField string) int {
	idx := strings.Index(usersField, "pid=")
	if idx < 0 {
		return 0
	}

	pidStr := usersField[idx+4:]

	// Find the end of the PID number
	end := strings.IndexAny(pidStr, ",)")
	if end > 0 {
		pidStr = pidStr[:end]
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0
	}

	return pid
}

func extractProcessName(usersField string) string {
	// Process name is between the first pair of quotes
	firstQuote := strings.Index(usersField, "\"")
	if firstQuote < 0 {
		return "unknown"
	}

	rest := usersField[firstQuote+1:]
	secondQuote := strings.Index(rest, "\"")
	if secondQuote <= 0 {
		return "unknown"
	}

	return rest[:secondQuote]
}
