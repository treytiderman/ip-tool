// Imports
import { writable } from 'svelte/store'
import { appWindow } from "@tauri-apps/api/window";
import { enable, disable, isEnabled } from "tauri-plugin-autostart-api";

// Export Store
export const settings = writable()

// Get / Init settings
export const default_settings = {
    "fontSize": 16,
    "theme": "dark",
    "themes": ["dark", "light"],
    "isAlwaysOnTop": false,
    "windowDecorations": false,
    "startOnBoot": false,
    "ipPollRate_ms": 1000,
    "colored_headers": true,
};
let settingsLocalStorage = localStorage.getItem("settings");
if (settingsLocalStorage == null || settingsLocalStorage === "") {
    settingsLocalStorage = default_settings
    console.log("settings: Nothing in localSorage. Setting to default", default_settings);
    settings.set(default_settings);
    localStorage.setItem("settings", JSON.stringify(default_settings));
}
else {
    const settingsLocalStorageParsed = JSON.parse(settingsLocalStorage);
    settings.set(settingsLocalStorageParsed);
    settings.update(obj => {
        Object.keys(default_settings).forEach(key => {
            if (obj[key] === undefined) {
                console.log(`setting: ${key} doesn't exist in localStorage. Setting to default ${default_settings[key]}`);
                obj[key] = default_settings[key];
            }
        })
        return obj
    });
}

// Every settings change set the values 
settings.subscribe(async settings => {
    if (settings !== undefined && settings != default_settings) {
        // console.log("settings: updated", settings);

        // Update Settings
        appWindow.setAlwaysOnTop(settings.isAlwaysOnTop);

        // Hide Title bar, don't hide in the Tauri config so they show unless this page loads correctly
        appWindow.setDecorations(settings.windowDecorations);


        // Theme
        document.documentElement.classList = settings.theme;

        // Font Size
        if (settings.fontSize >= 8 && settings.fontSize <= 36) {
            document.documentElement.style.fontSize = settings.fontSize + "px";
        }

        // Start on boot
        try {
            if (settings.startOnBoot) {
                await enable();
                settings.startOnBoot = await isEnabled();
            }
            else {
                await disable();
                settings.startOnBoot = await isEnabled();
            }
        }
        catch (error) {
            // console.log(error);
        }

        // Save settings to localStorage
        localStorage.setItem("settings", JSON.stringify(settings));

    }
})