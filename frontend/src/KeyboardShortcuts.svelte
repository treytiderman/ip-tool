<script>
    import * as app from "../wailsjs/go/main/App.js";
    import { getFontSize, setFontSize, getWindowSize, setWindowSize, toggleWindowAlwaysOnTop } from "./ts/window";

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
    let defaultFontSize = 16;
    let fontChangeFactor = 1;
    let zoomChangeFactor = 0.05;

    function onMouseDown(event) {
        isMouseDown = true;

        if (event.ctrlKey && event.which === 2) {
            console.log("KeyboardShortcut: Reset font size and window size");
            setFontSize(defaultFontSize);
            setWindowSize(300, 500);
        }
    }

    function onMouseUp(event) {
        isMouseDown = false;
    }

    async function onKeyDown(event) {
        isKeyDown = true;

        if (event.ctrlKey && event.key === "=") {
            console.log("KeyboardShortcut: Font size increase");
            const size = await getWindowSize();
            const newWidth = Math.round(size.w * (1 + zoomChangeFactor));
            const newHeight = Math.round(size.h * (1 + zoomChangeFactor));
            setWindowSize(newWidth, newHeight);
            setFontSize(getFontSize() + fontChangeFactor);
        } else if (event.ctrlKey && event.key === "-") {
            console.log("KeyboardShortcut: Font size decrease");
            const size = await getWindowSize();
            const newWidth = Math.round(size.w * (1 - zoomChangeFactor));
            const newHeight = Math.round(size.h * (1 - zoomChangeFactor));
            setFontSize(getFontSize() - fontChangeFactor);
            setWindowSize(newWidth, newHeight);
        }

        if (event.ctrlKey && event.key === "0") {
            console.log("KeyboardShortcut: Reset font size and window size");
            setWindowSize(300, 500);
            setFontSize(defaultFontSize);
        }

        if (event.ctrlKey && event.key === "t") {
            console.log("KeyboardShortcut: Toggle Always on Top");
            toggleWindowAlwaysOnTop();
        }

        if (event.ctrlKey && event.key === "w") {
            console.log("KeyboardShortcut: GetInterfaces");
            const nics = await app.GetInterfaces();
            console.log(nics);
        }
    }

    function onKeyUp(event) {
        isKeyDown = false;
    }

    async function onWheel(event) {
        if (event.ctrlKey && event.deltaY <= 0) {
            console.log("KeyboardShortcut: Font size increase");
            const size = await getWindowSize();
            const newWidth = Math.round(size.w * (1 + zoomChangeFactor));
            const newHeight = Math.round(size.h * (1 + zoomChangeFactor));
            setWindowSize(newWidth, newHeight);
            setFontSize(getFontSize() + fontChangeFactor);
        } else if (event.ctrlKey && event.deltaY >= 0) {
            console.log("KeyboardShortcut: Font size decrease");
            const size = await getWindowSize();
            const newWidth = Math.round(size.w * (1 - zoomChangeFactor));
            const newHeight = Math.round(size.h * (1 - zoomChangeFactor));
            setFontSize(getFontSize() - fontChangeFactor);
            setWindowSize(newWidth, newHeight);
        }
    }
</script>

<svelte:window
    on:keydown={onKeyDown}
    on:keyup={onKeyUp}
    on:pointerup={onMouseUp}
    on:pointerdown={onMouseDown}
    on:wheel={onWheel}
/>
