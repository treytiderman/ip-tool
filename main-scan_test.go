//go:build windows
// +build windows

package main

import (
	"net"
	"testing"

	"golang.org/x/net/dns/dnsmessage"
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

func TestMDNSDiscoveryBuildsServiceResults(t *testing.T) {
	const (
		instance = "_living-room._http._tcp.local."
		hostname = "speaker.local."
	)
	discovery := mdnsDiscovery{
		hosts:           make(map[string]*mdnsHost),
		instances:       make(map[string]mdnsServiceInstance),
		instanceDetails: make(map[string]map[string]struct{}),
	}
	message := dnsmessage.Message{
		Header: dnsmessage.Header{Response: true},
		Answers: []dnsmessage.Resource{
			{
				Header: dnsmessage.ResourceHeader{Name: dnsmessage.MustNewName("_http._tcp.local.")},
				Body:   &dnsmessage.PTRResource{PTR: dnsmessage.MustNewName(instance)},
			},
			{
				Header: dnsmessage.ResourceHeader{Name: dnsmessage.MustNewName(instance)},
				Body:   &dnsmessage.SRVResource{Target: dnsmessage.MustNewName(hostname), Port: 8080},
			},
			{
				Header: dnsmessage.ResourceHeader{Name: dnsmessage.MustNewName(instance)},
				Body:   &dnsmessage.TXTResource{TXT: []string{"model=example", "version=1"}},
			},
			{
				Header: dnsmessage.ResourceHeader{Name: dnsmessage.MustNewName(hostname)},
				Body:   &dnsmessage.AResource{A: [4]byte{192, 168, 1, 25}},
			},
		},
	}

	discovery.addMessage(message)
	results := discovery.results()
	if len(results) != 1 {
		t.Fatalf("got %d devices, want 1", len(results))
	}
	result := results[0]
	if result.IPAddresses == nil || result.Services == nil || result.Details == nil {
		t.Fatal("mDNS result collections must be empty arrays, not nil")
	}
	if result.Hostname != "speaker.local" {
		t.Errorf("hostname = %q, want %q", result.Hostname, "speaker.local")
	}
	if len(result.IPAddresses) != 1 || result.IPAddresses[0] != "192.168.1.25" {
		t.Errorf("addresses = %v, want [192.168.1.25]", result.IPAddresses)
	}
	if len(result.Services) != 1 || result.Services[0] != "_http._tcp.local (8080)" {
		t.Errorf("services = %v, want [_http._tcp.local (8080)]", result.Services)
	}
	if len(result.Details) != 3 {
		t.Errorf("details = %v, want SRV instance and 2 TXT values", result.Details)
	}
}

func TestPackMDNSQueryAcceptsNameWithoutTrailingDot(t *testing.T) {
	packet, err := packMDNSQuery("_nvstream._tcp.local")
	if err != nil {
		t.Fatalf("packMDNSQuery() error = %v", err)
	}

	var message dnsmessage.Message
	if err := message.Unpack(packet); err != nil {
		t.Fatalf("unpack query: %v", err)
	}
	if got, want := message.Questions[0].Name.String(), "_nvstream._tcp.local."; got != want {
		t.Errorf("query name = %q, want %q", got, want)
	}
}
