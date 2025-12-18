<script lang="ts">
    import * as app from "../../wailsjs/go/main/App.js";
    import { setPage } from "../ts/router";
    import { nics, setNic, currentNicIndex } from "../ts/nic";
    import { presets, setPresetToInterface, selectPreset } from "../ts/presets";
</script>

<div class="flex bottom">
    <div class="grow text-dark thin">Presets</div>
    <button class="transparent text-dark" title="Add New Preset" on:click={() => setPage("Create Preset")}>
        <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="var(--border-width)"
            stroke-linecap="round"
            stroke-linejoin="round"
        >
            <path d="M5 12h14" />
            <path d="M12 5v14" />
        </svg>
    </button>
</div>
<div
    class="bg border radius shadow pad-1 flex column gap-1 grow overflow"
    hidden={$nics[$currentNicIndex].ips[0]?.ip_address === ""}
>
    <div class="flex center-y gap-1">
        <button
            class="transparent"
            on:click={async () => await app.SetDhcp($nics[$currentNicIndex].interface_name)}
            title="Set Interface to DHCP"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="var(--warning-border)"
                stroke="var(--warning)"
                stroke-width="var(--border-width)"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <polygon points="6 3 20 12 6 21 6 3" />
            </svg>
        </button>
        <div class="grow">DHCP</div>
        {#if $nics[$currentNicIndex].ip_is_dhcp}
            <button
                class="transparent text-dark"
                on:click={async () => {
                    await app.ReleaseDhcp();
                }}
                title="Release DHCP
WARNING: affects all interfaces"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="var(--border-width)"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <circle cx="12" cy="12" r="10" />
                    <line x1="9" x2="15" y1="15" y2="9" />
                </svg>
            </button>
            <button
                class="transparent text-dark"
                on:click={async () => {
                    await app.RenewDhcp();
                }}
                title="Renew DHCP
WARNING: affects all interfaces"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="var(--border-width)"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
                    <path d="M3 3v5h5" />
                    <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16" />
                    <path d="M16 16h5v5" />
                </svg>
            </button>
        {/if}
    </div>

    {#each $presets as preset}
        <hr>
        <div class="flex center-y gap-1">
            <button
                class="transparent"
                on:click={() => {
                    setPresetToInterface($nics[$currentNicIndex].interface_name, preset.name);
                }}
                title="Set Interface
Preset: {preset.name}
IP: {preset.ips[0].ip_address}
Mask: {preset.ips[0].subnet_mask}
Gate: {preset.gateways[0]?.gateway_address}
DNS: {preset.dns_servers[0]}
"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="var(--warning-border)"
                    stroke="var(--warning)"
                    stroke-width="var(--border-width)"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <polygon points="6 3 20 12 6 21 6 3" />
                </svg>
            </button>
            <div class="grow">{preset.name}</div>
            <button
                class="transparent text-dark"
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
                    stroke-width="var(--border-width)"
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
