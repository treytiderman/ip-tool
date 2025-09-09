import { writable, get } from "svelte/store"
import * as app from "../../wailsjs/go/main/App.js";
import { main } from "../../wailsjs/go/models";

export {
    nics,
    nicTemp,
    nicStatus,
    currentNicIndex,
    lastSelectedNicName,
    initNics,
    updateNics,
    pollNics,
    setNic,
    setNicToInterface,
}

const blankNic: any = {
    interface_name: "xxxxxx",
    interface_metric: 999,
    interface_type: "Dedicated",
    connected: false,
    disabled: false,
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
}

const nics = writable<main.Interface[]>([blankNic])

const nicTemp = writable<main.Interface>(JSON.parse(JSON.stringify(blankNic)))

const nicStatus = writable("")

const currentNicIndex = writable(0)
const lastSelectedNicName = writable("")
let lastSelectedWaitCounter = 0

function saveLastSelectedNicName() {
    const temp = get(lastSelectedNicName)
    localStorage.setItem("lastSelectedNicName", temp)
}

function loadLastSelectedNicName() {
    const temp = localStorage.getItem("lastSelectedNicName")
    console.log("localStorage: lastSelectedNicName", temp);
    if (temp) {
        setNic(temp)
    }
}

async function initNics() {
    const tempNics = await app.GetInterfaces()
    console.log("api: Network Interfaces (nics)", tempNics)
    nics.set(tempNics)
    loadLastSelectedNicName()
}

function setNic(name: string) {
    const tempNics = get(nics)
    const foundIndex = tempNics.findIndex((nic) => nic.interface_name === name)

    if (foundIndex === -1) {
        console.log("ui: Network Interface (nic) Selected does not exist", name)
        return
    }

    console.log("ui: Network Interface (nic) Selected", tempNics, tempNics[foundIndex])
    currentNicIndex.set(foundIndex)
    nicTemp.set(JSON.parse(JSON.stringify(tempNics[foundIndex])))
    lastSelectedNicName.set(tempNics[foundIndex].interface_name)
    saveLastSelectedNicName()
}

function resetNic() {
    const tempNics = get(nics)
    if (tempNics.length > 0) {
        console.log("ui: Network Interface Resetting:", lastSelectedWaitCounter, tempNics[0].interface_name)

        if (lastSelectedWaitCounter < 10) {
            lastSelectedWaitCounter++
        } else {
            lastSelectedWaitCounter = 0
            lastSelectedNicName.set(tempNics[0].interface_name)
            // saveLastSelectedNicName()
        }

        currentNicIndex.set(0)
        nicTemp.set(JSON.parse(JSON.stringify(tempNics[0])))
        nics.set(tempNics)
    } else {
        console.log("ui: No Network Interfaces found")
        currentNicIndex.set(0)
        nicTemp.set(JSON.parse(JSON.stringify(tempNics[0])))
        nics.set(tempNics)
    }
}

async function updateNics() {
    const tempNics = await app.GetInterfaces()
    if (JSON.stringify(tempNics) !== JSON.stringify(get(nics))) {
        console.log("api: Network Interfaces (nics)", tempNics, get(nics))
        nics.set(tempNics)
        nicTemp.set(JSON.parse(JSON.stringify(tempNics[get(currentNicIndex)])))
    }

    const tempLastSelectedNicName = get(lastSelectedNicName)
    const foundIndex = tempNics.findIndex((nic) => nic.interface_name === tempLastSelectedNicName)
    if (foundIndex === -1) {
        console.log("api: Network Interfaces (nics) no longer found", tempLastSelectedNicName || undefined)
        resetNic()
        return
    }

    if (foundIndex !== get(currentNicIndex)) {
        console.log("api: Network Interface Selected at new index", foundIndex)
        currentNicIndex.set(foundIndex)
        nicTemp.set(JSON.parse(JSON.stringify(tempNics[foundIndex])))
        setNic(tempNics[foundIndex].interface_name)
        lastSelectedWaitCounter = 0
    }
}

function pollNics(ms = 2000) {
    return setInterval(updateNics, ms)
}

async function setNicToInterface(interface_name: string, nic: main.Interface) {
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
