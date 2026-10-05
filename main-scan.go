//go:build windows
// +build windows

package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/binary"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// The vendor assignments are from the IEEE MA-L, MA-M, and MA-S registries.
// Sources: https://standards-oui.ieee.org/oui/oui.csv,
// https://standards-oui.ieee.org/oui28/mam.csv, and
// https://standards-oui.ieee.org/oui36/oui36.csv.
//
//go:embed data/mac-vendors.csv
var macVendorData []byte

const (
	ethernetAddressLength = 6
	maxScanAddresses      = 65536
	maxConcurrentPings    = 256
	pingTimeout           = 300 * time.Millisecond
)

type ScanResult struct {
	IPAddress  string `json:"ip_address"`
	LatencyMs  int64  `json:"latency_ms"`
	IsLocal    bool   `json:"is_local"`
	MACAddress string `json:"mac_address"`
	Vendor     string `json:"vendor"`
}

type neighborEntry struct {
	IPAddress        string `json:"IPAddress"`
	LinkLayerAddress string `json:"LinkLayerAddress"`
}

var (
	macVendorsOnce sync.Once
	macVendors     map[string]string
	macVendorsErr  error
)

func ScanSubnet(interfaceName, ipAddress, subnetMask string) ([]ScanResult, error) {
	if interfaceName == "" {
		return nil, fmt.Errorf("interface name is required")
	}
	subnet, prefix, err := parseScanSubnet(ipAddress, subnetMask)
	if err != nil {
		return nil, err
	}

	configs, err := GetInterfaceConfigs()
	if err != nil {
		return nil, fmt.Errorf("read interface configuration: %w", err)
	}

	var selected *InterfaceConfig
	for i := range configs {
		if configs[i].InterfaceName == interfaceName {
			selected = &configs[i]
			break
		}
	}
	if selected == nil {
		return nil, fmt.Errorf("interface %q was not found", interfaceName)
	}

	var localIP net.IP
	for _, configuredIP := range selected.Ips {
		parsedIP := net.ParseIP(configuredIP.IpAddress).To4()
		if parsedIP == nil {
			continue
		}
		localIP = parsedIP
		break
	}
	if localIP == nil {
		return nil, fmt.Errorf("interface %q has no configured IPv4 address", interfaceName)
	}

	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("find network interface %q: %w", interfaceName, err)
	}

	subnetSize := uint64(1) << uint(32-prefix)
	networkAddress := binary.BigEndian.Uint32(subnet.IP.To4())
	firstOffset := uint64(0)
	addressCount := subnetSize
	if prefix <= 30 {
		firstOffset = 1
		addressCount -= 2
	}

	type target struct {
		ip      string
		isLocal bool
	}
	targets := make([]target, 0, addressCount)
	for offset := firstOffset; offset < firstOffset+addressCount; offset++ {
		address := make(net.IP, net.IPv4len)
		binary.BigEndian.PutUint32(address, networkAddress+uint32(offset))
		targets = append(targets, target{
			ip:      address.String(),
			isLocal: address.Equal(localIP),
		})
	}

	results := make([]ScanResult, 0)
	var resultsMu sync.Mutex
	var scanErr error
	var scanErrOnce sync.Once
	scanCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	semaphore := make(chan struct{}, maxConcurrentPings)
	var workers sync.WaitGroup

	for _, host := range targets {
		if scanCtx.Err() != nil {
			break
		}
		semaphore <- struct{}{}
		workers.Add(1)
		go func(host target) {
			defer workers.Done()
			defer func() { <-semaphore }()

			started := time.Now()
			ctx, stop := context.WithTimeout(scanCtx, 2*time.Second)
			defer stop()
			cmd := exec.CommandContext(ctx, "ping", "-n", "1", "-w", strconv.Itoa(int(pingTimeout/time.Millisecond)), host.ip)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			err := cmd.Run()
			if err != nil {
				var exitErr *exec.ExitError
				if errors.Is(err, context.DeadlineExceeded) || errors.As(err, &exitErr) {
					return
				}
				scanErrOnce.Do(func() {
					scanErr = fmt.Errorf("ping %s: %w", host.ip, err)
					cancel()
				})
				return
			}

			latency := time.Since(started).Milliseconds()
			resultsMu.Lock()
			results = append(results, ScanResult{
				IPAddress: host.ip,
				LatencyMs: latency,
				IsLocal:   host.isLocal,
			})
			resultsMu.Unlock()
		}(host)
	}

	workers.Wait()
	if scanErr != nil {
		return nil, scanErr
	}

	neighbors, err := getIPv4Neighbors(iface.Index)
	if err != nil {
		return nil, fmt.Errorf("read interface neighbors: %w", err)
	}
	localMAC := iface.HardwareAddr.String()
	for i := range results {
		macAddress := localMAC
		if !results[i].IsLocal {
			macAddress = neighbors[results[i].IPAddress]
		}
		if macAddress == "" {
			results[i].Vendor = "Unknown"
			continue
		}
		results[i].MACAddress = macAddress
		vendor, err := lookupMACVendor(macAddress)
		if err != nil {
			return nil, fmt.Errorf("look up vendor for %s: %w", macAddress, err)
		}
		results[i].Vendor = vendor
	}

	sort.Slice(results, func(i, j int) bool {
		return bytes.Compare(net.ParseIP(results[i].IPAddress).To4(), net.ParseIP(results[j].IPAddress).To4()) < 0
	})
	return results, nil
}

