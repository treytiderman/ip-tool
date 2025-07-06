// @ts-nocheck

import { writable } from "svelte/store"

export {
    state as store,
    alwaysOnTopStore,
    saveState as saveStore,
    loadState as loadStore,
    resetState as resetStore,
    getFontSize,
    setFontSize,
    getWindowSize,
    setWindowSize,
    toggleWindowMaximise,
    setWindowMinimise,
    getWindowAlwaysOnTop,
    setWindowAlwaysOnTop,
    toggleWindowAlwaysOnTop,
    quitApp
}

const defaultState = {
    fontSize: 16,
    windowSizeX: 300,
    windowSizeY: 500,
    windowAlwaysOnTop: false,
    // "ipPollRate_ms": 1000,
}

let state = JSON.parse(JSON.stringify(defaultState))

const alwaysOnTopStore = writable(false)


// Start Up

loadState()


// Functions

function saveState() {
    localStorage.setItem("store", JSON.stringify(state))
    return state
}

function loadState() {
    const savedStore = localStorage.getItem("store")
    if (savedStore) {
        state = JSON.parse(savedStore)
    }
    setFontSize(state.fontSize)
    setWindowSize(state.windowSizeX, state.windowSizeY)
    setWindowAlwaysOnTop(state.windowAlwaysOnTop)
    return state
}

function resetState() {
    state = JSON.parse(JSON.stringify(defaultState))
    saveState()
    return state
}

function getFontSize() {
    return state.fontSize
}

function setFontSize(size: number) {
    if (size <= 8 || size >= 36) return
    state.fontSize = size
    document.documentElement.style.setProperty('--font-size', `${state.fontSize}px`)
    document.documentElement.style.setProperty('--font-size-mono', `${state.fontSize - 1}px`)
    saveState()
}

async function getWindowSize() {
    return await window.runtime.WindowGetSize()
}

function setWindowSize(x: number, y: number) {
    state.windowSizeX = Math.round(x)
    state.windowSizeY = Math.round(y)
    window.runtime.WindowSetSize(state.windowSizeX, state.windowSizeY)
    saveState()
}

function toggleWindowMaximise() {
    window.runtime.WindowToggleMaximise()
}

function setWindowMinimise() {
    window.runtime.WindowMinimise()
}

function getWindowAlwaysOnTop() {
    return state.windowAlwaysOnTop
}

function setWindowAlwaysOnTop(alwaysOnTop: boolean) {
    state.windowAlwaysOnTop = alwaysOnTop
    alwaysOnTopStore.set(alwaysOnTop)
    window.runtime.WindowSetAlwaysOnTop(state.windowAlwaysOnTop)
    saveState()
}

function toggleWindowAlwaysOnTop() {
    setWindowAlwaysOnTop(!state.windowAlwaysOnTop)
}

function quitApp() {
    window.runtime.Quit()
}
