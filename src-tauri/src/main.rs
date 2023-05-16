#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

use tauri_plugin_autostart::MacosLauncher;
mod invoke;
mod netsh;
use tauri::Manager;

fn main() {
    tauri::Builder::default()
        .setup(|app| {
            let main_window = app.get_window("main").unwrap();
            tauri::async_runtime::spawn(async move {
                std::thread::sleep(std::time::Duration::from_millis(500));
                main_window.show().unwrap();
            });
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            invoke::greet,
            invoke::move_mouse_test,
            invoke::is_admin,
            invoke::get_interfaces,
            invoke::set_ip_dhcp,
            invoke::set_ip_static,
            invoke::add_ip_static,
            invoke::set_dns_dhcp,
            invoke::set_dns_static,
            invoke::set_dhcp,
            invoke::set_static,
            invoke::set_metric,
            invoke::set_metric_auto,
            invoke::set_interface_name,
            invoke::release,
            invoke::renew,
            invoke::flushdns,
        ])
        .plugin(tauri_plugin_autostart::init(
            MacosLauncher::LaunchAgent,
            Some(vec!["--flag1", "--flag2"]),
        ))
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
