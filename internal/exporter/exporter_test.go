package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Vinay-Madarkhandi/portview/internal/types"
)

func TestWriteCSV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ports.csv")
	ports := []types.PortInfo{
		{Protocol: "tcp", Port: 8080, Address: "127.0.0.1", Process: "nginx", PID: 1234},
	}

	if err := WriteCSV(path, ports); err != nil {
		t.Fatalf("WriteCSV() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	got := string(data)
	if !strings.Contains(got, "protocol,port,address,process,pid") {
		t.Fatalf("missing CSV header: %q", got)
	}
	if !strings.Contains(got, "tcp,8080,127.0.0.1,nginx,1234") {
		t.Fatalf("missing CSV row: %q", got)
	}
}

func TestWriteJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ports.json")
	ports := []types.PortInfo{
		{Protocol: "udp", Port: 53, Address: "0.0.0.0", Process: "dns", PID: 99},
	}

	if err := WriteJSON(path, ports); err != nil {
		t.Fatalf("WriteJSON() error = %v", err)
	}

	var got []types.PortInfo
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(got) != 1 || got[0].Port != 53 {
		t.Fatalf("decoded ports = %+v", got)
	}
}

func TestExportRejectsEmptyPorts(t *testing.T) {
	if _, err := Export(FormatCSV, nil); err == nil {
		t.Fatal("expected empty export error")
	}
}

func TestDefaultPath(t *testing.T) {
	got := defaultPath(FormatJSON, time.Date(2026, 7, 3, 12, 34, 56, 0, time.UTC))
	want := filepath.Join(".", "portview-ports-20260703-123456.json")
	if got != want {
		t.Fatalf("defaultPath() = %q, want %q", got, want)
	}
}
