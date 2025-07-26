<script lang="ts">
    import * as app from "../wailsjs/go/main/App.js";

    // Components / Pages
    import Header from "./Header.svelte";
    import KeyboardShortcuts from "./KeyboardShortcuts.svelte";
    import { pageStore } from "./router";
    import Presets from "./Presets.svelte";
    import CreatePreset from "./CreatePreset.svelte";
    import EditPreset from "./EditPreset.svelte";
    import EditInterface from "./EditInterface.svelte";
    import { onMount } from "svelte";
    // import Settings from "./Settings.svelte"

    // Check if the app is running as Administrator on Windows
    let isAdmin = true
    onMount(async () => {
        isAdmin = await app.IsAdmin()
        if (isAdmin) console.log("App is running as Administrator on Windows");
        else console.log("App is NOT running as Administrator");
    })

</script>

<KeyboardShortcuts />

<div class="flex column" style="height: 100svh; --default-contextmenu: show;">
    <Header title={$pageStore.name} />

    {#if isAdmin === false}
        <div class="not-admin">
            <div class="flex center-y gap-sm">
                <div>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="10"/>
                        <path d="M12 16v-4"/>
                        <path d="M12 8h.01"/>
                    </svg>
                </div>
                <div>Run app as Administrator to change network settings</div>
            </div>
        </div>
    {/if}

    <main class="grow flex column">
        <div class="h-full" hidden={$pageStore.name !== "IPv4 Presets"}>
            <Presets />
        </div>
        <div class="h-full" hidden={$pageStore.name !== "Create Preset"}>
            <CreatePreset />
        </div>
        <div class="h-full" hidden={$pageStore.name !== "Edit Preset"}>
            <EditPreset />
        </div>
        <div class="h-full" hidden={$pageStore.name !== "Edit Interface"}>
            <EditInterface />
        </div>
    </main>
</div>

<style>
    main {
        overflow-y: auto;
        /* padding: var(--gap-sm);
        padding-top: 1px; */
    }
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
