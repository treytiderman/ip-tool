<script lang="ts">
    import { settingsStore, windowSetSize } from "../ts/settings";
    import { presets, resetPresets } from "../ts/presets";

    async function fileUploaded(event: any) {
        const file = event.target.files?.[0];
        if (!file) return;

        const name = file.name
        const text = await file.text();
        if (!JSON.parse(text)) return;
        const obj = JSON.parse(text);

        if (
            window.confirm(
                `Are you sure you want to overwrite all presets with the contents of "${name}"? This cannot be undone.`,
            )
        ) {
            console.log("Re-write Presets", obj);
            resetPresets(obj);
        }
    }
</script>

<div class="grid gap-4 pad-1">
    <div class="grid gap-1">
        <div class="flex bottom">
            <div class="grow text-dark thin">Presets</div>
            <button
                class="transparent text-dark error-bg-hover"
                title="Reset ALL Presets to Default"
                on:click={() => {
                    if (window.confirm("Are you sure you want to reset presets? This cannot be undone.")) {
                        resetPresets();
                    }
                }}
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
                    <path d="M3 3v5h5" />
                </svg>
            </button>
        </div>

        <div class="flex wrap even gap-2 break-xs">
            <a
                class="border border-radius padding bg-light bg-light-hover text-decoration-none shadow flex center-y gap-2"
                download="presets.json"
                href={URL.createObjectURL(new Blob([JSON.stringify($presets, null, 4)], { type: "text/json" }))}
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <path
                        d="M6 22a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h8a2.4 2.4 0 0 1 1.704.706l3.588 3.588A2.4 2.4 0 0 1 20 8v12a2 2 0 0 1-2 2z"
                    />
                    <path d="M14 2v5a1 1 0 0 0 1 1h5" />
                    <path d="M12 18v-6" />
                    <path d="m9 15 3 3 3-3" />
                </svg>
                <div>Download</div>
            </a>
            <label
                for="file"
                class="padding border border-radius text-light bg-light bg-light-hover shadow width-100 pointer flex center-y gap-2"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <path
                        d="M6 22a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h8a2.4 2.4 0 0 1 1.704.706l3.588 3.588A2.4 2.4 0 0 1 20 8v12a2 2 0 0 1-2 2z"
                    />
                    <path d="M14 2v5a1 1 0 0 0 1 1h5" />
                    <path d="M12 12v6" />
                    <path d="m15 15-3-3-3 3" />
                </svg>
                Upload
            </label>
        </div>

        <input type="file" name="file" id="file" class="width-100" hidden on:change={fileUploaded} />
    </div>

    <div class="flex even center-y gap-2">
        <div class="grid gap-1">
            <label for="theme" class="text-dark thin">Theme</label>
            <select id="theme" class="border shadow width-100 height-7" bind:value={$settingsStore.theme}>
                <option value="auto">Auto</option>
                <option value="black">Black</option>
                <option value="dark">Dark</option>
                <option value="white">White</option>
            </select>
        </div>

        <div class="grid gap-1">
            <label for="fontSize" class="text-dark thin">Font size</label>
            <input
                type="number"
                name="fontSize"
                id="fontSize"
                class="border shadow-inset width-100 height-7"
                bind:value={$settingsStore.fontSize}
                placeholder="16"
            />
        </div>
    </div>

    <div class="flex even center-y gap-2">
        <div class="grid gap-1">
            <label for="windowSizeWidth" class="text-dark thin">Window Width</label>
            <input
                type="number"
                name="windowSizeWidth"
                id="windowSizeWidth"
                class="border shadow-inset width-100"
                on:blur={windowSetSize}
                bind:value={$settingsStore.windowSizeWidth}
                placeholder="300"
                min="100"
                step="10"
            />
        </div>
        <div class="grid gap-1">
            <label for="windowSizeHeight" class="text-dark thin">Window Height</label>
            <input
                type="number"
                name="windowSizeHeight"
                id="windowSizeHeight"
                class="border shadow-inset width-100"
                on:blur={windowSetSize}
                bind:value={$settingsStore.windowSizeHeight}
                placeholder="500"
                min="100"
                step="10"
            />
        </div>
    </div>

    <div class="grid gap-1">
        <label for="ip_poll_rate_ms" class="text-dark thin">Interface Poll Rate (ms)</label>
        <input
            type="number"
            name="ip_poll_rate_ms"
            id="ip_poll_rate_ms"
            class="border shadow-inset"
            bind:value={$settingsStore.interfacePollRateMs}
            placeholder="3000"
            min="500"
            max="20000"
            step="100"
        />
    </div>

    <div class="flex center-y gap-1">
        <input
            type="checkbox"
            name="show_poll_animation"
            id="show_poll_animation"
            class="border shadow-inset"
            bind:checked={$settingsStore.showPollAnimation}
        />
        <label for="show_poll_animation" class="text-dark thin">Show Poll Animation</label>
    </div>

    <div class="flex center-y gap-1">
        <input
            type="checkbox"
            name="poll_if_not_focued"
            id="poll_if_not_focued"
            class="border shadow-inset"
            bind:checked={$settingsStore.pollIfNotFocued}
        />
        <label for="poll_if_not_focued" class="text-dark thin">Poll If Not Focused</label>
    </div>

    <!-- <div class="flex center-y gap-1">
        <input
            type="checkbox"
            name="poll_if_minimised"
            id="poll_if_minimised"
            class="border shadow-inset"
            bind:checked={$settingsStore.pollIfMinimised}
        />
        <label for="poll_if_minimised" class="text-dark thin">Poll If Minimised</label>
    </div> -->
</div>
