// @ts-nocheck
import { writable } from "svelte/store"

export {
    windowStore as window,
    alwaysOnTopStore,
    themeStore,
    saveWindow,
    loadWindow,
    resetWindow,
    getFontSize,
    setFontSize,
    getWindowSize,
    setWindowSize,
    toggleWindowMaximise,
    setWindowMinimise,
    getWindowAlwaysOnTop,
    setWindowAlwaysOnTop,
    toggleWindowAlwaysOnTop,
    quitApp,
    setTheme,
}

const defaultState = {
    fontSize: 16,
    windowSizeX: 300,
    windowSizeY: 500,
    windowAlwaysOnTop: false,
    theme: "auto",
    // "ipPollRate_ms": 1000,
}

let windowStore = JSON.parse(JSON.stringify(defaultState))

const alwaysOnTopStore = writable(false)
const themeStore = writable("auto")

function saveWindow() {
    localStorage.setItem("window", JSON.stringify(windowStore))
    return windowStore
}

function loadWindow() {
    const savedStore = localStorage.getItem("window")
    if (savedStore) {
        windowStore = JSON.parse(savedStore)
    }
    console.log("localStorage: window", windowStore || undefined);
    setFontSize(windowStore.fontSize)
    setTheme(windowStore.theme)
    setWindowSize(windowStore.windowSizeX, windowStore.windowSizeY)
    setWindowAlwaysOnTop(windowStore.windowAlwaysOnTop)
    return windowStore
}

function resetWindow() {
    windowStore = JSON.parse(JSON.stringify(defaultState))
    saveWindow()
    return windowStore
}

function getFontSize() {
    return windowStore.fontSize
}

function setFontSize(size: number) {
    if (size <= 8 || size >= 36) return
    windowStore.fontSize = size
    document.documentElement.style.setProperty('--font-size', `${windowStore.fontSize}px`)
    document.documentElement.style.setProperty('--font-size-mono', `${windowStore.fontSize - 1}px`)
    saveWindow()
}

async function getWindowSize() {
    return await window.runtime.WindowGetSize()
}

function setWindowSize(x: number, y: number) {
    windowStore.windowSizeX = Math.round(x)
    windowStore.windowSizeY = Math.round(y)
    window.runtime.WindowSetSize(windowStore.windowSizeX, windowStore.windowSizeY)
    saveWindow()
}

function toggleWindowMaximise() {
    window.runtime.WindowToggleMaximise()
}

function setWindowMinimise() {
    window.runtime.WindowMinimise()
}

function getWindowAlwaysOnTop() {
    return windowStore.windowAlwaysOnTop
}

function setWindowAlwaysOnTop(alwaysOnTop: boolean) {
    windowStore.windowAlwaysOnTop = alwaysOnTop
    alwaysOnTopStore.set(alwaysOnTop)
    window.runtime.WindowSetAlwaysOnTop(windowStore.windowAlwaysOnTop)
    saveWindow()
}

function toggleWindowAlwaysOnTop() {
    setWindowAlwaysOnTop(!windowStore.windowAlwaysOnTop)
}

function quitApp() {
    window.runtime.Quit()
}

function setTheme(theme: string) {
    windowStore.theme = theme
    themeStore.set(theme)
    document.getElementsByTagName('html')[0].classList = theme
    saveWindow()
}
