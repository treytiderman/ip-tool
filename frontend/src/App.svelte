<script lang="ts">
    import { onMount } from "svelte";
    import { pageStore } from "./ts/router";
    import { loadSettings } from "./ts/settings";
    import { loadPresets } from "./ts/presets";
    import * as app from "../wailsjs/go/main/App.js";

    // Components
    import Header from "./Header.svelte";
    import KeyboardShortcuts from "./KeyboardShortcuts.svelte";

    // Pages
    import Presets from "./pages/Presets.svelte";
    import Interfaces from "./pages/Interfaces.svelte";
    import EditPreset from "./pages/EditPreset.svelte";
    import CreatePreset from "./pages/CreatePreset.svelte";
    import EditInterface from "./pages/EditInterface.svelte";
    import Settings from "./pages/Settings.svelte";

    // App start up
    let isAdmin = true;
    onMount(async () => {

        // Check if the app is running as Administrator on Windows
        isAdmin = await app.IsAdmin();
        if (isAdmin) console.log("App is running as Administrator on Windows");
        else console.log("App is NOT running as Administrator");

        // Load localStorage
        loadSettings();
        loadPresets();

    });
</script>

<KeyboardShortcuts />

<div class="flex column pad-1 gap-1" style="height: 100svh; --default-contextmenu: show; overflow: hidden;">
    <Header title={$pageStore.name} />

    <hr>

    {#if isAdmin === false}
        <div class="pad-2 radius error error-bg error-border shadow">
            <div class="flex center-y gap-2">
                <div>
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="var(--border-width)"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <circle cx="12" cy="12" r="10" />
                        <path d="M12 16v-4" />
                        <path d="M12 8h.01" />
                    </svg>
                </div>
                <div>Run app as Administrator to change network settings</div>
            </div>
        </div>
    {/if}

    <main class="grow flex column overflow height-100">
        {#if $pageStore.name === "IPv4 Presets"}
            <Presets />
        {:else if $pageStore.name === "Interfaces"}
            <Interfaces />
        {:else if $pageStore.name === "Edit Preset"}
            <EditPreset />
        {:else if $pageStore.name === "Create Preset"}
            <CreatePreset />
        {:else if $pageStore.name === "Edit Interface"}
            <EditInterface />
        {:else if $pageStore.name === "Settings"}
            <Settings />
        {/if}
    </main>
</div>
