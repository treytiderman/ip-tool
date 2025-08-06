import { writable, get } from "svelte/store"
import * as app from "../../wailsjs/go/main/App.js";

export type {
    PresetIp,
    PresetGateway,
    Preset,
}
export {
    presets,
    presetTemp,
    presetSelectedIndex,
    savePresets,
    loadPresets,
    resetPresets,
    addPreset,
    editPreset,
    removePreset,
    selectPreset,
    setPresetToInterface
}

type PresetIp = {
    ip_address: string
    subnet_mask: string
    cidr: number
}

type PresetGateway = {
    gateway_address: string
}

type Preset = {
    name: string
    ip_is_dhcp: boolean
    ips: PresetIp[]
    gateways: PresetGateway[]
    dns_is_dhcp: boolean
    dns_servers: string[]
}

const defaultPresets: Preset[] = [
    {
        name: "192.168.0.7/24",
        ip_is_dhcp: false,
        ips: [
            {
                ip_address: "192.168.0.7",
                subnet_mask: "255.255.255.0",
                cidr: 24,
            },
        ],
        gateways: [
            {
                gateway_address: "192.168.0.1",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["192.168.0.1"],
    },
    {
        name: "192.168.1.7/24",
        ip_is_dhcp: false,
        ips: [
            {
                ip_address: "192.168.1.7",
                subnet_mask: "255.255.255.0",
                cidr: 24,
            },
        ],
        gateways: [
            {
                gateway_address: "192.168.1.1",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["192.168.1.1"],
    },
    {
        name: "192.168.2.7/24",
        ip_is_dhcp: false,
        ips: [
            {
                ip_address: "192.168.2.7",
                subnet_mask: "255.255.255.0",
                cidr: 24,
            },
        ],
        gateways: [
            {
                gateway_address: "192.168.2.1",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["192.168.2.1"],
    },
    {
        name: "169.254.0.7/24",
        ip_is_dhcp: false,
        ips: [
            {
                ip_address: "169.254.0.7",
                subnet_mask: "255.255.255.0",
                cidr: 24,
            },
        ],
        gateways: [
            {
                gateway_address: "169.254.0.1",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["169.254.0.1"],
    },
    {
        name: "172.22.0.7/16",
        ip_is_dhcp: false,
        ips: [
            {
                ip_address: "172.22.0.7",
                subnet_mask: "255.255.0.0",
                cidr: 16,
            },
        ],
        gateways: [
            {
                gateway_address: "172.22.0.2",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["172.22.0.2"],
    },
    {
        name: "198.18.0.99/16",
        ip_is_dhcp: false,
        ips: [
            {
                ip_address: "198.18.0.99",
                subnet_mask: "255.255.0.0",
                cidr: 16,
            },
        ],
        gateways: [
            {
                gateway_address: "198.18.0.1",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["198.18.0.1"],
    },
];

const presets = writable<Preset[]>(JSON.parse(JSON.stringify(defaultPresets)))

const presetTemp = writable<Preset>(JSON.parse(JSON.stringify(defaultPresets[0])))

let presetSelectedIndex = 0

function savePresets() {
    const tempPresets = get(presets)
    localStorage.setItem("presets", JSON.stringify(tempPresets))
}

function loadPresets() {
    const savedPresetsJson = localStorage.getItem("presets")
    if (savedPresetsJson) {
        presets.set(JSON.parse(savedPresetsJson))
    }
    console.log("localStorage: presets", get(presets) || undefined);
}

function resetPresets() {
    presets.set(JSON.parse(JSON.stringify(defaultPresets)))
    savePresets()
}

function addPreset(preset: Preset) {
    presets.update(ps => {
        ps.push(preset)
        return ps
    })
    savePresets()
}

function editPreset(preset_name: string, preset: Preset) {
    presets.update(ps => {
        const pIndex = ps.findIndex(p => p.name === preset_name)
        if (pIndex !== -1) {
            ps[pIndex] = preset
        }
        return ps
    })
    savePresets()
}

function removePreset(name: string) {
    presetSelectedIndex = 0
    presets.update(ps => {
        return ps.filter(p => p.name !== name)
    })
    savePresets()
}

function selectPreset(preset_name: string) {
    const ps = get(presets)
    presetSelectedIndex = ps.findIndex(p => p.name === preset_name)
    if (presetSelectedIndex !== -1) {
        console.log("Preset Selected", preset_name, ps[presetSelectedIndex])
        presetTemp.set(ps[presetSelectedIndex])
    } else {
        console.log("Preset Selected Reset to index 0", preset_name, ps[presetSelectedIndex])
        presetSelectedIndex = 0
    }
}

async function setPresetToInterface(interface_name: string, preset_name: string) {
    const ps = get(presets)
    const p = ps.find(p => p.name === preset_name)
    console.log("Set Preset To Interface", interface_name, p);
    await app.SetStatic(
        interface_name,
        p.ips[0].ip_address,
        p.ips[0].subnet_mask,
        p.gateways[0].gateway_address,
        p.dns_servers[0]
    )
    if (p.ips.length > 1) {
        for (let index = 1; index < p.ips.length; index++) {
            const ip = p.ips[index]
            await app.AddIpStatic(interface_name, ip.ip_address, ip.subnet_mask)
        }
    }
    if (p.dns_servers.length > 1) {
        for (let index = 1; index < p.dns_servers.length; index++) {
            const dns = p.dns_servers[index]
            await app.AddDnsStatic(interface_name, dns, index.toString())
        }
    }
}

