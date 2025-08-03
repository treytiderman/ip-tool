//go:build windows
// +build windows

package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type Ip struct {
	IpAddress  string `json:"ip_address"`
	SubnetMask string `json:"subnet_mask"`
	Cidr       int    `json:"cidr"`
}

type Gateway struct {
	GatewayAddress string `json:"gateway_address"`
	GatewayMetric  int    `json:"gateway_metric"`
}

type Interface struct {
	InterfaceName   string    `json:"interface_name"`
	InterfaceMetric int       `json:"interface_metric"`
	IpIsDhcp        bool      `json:"ip_is_dhcp"`
	Ips             []Ip      `json:"ips"`
	Gateways        []Gateway `json:"gateways"`
	DnsIsDhcp       bool      `json:"dns_is_dhcp"`
	DnsServers      []string  `json:"dns_servers"`
}

type Route struct {
	Destination string `json:"destination"`
	Netmask     string `json:"netmask"`
	Gateway     string `json:"gateway"`
	Interface   string `json:"interface"`
	Metric      int    `json:"metric"`
	Persistent  bool   `json:"persistent"`
}

func Cmd(prog string, args ...string) ([]string, error) {
	// "cmd.exe", "/C",
	cmd := exec.Command(prog, args...)
	// cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func IsAdmin() bool {
	lines, err := Cmd("net", "session")
	return err == nil && len(lines) >= 0
}

func IsIPv4Addr(ip string) bool {
	parsed := net.ParseIP(ip)
	return parsed != nil && parsed.To4() != nil
}

func CidrToMask(cidrString string) (string, error) {
	cidrInt, err := strconv.Atoi(strings.TrimPrefix(cidrString, "/"))
	if err != nil || cidrInt < 1 || cidrInt > 32 {
		return "", fmt.Errorf("invalid CIDR prefix: %s", cidrString)
	}
	mask := net.CIDRMask(cidrInt, 32)
	return net.IP(mask).String(), nil
}

func ParseInterfaces(lines []string) ([]Interface, error) {
	nics := []Interface{}
	curr := Interface{}
	currIP := Ip{}
	currGateway := Gateway{}
	newNic := false
	anotherDnsServerIsPossible := false

	reset := func() {
		curr = Interface{InterfaceMetric: -1, Ips: []Ip{}, Gateways: []Gateway{}, DnsServers: []string{}}
		currIP = Ip{}
	}
	reset()

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		switch {
		case strings.HasPrefix(line, "Configuration for interface"):
			newNic = true
			reset()
			curr.InterfaceName = strings.Trim(strings.TrimPrefix(line, "Configuration for interface "), "\"")

		case strings.HasPrefix(line, "DHCP enabled:"):
			val := strings.TrimSpace(strings.TrimPrefix(line, "DHCP enabled:"))
			curr.IpIsDhcp = (val == "Yes")

		case strings.HasPrefix(line, "InterfaceMetric:"):
			val := strings.TrimSpace(strings.TrimPrefix(line, "InterfaceMetric:"))
			i, err := strconv.Atoi(val)
			if err == nil {
				curr.InterfaceMetric = i
			}

		case strings.HasPrefix(line, "IP Address:"):
			currIP.IpAddress = strings.TrimSpace(strings.TrimPrefix(line, "IP Address:"))

		case strings.HasPrefix(line, "Subnet Prefix:"):
			parts := strings.Split(line, "(mask ")
			if len(parts) > 1 {
				currIP.SubnetMask = strings.TrimSuffix(strings.TrimSpace(parts[1]), ")")
			}
			cidrPart := strings.Split(line, "/")
			if len(cidrPart) > 1 {
				cidrVal, err := strconv.Atoi(strings.Split(cidrPart[1], " ")[0])
				if err == nil {
					currIP.Cidr = cidrVal
				}
			}
			curr.Ips = append(curr.Ips, currIP)
			currIP = Ip{}

		case strings.HasPrefix(line, "Default Gateway:"):
			currGateway = Gateway{}
			currGateway.GatewayAddress = strings.TrimSpace(strings.TrimPrefix(line, "Default Gateway:"))

		case strings.HasPrefix(line, "Gateway Metric:"):
			val := strings.TrimSpace(strings.TrimPrefix(line, "Gateway Metric:"))
			i, err := strconv.Atoi(val)
			if err == nil {
				currGateway.GatewayMetric = i
			}
			curr.Gateways = append(curr.Gateways, currGateway)

		case strings.Contains(line, "DNS servers configured through DHCP:"):
			curr.DnsIsDhcp = true
			val := strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			if IsIPv4Addr(val) {
				curr.DnsServers = append(curr.DnsServers, val)
				anotherDnsServerIsPossible = true
			}

		case strings.Contains(line, "Statically Configured DNS Servers:"):
			curr.DnsIsDhcp = false
			val := strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			if IsIPv4Addr(val) {
				curr.DnsServers = append(curr.DnsServers, val)
				anotherDnsServerIsPossible = true
			}

		case line == "" && newNic:
			nics = append(nics, curr)
			newNic = false
			reset()
			anotherDnsServerIsPossible = false

		default:
			if anotherDnsServerIsPossible && IsIPv4Addr(line) {
				curr.DnsServers = append(curr.DnsServers, line)
			}
		}
	}

	return nics, nil
}

