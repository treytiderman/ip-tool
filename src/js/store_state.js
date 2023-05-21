import { writable, get } from 'svelte/store'
import { appWindow, LogicalPosition, LogicalSize } from "@tauri-apps/api/window"
import { Command } from "@tauri-apps/api/shell"

// Variables
const default_state = {
    is_administrator: true,
    nav_tab_active: "IPv4",
    nav_tabs: ["IPv4", "Routes", "Settings"],
    window_size: { x: 922, y: 800 },
    window_position: { x: 0, y: 0 },
    window_scale_factor: 1,
    mouse: {
        which: 1,
        isDown: false,
        ctrlKey: false,
        altKey: false,
        shiftKey: false,
        clientX: 0,
        clientY: 0,
    },
    keys: {
        key: "",
        isDown: false,
        ctrlKey: false,
        altKey: false,
        shiftKey: false,
        keyCode: 0,
        clientX: 0,
        clientY: 0,
    },
    wheel: {
        scroll: 1,
        ctrlKey: false,
        altKey: false,
        shiftKey: false,
        deltaY: 150,
        clientX: 0,
        clientY: 0,
    },
}
export const state = createStore()

// Functions
async function check_administrator() {
    const cmd = new Command("net", "session")
    cmd.stdout.on("data", (line) => {
        // console.log("line", line);
    })
    cmd.stderr.on("data", (line) => {
        // console.log("line", line);
        if (line.startsWith("System error 5 has occurred")) {
            console.log("App is not running as Administrator")
            state.update(state => {
                state.is_administrator = false
                return state
            })
        }
    })
    await cmd.spawn()
}
async function window_get_scale() {
    let scale = await appWindow.scaleFactor()
    state.update(state => {
        state.window_scale_factor = scale
        return state
    })
}
async function window_set_size(height, width) {
    let state_now = get(state)
    let scaled_height = height / state_now.window_scale_factor
    let scaled_width = width / state_now.window_scale_factor
    return appWindow.setSize(new LogicalSize(scaled_height, scaled_width))
}
async function window_set_position(x, y) {
    let state_now = get(state)
    let scaled_x = x / state_now.window_scale_factor
    let scaled_y = y / state_now.window_scale_factor
    return appWindow.setPosition(new LogicalPosition(scaled_x, scaled_y))
}
function state_initialize() {
    if (localStorage.getItem("state") === null) state_default()
    else state_from_localStorage()
}
function state_default() {
    state.set(default_state)
    localStorage.setItem("state", JSON.stringify(default_state))
}
function state_from_localStorage() {
    try {
        const localStorage_json = localStorage.getItem("state")
        const localStorage_parsed = JSON.parse(localStorage_json)
        const state_new = {}
        Object.keys(default_state).forEach(key => {
            if (localStorage_parsed[key] == null) state_new[key] = default_state[key]
            else state_new[key] = localStorage_parsed[key]
        })
        state_new.is_administrator = true;
        state.set(state_new)
        let state_now = get(state)
        // window_set_size(state_now.window_size.height, state_now.window_size.width)
        // window_set_position(state_now.window_position.x, state_now.window_position.y)
    }
    catch (error) {
        console.log("state: localStorage not a JSON. Setting to default")
        state_default()
    }
}
function state_on_change() {
    state.subscribe(state => {
        localStorage.setItem("state", JSON.stringify(state))
    })
}
function createStore() {
    const { subscribe, set, update } = writable(default_state)

    // Store functions
    return {
        set,
        update,
        subscribe,

        check_administrator,

        state_initialize,
        state_from_localStorage,
    }
}

// Window Events
const unlisten_moved = await appWindow.onMoved(({ payload: position }) => {
    state.update(state => {
        if (state.window_position !== position) {
            state.window_position = position
        }
        return state
    })
})
const unlisten_resized = await appWindow.onResized(({ payload: size }) => {
    state.update(state => {
        if (state.window_size !== size) {
            state.window_size = size
        }
        return state
    })
})
const unlisten_scaled = await appWindow.onScaleChanged(({ payload }) => {
    state.update(state => {
        if (state.window_scale_factor !== payload.scaleFactor) {
            state.window_scale_factor = payload.scaleFactor
        }
        return state
    })
})

// Startup
state_initialize()
state_on_change()
check_administrator()
window_get_scale()