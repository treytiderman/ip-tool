<script lang="ts">
    import { onMount } from "svelte";
    import { pageStore } from "./ts/router";
    import { loadWindow } from "./ts/window";
    import { loadPresets } from "./ts/presets";
    import { initNics, pollNics } from "./ts/nic";
    import * as app from "../wailsjs/go/main/App.js";

    // Components
    import Header from "./Header.svelte";
    import KeyboardShortcuts from "./KeyboardShortcuts.svelte";

    // Pages
    import Presets from "./pages/Presets.svelte";
    import EditPreset from "./pages/EditPreset.svelte";
    import CreatePreset from "./pages/CreatePreset.svelte";
    import EditInterface from "./pages/EditInterface.svelte";

    // App start up
    let isAdmin = true;
    onMount(async () => {
        // Check if the app is running as Administrator on Windows
        isAdmin = await app.IsAdmin();
        if (isAdmin) console.log("App is running as Administrator on Windows");
        else console.log("App is NOT running as Administrator");

        // Load presets
        loadWindow();

        // Load presets
        loadPresets();

        // Poll for changes in network interfaces
        await initNics();
        const interval = pollNics();
    });
</script>

<KeyboardShortcuts />

<div class="flex column" style="height: 100svh; --default-contextmenu: show;">
    <Header title={$pageStore.name} />

    {#if isAdmin === false}
        <div class="not-admin">
            <div class="flex center-y gap-sm">
                <div>
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
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

    <main class="grow flex column overflow">
        <div class="h-full" hidden={$pageStore.name !== "IPv4 Presets"}>
            <Presets />
        </div>
        <div class="h-full" hidden={$pageStore.name !== "Edit Preset"}>
            <EditPreset />
        </div>
        <div class="h-full" hidden={$pageStore.name !== "Create Preset"}>
            <CreatePreset />
        </div>
        <div class="h-full" hidden={$pageStore.name !== "Edit Interface"}>
            <EditInterface />
        </div>
    </main>
</div>

<style>
    .not-admin {
        padding: var(--gap-sm);
        padding-top: 1px;
    }
    .not-admin > div {
        padding: var(--gap-sm);
        color: var(--color-bg-red);
        background-color: var(--color-text-red);
        border-radius: var(--border-radius);
    }
</style>