// run and parse command: netsh interface ipv4 show config
func GetInterfaces() ([]Interface, error) {
	lines, err := Cmd("netsh", "interface", "ipv4", "show", "config")
	if err != nil {
		return nil, err
	}
	return ParseInterfaces(lines)
}

// returns true if first line is blank
func NetshCmdBool(args ...string) bool {
	lines, err := Cmd("netsh", args...)
	if err != nil {
		return false
	}
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return true
	}
	return false
}

// run command: netsh interface ipv4 set address name="Ethernet" source=dhcp
func SetIpDhcp(iface string) bool {
	return NetshCmdBool("interface", "ipv4", "set", "address", "name="+iface, "source=dhcp")
}

// run command: netsh interface ipv4 set address name="Ethernet" static 192.168.1.9 255.255.255.0 192.168.1.254
func SetIpStatic(iface, ip, mask, gateway string) bool {
	if strings.HasPrefix(mask, "/") {
		var err error
		mask, err = CidrToMask(mask)
		if err != nil {
			return false
		}
	}
	return NetshCmdBool("interface", "ipv4", "set", "address", "name="+iface, "static", ip, mask, gateway)
}

// run command: netsh interface ipv4 add address name="Ethernet" 192.168.2.9 255.255.255.0
func AddIpStatic(iface, ip, mask string) bool {
	if strings.HasPrefix(mask, "/") {
		var err error
		mask, err = CidrToMask(mask)
		if err != nil {
			return false
		}
	}
	return NetshCmdBool("interface", "ipv4", "add", "address", "name="+iface, ip, mask)
}

// run command: netsh interface ipv4 set dns name="Ethernet" source=dhcp
func SetDnsDhcp(iface string) bool {
	return NetshCmdBool("interface", "ipv4", "set", "dns", "name="+iface, "source=dhcp")
}

// run command: netsh interface ipv4 set dns name="Ethernet" static 192.168.1.1
func SetDnsStatic(iface, dns string) bool {
	return NetshCmdBool("interface", "ipv4", "set", "dns", "name="+iface, "static", dns)
}

// run command: netsh interface ipv4 add dns name="Ethernet" 1.1.1.1 index=2
func AddDnsStatic(iface, dns, index string) bool {
	return NetshCmdBool("interface", "ipv4", "add", "dns", "name="+iface, dns, "index="+index)
}

// run command: netsh interface ipv4 set interface "Ethernet" metric=25
func SetInterfaceMetric(iface, metric string) bool {
	return NetshCmdBool("interface", "ipv4", "set", "interface", iface, "metric="+metric)
}

// run command: netsh interface ipv4 set interface "Ethernet" metric=auto
func SetInterfaceMetricAuto(iface string) bool {
	return NetshCmdBool("interface", "ipv4", "set", "interface", iface, "metric=auto")
}

// run command: netsh interface set interface name="Ethernet" newname="New Ethernet"
func SetInterfaceName(oldName, newName string) bool {
	return NetshCmdBool("interface", "set", "interface", "name="+oldName, "newname="+newName)
}

// run command: ipconfig /release
func ReleaseDhcp() error {
	_, err := Cmd("ipconfig", "/release")
	return err
}

// run command: ipconfig /renew &
// "&" runs in background. cheap way to not wait for the command to exit
func RenewDhcp() error {
	_, err := Cmd("ipconfig", "/renew")
	return err
}

// run command: ipconfig /flushdns
func FlushDns() error {
	_, err := Cmd("ipconfig", "/flushdns")
	return err
}

// set ip and dns to dhcp
func SetDhcp(iface string) bool {
	return SetIpDhcp(iface) && SetDnsDhcp(iface)
}

// set ip and dns to static
func SetStatic(iface, ip, mask, gateway, dns string) bool {
	return SetIpStatic(iface, ip, mask, gateway) && SetDnsStatic(iface, dns)
}
