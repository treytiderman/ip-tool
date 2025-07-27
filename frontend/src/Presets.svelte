<script lang="ts">
    import { onMount } from "svelte";
    import { setPage } from "./router";
    import { nic, nics, setNic, initNics, pollNics } from "./nic";
    import {
        presets,
        savePresets,
        loadPresets,
        resetPresets,
        addPreset,
        removePreset,
        setPresetToInterface,
        selectPreset,
    } from "./presets";
    import * as app from "../wailsjs/go/main/App.js";

    function setNicFromEvent(event: any) {
        const val = event.target.value;
        setNic(val);
    }

    onMount(async () => {
        loadPresets();
        await initNics();
        const interval = pollNics();
    });
</script>

<div class="flex column gap-sm pad-sm h-full" style="padding-top: 1px;">
    <div class="ip-nic flex column shadow">
        <select
            value={$nic.interface_name}
            id="networkInterface"
            name="networkInterface"
            title="Select Network Interface"
            on:change={setNicFromEvent}
        >
            <button>
                <selectedcontent></selectedcontent>
            </button>

            {#each $nics as nic, index}
                <option
                    selected
                    value={nic.interface_name}
                    title={`${nic.interface_name} (${nic.interface_metric})`}
                    hidden={nic.interface_name.includes("Loop")}
                >
                    <div>
                        <div class="flex center-y gap-sm">
                            <div>{nic.interface_name}</div>
                            <div class="color-dim" style="font-size: 0.6rem;">({nic.interface_metric})</div>
                        </div>
                        <div class="mono">{nic.ips[0]?.ip_address}</div>
                    </div>
                </option>
            {/each}
        </select>

        <div class="grid pad-sm">
            {#each $nic.ips as ipMask, index}
                <div class="grid center-y gap-sm" style="grid-template-columns: 2.75rem 1fr 2rem;">
                    {#if $nic.ip_is_dhcp}
                        <span class="color-dim grid right" style="font-size: 0.875rem;" title="DHCP">
                            <div class="flex center-y gap-xs">
                                <svg
                                    style="height: 0.75rem;"
                                    xmlns="http://www.w3.org/2000/svg"
                                    width="16"
                                    height="16"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="var(--color-text-heading)"
                                    stroke-width="2"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                >
                                    <path
                                        fill="var(--color-text-heading)"
                                        d="m21.64 3.64-1.28-1.28a1.21 1.21 0 0 0-1.72 0L2.36 18.64a1.21 1.21 0 0 0 0 1.72l1.28 1.28a1.2 1.2 0 0 0 1.72 0L21.64 5.36a1.2 1.2 0 0 0 0-1.72"
                                    />
                                    <path d="m14 7 3 3" />
                                    <path d="M5 6v4" />
                                    <path d="M19 14v4" />
                                    <path d="M10 2v2" />
                                    <path d="M7 8H3" />
                                    <path d="M21 16h-4" />
                                    <path d="M11 3H9" />
                                </svg>
                                <span>IP:</span>
                            </div>
                        </span>
                    {:else}
                        <span class="color-dim grid right" style="font-size: 0.875rem;" title="Static">IP:</span>
                    {/if}
                    <span class="mono" style="color: var(--color-text-heading); font-weight: bold;"
                        >{ipMask.ip_address}</span
                    >
                    {#if index === 0}
                        <div class="flex right">
                            <button
                                class="ip-icon-button"
                                on:click={() => setPage("Edit Interface")}
                                title="Edit Interface"
                            >
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
                                    <path
                                        d="M21.174 6.812a1 1 0 0 0-3.986-3.987L3.842 16.174a2 2 0 0 0-.5.83l-1.321 4.352a.5.5 0 0 0 .623.622l4.353-1.32a2 2 0 0 0 .83-.497z"
                                    />
                                    <path d="m15 5 4 4" />
                                </svg>
                            </button>
                        </div>
                    {/if}
                </div>
                <div class="grid center-y gap-sm" style="grid-template-columns: 2.75rem 1fr;">
                    <span class="color-dim grid right" style="font-size: 0.875rem;">Mask:</span>
                    <span class="mono">{ipMask.subnet_mask}</span>
                </div>
            {/each}

            {#each $nic.gateways as gateway}
                <div class="grid center-y gap-sm" style="grid-template-columns: 2.75rem 1fr;">
                    <span class="color-dim grid right" style="font-size: 0.875rem;">Gate:</span>
                    <span class="mono">{gateway.gateway_address}</span>
                </div>
            {/each}

            {#each $nic.dns_servers as dns_server}
                <div class="grid center-y gap-sm" style="grid-template-columns: 2.75rem 1fr;">
                    <span class="color-dim grid right" style="font-size: 0.875rem;">DNS:</span>
                    <span class="mono">{dns_server}</span>
                </div>
            {/each}
        </div>
    </div>

    <div class="presets grid overflow grow radius shadow" style="background-color: var(--color-bg-1);">
        <div class="flex center-y gap-sm pad-sm">
            <button
                class="ip-icon-button shadow"
                style="background-color: var(--color-bg-3);"
                on:click={async () => await app.SetDhcp($nic.interface_name)}
                title="Set Interface to DHCP"
            >
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
                    <polygon points="6 3 20 12 6 21 6 3" />
                </svg>
            </button>
            <div class="grow">
                <span>DHCP</span>
            </div>
        </div>

        {#each $presets as preset}
            <div class="flex center-y gap-sm pad-sm" style="border-top: var(--border);">
                <button
                    class="ip-icon-button shadow"
                    style="background-color: var(--color-bg-3);"
                    on:click={() => {
                        setPresetToInterface($nic.interface_name, preset.name);
                    }}
                    title="Set Interface
Preset: {preset.name}
IP: {preset.ips[0].ip_address}
Mask: {preset.ips[0].subnet_mask}
Gate: {preset.gateways[0].gateway_address}
DNS: {preset.dns_servers[0]}
"
                >
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
                        <polygon points="6 3 20 12 6 21 6 3" />
                    </svg>
                </button>
                <div class="grow">{preset.name}</div>
                <button
                    class="ip-icon-button"
                    on:click={async () => {
                        setPage("Edit Preset");
                        selectPreset(preset.name);
                    }}
                    title="Edit Preset"
                >
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
                        <path
                            d="M21.174 6.812a1 1 0 0 0-3.986-3.987L3.842 16.174a2 2 0 0 0-.5.83l-1.321 4.352a.5.5 0 0 0 .623.622l4.353-1.32a2 2 0 0 0 .83-.497z"
                        />
                        <path d="m15 5 4 4" />
                    </svg>
                </button>
            </div>
        {/each}
    </div>
    <div>
        <button
            class="ip-icon-button color-dim m-auto"
            style="font-size: 0.875rem;"
            on:click={() => setPage("Create Preset")}
            title="Add Preset"
        >
            <div>Add Preset</div>
        </button>
    </div>
</div>

<style>
    .ip-nic {
        border-radius: var(--border-radius);
        background-color: var(--color-bg-1);
    }

    .ip-nic select,
    .ip-nic ::picker(select) {
        appearance: base-select;
    }

    .ip-nic select {
        width: 100%;
        border-radius: var(--border-radius) var(--border-radius) 0 0;
    }

    @media (hover: hover) {
        .ip-nic select:hover:not(:disabled) {
            color: var(--color-text-button);
            background-color: var(--color-bg-hover);
        }
    }

    .ip-nic ::picker(select) {
        border: none;
        border-radius: var(--border-radius);
        box-shadow: 4px 8px 1px 1000px #000000cc;
    }

    .ip-nic select::picker-icon {
        color: var(--color-text);
    }

    .ip-nic select:open::picker-icon {
        rotate: 180deg;
    }

    .ip-nic option {
        border: none;
        outline: none;
        color: var(--color-text);
        background-color: var(--color-bg-1);
        padding: var(--gap-sm);
    }

    .ip-nic option div:nth-child(1) {
        white-space: nowrap;
        max-width: 75vw;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    selectedcontent div div:nth-child(1) {
        white-space: nowrap;
        max-width: 75vw !important;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    selectedcontent div div.mono:nth-child(2) {
        display: none !important;
        background-color: red;
    }

    .ip-nic option:hover,
    .ip-nic option:focus-visible {
        background: var(--color-bg-hover, #777);
    }

    .ip-nic option:focus-visible {
        outline-offset: -2px;
        outline: var(--outline, 2px solid #fff);
    }

    .ip-nic option:checked {
        font-weight: bold;
    }

    .ip-nic option:not(option:first-of-type) {
        border-top: var(--border, 1px solid #444);
    }

    .ip-nic option:nth-of-type(2) {
        border-radius: var(--border-radius) var(--border-radius) 0 0;
    }

    .ip-nic option:last-of-type {
        border-radius: 0 0 var(--border-radius) var(--border-radius);
    }

    .ip-nic option:not(option:last-of-type) {
        border-bottom: none;
    }

    .ip-nic option:nth-child(even) {
        filter: brightness(1.15);
    }

    .ip-nic option::checkmark {
        width: 1.25rem;
        height: 1.25rem;
        display: grid;
        place-items: center;
    }
</style>
