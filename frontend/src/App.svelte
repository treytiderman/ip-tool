<script lang="ts">
    import Header from "./Header.svelte";
    import KeyboardShortcuts from "./KeyboardShortcuts.svelte";
    import { pageStore } from "./router";
    import Presets from "./Presets.svelte";
    import CreatePreset from "./CreatePreset.svelte";
    import EditPreset from "./EditPreset.svelte";
    import EditInterface from "./EditInterface.svelte";
    // import Settings from "./Settings.svelte"
</script>

<KeyboardShortcuts />

<div class="flex column" style="height: 100svh;">
    <Header title={$pageStore.name} />

    <main class="grow flex column">
        <div class="h-full" hidden={$pageStore.name !== "Presets"}>
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
        padding: var(--gap-sm);
        padding-top: 1px;
    }
</style>
