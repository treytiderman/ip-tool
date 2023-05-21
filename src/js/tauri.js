import { invoke } from "@tauri-apps/api/tauri";

export async function is_admin() {
    const output = await invoke("is_admin");
    return output;
}
export async function get_interfaces() {
    const output = await invoke("get_interfaces");
    const interfaces_raw = JSON.parse(output);

    let interfaces = [];
    interfaces_raw.forEach((nic) => {
        if (
            nic.ip_and_masks[0]?.ip_address &&
            nic.ip_and_masks[0]?.ip_address !== "127.0.0.1"
        ) {
            if (nic.dns_server_2 === "") {
                nic.dns_servers = [nic.dns_server_1];
            }
            else {
                nic.dns_servers = [nic.dns_server_1, nic.dns_server_2];
            }
            interfaces.push(nic);
        }
    });

    // console.log("get_interfaces", interfaces);
    return interfaces;
}
export async function set_ip_dhcp(nic) {
    const output = await invoke("set_ip_dhcp", {
        interface: nic,
    });
    return output;
}
export async function set_ip_static(nic, ip, mask, gateway) {
    console.log("set_ip_static", nic, ip, mask, gateway);
    const output = await invoke("set_ip_static", {
        interface: nic,
        ip: ip,
        mask: mask,
        gateway: gateway,
    });
    return output;
}
export async function add_ip_static(nic, ip, mask) {
    const output = await invoke("add_ip_static", {
        interface: nic,
        ip: ip,
        mask: mask,
    });
    return output;
}
export async function set_dns_dhcp(nic) {
    const output = await invoke("set_dns_dhcp", {
        interface: nic,
    });
    return output;
}
export async function set_dns_static(nic, dns_1, dns_2 = "") {
    console.log("set_dns_static", nic, dns_1, dns_2 || "");
    const output = await invoke("set_dns_static", {
        interface: nic,
        dns1: dns_1,
        dns2: dns_2 || "",
    });
    return output;
}
export async function set_dhcp(nic) {
    const output = await invoke("set_dhcp", {
        interface: nic,
    });
    return output;
}
export async function set_static(nic, ip, mask, gateway, dns_1, dns_2) {
    const output = await invoke("set_static", {
        interface: nic,
        ip: ip,
        mask: mask,
        gateway: gateway,
        dns_1: dns_1,
        dns_2: dns_2,
    });
    return output;
}
export async function set_preset(interface_name, preset) {
    for (let i = 0; i < preset.ip_and_masks.length; i++) {
        const ip_and_mask = preset.ip_and_masks[i];
        if (i === 0) {
            set_ip_static(
                interface_name,
                ip_and_mask.ip_address,
                ip_and_mask.subnet_mask,
                preset.gateway
            );
        } else {
            add_ip_static(
                interface_name,
                ip_and_mask.ip_address,
                ip_and_mask.subnet_mask
            );
        }
    }
    const output = set_dns_static(
        interface_name,
        preset.dns_servers[0],
        preset.dns_servers[1]
    );
    return output;
}
export async function set_interface_name(nic, new_name) {
    const output = await invoke("set_interface_name", {
        interface: nic,
        interfaceNewName: new_name
    });
    return output;
}
export async function set_metric(nic, metric) {
    const output = await invoke("set_metric", {
        interface: nic,
        metric: metric
    });
    return output;
}
export async function set_metric_auto(nic) {
    const output = await invoke("set_metric_auto", {
        interface: nic,
    });
    return output;
}

export async function release_dhcp() {
    await invoke("release")
}
export async function renew_dhcp() {
    await invoke("renew")
}
export async function flushdns() {
    await invoke("flushdns")
}

export async function route_print_v4() {
    const output = await invoke("route_print_v4");
    const routes = JSON.parse(output);
    return routes;
}
export async function route_add(destination, mask, gateway, metric) {
    const output = await invoke("route_add", {
        destination: destination,
        mask: mask,
        gateway: gateway,
        metric: metric,
    });
    return output;
}
export async function route_add_persistent(destination, mask, gateway, metric) {
    const output = await invoke("route_add_persistent", {
        destination: destination,
        mask: mask,
        gateway: gateway,
        metric: metric,
    });
    return output;
}
export async function route_change(destination, mask, gateway, metric) {
    const output = await invoke("route_change", {
        destination: destination,
        mask: mask,
        gateway: gateway,
        metric: metric,
    });
    return output;
}
export async function route_change_persistent(destination, mask, gateway, metric) {
    const output = await invoke("route_change_persistent", {
        destination: destination,
        mask: mask,
        gateway: gateway,
        metric: metric,
    });
    return output;
}
export async function route_delete(destination) {
    const output = await invoke("route_delete", {
        destination: destination,
    });
    return output;
}