func parseScanSubnet(ipAddress, subnetMask string) (*net.IPNet, int, error) {
	scanIP := net.ParseIP(ipAddress).To4()
	if scanIP == nil {
		return nil, 0, fmt.Errorf("invalid IPv4 address: %q", ipAddress)
	}
	parsedMask := net.ParseIP(subnetMask).To4()
	if parsedMask == nil {
		return nil, 0, fmt.Errorf("invalid IPv4 subnet mask: %q", subnetMask)
	}
	mask := net.IPMask(parsedMask)
	prefix, bits := mask.Size()
	if bits != 32 || prefix < 1 || prefix > 32 {
		return nil, 0, fmt.Errorf("invalid IPv4 subnet mask: %q", subnetMask)
	}
	subnetSize := uint64(1) << uint(32-prefix)
	if subnetSize > maxScanAddresses {
		return nil, 0, fmt.Errorf("subnet is too large to scan (%d addresses; maximum is %d)", subnetSize, maxScanAddresses)
	}
	return &net.IPNet{IP: scanIP.Mask(mask), Mask: mask}, prefix, nil
}

func getIPv4Neighbors(interfaceIndex int) (map[string]string, error) {
	script := fmt.Sprintf(
		`$ErrorActionPreference = 'Stop'; $entries = @(Get-NetNeighbor -AddressFamily IPv4 -InterfaceIndex %d | Select-Object IPAddress, LinkLayerAddress); ConvertTo-Json -InputObject $entries -Compress`,
		interfaceIndex,
	)
	lines, err := Cmd("powershell", "-NoLogo", "-NoProfile", "-NonInteractive", "-Command", script)
	if err != nil {
		return nil, err
	}

	var entries []neighborEntry
	if err := json.Unmarshal([]byte(strings.Join(lines, "\n")), &entries); err != nil {
		return nil, fmt.Errorf("parse neighbor table: %w", err)
	}

	neighbors := make(map[string]string, len(entries))
	for _, entry := range entries {
		ip := net.ParseIP(entry.IPAddress).To4()
		mac, err := net.ParseMAC(entry.LinkLayerAddress)
		if ip == nil || err != nil || len(mac) != ethernetAddressLength || mac[0]&1 != 0 {
			continue
		}
		neighbors[ip.String()] = mac.String()
	}
	return neighbors, nil
}

func lookupMACVendor(macAddress string) (string, error) {
	mac, err := net.ParseMAC(macAddress)
	if err != nil || len(mac) != ethernetAddressLength {
		return "Unknown", nil
	}
	if mac[0]&0x02 != 0 {
		return "Locally administered", nil
	}

	macVendorsOnce.Do(loadMACVendors)
	if macVendorsErr != nil {
		return "", macVendorsErr
	}

	address := make([]byte, ethernetAddressLength*2)
	hex.Encode(address, mac)
	for _, prefixLength := range []int{9, 7, 6} {
		if vendor, ok := macVendors[string(bytes.ToUpper(address[:prefixLength]))]; ok {
			return vendor, nil
		}
	}
	return "Unknown", nil
}

func loadMACVendors() {
	reader := csv.NewReader(bytes.NewReader(macVendorData))
	if _, err := reader.Read(); err != nil {
		macVendorsErr = fmt.Errorf("read IEEE vendor database header: %w", err)
		return
	}

	macVendors = make(map[string]string)
	for {
		row, err := reader.Read()
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			macVendorsErr = fmt.Errorf("read IEEE vendor database: %w", err)
			return
		}
		if len(row) != 2 {
			macVendorsErr = fmt.Errorf("invalid IEEE vendor database row")
			return
		}
		if !validMACPrefix(row[0]) {
			macVendorsErr = fmt.Errorf("invalid IEEE vendor prefix %q", row[0])
			return
		}
		macVendors[strings.ToUpper(row[0])] = row[1]
	}
}

func validMACPrefix(prefix string) bool {
	if len(prefix) != 6 && len(prefix) != 7 && len(prefix) != 9 {
		return false
	}
	for _, digit := range prefix {
		if (digit < '0' || digit > '9') && (digit < 'A' || digit > 'F') && (digit < 'a' || digit > 'f') {
			return false
		}
	}
	return true
}
