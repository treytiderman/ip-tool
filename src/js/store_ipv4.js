import { writable } from 'svelte/store'
import { get_interfaces } from "../js/tauri";

// Variables
const default_interface = {
    interface_name: "Ethernet",
    interface_metric: 25,
    gateway: "192.168.1.254",
    gateway_metric: 1,
    ip_is_dhcp: false,
    ip_and_masks: [
        {
            ip_address: "192.168.1.9",
            subnet_mask: "255.255.255.0",
        },
    ],
    dns_is_dhcp: false,
    dns_servers: ["192.168.1.1"],
}
const default_interfaces = [
    default_interface,
    {
        interface_name: "WiFi",
        interface_metric: 50,
        gateway: "10.10.1.1",
        gateway_metric: 1,
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "10.10.1.9",
                subnet_mask: "255.255.0.0",
            },
            {
                ip_address: "10.10.2.9",
                subnet_mask: "255.255.0.0",
            },
            {
                ip_address: "10.10.3.9",
                subnet_mask: "255.255.0.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["10.10.1.1"],
    },
]
const default_interface_table_sort = {
    interface: "NONE",
    ip_addresses: "NONE",
    subnet_masks: "NONE",
    gateway: "NONE",
    dns_servers: "NONE",
}
export const default_preset = {
    name: "Home 🏠",
    gateway: "192.168.1.254",
    ip_is_dhcp: false,
    ip_and_masks: [
        {
            ip_address: "192.168.1.9",
            subnet_mask: "255.255.255.0",
        },
    ],
    dns_is_dhcp: false,
    dns_servers: ["192.168.1.1"],
}
export const default_presets = [
    default_preset,
    {
        name: "default_0",
        gateway: "192.168.0.1",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "192.168.0.7",
                subnet_mask: "255.255.255.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["192.168.0.1"],
    },
    {
        name: "default_1",
        gateway: "192.168.1.1",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "192.168.1.7",
                subnet_mask: "255.255.255.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["192.168.1.1"],
    },
    {
        name: "default_2",
        gateway: "192.168.2.1",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "192.168.2.7",
                subnet_mask: "255.255.255.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["192.168.2.1"],
    },
    {
        name: "default_3",
        gateway: "192.168.3.1",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "192.168.3.7",
                subnet_mask: "255.255.255.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["192.168.3.1"],
    },
    {
        name: "default_4",
        gateway: "192.168.4.1",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "192.168.4.7",
                subnet_mask: "255.255.255.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["192.168.4.1"],
    },
    {
        name: "default_10",
        gateway: "10.0.0.1",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "10.0.0.7",
                subnet_mask: "255.0.0.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["10.0.0.1"],
    },
    {
        name: "default_169",
        gateway: "169.254.0.1",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "169.254.0.9",
                subnet_mask: "255.255.0.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["169.254.0.1"],
    },
    {
        name: "default_172",
        gateway: "172.22.0.2",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "172.22.0.9",
                subnet_mask: "255.255.0.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["172.22.0.2"],
    },
    {
        name: "default_198",
        gateway: "198.18.0.1",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "198.18.0.9",
                subnet_mask: "255.255.0.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["198.18.0.1"],
    },
    {
        name: "default_all",
        gateway: "192.168.1.1",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "169.254.0.9",
                subnet_mask: "255.255.0.0",
            },
            {
                ip_address: "172.22.0.9",
                subnet_mask: "255.255.0.0",
            },
            {
                ip_address: "192.168.0.9",
                subnet_mask: "255.255.0.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["8.8.8.8", "8.8.4.4"],
    },
];
const default_preset_table_sort = {
    preset: "NONE",
    ip_addresses: "NONE",
    subnet_masks: "NONE",
    gateway: "NONE",
    dns_servers: "NONE",
}
const default_state = {
    interfaces: default_interfaces,
    interface_active: default_interface,
    interface_editing: false,
    interface_table_sort: default_interface_table_sort,
    presets: default_presets,
    preset_active: default_preset,
    preset_editing: false,
    preset_table_sort: default_preset_table_sort,
}
export const ipv4 = create_store()
export const DIRECTIONS = {
    NONE: "NONE",
    DOWN: "DOWN",
    UP: "UP",
};

// Functions
function state_initialize() {
    if (localStorage.getItem("ipv4") === null) state_default()
    else state_from_localStorage()
}
function state_default() {
    ipv4.set(default_state)
    localStorage.setItem("ipv4", JSON.stringify(default_state))
}
function state_from_localStorage() {
    try {
        const localStorage_json = localStorage.getItem("ipv4")
        const localStorage_parsed = JSON.parse(localStorage_json)
        const state_new = {}
        Object.keys(default_state).forEach(key => {
            if (localStorage_parsed[key] == null) state_new[key] = default_state[key];
            else state_new[key] = localStorage_parsed[key]
        })
        ipv4.set(state_new)
    }
    catch (error) {
        console.log("ipv4: localStorage not a JSON. Setting to default")
        state_initialize()
    }
}
function state_on_change() {
    ipv4.subscribe(state => {
        localStorage.setItem("ipv4", JSON.stringify(state))
    })
}
function create_store() {
    const { subscribe, set, update } = writable(default_state)

    // Store functions
    return {
        set,
        update,
        subscribe,

        state_initialize,
        state_from_localStorage,
    }
}

// Startup
state_initialize()
state_on_change()