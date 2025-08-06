package main

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Cmd(t *testing.T) {
	lines, err := Cmd("cmd", "echo", "test")
	if err != nil {
		t.Error(err)
	}
	for _, line := range lines {
		t.Log(line)
		fmt.Println(line)
	}
}

func Test_IsAdmin(t *testing.T) {
	adminStatus := IsAdmin()
	fmt.Println("adminStatus", adminStatus)
}

func Test_IsIPv4Addr(t *testing.T) {
	var tests = []struct {
		expected bool
		ip       string
	}{
		{true, "0.0.0.0"},
		{true, "1.1.1.1"},
		{true, "10.10.0.100"},
		{true, "0.0.0.255"},
		{true, "255.255.255.255"},
		{false, "192.168.1"},
		{false, "192.168.1.300"},
		{false, "192.168.x.1"},
		{false, "192.168.010.001"},
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			result := IsIPv4Addr(tt.ip)
			if result != tt.expected {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func Test_ParseInterfaces(t *testing.T) {

	// --- Example 1 ---

	example1 := `
Configuration for interface "Ethernet"
    DHCP enabled:                         Yes
    IP Address:                           192.168.1.9
    Subnet Prefix:                        192.168.1.0/24 (mask 255.255.255.0)
    Default Gateway:                      192.168.1.254
    Gateway Metric:                       0
    InterfaceMetric:                      25
    DNS servers configured through DHCP:  192.168.1.1
    Register with which suffix:           Primary only
    WINS servers configured through DHCP: None

Configuration for interface "Loopback Pseudo-Interface 1"
    DHCP enabled:                         No
    IP Address:                           127.0.0.1
    Subnet Prefix:                        127.0.0.0/8 (mask 255.0.0.0)
    InterfaceMetric:                      75
    Statically Configured DNS Servers:    None
    Register with which suffix:           Primary only
    Statically Configured WINS Servers:   None

Configuration for interface "vEthernet (WSL (Hyper-V firewall))"
    DHCP enabled:                         No
    IP Address:                           172.30.96.1
    Subnet Prefix:                        172.30.96.0/20 (mask 255.255.240.0)
    InterfaceMetric:                      5000
    Statically Configured DNS Servers:    None
    Register with which suffix:           None
    Statically Configured WINS Servers:   None
`

	split1 := strings.Split(example1, "\n")
	interfaces1, err := ParseInterfaces(split1)
	if err != nil {
		t.Error(err)
	}

	// b1, err := json.MarshalIndent(interfaces1, "", "  ")
	// if err != nil {
	// 	t.Error(err)
	// }
	// fmt.Println(string(b1))

	expected1 := []Interface{
		{
			InterfaceName:   "Ethernet",
			InterfaceMetric: 25,
			IpIsDhcp:        true,
			Ips: []Ip{
				{
					IpAddress:  "192.168.1.9",
					SubnetMask: "255.255.255.0",
					Cidr:       24,
				},
			},
			Gateways: []Gateway{
				{
					GatewayAddress: "192.168.1.254",
					GatewayMetric:  0,
				},
			},
			DnsIsDhcp: true,
			DnsServers: []string{
				"192.168.1.1",
			},
		},
		{
			InterfaceName:   "Loopback Pseudo-Interface 1",
			InterfaceMetric: 75,
			IpIsDhcp:        false,
			Ips: []Ip{
				{
					IpAddress:  "127.0.0.1",
					SubnetMask: "255.0.0.0",
					Cidr:       8,
				},
			},
			Gateways: []Gateway{
				{
					GatewayAddress: "",
					GatewayMetric:  -1,
				},
			},
			DnsIsDhcp:  false,
			DnsServers: []string{},
		},
		{
			InterfaceName:   "vEthernet (WSL (Hyper-V firewall))",
			InterfaceMetric: 5000,
			IpIsDhcp:        false,
			Ips: []Ip{
				{
					IpAddress:  "172.30.96.1",
					SubnetMask: "255.255.240.0",
					Cidr:       20,
				},
			},
			Gateways: []Gateway{
				{
					GatewayAddress: "",
					GatewayMetric:  -1,
				},
			},
			DnsIsDhcp:  false,
			DnsServers: []string{},
		},
	}

	for i, iface := range interfaces1 {
		if iface.InterfaceName != expected1[i].InterfaceName {
			t.Errorf("expected InterfaceName %q, got %q", expected1[i].InterfaceName, iface.InterfaceName)
		}
		if iface.InterfaceMetric != expected1[i].InterfaceMetric {
			t.Errorf("expected InterfaceMetric %d, got %d", expected1[i].InterfaceMetric, iface.InterfaceMetric)
		}
		if iface.IpIsDhcp != expected1[i].IpIsDhcp {
			t.Errorf("expected IpIsDhcp %v, got %v", expected1[i].IpIsDhcp, iface.IpIsDhcp)
		}
		if iface.DnsIsDhcp != expected1[i].DnsIsDhcp {
			t.Errorf("expected DnsIsDhcp %v, got %v", expected1[i].DnsIsDhcp, iface.DnsIsDhcp)
		}
		for j, ipAndMask := range iface.Ips {
			if ipAndMask.IpAddress != expected1[i].Ips[j].IpAddress {
				t.Errorf("expected IpAddress %q, got %q", expected1[i].Ips[j].IpAddress, ipAndMask.IpAddress)
			}
			if ipAndMask.SubnetMask != expected1[i].Ips[j].SubnetMask {
				t.Errorf("expected SubnetMask %q, got %q", expected1[i].Ips[j].SubnetMask, ipAndMask.SubnetMask)
			}
			if ipAndMask.Cidr != expected1[i].Ips[j].Cidr {
				t.Errorf("expected Cidr %d, got %d", expected1[i].Ips[j].Cidr, ipAndMask.Cidr)
			}
		}
		for j, Gateway := range iface.Gateways {
			if Gateway.GatewayAddress != expected1[i].Gateways[j].GatewayAddress {
				t.Errorf("expected GatewayAddress %q, got %q", expected1[i].Gateways[j].GatewayAddress, Gateway.GatewayAddress)
			}
			if Gateway.GatewayMetric != expected1[i].Gateways[j].GatewayMetric {
				t.Errorf("expected GatewayMetric %d, got %d", expected1[i].Gateways[j].GatewayMetric, Gateway.GatewayMetric)
			}
		}
		for j, dnsServer := range iface.DnsServers {
			if dnsServer != expected1[i].DnsServers[j] {
				t.Errorf("expected DnsServer %q, got %q", expected1[i].DnsServers[j], dnsServer)
			}
		}
	}

	// --- Example 2 ---

	example2 := `
Configuration for interface "Ethernet"
    DHCP enabled:                         No
    IP Address:                           192.168.1.9
    Subnet Prefix:                        192.168.1.0/24 (mask 255.255.255.0)
    IP Address:                           192.168.1.63
    Subnet Prefix:                        192.168.0.0/16 (mask 255.255.0.0)
    IP Address:                           192.168.1.67
    Subnet Prefix:                        192.0.0.0/8 (mask 255.0.0.0)
    Default Gateway:                      192.168.1.254
    Gateway Metric:                       1
    Default Gateway:                      192.168.1.1
    Gateway Metric:                       99
    InterfaceMetric:                      9
    Statically Configured DNS Servers:    192.168.1.1
                                          1.1.1.1
                                          8.8.8.8
    Register with which suffix:           Primary only
    Statically Configured WINS Servers:   None
`

	split2 := strings.Split(example2, "\n")
	interfaces2, err := ParseInterfaces(split2)
	if err != nil {
		t.Error(err)
	}

	// b2, err := json.MarshalIndent(interfaces2, "", "  ")
	// if err != nil {
	// 	t.Error(err)
	// }
	// fmt.Println(string(b2))

	expected2 := []Interface{
		{
			InterfaceName:   "Ethernet",
			InterfaceMetric: 9,
			IpIsDhcp:        false,
			Ips: []Ip{
				{
					IpAddress:  "192.168.1.9",
					SubnetMask: "255.255.255.0",
					Cidr:       24,
				},
				{
					IpAddress:  "192.168.1.63",
					SubnetMask: "255.255.0.0",
					Cidr:       16,
				},
				{
					IpAddress:  "192.168.1.67",
					SubnetMask: "255.0.0.0",
					Cidr:       8,
				},
			},
			Gateways: []Gateway{
				{
					GatewayAddress: "192.168.1.254",
					GatewayMetric:  1,
				},
				{
					GatewayAddress: "192.168.1.1",
					GatewayMetric:  99,
				},
			},
			DnsIsDhcp: false,
			DnsServers: []string{
				"192.168.1.1",
				"1.1.1.1",
			},
		},
	}

	for i, iface := range interfaces2 {
		if iface.InterfaceName != expected2[i].InterfaceName {
			t.Errorf("expected InterfaceName %q, got %q", expected2[i].InterfaceName, iface.InterfaceName)
		}
		if iface.InterfaceMetric != expected2[i].InterfaceMetric {
			t.Errorf("expected InterfaceMetric %d, got %d", expected2[i].InterfaceMetric, iface.InterfaceMetric)
		}
		if iface.IpIsDhcp != expected2[i].IpIsDhcp {
			t.Errorf("expected IpIsDhcp %v, got %v", expected2[i].IpIsDhcp, iface.IpIsDhcp)
		}
		if iface.DnsIsDhcp != expected2[i].DnsIsDhcp {
			t.Errorf("expected DnsIsDhcp %v, got %v", expected2[i].DnsIsDhcp, iface.DnsIsDhcp)
		}
		for j, ipAndMask := range iface.Ips {
			if ipAndMask.IpAddress != expected2[i].Ips[j].IpAddress {
				t.Errorf("expected IpAddress %q, got %q", expected2[i].Ips[j].IpAddress, ipAndMask.IpAddress)
			}
			if ipAndMask.SubnetMask != expected2[i].Ips[j].SubnetMask {
				t.Errorf("expected SubnetMask %q, got %q", expected2[i].Ips[j].SubnetMask, ipAndMask.SubnetMask)
			}
			if ipAndMask.Cidr != expected2[i].Ips[j].Cidr {
				t.Errorf("expected Cidr %d, got %d", expected2[i].Ips[j].Cidr, ipAndMask.Cidr)
			}
		}
		for j, Gateway := range iface.Gateways {
			if Gateway.GatewayAddress != expected2[i].Gateways[j].GatewayAddress {
				t.Errorf("expected GatewayAddress %q, got %q", expected2[i].Gateways[j].GatewayAddress, Gateway.GatewayAddress)
			}
			if Gateway.GatewayMetric != expected2[i].Gateways[j].GatewayMetric {
				t.Errorf("expected GatewayMetric %d, got %d", expected2[i].Gateways[j].GatewayMetric, Gateway.GatewayMetric)
			}
		}
		for k, dnsServer := range expected2[i].DnsServers {
			if dnsServer != iface.DnsServers[k] {
				t.Errorf("expected DnsServer %q, got %q", dnsServer, iface.DnsServers[k])
			}
		}
	}

}

func Test_CidrToMask(t *testing.T) {
	var tests = []struct {
		cidr     string
		expected string
		hasError bool
	}{
		{"/1", "128.0.0.0", false},
		{"/8", "255.0.0.0", false},
		{"/16", "255.255.0.0", false},
		{"/20", "255.255.240.0", false},
		{"/24", "255.255.255.0", false},
		{"/32", "255.255.255.255", false},
		{"24", "255.255.255.0", false},
		{"/0", "", true},
		{"/33", "", true},
		{"/-1", "", true},
		{"abc", "", true},
		{"/abc", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.cidr, func(t *testing.T) {
			mask, err := CidrToMask(tt.cidr)
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error for input %q, got none", tt.cidr)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", tt.cidr, err)
				}
				if mask != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, mask)
				}
			}
		})
	}
}

