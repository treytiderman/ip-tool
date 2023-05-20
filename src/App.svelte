<script>
    // Imports
    import { state } from "./js/store_state";
    import { settings } from "./js/store_settings";

    // Components
    import Header from "./layout/Header.svelte";
    import Main from "./layout/Main.svelte";
    import KeyboardShortcuts from "./layout/KeyboardShortcuts.svelte";

    // Debug
    // $: console.log("$state", $state)
</script>

<div
    id="window"
    style="
        border: {$settings.hasDecorations ? 'none' : 'var(--border)'};
        border-color: {$state.is_administrator ? ' ' : 'var(--color-bg-red)'};
    "
>
    <Header />
    <div hidden={$state.is_administrator}>
        <i class="fa-solid fa-circle-exclamation" />
        Run as Administrator in order to change network settings
    </div>
    <Main />
</div>
<KeyboardShortcuts />

<style>
    #window {
        width: calc(100% - 1px);
        height: calc(100% - 1px);
        display: flex;
        flex-direction: column;
    }
    #window > div {
        background-color: var(--color-bg-red);
        color: var(--color-text-red);
        padding: var(--pad);
        margin: var(--pad);
    }
</style>
