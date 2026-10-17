//go:build windows
// +build windows

package main

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/dns/dnsmessage"
	"golang.org/x/net/ipv4"
)

const (
	mdnsMulticastAddress = "224.0.0.251:5353"
	mdnsListenDuration   = 5 * time.Second
	mdnsReadInterval     = 250 * time.Millisecond
)

var mdnsBrowseServices = []string{
	"_services._dns-sd._udp.local.",
	"_workstation._tcp.local.",
	"_http._tcp.local.",
	"_https._tcp.local.",
	"_ssh._tcp.local.",
	"_smb._tcp.local.",
	"_ipp._tcp.local.",
	"_printer._tcp.local.",
	"_airplay._tcp.local.",
	"_raop._tcp.local.",
	"_googlecast._tcp.local.",
	"_hap._tcp.local.",
	"_homekit._tcp.local.",
	"_companion-link._tcp.local.",
	"_device-info._tcp.local.",
}

type MDNSResult struct {
	Hostname    string   `json:"hostname"`
	IPAddresses []string `json:"ip_addresses"`
	Services    []string `json:"services"`
	Details     []string `json:"details"`
}

type mdnsHost struct {
	hostname  string
	addresses map[string]struct{}
	services  map[string]struct{}
	details   map[string]struct{}
}

type mdnsServiceInstance struct {
	service string
	host    string
	port    uint16
}

type mdnsDiscovery struct {
	hosts           map[string]*mdnsHost
	instances       map[string]mdnsServiceInstance
	instanceDetails map[string]map[string]struct{}
}

func ScanMDNS(ctx context.Context, interfaceName string, progress func(int)) ([]MDNSResult, error) {
	if interfaceName == "" {
		return nil, fmt.Errorf("interface name is required")
	}
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("find network interface %q: %w", interfaceName, err)
	}

	group := &net.UDPAddr{IP: net.ParseIP("224.0.0.251"), Port: 5353}
	conn, err := net.ListenMulticastUDP("udp4", iface, group)
	if err != nil {
		return nil, fmt.Errorf("listen for mDNS on interface %q: %w", interfaceName, err)
	}
	defer conn.Close()
	packetConn := ipv4.NewPacketConn(conn)
	if err := packetConn.SetMulticastInterface(iface); err != nil {
		return nil, fmt.Errorf("select mDNS interface %q: %w", interfaceName, err)
	}
	if err := packetConn.SetMulticastTTL(255); err != nil {
		return nil, fmt.Errorf("set mDNS multicast TTL: %w", err)
	}

	discovery := mdnsDiscovery{
		hosts:           make(map[string]*mdnsHost),
		instances:       make(map[string]mdnsServiceInstance),
		instanceDetails: make(map[string]map[string]struct{}),
	}
	queried := make(map[string]struct{})
	started := time.Now()
	nextQuery := 0
	queryRounds := []time.Duration{0, time.Second, 3 * time.Second}
	buffer := make([]byte, 9000)

	for time.Since(started) < mdnsListenDuration {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		elapsed := time.Since(started)
		if nextQuery < len(queryRounds) && elapsed >= queryRounds[nextQuery] {
			for _, service := range mdnsBrowseServices {
				if err := sendMDNSQuery(conn, group, service); err != nil {
					return nil, fmt.Errorf("query mDNS service %s: %w", service, err)
				}
				queried[normalizeMDNSName(service)] = struct{}{}
			}
			nextQuery++
		}

		deadline := time.Now().Add(mdnsReadInterval)
		if remaining := time.Until(started.Add(mdnsListenDuration)); remaining < mdnsReadInterval {
			deadline = time.Now().Add(remaining)
		}
		if err := conn.SetReadDeadline(deadline); err != nil {
			return nil, fmt.Errorf("set mDNS read deadline: %w", err)
		}
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, fmt.Errorf("read mDNS response: %w", err)
		}

		var message dnsmessage.Message
		if err := message.Unpack(buffer[:n]); err != nil || !message.Header.Response {
			continue
		}
		newServices := discovery.addMessage(message)
		for _, service := range newServices {
			service = normalizeMDNSName(service)
			if service == "" {
				continue
			}
			if _, alreadyQueried := queried[service]; alreadyQueried {
				continue
			}
			if err := sendMDNSQuery(conn, group, service); err != nil {
				return nil, fmt.Errorf("query discovered mDNS service %s: %w", service, err)
			}
			queried[service] = struct{}{}
		}
		if progress != nil {
			progress(len(discovery.results()))
		}
	}

	return discovery.results(), nil
}

