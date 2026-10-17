//go:build windows
// +build windows

package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

type CDPInfo struct {
	DeviceID          string `json:"device_id"`
	PortID            string `json:"port_id"`
	Platform          string `json:"platform"`
	SoftwareVersion   string `json:"software_version"`
	Capabilities      string `json:"capabilities"`
	NativeVLAN        string `json:"native_vlan"`
	ManagementAddress string `json:"management_address"`
	Raw               string `json:"raw"`
	Status            string `json:"status"`
}

type cdpMonitor struct {
	mu   sync.RWMutex
	info CDPInfo
	stop chan struct{}
}

var (
	cdpMonitorsMu sync.Mutex
	cdpMonitors   = make(map[string]*cdpMonitor)
)

func GetInterfaceCDP(interfaceName string) CDPInfo {
	cdpMonitorsMu.Lock()
	monitor, exists := cdpMonitors[interfaceName]
	if !exists {
		monitor = &cdpMonitor{
			info: CDPInfo{Status: "Listening for CDP advertisements"},
			stop: make(chan struct{}),
		}
		cdpMonitors[interfaceName] = monitor
		go monitor.listen(interfaceName)
	}
	cdpMonitorsMu.Unlock()

	monitor.mu.RLock()
	defer monitor.mu.RUnlock()
	return monitor.info
}

func StopCDPMonitors() {
	cdpMonitorsMu.Lock()
	for _, monitor := range cdpMonitors {
		close(monitor.stop)
		monitor.mu.Lock()
		monitor.info.Status = "CDP listener stopped"
		monitor.mu.Unlock()
	}
	cdpMonitors = make(map[string]*cdpMonitor)
	cdpMonitorsMu.Unlock()
}

func stopInterfaceCDPMonitor(interfaceName string) {
	cdpMonitorsMu.Lock()
	if monitor, exists := cdpMonitors[interfaceName]; exists {
		close(monitor.stop)
		delete(cdpMonitors, interfaceName)
	}
	cdpMonitorsMu.Unlock()
}

func (m *cdpMonitor) listen(interfaceName string) {
	select {
	case <-m.stop:
		return
	default:
	}

	dll := syscall.NewLazyDLL("wpcap.dll")
	openLive := dll.NewProc("pcap_open_live")
	nextEx := dll.NewProc("pcap_next_ex")
	closeCapture := dll.NewProc("pcap_close")
	dataLink := dll.NewProc("pcap_datalink")
	compileFilter := dll.NewProc("pcap_compile")
	setFilter := dll.NewProc("pcap_setfilter")
	freeCode := dll.NewProc("pcap_freecode")
	if err := dll.Load(); err != nil {
		m.setStatus("Npcap is not installed")
		return
	}
	if err := openLive.Find(); err != nil {
		m.setStatus("Npcap capture API is unavailable")
		return
	}
	if err := nextEx.Find(); err != nil {
		m.setStatus("Npcap capture API is unavailable")
		return
	}
	if err := closeCapture.Find(); err != nil {
		m.setStatus("Npcap capture API is unavailable")
		return
	}
	if err := dataLink.Find(); err != nil {
		m.setStatus("Npcap capture API is unavailable")
		return
	}
	if err := compileFilter.Find(); err != nil {
		m.setStatus("Npcap filtering API is unavailable")
		return
	}
	if err := setFilter.Find(); err != nil {
		m.setStatus("Npcap filtering API is unavailable")
		return
	}
	if err := freeCode.Find(); err != nil {
		m.setStatus("Npcap filtering API is unavailable")
		return
	}

	device, err := getNpcapDeviceName(interfaceName)
	if err != nil {
		m.setStatus(err.Error())
		return
	}
	deviceName, err := syscall.BytePtrFromString(device)
	if err != nil {
		m.setStatus(fmt.Sprintf("Invalid Npcap device name: %v", err))
		return
	}
	errBuffer := make([]byte, 256)
	handle, _, _ := openLive.Call(
		uintptr(unsafe.Pointer(deviceName)),
		65536,
		1,
		1000,
		uintptr(unsafe.Pointer(&errBuffer[0])),
	)
	if handle == 0 {
		reason := strings.TrimSpace(syscallByteString(errBuffer))
		if reason == "" {
			reason = "could not open interface for CDP capture"
		}
		m.setStatus(reason)
		return
	}
	defer closeCapture.Call(handle)

	linkType, _, _ := dataLink.Call(handle)
	if linkType != 1 {
		m.setStatus("Unsupported Npcap link type")
		return
	}
	filterText, err := syscall.BytePtrFromString("ether dst 01:00:0c:cc:cc:cc")
	if err != nil {
		m.setStatus(fmt.Sprintf("Invalid CDP capture filter: %v", err))
		return
	}
	var filter bpfProgram
	result, _, _ := compileFilter.Call(
		handle,
		uintptr(unsafe.Pointer(&filter)),
		uintptr(unsafe.Pointer(filterText)),
		1,
		uintptr(0xffffffff),
	)
	if int32(result) != 0 {
		m.setStatus("Npcap could not compile the CDP capture filter")
		return
	}
	defer freeCode.Call(uintptr(unsafe.Pointer(&filter)))
	result, _, _ = setFilter.Call(handle, uintptr(unsafe.Pointer(&filter)))
	if int32(result) != 0 {
		m.setStatus("Npcap could not apply the CDP capture filter")
		return
	}

	m.setStatus("Listening for CDP advertisements")
	for {
		select {
		case <-m.stop:
			return
		default:
		}
		var headerPointer uintptr
		var packetPointer uintptr
		result, _, _ := nextEx.Call(
			handle,
			uintptr(unsafe.Pointer(&headerPointer)),
			uintptr(unsafe.Pointer(&packetPointer)),
		)
		switch int32(result) {
		case 1:
			if headerPointer == 0 || packetPointer == 0 {
				continue
			}
			capturedLength := *(*uint32)(unsafe.Pointer(headerPointer + 8))
			if capturedLength < 14 || capturedLength > 65536 {
				continue
			}
			packet := unsafe.Slice((*byte)(unsafe.Pointer(packetPointer)), int(capturedLength))
			info, ok := ParseCDPFrame(packet)
			if ok {
				m.setInfo(info)
			}
		case 0:
			continue
		default:
			m.setStatus("Npcap stopped capturing CDP frames")
			return
		}
	}
}

