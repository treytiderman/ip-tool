package main

import (
	"encoding/json"
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

func Test_ParseInterfaces(t *testing.T) {

	var tests = []struct {
		input    string
		output   []InterfaceConfig
		hasError error
	}{
		{
			`
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
`,
			[]InterfaceConfig{
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
					Gateways:   []Gateway{},
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
					Gateways:   []Gateway{},
					DnsIsDhcp:  false,
					DnsServers: []string{},
				},
			},
			nil,
		},
		{
			`
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
`,
			[]InterfaceConfig{
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
						"8.8.8.8",
					},
				},
			},
			nil,
		},
		// 		{
		// 			`
		// Configuration for interface "Ethernet"
		//     DHCP enabled:                         Yes
		//     IP Address:                           192.168.1.109
		//     Subnet Prefix:                        192.168.1.0/24 (mask 255.255.255.0)
		//     Default Gateway:                      192.168.1.1
		//     Gateway Metric:                       0
		//     InterfaceMetric:                      25
		//     DNS servers configured through DHCP:  None
		//     Register with which suffix:           Primary only
		//     WINS servers configured through DHCP: None

		// Configuration for interface "Local Area Connection* 1"
		//     DHCP enabled:                         Yes
		//     InterfaceMetric:                      25
		//     DNS servers configured through DHCP:  None
		//     Register with which suffix:           Primary only
		//     WINS servers configured through DHCP: None

		// Configuration for interface "Local Area Connection* 2"
		//     DHCP enabled:                         Yes
		//     InterfaceMetric:                      25
		//     DNS servers configured through DHCP:  None
		//     Register with which suffix:           Primary only
		//     WINS servers configured through DHCP: None

		// Configuration for interface "Wi-Fi"
		//     DHCP enabled:                         Yes
		//     IP Address:                           10.52.214.89
		//     Subnet Prefix:                        10.52.208.0/20 (mask 255.255.240.0)
		//     Default Gateway:                      10.52.208.20
		//     Gateway Metric:                       0
		//     InterfaceMetric:                      35
		//     DNS servers configured through DHCP:  10.52.32.10
		//                                           172.16.32.10
		//     Register with which suffix:           Primary only
		//     WINS servers configured through DHCP: None

		// Configuration for interface "Bluetooth Network Connection"
		// 	DHCP enabled:                         Yes
		// 	InterfaceMetric:                      65
		// 	DNS servers configured through DHCP:  None
		// 	Register with which suffix:           Primary only
		// 	WINS servers configured through DHCP: None

		// Configuration for interface "Loopback Pseudo-Interface 1"
		//     DHCP enabled:                         No
		//     IP Address:                           127.0.0.1
		//     Subnet Prefix:                        127.0.0.0/8 (mask 255.0.0.0)
		//     InterfaceMetric:                      75
		//     Statically Configured DNS Servers:    None
		//     Register with which suffix:           None
		//     Statically Configured WINS Servers:   None
		// `,
		// 			[]Interface{
		// 				{
		// 					InterfaceName:   "Ethernet",
		// 					InterfaceMetric: 9,
		// 					IpIsDhcp:        false,
		// 					Ips: []Ip{
		// 						{
		// 							IpAddress:  "192.168.1.9",
		// 							SubnetMask: "255.255.255.0",
		// 							Cidr:       24,
		// 						},
		// 						{
		// 							IpAddress:  "192.168.1.63",
		// 							SubnetMask: "255.255.0.0",
		// 							Cidr:       16,
		// 						},
		// 						{
		// 							IpAddress:  "192.168.1.67",
		// 							SubnetMask: "255.0.0.0",
		// 							Cidr:       8,
		// 						},
		// 					},
		// 					Gateways: []Gateway{
		// 						{
		// 							GatewayAddress: "192.168.1.254",
		// 							GatewayMetric:  1,
		// 						},
		// 						{
		// 							GatewayAddress: "192.168.1.1",
		// 							GatewayMetric:  99,
		// 						},
		// 					},
		// 					DnsIsDhcp: false,
		// 					DnsServers: []string{
		// 						"192.168.1.1",
		// 						"1.1.1.1",
		// 						"8.8.8.8",
		// 					},
		// 				},
		// 			},
		// 			nil,
		// 		},
	}

	for _, test := range tests {

		split := strings.Split(test.input, "\n")
		interfaces, err := ParseInterfaces(split)
		if err != nil {
			t.Error(err)
		}

		jsonToPrint, err := json.MarshalIndent(interfaces, "", "  ")
		if err != nil {
			t.Error(err)
		}
		fmt.Println(test.input)
		fmt.Println(string(jsonToPrint))

		for i, iface := range interfaces {
			if test.hasError != nil {
				t.Errorf("expected no Error, got Error %v", test.hasError)
			}
			if iface.InterfaceName != test.output[i].InterfaceName {
				t.Errorf("expected InterfaceName %q, got %q", test.output[i].InterfaceName, iface.InterfaceName)
			}
			if iface.InterfaceMetric != test.output[i].InterfaceMetric {
				t.Errorf("expected InterfaceMetric %d, got %d", test.output[i].InterfaceMetric, iface.InterfaceMetric)
			}
			if iface.IpIsDhcp != test.output[i].IpIsDhcp {
				t.Errorf("expected IpIsDhcp %v, got %v", test.output[i].IpIsDhcp, iface.IpIsDhcp)
			}
			if iface.DnsIsDhcp != test.output[i].DnsIsDhcp {
				t.Errorf("expected DnsIsDhcp %v, got %v", test.output[i].DnsIsDhcp, iface.DnsIsDhcp)
			}

			if len(iface.Ips) != len(test.output[i].Ips) {
				t.Errorf("expected number Ips %d, got %d", len(test.output[i].Ips), len(iface.Ips))
			}
			for j, ipAndMask := range iface.Ips {
				if ipAndMask.IpAddress != test.output[i].Ips[j].IpAddress {
					t.Errorf("expected IpAddress %q, got %q", test.output[i].Ips[j].IpAddress, ipAndMask.IpAddress)
				}
				if ipAndMask.SubnetMask != test.output[i].Ips[j].SubnetMask {
					t.Errorf("expected SubnetMask %q, got %q", test.output[i].Ips[j].SubnetMask, ipAndMask.SubnetMask)
				}
				if ipAndMask.Cidr != test.output[i].Ips[j].Cidr {
					t.Errorf("expected Cidr %d, got %d", test.output[i].Ips[j].Cidr, ipAndMask.Cidr)
				}
			}

			if len(iface.Gateways) != len(test.output[i].Gateways) {
				t.Errorf("expected number Gateways %d, got %d", len(test.output[i].Gateways), len(iface.Gateways))
			}
			for j, Gateway := range iface.Gateways {
				if Gateway.GatewayAddress != test.output[i].Gateways[j].GatewayAddress {
					t.Errorf("expected GatewayAddress %q, got %q", test.output[i].Gateways[j].GatewayAddress, Gateway.GatewayAddress)
				}
				if Gateway.GatewayMetric != test.output[i].Gateways[j].GatewayMetric {
					t.Errorf("expected GatewayMetric %d, got %d", test.output[i].Gateways[j].GatewayMetric, Gateway.GatewayMetric)
				}
			}

			if len(iface.DnsServers) != len(test.output[i].DnsServers) {
				t.Errorf("expected number DnsServers %d, got %d", len(test.output[i].DnsServers), len(iface.DnsServers))
			}
			for j, dnsServer := range iface.DnsServers {
				if dnsServer != test.output[i].DnsServers[j] {
					t.Errorf("expected DnsServer %q, got %q", test.output[i].DnsServers[j], dnsServer)
				}
			}
		}
	}
}

