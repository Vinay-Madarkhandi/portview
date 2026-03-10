// Package types defines shared data types for the PortView application.
package types

import "fmt"

type PortInfo struct {
	Protocol string
	Port     int
	Address  string
	Process  string
	PID      int
}

func (p PortInfo) PortString() string {
	return fmt.Sprintf("%d", p.Port)
}

func (p PortInfo) PIDString() string {
	if p.PID == 0 {
		return "-"
	}
	return fmt.Sprintf("%d", p.PID)
}

func (p PortInfo) AddressPort() string {
	return fmt.Sprintf("%s:%d", p.Address, p.Port)
}
