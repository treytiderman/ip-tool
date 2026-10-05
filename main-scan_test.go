//go:build windows
// +build windows

package main

import (
	"net"
	"testing"
)

func TestParseScanSubnet(t *testing.T) {
	subnet, prefix, err := parseScanSubnet("192.168.1.42", "255.255.255.0")
	if err != nil {
		t.Fatalf("parseScanSubnet() error = %v", err)
	}
	if prefix != 24 {
		t.Errorf("parseScanSubnet() prefix = %d, want 24", prefix)
	}
	if got, want := subnet.String(), "192.168.1.0/24"; got != want {
		t.Errorf("parseScanSubnet() subnet = %q, want %q", got, want)
	}
	if !subnet.Contains(net.ParseIP("192.168.1.200")) {
		t.Error("parsed subnet does not include addresses in the selected network")
	}
}

func TestParseScanSubnetAllowsMaximumSize(t *testing.T) {
	_, prefix, err := parseScanSubnet("192.168.1.42", "255.255.0.0")
	if err != nil {
		t.Fatalf("parseScanSubnet() rejected a /16 subnet: %v", err)
	}
	if prefix != 16 {
		t.Errorf("parseScanSubnet() prefix = %d, want 16", prefix)
	}
}

func TestParseScanSubnetAllowsSlash17(t *testing.T) {
	_, prefix, err := parseScanSubnet("192.168.1.42", "255.255.128.0")
	if err != nil {
		t.Fatalf("parseScanSubnet() rejected a /17 subnet: %v", err)
	}
	if prefix != 17 {
		t.Errorf("parseScanSubnet() prefix = %d, want 17", prefix)
	}
}

func TestParseScanSubnetRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		mask string
	}{
		{name: "invalid address", ip: "not-an-ip", mask: "255.255.255.0"},
		{name: "non-contiguous mask", ip: "192.168.1.42", mask: "255.0.255.0"},
		{name: "subnet exceeds scan limit", ip: "10.0.0.1", mask: "255.254.0.0"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, _, err := parseScanSubnet(test.ip, test.mask); err == nil {
				t.Fatal("parseScanSubnet() expected an error")
			}
		})
	}
}

func TestLookupMACVendor(t *testing.T) {
	tests := []struct {
		name     string
		mac      string
		expected string
	}{
		{name: "IEEE assignment", mac: "00:1b:63:00:00:01", expected: "Apple, Inc."},
		{name: "MA-M assignment", mac: "c8:5c:e2:70:00:00", expected: "SYNERGY SYSTEMS AND SOLUTIONS"},
		{name: "MA-S assignment", mac: "8c:1f:64:af:a0:00", expected: "DATA ELECTRONIC DEVICES, INC"},
		{name: "locally administered", mac: "02:00:00:00:00:01", expected: "Locally administered"},
		{name: "unknown assignment", mac: "00:08:33:00:00:01", expected: "Unknown"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vendor, err := lookupMACVendor(test.mac)
			if err != nil {
				t.Fatalf("lookupMACVendor() error = %v", err)
			}
			if vendor != test.expected {
				t.Errorf("lookupMACVendor() = %q, want %q", vendor, test.expected)
			}
		})
	}
}
