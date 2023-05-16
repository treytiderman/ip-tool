use crate::netsh;
use enigo::{Enigo, MouseControllable};
use serde_json;

#[tauri::command]
pub fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted the by Rust backend!", name)
}

#[tauri::command]
pub fn move_mouse_test() {
    let mut enigo = Enigo::new();
    enigo.mouse_move_to(500, 200);
}

#[tauri::command]
pub fn is_admin() -> bool {
    println!("is_admin");
    let output = netsh::is_admin();
    output
}

#[tauri::command]
pub fn get_interfaces() -> String {
    // println!("get_interfaces");
    let output = netsh::get_interfaces();
    serde_json::to_string(&output).unwrap()
}

#[tauri::command]
pub fn set_ip_dhcp(interface: &str) -> bool {
    println!("set_ip_dhcp");
    let output = netsh::set_ip_dhcp(interface);
    output
}

#[tauri::command]
pub fn set_ip_static(interface: &str, ip: &str, mask: &str, gateway: &str) -> bool {
    println!("set_ip_static");
    let output = netsh::set_ip_static(interface, ip, mask, gateway);
    output
}

#[tauri::command]
pub fn add_ip_static(interface: &str, ip: &str, mask: &str) -> bool {
    println!("add_ip_static");
    let output = netsh::add_ip_static(interface, ip, mask);
    output
}

#[tauri::command]
pub fn set_dns_dhcp(interface: &str) -> bool {
    println!("set_dns_dhcp");
    let output = netsh::set_dns_dhcp(interface);
    output
}

#[tauri::command]
pub fn set_dns_static(interface: &str, dns1: &str, dns2: &str) -> bool {
    println!("set_dns_static");
    let output = netsh::set_dns_static(interface, dns1, dns2);
    output
}

#[tauri::command]
pub fn set_dhcp(interface: &str) -> bool {
    println!("set_dhcp");
    let output = netsh::set_dhcp(interface);
    output
}

#[tauri::command]
pub fn set_static(
    interface: &str,
    ip: &str,
    mask: &str,
    gateway: &str,
    dns1: &str,
    dns2: &str,
) -> bool {
    println!("set_static");
    let output = netsh::set_static(interface, ip, mask, gateway, dns1, dns2);
    output
}

#[tauri::command]
pub fn set_metric(interface: &str, metric: &str) -> bool {
    println!("set_metric");
    let output = netsh::set_metric(interface, metric);
    output
}

#[tauri::command]
pub fn set_metric_auto(interface: &str) -> bool {
    println!("set_metric_auto");
    let output = netsh::set_metric_auto(interface);
    output
}

#[tauri::command]
pub fn set_interface_name(interface: &str, interface_new_name: &str) -> bool {
    println!("set_interface_name {interface} {interface_new_name}");
    let output = netsh::set_interface_name(interface, interface_new_name);
    output
}

#[tauri::command]
pub fn release() {
    println!("release");
    netsh::release();
}

#[tauri::command]
pub fn renew() {
    println!("renew");
    netsh::renew();
}

#[tauri::command]
pub fn flushdns() {
    println!("flushdns");
    netsh::flushdns();
}
