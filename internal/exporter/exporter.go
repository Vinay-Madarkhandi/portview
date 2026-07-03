// Package exporter writes PortView port data to disk.
package exporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Vinay-Madarkhandi/portview/internal/types"
)

type Format string

const (
	FormatCSV  Format = "csv"
	FormatJSON Format = "json"
)

func Export(format Format, ports []types.PortInfo) (string, error) {
	if len(ports) == 0 {
		return "", fmt.Errorf("no ports to export")
	}

	path := defaultPath(format, time.Now())
	switch format {
	case FormatCSV:
		return path, WriteCSV(path, ports)
	case FormatJSON:
		return path, WriteJSON(path, ports)
	default:
		return "", fmt.Errorf("unsupported export format %q", format)
	}
}

func WriteCSV(path string, ports []types.PortInfo) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create CSV: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"protocol", "port", "address", "process", "pid"}); err != nil {
		return fmt.Errorf("write CSV header: %w", err)
	}
	for _, port := range ports {
		if err := writer.Write([]string{
			port.Protocol,
			strconv.Itoa(port.Port),
			port.Address,
			port.Process,
			strconv.Itoa(port.PID),
		}); err != nil {
			return fmt.Errorf("write CSV row: %w", err)
		}
	}
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush CSV: %w", err)
	}
	return nil
}

func WriteJSON(path string, ports []types.PortInfo) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create JSON: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(ports); err != nil {
		return fmt.Errorf("write JSON: %w", err)
	}
	return nil
}

func defaultPath(format Format, now time.Time) string {
	ext := string(format)
	if ext == "" {
		ext = "txt"
	}
	return filepath.Join(".", fmt.Sprintf("portview-ports-%s.%s", now.Format("20060102-150405"), ext))
}
