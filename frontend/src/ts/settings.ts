// @ts-nocheck
import { writable, get } from "svelte/store";

export {
    defaultSettings,
    settingsStore,
    saveSettings,
    loadSettings,
    resetSettings,
    quitApp,
    windowSetMaximise,
    windowSetUnmaximise,
    windowSetMinimise,
    windowGetSize,
    windowUpdateSize,
    windowSetSize,
};

const defaultSettings = {
    theme: "auto",
    fontSize: 16,
    windowSizeWidth: 300,
    windowSizeHeight: 500,
    windowAlwaysOnTop: false,
    interfacePollRateMs: 3000,
    showPollAnimation: true,
    foldIfNotFocued: false,
    pollIfMinimised: false,
    pollIfNotFocued: false,
};

// Store
const settingsStore = writable(JSON.parse(JSON.stringify(defaultSettings)));

function saveSettings() {
    localStorage.setItem("settings", JSON.stringify(get(settingsStore)));
}

function loadSettings() {
    const localStorageSettings = localStorage.getItem("settings");
    console.log(
        "localStorage: settings",
        JSON.parse(localStorageSettings) || undefined,
    );
    if (localStorageSettings) {
        settingsStore.set(JSON.parse(localStorageSettings));
        windowSetSize();
    }
    settingsStore.subscribe(onSettingsChange);
}

function resetSettings() {
    settingsStore.set(JSON.parse(JSON.stringify(defaultSettings)));
    saveSettings();
}

// Window Functions
function quitApp() {
    window.runtime.Quit();
}

function windowSetMaximise() {
    window.runtime.WindowMaximise();
}

function windowSetUnmaximise() {
    window.runtime.WindowUnmaximise();
}

function windowSetMinimise() {
    window.runtime.WindowMinimise();
}

async function windowGetSize(): Promise<{ w: number; h: number }> {
    return await window.runtime.WindowGetSize();
}

function windowUpdateSize(w: number, h: number) {
    settingsStore.update((store) => {
        store.windowSizeWidth = w;
        store.windowSizeHeight = h;
        return store;
    });
}

function windowSetSize() {
    const settings = get(settingsStore);
    window.runtime.WindowSetSize(
        settings.windowSizeWidth,
        settings.windowSizeHeight,
    );
}

// Listen for updates
function onSettingsChange() {
    const settings = get(settingsStore);
    console.log("settings updated", settings);

    // Font Size
    const FONT_SIZE_MIN = 8;
    const FONT_SIZE_MAX = 36;
    settings.fontSize = Math.min(
        Math.max(settings.fontSize, FONT_SIZE_MIN),
        FONT_SIZE_MAX,
    );
    document.documentElement.style.setProperty(
        "--font-size",
        `${settings.fontSize}px`,
    );
    document.documentElement.style.setProperty(
        "--font-size-mono",
        `${settings.fontSize - 1}px`,
    );

    // Theme
    document.getElementsByTagName("html")[0].classList = settings.theme;

    // Always On Top
    window.runtime.WindowSetAlwaysOnTop(settings.windowAlwaysOnTop);

    // Window Size
    // handled manually with windowGetSize(), windowSetSize(), windowUpdateSize()

    saveSettings();
}

// Shrink Window on Blur. Go back to normal on Hover'
let blurTimer = null;
let lastWindowSize = {
    windowSizeWidth: settingsStore.windowSizeWidth,
    windowSizeHeight: settingsStore.windowSizeHeight,
};

async function shrinkWindow() {
    const s = get(settingsStore);

    if (s.foldIfNotFocued) {
        lastWindowSize.windowSizeWidth = s.windowSizeWidth;
        lastWindowSize.windowSizeHeight = s.windowSizeHeight;
        console.log("Window blurred", lastWindowSize);
        window.runtime.WindowSetSize(lastWindowSize.windowSizeWidth, 10);
    }
}

async function restoreWindow() {
    const s = get(settingsStore);

    if (s.foldIfNotFocued) {
        console.log("Window focused", s.windowSizeWidth, s.windowSizeHeight);
        window.runtime.WindowSetSize(
            lastWindowSize.windowSizeWidth,
            lastWindowSize.windowSizeHeight,
        );
    }
}

window.addEventListener("blur", () => {
    blurTimer = setTimeout(() => {
        blurTimer = null;
        shrinkWindow();
    }, 10);
});

window.addEventListener("focus", () => {
    if (blurTimer) {
        clearTimeout(blurTimer);
        blurTimer = null;
        return;
    }

    restoreWindow();
});
