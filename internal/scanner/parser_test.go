package scanner

import (
	"testing"
)

func TestParseSSOutput_TCP(t *testing.T) {
	input := `tcp   LISTEN 0      128          0.0.0.0:22        0.0.0.0:*    users:(("sshd",pid=1234,fd=3))
tcp   LISTEN 0      128             [::]:22           [::]:*    users:(("sshd",pid=1234,fd=4))`

	ports := ParseSSOutput(input)

	if len(ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(ports))
	}

	p := ports[0]
	if p.Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", p.Protocol)
	}
	if p.Port != 22 {
		t.Errorf("expected port 22, got %d", p.Port)
	}
	if p.Process != "sshd" {
		t.Errorf("expected process sshd, got %s", p.Process)
	}
	if p.PID != 1234 {
		t.Errorf("expected PID 1234, got %d", p.PID)
	}
	if p.Address != "0.0.0.0" {
		t.Errorf("expected address 0.0.0.0, got %s", p.Address)
	}
}

func TestParseSSOutput_UDP(t *testing.T) {
	input := `udp   UNCONN 0      0            0.0.0.0:5353      0.0.0.0:*    users:(("avahi-daemon",pid=567,fd=12))`

	ports := ParseSSOutput(input)

	if len(ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(ports))
	}

	p := ports[0]
	if p.Protocol != "udp" {
		t.Errorf("expected protocol udp, got %s", p.Protocol)
	}
	if p.Port != 5353 {
		t.Errorf("expected port 5353, got %d", p.Port)
	}
	if p.Process != "avahi-daemon" {
		t.Errorf("expected process avahi-daemon, got %s", p.Process)
	}
	if p.PID != 567 {
		t.Errorf("expected PID 567, got %d", p.PID)
	}
}

func TestParseSSOutput_NoProcessInfo(t *testing.T) {
	input := `tcp   LISTEN 0      128          0.0.0.0:8080      0.0.0.0:*`

	ports := ParseSSOutput(input)

	if len(ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(ports))
	}

	p := ports[0]
	if p.Process != "unknown" {
		t.Errorf("expected process unknown, got %s", p.Process)
	}
	if p.PID != 0 {
		t.Errorf("expected PID 0, got %d", p.PID)
	}
}

func TestParseSSOutput_EmptyInput(t *testing.T) {
	ports := ParseSSOutput("")
	if len(ports) != 0 {
		t.Errorf("expected 0 ports, got %d", len(ports))
	}
}

func TestParseSSOutput_SkipsNonListening(t *testing.T) {
	input := `tcp   ESTAB  0      0      192.168.1.100:44556  142.250.80.78:443`

	ports := ParseSSOutput(input)
	if len(ports) != 0 {
		t.Errorf("expected 0 ports (ESTAB should be skipped), got %d", len(ports))
	}
}

func TestParseSSOutput_ShortLine(t *testing.T) {
	input := `tcp   LISTEN 0`

	ports := ParseSSOutput(input)
	if len(ports) != 0 {
		t.Errorf("expected 0 ports for short line, got %d", len(ports))
	}
}

func TestExtractPID(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`users:(("sshd",pid=1234,fd=3))`, 1234},
		{`users:(("nginx",pid=99,fd=6))`, 99},
		{`nopidhere`, 0},
		{`pid=abc,fd=1`, 0},
	}

	for _, tt := range tests {
		got := extractPID(tt.input)
		if got != tt.expected {
			t.Errorf("extractPID(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestExtractProcessName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`users:(("sshd",pid=1234,fd=3))`, "sshd"},
		{`users:(("my-app",pid=1,fd=1))`, "my-app"},
		{`noquoteshere`, "unknown"},
	}

	for _, tt := range tests {
		got := extractProcessName(tt.input)
		if got != tt.expected {
			t.Errorf("extractProcessName(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		input string
		addr  string
		port  int
	}{
		{"0.0.0.0:8080", "0.0.0.0", 8080},
		{"127.0.0.1:3000", "127.0.0.1", 3000},
		{"[::]:22", "::", 22},
		{"[::1]:443", "::1", 443},
		{"*:80", "*", 80},
	}

	for _, tt := range tests {
		addr, port := parseAddress(tt.input)
		if addr != tt.addr || port != tt.port {
			t.Errorf("parseAddress(%q) = (%q, %d), want (%q, %d)",
				tt.input, addr, port, tt.addr, tt.port)
		}
	}
}
