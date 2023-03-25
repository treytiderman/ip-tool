# IP Settings / Presets (Windows)

Quickly change and save network settings.
For people that plug into a lot of networks

### Download

[Installer](https://github.com/TreyTiderman/IP-Tool/releases/download/v0.2.0/IP-Tool_0.2.1_Installer.msi)
/
[Portable](https://github.com/TreyTiderman/IP-Tool/releases/download/v0.2.0/IP-Tool_0.2.1_Portable.exe)

![Preview UI](./public/preview.png)

# Source Code

### Setup

https://tauri.app/v1/guides/getting-started/prerequisites

```
npm install
```

### Run

```
npm run tauri dev
```

### Build

```
npm run tauri build
```

# Notes

### TODO

- [ ] alternating colors shift with contextmenu
- [ ] row moving
- [ ] table sorting
- [ ] routes page
- [ ] Hide interface
- [ ] More options to add presets
- [ ] Network Scan page
- [x] Interface friendly name change

### Network Scan

https://crates.io/crates/netscan

### Force Admin

https://github.com/nabijaczleweli/rust-embed-resource#example-embedding-a-windows-manifest

### Logging?

https://github.com/tauri-apps/tauri-plugin-log

### Add to `tauri.conf.json` for offline install of WebView2

```json
{
  "tauri": {
    "bundle": {
      "windows": {
        "webviewInstallMode": {
          "type": "offlineInstaller"
        }
      }
    }
  }
}
```