// func TestNotReal()  {
// 	example1 := `
// Admin State    State          Type             Interface Name
// -------------------------------------------------------------------------
// Enabled        Connected      Dedicated        Ethernet
// Enabled        Disconnected   Dedicated        Ethernet 3
// `
// }

func Test_Manual(t *testing.T) {
	adminStatus := IsAdmin()
	if adminStatus == false {
		t.Error("Not administrator")
	}

	// Static Series
	// SetIpStatic("Ethernet", "192.168.1.27", "255.255.255.0", "192.168.1.254")
	// SetDnsStatic("Ethernet", "192.168.1.1")
	// AddIpStatic("Ethernet", "192.168.1.26", "255.255.0.0")
	// AddIpStatic("Ethernet", "192.168.1.28", "/8")
	// AddDnsStatic("Ethernet", "1.1.1.1", "2")
	// AddDnsStatic("Ethernet", "8.8.8.8", "3")
	// SetInterfaceMetric("Ethernet", "9")

	// DHCP
	// SetIpDhcp("Ethernet")
	// SetDnsDhcp("Ethernet")
	// SetInterfaceMetricAuto("Ethernet")

	// Rename Interface
	// SetInterfaceName("Ethernet", "Ethernet 2")
	// SetInterfaceName("Ethernet 2", "Ethernet")

}