type bpfProgram struct {
	length uint32
	_      uint32
	instr  uintptr
}

func (m *cdpMonitor) setStatus(status string) {
	m.mu.Lock()
	m.info.Status = status
	m.mu.Unlock()
}

func (m *cdpMonitor) setInfo(info CDPInfo) {
	m.mu.Lock()
	m.info = info
	m.mu.Unlock()
}

func getNpcapDeviceName(interfaceName string) (string, error) {
	escapedName := strings.ReplaceAll(interfaceName, "'", "''")
	script := fmt.Sprintf(
		`Get-NetAdapter -Name '%s' -ErrorAction Stop | Select-Object -ExpandProperty InterfaceGuid`,
		escapedName,
	)
	lines, err := Cmd("powershell", "-NoLogo", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	if err != nil {
		return "", fmt.Errorf("find Windows interface identifier for CDP capture: %w", err)
	}
	for _, line := range lines {
		guid := strings.TrimSpace(line)
		if strings.HasPrefix(guid, "{") && strings.HasSuffix(guid, "}") {
			return `\Device\NPF_` + guid, nil
		}
	}
	return "", fmt.Errorf("could not resolve Windows interface identifier for CDP capture")
}

func ParseCDPFrame(frame []byte) (CDPInfo, bool) {
	var info CDPInfo
	if len(frame) < 26 || !bytesEqual(frame[14:22], []byte{0xaa, 0xaa, 0x03, 0x00, 0x00, 0x0c, 0x20, 0x00}) {
		return info, false
	}

	offset := 26
	rawLines := make([]string, 0, 8)
	for offset+4 <= len(frame) {
		tlvType := binary.BigEndian.Uint16(frame[offset : offset+2])
		tlvLength := int(binary.BigEndian.Uint16(frame[offset+2 : offset+4]))
		if tlvLength < 4 || offset+tlvLength > len(frame) {
			return CDPInfo{}, false
		}
		value := frame[offset+4 : offset+tlvLength]
		text := strings.TrimSpace(string(value))
		switch tlvType {
		case 0x0001:
			info.DeviceID = text
			rawLines = append(rawLines, "Device ID: "+text)
		case 0x0002:
			info.ManagementAddress = parseCDPAddresses(value)
			if info.ManagementAddress != "" {
				rawLines = append(rawLines, "Management Address: "+info.ManagementAddress)
			}
		case 0x0003:
			info.PortID = text
			rawLines = append(rawLines, "Port ID: "+text)
		case 0x0004:
			if len(value) == 4 {
				info.Capabilities = fmt.Sprintf("0x%08x", binary.BigEndian.Uint32(value))
				rawLines = append(rawLines, "Capabilities: "+info.Capabilities)
			}
		case 0x0005:
			info.SoftwareVersion = text
			rawLines = append(rawLines, "Software Version: "+text)
		case 0x0006:
			info.Platform = text
			rawLines = append(rawLines, "Platform: "+text)
		case 0x000a:
			if len(value) == 2 {
				info.NativeVLAN = fmt.Sprintf("%d", binary.BigEndian.Uint16(value))
				rawLines = append(rawLines, "Native VLAN: "+info.NativeVLAN)
			}
		}
		offset += tlvLength
	}
	if info.DeviceID == "" && info.PortID == "" && info.Platform == "" {
		return CDPInfo{}, false
	}
	info.Raw = strings.Join(rawLines, "\n")
	info.Status = "CDP neighbor discovered"
	return info, true
}

func parseCDPAddresses(value []byte) string {
	if len(value) < 4 {
		return ""
	}
	count := uint64(binary.BigEndian.Uint32(value[:4]))
	offset := 4
	addresses := make([]string, 0)
	for index := uint64(0); index < count && offset < len(value); index++ {
		if offset+2 > len(value) {
			return strings.Join(addresses, ", ")
		}
		protocolLength := int(value[offset+1])
		offset += 2
		if offset+protocolLength+2 > len(value) {
			return strings.Join(addresses, ", ")
		}
		protocol := value[offset : offset+protocolLength]
		offset += protocolLength
		addressLength := int(binary.BigEndian.Uint16(value[offset : offset+2]))
		offset += 2
		if offset+addressLength > len(value) {
			return strings.Join(addresses, ", ")
		}
		address := value[offset : offset+addressLength]
		offset += addressLength
		if len(protocol) == 1 && protocol[0] == 0xcc && len(address) == net.IPv4len {
			addresses = append(addresses, net.IP(address).String())
		}
	}
	return strings.Join(addresses, ", ")
}

func bytesEqual(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func syscallByteString(buffer []byte) string {
	for index, value := range buffer {
		if value == 0 {
			return string(buffer[:index])
		}
	}
	return string(buffer)
}