func Test_ParseInterfaceStates(t *testing.T) {

	var tests = []struct {
		input    string
		output   []InterfaceState
		hasError error
	}{
		{
			`
Admin State    State          Type             Interface Name
-------------------------------------------------------------------------
Disabled       Disconnected   Dedicated        Ethernet 2
Enabled        Connected      Dedicated        Wi-Fi
Disabled       Disconnected   Dedicated        Ethernet     5
Enabled        Connected      Dedicated        Ethernet
`,
			[]InterfaceState{
				{
					InterfaceName: "Ethernet 2",
					InterfaceType: "Dedicated",
					Connected:     false,
					Disabled:      true,
				},
				{
					InterfaceName: "Wi-Fi",
					InterfaceType: "Dedicated",
					Connected:     true,
					Disabled:      false,
				},
				{
					InterfaceName: "Ethernet     5",
					InterfaceType: "Dedicated",
					Connected:     false,
					Disabled:      true,
				},
				{
					InterfaceName: "Ethernet",
					InterfaceType: "Dedicated",
					Connected:     true,
					Disabled:      false,
				},
			},
			nil,
		},
		{
			`
Admin State    State          Type             Interface Name
-------------------------------------------------------------------------
Enabled        Connected      Dedicated        Ethernet
Enabled        Disconnected   Dedicated        Ethernet 3
		`,
			[]InterfaceState{
				{
					InterfaceName: "Ethernet",
					InterfaceType: "Dedicated",
					Connected:     true,
					Disabled:      false,
				},
				{
					InterfaceName: "Ethernet 3",
					InterfaceType: "Dedicated",
					Connected:     false,
					Disabled:      false,
				},
			},
			nil,
		},
	}

	for _, test := range tests {

		split := strings.Split(test.input, "\n")
		interfaceStates, err := ParseInterfaceStates(split)
		if err != nil {
			t.Error(err)
		}

		jsonToPrint, err := json.MarshalIndent(interfaceStates, "", "  ")
		if err != nil {
			t.Error(err)
		}
		fmt.Println(test.input)
		fmt.Println(string(jsonToPrint))

		if len(test.output) != len(interfaceStates) {
			t.Errorf("expected number interfaces %d, got %d", len(test.output), len(interfaceStates))
		}
		for i, iface := range interfaceStates {
			if test.hasError != nil {
				t.Errorf("expected no Error, got Error %v", test.hasError)
			}
			if iface.InterfaceName != test.output[i].InterfaceName {
				t.Errorf("expected InterfaceName %q, got %q", test.output[i].InterfaceName, iface.InterfaceName)
			}
			if iface.InterfaceType != test.output[i].InterfaceType {
				t.Errorf("expected InterfaceType %q, got %q", test.output[i].InterfaceType, iface.InterfaceType)
			}
			if iface.Connected != test.output[i].Connected {
				t.Errorf("expected Connected %v, got %v", test.output[i].Connected, iface.Connected)
			}
			if iface.Disabled != test.output[i].Disabled {
				t.Errorf("expected Disabled %v, got %v", test.output[i].Disabled, iface.Disabled)
			}
		}
	}
}

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

	// Enable/Disable Interface
	// EnableInterface("Ethernet Anker")
	// DisableInterface("Ethernet Anker")

}
