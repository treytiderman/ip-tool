import { writable, get } from "svelte/store"
import * as app from "../wailsjs/go/main/App.js";

export type {
    NicIp,
    NicGateway,
    Nic,
}
export {
    nic,
    nics,
    nicTemp,
    initNics,
    updateNics,
    pollNics,
    setNic,
    setNicToInterface,
}

type NicIp = {
    ip_address: string;
    subnet_mask: string;
    cidr: number;
};

type NicGateway = {
    gateway_address: string;
    gateway_metric: number;
};

type Nic = {
    interface_name: string;
    interface_metric: number;
    ip_is_dhcp: boolean;
    ips: NicIp[];
    gateways: NicGateway[];
    dns_is_dhcp: boolean;
    dns_servers: string[];
};

const nic = writable<Nic>({
    interface_name: "xxxxxx",
    interface_metric: 999,
    ip_is_dhcp: false,
    ips: [
        {
            ip_address: "xxx.xxx.xxx.xxx",
            subnet_mask: "xxx.xxx.xxx.xxx",
            cidr: 24,
        },
    ],
    gateways: [
        {
            gateway_address: "xxx.xxx.xxx.xxx",
            gateway_metric: 999,
        },
    ],
    dns_is_dhcp: false,
    dns_servers: ["xxx.xxx.xxx.xxx"],
})

const nics = writable<Nic[]>([]) // TODO: change this to an index of the main nics array

const nicTemp = writable<Nic>(JSON.parse(JSON.stringify(get(nic))))

async function initNics() {
    const tempNics = await app.GetInterfaces()
    console.log("Interfaces (nics)", tempNics)
    nics.set(tempNics)
    setNic(tempNics[0].interface_name)
}

async function updateNics() {
    const tempNics = await app.GetInterfaces()
    nics.set(tempNics)
    if (JSON.stringify(tempNics) !== JSON.stringify(get(nics))) {
        console.log("Interfaces (nics)", tempNics)
    }
    if (!tempNics.find(n => n.interface_name === get(nic).interface_name)) {
        setNic(get(nics)[0].interface_name)
    }
    setNic(get(nic).interface_name)
}

function pollNics(ms = 2000) {
    return setInterval(updateNics, ms)
}

function setNic(name: string) {
    const tempNics = get(nics)
    const found = tempNics.find((nic) => nic.interface_name === name)
    if (found) {
        if (JSON.stringify(found) !== JSON.stringify(get(nic))) {
            console.log("Interface Selected (nic)", found)
            nic.set(found)
            nicTemp.set(JSON.parse(JSON.stringify(found)))
        }
    } else {
        console.error("Interface Selected (nic)", name, found)
    }
}

async function setNicToInterface(interface_name: string, nic: Nic) {
    console.log("Set Nic To Interface", interface_name, nic);
    await app.SetStatic(
        interface_name,
        nic.ips[0].ip_address,
        nic.ips[0].subnet_mask,
        nic.gateways[0].gateway_address,
        nic.dns_servers[0]
    )
    if (nic.ips.length > 1) {
        for (let index = 1; index < nic.ips.length; index++) {
            const ip = nic.ips[index]
            await app.AddIpStatic(interface_name, ip.ip_address, ip.subnet_mask)
        }
    }
    if (nic.dns_servers.length > 1) {
        for (let index = 1; index < nic.dns_servers.length; index++) {
            const dns = nic.dns_servers[index]
            await app.AddDnsStatic(interface_name, dns, index.toString())
        }
    }
    if (interface_name !== nic.interface_name) {
        await app.SetInterfaceName(interface_name, nic.interface_name)
    }
}
