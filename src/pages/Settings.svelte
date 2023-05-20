<script>
    import { settings, default_settings } from "../js/store_settings";
    let showReset = false;
</script>

<article class="mono" class:coloredHeaders={$settings.coloredHeaders}>
    <table style="width: 100%;">
        <thead>
            <tr>
                <th>Key</th>
                <th>Value</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td><div>defaultSettings</div></td>
                <td>
                    <button
                        hidden={showReset}
                        on:click={() => {
                            showReset = true;
                        }}
                    >
                        set
                    </button>
                    <button
                        hidden={!showReset}
                        style="color: var(--color-bg-red)"
                        on:click={() => {
                            showReset = false;
                        }}
                    >
                        no cancel!
                    </button>
                    <button
                        hidden={!showReset}
                        style="color: var(--color-bg-green)"
                        on:click={() => {
                            showReset = false;
                            $settings = default_settings;
                            localStorage.setItem("settings", JSON.stringify(default_settings));
                            location.reload(true);
                        }}
                    >
                        reset settings to default?
                    </button>
                </td>
            </tr>
            <tr>
                <td><div>isAlwaysOnTop</div></td>
                <td>
                    <button
                        on:click={() => {
                            $settings.isAlwaysOnTop = !$settings.isAlwaysOnTop;
                        }}
                    >
                        {$settings.isAlwaysOnTop}
                    </button>
                </td>
            </tr>
            <tr>
                <td><div>windowDecorations</div></td>
                <td>
                    <button
                        on:click={() => {
                            $settings.windowDecorations = !$settings.windowDecorations;
                        }}
                    >
                        {$settings.windowDecorations}
                    </button>
                </td>
            </tr>
            <tr>
                <td><div>startOnBoot</div></td>
                <td>
                    <button
                        on:click={() => {
                            $settings.startOnBoot = !$settings.startOnBoot;
                        }}
                    >
                        {$settings.startOnBoot}
                    </button>
                </td>
            </tr>
            <tr>
                <td><div>theme</div></td>
                <td>
                    <select
                        on:input={(event) => {
                            $settings.theme = event.target.value;
                        }}
                    >
                        <option>{$settings.theme}</option>
                        {#each $settings.themes as theme}
                            {#if theme !== $settings.theme}
                                <option>{theme}</option>
                            {/if}
                        {/each}
                    </select>
                </td>
            </tr>
            <tr>
                <td><div>coloredHeaders</div></td>
                <td>
                    <button
                        on:click={() => {
                            $settings.coloredHeaders = !$settings.coloredHeaders;
                        }}
                    >
                        {$settings.coloredHeaders}
                    </button>
                </td>
            </tr>
            <tr>
                <td><div>fontSize (zoom)</div></td>
                <td>
                    <input
                        type="number"
                        max="36"
                        min="8"
                        placeholder="16"
                        value={$settings.fontSize}
                        on:input={(event) => {
                            $settings.fontSize = event.target.value;
                        }}
                    />
                </td>
            </tr>
            <tr>
                <td><div>ipPollRate_ms</div></td>
                <td>
                    <input
                        type="number"
                        max="10000"
                        min="200"
                        placeholder="500"
                        value={$settings.ipPollRate_ms}
                        on:input={(event) => {
                            $settings.ipPollRate_ms = event.target.value;
                        }}
                    />
                </td>
            </tr>
        </tbody>
    </table>
<!-- </article>

<article class="mono"> -->
    <table style="width: 100%;">
        <thead>
            <tr>
                <th>Shortcut</th>
                <th>Key(s)</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td>
                    <div>Zoom in</div>
                </td>
                <td>
                    <div>
                        <code>Control</code>
                        <code>=</code>
                    </div>
                    <div>
                        <code>Control</code>
                        <code>Scroll Wheel Forward</code>
                    </div>
                </td>
            </tr>
            <tr>
                <td>
                    <div>Zoom out</div>
                </td>
                <td>
                    <div>
                        <code>Control</code>
                        <code>-</code>
                    </div>
                    <div>
                        <code>Control</code>
                        <code>Scroll Wheel Back</code>
                    </div>
                </td>
            </tr>
            <tr>
                <td>
                    <div>Default Zoom & Size</div>
                </td>
                <td>
                    <div>
                        <code>Control</code>
                        <code>Middle Mouse</code>
                    </div>
                </td>
            </tr>
            <tr>
                <td>
                    <div>Toggle isAlwaysOnTop</div>
                </td>
                <td>
                    <div>
                        <code>Control</code>
                        <code>t</code>
                    </div>
                </td>
            </tr>
            <tr>
                <td>
                    <div>Toggle lightMode</div>
                </td>
                <td>
                    <div>
                        <code>Control</code>
                        <code>l</code>
                    </div>
                </td>
            </tr>
        </tbody>
    </table>
</article>

<style>
    article {
        display: grid;
        padding: var(--pad);
        gap: var(--gap);
        max-width: fit-content;
    }
    td {
        padding: 0;
    }
    /* th {
        color: var(--color-text-bright);
    } */
    div,
    button,
    select,
    input {
        background-color: transparent;
        min-width: 6rem;
        text-align: left;
        padding: var(--pad);
        border-radius: 0;
    }
    code {
        font-size: 1rem;
        background-color: var(--color-bg-dim);
        color: var(--color-text-input);
        padding: 0.3rem var(--pad);
        border: var(--border);
    }
</style>
