<script>
    import { settingsStore, defaultSettings, windowGetSize, windowUpdateSize, windowSetSize } from "./ts/settings";

    // Helper
    // event.which, // 1 = left mouse button, 2 = middle, 3 = right
    // event.ctrlKey, // true if the control key is pressed
    // event.altKey, // true if the alt key is pressed
    // event.shiftKey, // true if the shift key is pressed
    // event.clientX, // X coordinate of the mouse pointer relative to the viewport
    // event.clientY, // Y coordinate of the mouse pointer relative to the viewport
    // event.keyCode, // Key code of the pressed key
    // event.deltaY, // Negative value means scrolling up, positive means scrolling down

    let isKeyDown = false;
    let isMouseDown = false;
    let fontChangeFactor = 1;

    function onMouseDown(event) {
        isMouseDown = true;

        if (event.ctrlKey && event.which === 2) {
            console.log("KeyboardShortcut: Reset font size and window size");
            settingsStore.update(store => {
                store.fontSize = defaultSettings.fontSize
                store.windowSizeWidth = defaultSettings.windowSizeWidth
                store.windowSizeHeight = defaultSettings.windowSizeHeight
                windowSetSize()
                return store
            })
        }
    }

    function onMouseUp(event) {
        isMouseDown = false;
    }

    async function onKeyDown(event) {
        isKeyDown = true;

        if (event.ctrlKey && event.key === "=") {
            console.log("KeyboardShortcut: Font size increase");
            settingsStore.update(store => {
                store.fontSize = store.fontSize + fontChangeFactor
                return store
            })
        } else if (event.ctrlKey && event.key === "-") {
            console.log("KeyboardShortcut: Font size decrease");
            settingsStore.update(store => {
                store.fontSize = store.fontSize - fontChangeFactor
                return store
            })
        }

        if (event.ctrlKey && event.key === "0") {
            console.log("KeyboardShortcut: Reset font size and window size");
            settingsStore.update(store => {
                store.fontSize = defaultSettings.fontSize
                store.windowSizeWidth = defaultSettings.windowSizeWidth
                store.windowSizeHeight = defaultSettings.windowSizeHeight
                windowSetSize()
                return store
            })
        }

        if (event.ctrlKey && event.key === "t") {
            console.log("KeyboardShortcut: Toggle Always on Top");
            settingsStore.update(store => {
                store.windowAlwaysOnTop = !store.windowAlwaysOnTop
                return store
            })
        }

        if (event.ctrlKey && event.key === "r") {
            location.reload()
        }
    }

    function onKeyUp(event) {
        isKeyDown = false;
    }

    async function onWheel(event) {
        if (event.ctrlKey && event.deltaY <= 0) {
            console.log("KeyboardShortcut: Font size increase");
            settingsStore.update(store => {
                store.fontSize = store.fontSize + fontChangeFactor
                return store
            })
        } else if (event.ctrlKey && event.deltaY >= 0) {
            console.log("KeyboardShortcut: Font size decrease");
            settingsStore.update(store => {
                store.fontSize = store.fontSize - fontChangeFactor
                return store
            })
        }
    }

    async function onResize(event) {
        const size = await windowGetSize()
        windowUpdateSize(size.w, size.h)
    }
</script>

<svelte:window
    on:keydown={onKeyDown}
    on:keyup={onKeyUp}
    on:pointerup={onMouseUp}
    on:pointerdown={onMouseDown}
    on:wheel={onWheel}
    on:resize={onResize}
/>
