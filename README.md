# README
 
## Dev

Run `wails dev`


## Build

Run `wails build`
Run `wails build -clean -o ip-tool.exe`


## Go Tests

Run `go test`


## Ideas

- [x] Add DHCP functions renew, release
- [ ] Add CIDR input of Subnet Mask if starts with "/"
- [ ] Add disable/enable interfaces
- [ ] Add css themes
- [ ] Add logging for each netsh command
- [ ] Add settings page
    - [ ] Be able to change/add filters to interfaces. Currently if contain "Loop" it does not show
    - [ ] Add poll interfaces time control
    - [ ] Add setting for when to poll interfaces (focused, minimized)
    - [ ] Add reset presets button
    - [ ] Add export presets button


## Always start as Administrator

Add the "trustInfo" section to wails.exe.manifest

```
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly manifestVersion="1.0" xmlns="urn:schemas-microsoft-com:asm.v1" xmlns:asmv3="urn:schemas-microsoft-com:asm.v3">
    <assemblyIdentity type="win32" name="com.wails.{{.Name}}" version="{{.Info.ProductVersion}}.0" processorArchitecture="*"/>
    <dependency>
        <dependentAssembly>
            <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
        </dependentAssembly>
    </dependency>
    <asmv3:application>
        <asmv3:windowsSettings>
            <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true/pm</dpiAware> <!-- fallback for Windows 7 and 8 -->
            <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">permonitorv2,permonitor</dpiAwareness> <!-- falls back to per-monitor if per-monitor v2 is not supported -->
        </asmv3:windowsSettings>
    </asmv3:application>
    <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
        <security>
            <requestedPrivileges>
                <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
            </requestedPrivileges>
        </security>
    </trustInfo>
</assembly>
```

Adding the "windows" key to wails.json did not work

```
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "ip-tool",
  "outputfilename": "ip-tool",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "Trey Tiderman",
    "email": "tiderman.trey@gmail.com"
  },
  "windows": {
    "webview2": "browser",
    "rtl": false,
    "disableWindowIcon": false,
    "darkMode": true,
    "theme": "system",
    "messages": {
      "runtimeNotInstalled": "You need to install the WebView2 runtime to run this application. Would you like to download it now?",
      "webview2DownloadError": "Failed to download the WebView2 runtime. Please check your internet connection.",
      "webview2InstallationError": "Failed to install the WebView2 runtime. Please try again or install it manually."
    },
    "executionLevel": "requireAdministrator"
  }
}
```

