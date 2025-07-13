#![allow(unused)]
#![allow(non_snake_case)]

use serde::{Deserialize, Serialize};
use std::env;
use std::io::{BufRead, BufReader, Error, ErrorKind};
use std::net::Ipv4Addr;
use std::os::windows::process::CommandExt;
use std::process::{Command, Stdio};
use std::thread;
use std::time::Duration;

// Structs
#[derive(Serialize, Deserialize, Debug)]
pub struct Ip {
    ip_address: String,
    subnet_mask: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Interface {
    interface_name: String,
    interface_metric: u16,
    gateway: String,
    gateway_metric: u16,
    ip_is_dhcp: bool,
    ip_and_masks: Vec<Ip>,
    dns_is_dhcp: bool,
    dns_server_1: String,
    dns_server_2: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Route {
    destination: String,
    netmask: String,
    gateway: String,
    interface: String,
    metric: u16,
    persistent: bool,
}

// Private Functions
fn cmd(program: &str, args: &str) -> Result<Vec<String>, Error> {
    const CREATE_NO_WINDOW: u32 = 0x08000000;
    const DETACHED_PROCESS: u32 = 0x00000008;
    let args_vec: Vec<&str> = args.split_whitespace().collect();
    let stdout = Command::new(program)
        .args(args_vec)
        .creation_flags(DETACHED_PROCESS)
        .stdout(Stdio::piped())
        .spawn()?
        .stdout
        .ok_or_else(|| Error::new(ErrorKind::Other, "Could not capture standard output."))?;
    let reader = BufReader::new(stdout);
    let lines: Vec<String> = reader.lines().filter_map(|line| line.ok()).collect();
    Ok(lines)
}
fn windows_cmd(program: &str, args: &str) -> Result<Vec<String>, Error> {
    if env::consts::OS == "windows" {
        let lines = cmd(program, args);
        return lines;
    } else {
        panic!("only run on windows");
    }
}
fn netsh_cmd_bool(args: &str) -> bool {
    let lines = windows_cmd("netsh", args).unwrap();
    let first_line = lines.get(0).unwrap();

    let mut success = false;
    if first_line.trim() == "" {
        success = true;
    } else {
        success = false;
    }
    success
}

// Public Functions
pub fn is_admin() -> bool {
    if env::consts::OS == "windows" {
        // cmd net session
        let lines = cmd("net", "").unwrap();
        // let lines = cmd("cmd", "net session").unwrap();
        println!("cmd net session");
        dbg!(&lines);
        for line in lines {
            // println!("{}", line);
            // if line.trim().starts_with("Configuration for interface") {
            // }
        }
        true
    } else {
        panic!("only run on windows");
    }
}
pub fn is_Ipv4Addr(ip: &str) -> bool {
    ip.parse::<Ipv4Addr>().is_ok()
}
pub fn get_interfaces() -> Vec<Interface> {
    // netsh interface ipv4 show config
    let lines = windows_cmd("netsh", "interface ipv4 show config").unwrap();

    let mut nics: Vec<Interface> = vec![];
    let mut new_interface = false;
    let mut is_dns_server_2 = false;
    let mut nic = Interface {
        interface_name: String::new(),
        interface_metric: 9999,
        gateway: String::new(),
        gateway_metric: 9999,
        ip_is_dhcp: true,
        ip_and_masks: vec![],
        dns_is_dhcp: true,
        dns_server_1: String::new(),
        dns_server_2: String::new(),
    };
    let mut ip_and_mask = Ip {
        ip_address: String::new(),
        subnet_mask: String::new(),
    };

    for line in lines {
        if line.trim().starts_with("Configuration for interface") {
            new_interface = true;
            let line_replace = line
                .replace("\"", "")
                .replace("Configuration for interface ", "");
            nic.interface_name = line_replace;
        } else if line.trim().starts_with("InterfaceMetric:") {
            let line_replace = line.replace("InterfaceMetric:", "");
            nic.interface_metric = line_replace.trim().parse().unwrap();
        } else if line.trim().starts_with("DHCP enabled:") {
            let line_replace = line.replace("DHCP enabled:", "");
            match line_replace.trim() {
                "No" => nic.ip_is_dhcp = false,
                "Yes" => nic.ip_is_dhcp = true,
                _ => nic.ip_is_dhcp = false,
            }
        } else if line.trim().starts_with("IP Address:") {
            let line_replace = line.trim().replace("IP Address:", "");
            ip_and_mask.ip_address = line_replace.trim().to_string();
        } else if line.trim().starts_with("Subnet Prefix:") {
            let clone = line.replace(")", "").clone();
            let line_replace: Vec<&str> = clone.split("(mask ").collect();
            ip_and_mask.subnet_mask = line_replace.get(1).unwrap().trim().to_string();
            nic.ip_and_masks.push(ip_and_mask);
            ip_and_mask = Ip {
                ip_address: String::new(),
                subnet_mask: String::new(),
            };
        } else if line.trim().starts_with("Default Gateway:") {
            let line_replace = line.trim().replace("Default Gateway:", "");
            nic.gateway = line_replace.trim().to_string();
        } else if line.trim().starts_with("Gateway Metric:") {
            let line_replace = line.trim().replace("Gateway Metric:", "");
            nic.gateway_metric = line_replace.trim().parse().unwrap();
        } else if line
            .trim()
            .starts_with("DNS servers configured through DHCP:")
        {
            nic.dns_is_dhcp = true;
            let line_replace = line
                .trim()
                .replace("DNS servers configured through DHCP:", "");
            if is_Ipv4Addr(line_replace.trim()) {
                nic.dns_server_1 = line_replace.trim().to_string();
                is_dns_server_2 = true;
            }
        } else if line
            .trim()
            .starts_with("Statically Configured DNS Servers:")
        {
            nic.dns_is_dhcp = false;
            let line_replace = line
                .trim()
                .replace("Statically Configured DNS Servers:", "");
            if is_Ipv4Addr(line_replace.trim()) {
                nic.dns_server_1 = line_replace.trim().to_string();
                is_dns_server_2 = true;
            }
        } else if line.trim() == "" && new_interface == true {
            nics.push(nic);
            nic = Interface {
                interface_name: String::new(),
                interface_metric: 9999,
                gateway: String::new(),
                gateway_metric: 9999,
                ip_is_dhcp: true,
                ip_and_masks: vec![],
                dns_is_dhcp: true,
                dns_server_1: String::new(),
                dns_server_2: String::new(),
            };
            ip_and_mask = Ip {
                ip_address: String::new(),
                subnet_mask: String::new(),
            };
        } else {
            if is_dns_server_2 && is_Ipv4Addr(line.trim()) {
                nic.dns_server_2 = line.trim().to_string();
            } else {
                // println!("{}", line);
            }
            is_dns_server_2 = false;
        }
    }
    nics
}
pub fn set_ip_dhcp(interface: &str) -> bool {
    // netsh interface ipv4 set address name="Ethernet" source=dhcp
    let args = format!("interface ipv4 set address name=\"{interface}\" source=dhcp");
    let success = netsh_cmd_bool(args.as_str());
    success
}
pub fn set_ip_static(interface: &str, ip: &str, mask: &str, gateway: &str) -> bool {
    // netsh interface ipv4 set address name="Ethernet" static 192.168.1.9 255.255.255.0 192.168.1.254
    let args =
        format!("interface ipv4 set address name=\"{interface}\" static {ip} {mask} {gateway}");
    let success = netsh_cmd_bool(args.as_str());
    success
}
pub fn add_ip_static(interface: &str, ip: &str, mask: &str) -> bool {
    // netsh interface ipv4 add address name="Ethernet" 192.168.2.9 255.255.255.0
    let args = format!("interface ipv4 add address name=\"{interface}\" {ip} {mask}");
    let success = netsh_cmd_bool(args.as_str());
    success
}
pub fn set_dns_dhcp(interface: &str) -> bool {
    // netsh interface ipv4 set dns name="Ethernet" source=dhcp
    let args = format!("interface ipv4 set dns name=\"{interface}\" source=dhcp");
    let success = netsh_cmd_bool(args.as_str());
    success
}
pub fn set_dns_static(interface: &str, dns_1: &str, dns_2: &str) -> bool {
    // netsh interface ipv4 set dns name="Ethernet" static 192.168.1.1
    // netsh interface ipv4 add dns name="Ethernet" 1.1.1.1 index=2
    let args_1 = format!("interface ipv4 set dns name=\"{interface}\" static {dns_1}");
    let success_1 = netsh_cmd_bool(args_1.as_str());
    let args_2 = format!("interface ipv4 add dns name=\"{interface}\" {dns_2} index=2");
    let success_2 = netsh_cmd_bool(args_2.as_str());
    success_1 && success_2
}
pub fn set_dhcp(interface: &str) -> bool {
    let success_dns = set_dns_dhcp(interface);
    let success_ip = set_ip_dhcp(interface);
    success_ip && success_dns
}
pub fn set_static(
    interface: &str,
    ip: &str,
    mask: &str,
    gateway: &str,
    dns_1: &str,
    dns_2: &str,
) -> bool {
    let success_ip = set_ip_static(interface, ip, mask, gateway);
    let success_dns = set_dns_static(interface, dns_1, dns_2);
    success_ip && success_dns
}
pub fn set_metric(interface: &str, metric: &str) -> bool {
    // netsh interface ipv4 set interface "Ethernet" metric=25
    let args = format!("interface ipv4 set interface \"{interface}\" metric={metric}");
    let success = netsh_cmd_bool(args.as_str());
    success
}
pub fn set_metric_auto(interface: &str) -> bool {
    // netsh interface ipv4 set interface "Ethernet" metric=auto
    let args = format!("interface ipv4 set interface \"{interface}\" metric=auto");
    let success = netsh_cmd_bool(args.as_str());
    success
}
pub fn set_interface_name(interface: &str, interface_new_name: &str) -> bool {
    // netsh interface set interface name="Ethernet" newname="New Ethernet"
    let args =
        format!("interface set interface name=\"{interface}\" newname=\"{interface_new_name}\"");
    let success = netsh_cmd_bool(args.as_str());
    success
}
pub fn release() {
    // ipconfig /release
    let args = format!("/release");
    let lines = windows_cmd("ipconfig", args.as_str()).unwrap();
}
pub fn renew() {
    // ipconfig /renew &
    // "&" runs in background. cheap way to not wait for the command to exit
    let args = format!("/renew &");
    let lines = windows_cmd("ipconfig", args.as_str()).unwrap();
}
pub fn flushdns() {
    // ipconfig /flushdns
    let args = format!("/flushdns");
    let lines = windows_cmd("ipconfig", args.as_str()).unwrap();
}
pub fn route_print_v4() -> Vec<Route> {
    // route PRINT -4
    let args = format!("PRINT -4");
    let lines = windows_cmd("route", args.as_str()).unwrap();

    let mut routes: Vec<Route> = vec![];
    let mut network_destination_flag = false;
    let mut network_address_flag = false;
    let mut route = Route {
        destination: String::new(),
        netmask: String::new(),
        gateway: String::new(),
        interface: String::new(),
        metric: 9999,
        persistent: false,
    };

    for line in lines {
        if line.trim().starts_with("Network Destination") {
            network_destination_flag = true;
            network_address_flag = false;
        } else if line.trim().starts_with("===") {
            network_destination_flag = false;
            network_address_flag = false;
        } else if line.trim().starts_with("Network Address") {
            network_destination_flag = false;
            network_address_flag = true;
        } else {
            if (network_destination_flag) {
                let clone = line.trim().clone();
                let route_string_array: Vec<&str> = clone.split_whitespace().collect();
                let destination = route_string_array.get(0).unwrap().trim().to_string();
                let netmask = route_string_array.get(1).unwrap().trim().to_string();
                let gateway = route_string_array.get(2).unwrap().trim().to_string();
                let interface = route_string_array.get(3).unwrap().trim().to_string();
                let metric: u16 = route_string_array.get(4).unwrap().trim().parse().unwrap();
                route = Route {
                    destination: destination,
                    netmask: netmask,
                    gateway: gateway,
                    interface: interface,
                    metric: metric,
                    persistent: false,
                };
                routes.push(route);
            } else if (network_address_flag) {
                let clone = line.trim().clone();
                let route_string_array: Vec<&str> = clone.split_whitespace().collect();
                let destination = route_string_array.get(0).unwrap().trim().to_string();
                let netmask = route_string_array.get(1).unwrap().trim().to_string();
                let gateway = route_string_array.get(2).unwrap().trim().to_string();
                let metric: u16 = route_string_array.get(3).unwrap().trim().parse().unwrap();
                let route_index = routes
                    .iter()
                    .position(|r| r.destination == destination)
                    .unwrap();
                routes[route_index].persistent = true;
            }
        }
    }
    routes
}
pub fn route_add(destination: &str, mask: &str, gateway: &str, metric: &str) {
    // route ADD 192.168.2.0 MASK 255.255.255.0 192.168.1.254 METRIC 1
    let args = format!("ADD {destination} MASK {mask} {gateway} METRIC {metric}");
    let lines = windows_cmd("route", args.as_str()).unwrap();
}
pub fn route_add_persistent(destination: &str, mask: &str, gateway: &str, metric: &str) {
    // route -p ADD 192.168.2.0 MASK 255.255.255.0 192.168.1.254 METRIC 1
    let args =
        format!("-p ADD {destination} MASK {mask} {gateway} METRIC {metric}");
    let lines = windows_cmd("route", args.as_str()).unwrap();
}
pub fn route_change(destination: &str, mask: &str, gateway: &str, metric: &str) {
    // route CHANGE 192.168.2.0 MASK 255.255.255.0 192.168.1.254 METRIC 1
    let args =
        format!("CHANGE {destination} MASK {mask} {gateway} METRIC {metric}");
    let lines = windows_cmd("route", args.as_str()).unwrap();
}
pub fn route_change_persistent(destination: &str, mask: &str, gateway: &str, metric: &str) {
    // route -p CHANGE 192.168.2.0 MASK 255.255.255.0 192.168.1.254 METRIC 1
    // can only change gateway, metric, persistent
    let args = format!("-p CHANGE {destination} MASK {mask} {gateway} METRIC {metric}");
    let lines = windows_cmd("route", args.as_str()).unwrap();
}
pub fn route_delete(destination: &str) {
    // route DELETE 157.0.0.0
    // deletes persistent routes as well
    let args = format!("-p DELETE {destination}");
    let lines = windows_cmd("route", args.as_str()).unwrap();
}

#[cfg(test)]
mod tests {
    use super::*;

    const GOOD_IPS: [&str; 5] = [
        "0.0.0.0",
        "1.1.1.1",
        "10.10.0.100",
        "0.0.0.255",
        "255.255.255.255",
    ];

    const BAD_IPS: [&str; 4] = [
        "192.168.1",
        "192.168.1.300",
        "192.168.x.1",
        "192.168.010.001",
    ];

    #[test]
    fn test_is_Ipv4Addr() {
        for ip in GOOD_IPS {
            assert_eq!(true, is_Ipv4Addr(ip));
        }
        for ip in BAD_IPS {
            assert_eq!(false, is_Ipv4Addr(ip));
        }
    }

    #[test]
    fn test_route_print_v4() {
        route_print_v4();
        // route_add("192.168.2.0", "255.255.255.0", "192.168.1.254", "1")
    }
}
