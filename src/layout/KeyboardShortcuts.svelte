<script>
    // Imports
    import { state } from "../js/store_state";
    import { settings } from "../js/store_settings";
    import { numMinMax } from "../js/helper";
    import { appWindow, LogicalSize } from "@tauri-apps/api/window";

    // Functions
    function onKeyDown(event) {
        $state.keys = {
            key: event.key,
            ctrlKey: event.ctrlKey,
            altKey: event.altKey,
            shiftKey: event.shiftKey,
            keyCode: event.keyCode,
            clientX: event.clientX,
            clientY: event.clientY,
        };
        if (event.ctrlKey && event.key === "=") {
            $settings.fontSize = numMinMax(($settings.fontSize += 1), 8, 36);
        } else if (event.ctrlKey && event.key === "-") {
            $settings.fontSize = numMinMax(($settings.fontSize -= 1), 8, 36);
        } else if (event.ctrlKey && event.key === "t") {
            $settings.isAlwaysOnTop = !$settings.isAlwaysOnTop;
        } else if (event.ctrlKey && event.key === "t") {
            $settings.isAlwaysOnTop = !$settings.isAlwaysOnTop;
        } else if (event.ctrlKey && event.key === "l") {
            $settings.theme = $settings.theme === "dark" ? "light" : "dark";
        }
        // console.log($state.keys);
    }
    function onWheel(event) {
        $state.wheel = {
            ctrlKey: event.ctrlKey,
            altKey: event.altKey,
            shiftKey: event.shiftKey,
            scroll: Math.sign(-1 * event.deltaY),
            deltaY: event.deltaY,
            clientX: event.clientX,
            clientY: event.clientY,
        };
        if (event.ctrlKey && $state.wheel.scroll === 1) {
            $settings.fontSize = numMinMax(($settings.fontSize += 1), 8, 36);
        } else if (event.ctrlKey && $state.wheel.scroll === -1) {
            $settings.fontSize = numMinMax(($settings.fontSize -= 1), 8, 36);
        }
        // console.log($state.wheel);
    }
    function onClick(event) {
        $state.click = {
            ctrlKey: event.ctrlKey,
            altKey: event.altKey,
            shiftKey: event.shiftKey,
            which: event.which,
            clientX: event.clientX,
            clientY: event.clientY,
        };
        // which: 1 = left mouse button, 2 = middle, 3 = right
        if (event.ctrlKey && $state.click.which === 2) {
            $settings.fontSize = 16;
            appWindow.setSize(new LogicalSize(900, 900))
        }
        // console.log($state.click);
    }
</script>

<svelte:window on:keydown={onKeyDown} on:wheel={onWheel} on:mousewheel={onWheel} on:mouseup={onClick} />
