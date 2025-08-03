package main

import (
	"context"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
}

// IPv4 functions

func (a *App) IsAdmin() bool {
	return IsAdmin()
}

func (a *App) GetInterfaces() []Interface {
	nics, err := GetInterfaces()
	if err != nil {
		return []Interface{}
	}
	return nics
}

func (a *App) SetIpDhcp(iface string) bool {
    return SetIpDhcp(iface)
}

func (a *App) SetIpStatic(iface, ip, mask, gateway string) bool {
    return SetIpStatic(iface, ip, mask, gateway)
}

func (a *App) AddIpStatic(iface, ip, mask string) bool {
    return AddIpStatic(iface, ip, mask)
}

func (a *App) SetDnsDhcp(iface string) bool {
    return SetDnsDhcp(iface)
}

func (a *App) SetDnsStatic(iface, dns string) bool {
    return SetDnsStatic(iface, dns)
}

func (a *App) AddDnsStatic(iface, dns, index string) bool {
    return AddDnsStatic(iface, dns, index)
}

func (a *App) SetInterfaceMetric(iface, metric string) bool {
    return SetInterfaceMetric(iface, metric)
}

func (a *App) SetInterfaceMetricAuto(iface string) bool {
    return SetInterfaceMetricAuto(iface)
}

func (a *App) SetInterfaceName(oldName, newName string) bool {
    return SetInterfaceName(oldName, newName)
}

func (a *App) ReleaseDhcp() bool {
	return ReleaseDhcp() == nil
}

func (a *App) RenewDhcp() bool {
	return RenewDhcp() == nil
}

func (a *App) FlushDns() bool {
	return FlushDns() == nil
}

func (a *App) SetDhcp(iface string) bool {
    return SetDhcp(iface)
}

func (a *App) SetStatic(iface, ip, mask, gateway, dns string) bool {
    return SetStatic(iface, ip, mask, gateway, dns)
}
