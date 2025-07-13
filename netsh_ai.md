AI response to netsh.rs to netsh.go

```go
// +build windows

package main

import (
    "bufio"
    // "encoding/json"
    "errors"
    "net"
    "os/exec"
    "strings"
)

// Structs
type Ip struct {
    IpAddress   string `json:"ip_address"`
    SubnetMask  string `json:"subnet_mask"`
}

type Interface struct {
    InterfaceName   string `json:"interface_name"`
    InterfaceMetric uint16 `json:"interface_metric"`
    Gateway         string `json:"gateway"`
    GatewayMetric   uint16 `json:"gateway_metric"`
    IpIsDhcp        bool   `json:"ip_is_dhcp"`
    IpAndMasks      []Ip   `json:"ip_and_masks"`
    DnsIsDhcp       bool   `json:"dns_is_dhcp"`
    DnsServer1      string `json:"dns_server_1"`
    DnsServer2      string `json:"dns_server_2"`
}

type Route struct {
    Destination string `json:"destination"`
    Netmask     string `json:"netmask"`
    Gateway     string `json:"gateway"`
    Interface   string `json:"interface"`
    Metric      uint16 `json:"metric"`
    Persistent  bool   `json:"persistent"`
}

func cmd(prog string, args ...string) ([]string, error) {
    out, err := exec.Command(prog, args...).Output()
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

func RoutePrintV4() ([]Route, error) {
    lines, err := cmd("route", "PRINT", "-4")
    if err != nil {
        return nil, err
    }

    var routes []Route
    var inDest, inPersistent bool

    for _, raw := range lines {
        line := strings.TrimSpace(raw)
        if strings.HasPrefix(line, "Network Destination") {
            inDest = true
            continue
        }
        if strings.HasPrefix(line, "===") {
            inDest = false
            inPersistent = false
            continue
        }
        if strings.HasPrefix(line, "Network Address") {
            inPersistent = true
            continue
        }
        if inDest {
            parts := strings.Fields(line)
            if len(parts) >= 5 {
                m, _ := parseUint(parts[4])
                routes = append(routes, Route{
                    Destination: parts[0],
                    Netmask:     parts[1],
                    Gateway:     parts[2],
                    Interface:   parts[3],
                    Metric:      uint16(m),
                    Persistent:  false,
                })
            }
        } else if inPersistent {
            parts := strings.Fields(line)
            if len(parts) >= 4 {
                for i := range routes {
                    if routes[i].Destination == parts[0] {
                        routes[i].Persistent = true
                        break
                    }
                }
            }
        }
    }

    return routes, nil
}

func RouteAdd(dest, mask, gateway, metric string) error {
    _, err := cmd("route", "ADD", dest, "MASK", mask, gateway, "METRIC", metric)
    return err
}

func RouteAddPersistent(dest, mask, gateway, metric string) error {
    _, err := cmd("route", "-p", "ADD", dest, "MASK", mask, gateway, "METRIC", metric)
    return err
}

func RouteChange(dest, mask, gateway, metric string) error {
    _, err := cmd("route", "CHANGE", dest, "MASK", mask, gateway, "METRIC", metric)
    return err
}

func RouteChangePersistent(dest, mask, gateway, metric string) error {
    _, err := cmd("route", "-p", "CHANGE", dest, "MASK", mask, gateway, "METRIC", metric)
    return err
}

func RouteDelete(dest string) error {
    _, err := cmd("route", "-p", "DELETE", dest)
    return err
}

```
