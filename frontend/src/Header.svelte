<script lang="ts">
    export let title: string = "Title";
    import { setHomePage } from "./ts/router";
    import { quitApp, setWindowMinimise, alwaysOnTopStore, toggleWindowAlwaysOnTop } from "./ts/window";
</script>

<header class="flex center-y">
    {#if title !== "IPv4 Presets"}
        <button class="ip-icon-button" on:click={setHomePage} title="Back to IPv4 Presets">
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
                <path d="m12 19-7-7 7-7" />
                <path d="M19 12H5" />
            </svg>
        </button>
    {/if}

    <div class="grow" style="--wails-draggable:drag">{title}</div>
    {#if $alwaysOnTopStore}
        <button class="ip-icon-button" on:click={toggleWindowAlwaysOnTop} title="Disable Always on Top">
            <svg
                style="height: 0.875rem;"
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
                <path d="M12 17v5" />
                <path
                    fill="currentColor"
                    d="M9 10.76a2 2 0 0 1-1.11 1.79l-1.78.9A2 2 0 0 0 5 15.24V16a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-.76a2 2 0 0 0-1.11-1.79l-1.78-.9A2 2 0 0 1 15 10.76V7a1 1 0 0 1 1-1 2 2 0 0 0 0-4H8a2 2 0 0 0 0 4 1 1 0 0 1 1 1z"
                />
            </svg>
        </button>
    {:else}
        <button class="ip-icon-button color-dim" on:click={toggleWindowAlwaysOnTop} title="Enable Always on Top">
            <svg
                style="height: 0.875rem;"
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
                <path d="M12 17v5" />
                <path fill="s" d="M15 9.34V7a1 1 0 0 1 1-1 2 2 0 0 0 0-4H7.89" />
                <path d="m2 2 20 20" />
                <path d="M9 9v1.76a2 2 0 0 1-1.11 1.79l-1.78.9A2 2 0 0 0 5 15.24V16a1 1 0 0 0 1 1h11" />
            </svg>
        </button>
    {/if}
    <button class="ip-icon-button" on:click={setWindowMinimise} title="Minimise App">
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
            <path d="M5 12h14" />
        </svg>
    </button>
    <button class="ip-icon-button red-hover" on:click={quitApp} title="Quit App">
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
            <path d="M18 6 6 18" />
            <path d="m6 6 12 12" />
        </svg>
    </button>
</header>

<style>
    header {
        padding: var(--gap-xs, 0.25rem);
    }
    header div {
        cursor: pointer;
        user-select: none;
        padding-inline: var(--gap-xs, 0.25rem);
    }
</style>