func sendMDNSQuery(conn *net.UDPConn, group *net.UDPAddr, service string) error {
	packet, err := packMDNSQuery(service)
	if err != nil {
		return err
	}
	if _, err := conn.WriteToUDP(packet, group); err != nil {
		return err
	}
	return nil
}

func packMDNSQuery(service string) ([]byte, error) {
	if !strings.HasSuffix(service, ".") {
		service += "."
	}
	name, err := dnsmessage.NewName(service)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}
	message := dnsmessage.Message{
		Header: dnsmessage.Header{RecursionDesired: false},
		Questions: []dnsmessage.Question{{
			Name:  name,
			Type:  dnsmessage.TypePTR,
			Class: dnsmessage.ClassINET,
		}},
	}
	packet, err := message.Pack()
	if err != nil {
		return nil, fmt.Errorf("pack DNS query: %w", err)
	}
	return packet, nil
}

func (d *mdnsDiscovery) addMessage(message dnsmessage.Message) []string {
	var discoveredServices []string
	records := make([]dnsmessage.Resource, 0, len(message.Answers)+len(message.Additionals))
	records = append(records, message.Answers...)
	records = append(records, message.Additionals...)
	records = append(records, message.Authorities...)

	for _, record := range records {
		owner := normalizeMDNSName(record.Header.Name.String())
		switch body := record.Body.(type) {
		case *dnsmessage.PTRResource:
			target := normalizeMDNSName(body.PTR.String())
			if owner == "_services._dns-sd._udp.local" {
				discoveredServices = append(discoveredServices, target)
				continue
			}
			if target != "" {
				instance := d.instances[target]
				instance.service = owner
				d.instances[target] = instance
			}
		case *dnsmessage.SRVResource:
			target := normalizeMDNSName(body.Target.String())
			instance := d.instances[owner]
			instance.host = target
			instance.port = body.Port
			d.instances[owner] = instance
			if target != "" {
				d.host(target).hostname = target
				d.host(target).details[owner] = struct{}{}
			}
		case *dnsmessage.TXTResource:
			if _, ok := d.instanceDetails[owner]; !ok {
				d.instanceDetails[owner] = make(map[string]struct{})
			}
			for _, detail := range body.TXT {
				if detail != "" {
					d.instanceDetails[owner][detail] = struct{}{}
				}
			}
		case *dnsmessage.AResource:
			address := net.IP(body.A[:]).String()
			if owner != "" {
				d.host(owner).addresses[address] = struct{}{}
			}
		case *dnsmessage.AAAAResource:
			address := net.IP(body.AAAA[:]).String()
			if owner != "" {
				d.host(owner).addresses[address] = struct{}{}
			}
		}
	}

	return discoveredServices
}

func (d *mdnsDiscovery) host(name string) *mdnsHost {
	name = normalizeMDNSName(name)
	host, ok := d.hosts[name]
	if !ok {
		host = &mdnsHost{
			hostname:  name,
			addresses: make(map[string]struct{}),
			services:  make(map[string]struct{}),
			details:   make(map[string]struct{}),
		}
		d.hosts[name] = host
	}
	return host
}

func (d *mdnsDiscovery) results() []MDNSResult {
	results := make([]MDNSResult, 0, len(d.hosts))
	for instanceName, instance := range d.instances {
		if instance.host == "" {
			continue
		}
		host := d.host(instance.host)
		if instance.service != "" {
			host.services[fmt.Sprintf("%s (%d)", instance.service, instance.port)] = struct{}{}
		}
		for detail := range d.instanceDetails[instanceName] {
			host.details[detail] = struct{}{}
		}
	}
	for _, host := range d.hosts {
		result := MDNSResult{
			Hostname:    host.hostname,
			IPAddresses: []string{},
			Services:    []string{},
			Details:     []string{},
		}
		for address := range host.addresses {
			result.IPAddresses = append(result.IPAddresses, address)
		}
		for service := range host.services {
			result.Services = append(result.Services, service)
		}
		for detail := range host.details {
			result.Details = append(result.Details, detail)
		}
		sort.Strings(result.IPAddresses)
		sort.Strings(result.Services)
		sort.Strings(result.Details)
		results = append(results, result)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Hostname < results[j].Hostname
	})
	return results
}

func normalizeMDNSName(name string) string {
	name = strings.TrimSuffix(strings.TrimSpace(name), ".")
	name = strings.ToLower(name)
	if strings.HasSuffix(strings.ToLower(name), ".local") {
		return name
	}
	if name == "" {
		return ""
	}
	return name + ".local"
}
