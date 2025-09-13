<script lang="ts">
    import { setPage } from "../ts/router";
    import { nics, setNic, currentNicIndex } from "../ts/nic";
    import { presets, setPresetToInterface, selectPreset } from "../ts/presets";
    import * as app from "../../wailsjs/go/main/App.js";
</script>

<div class="flex bottom gap-xs pad-inline-sm">
    <label>Interface</label>
    <button
        class="ip-icon-button color-dim"
        style="margin-left: auto; z color: var(--color-text-blue)"
        title="Change Interface"
        on:click={() => setPage("Interfaces")}
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
            <path d="M12 12h.01" />
            <path d="M16 12h.01" />
            <path d="m17 7 5 5-5 5" />
            <path d="m7 7-5 5 5 5" />
            <path d="M8 12h.01" />
        </svg>
    </button>
</div>

<section class="radius shadow" style="margin-inline: var(--gap-sm); background-color: var(--color-bg-1);">
    <div class="flex center-y gap-sm pad-sm">
        {#if $nics[$currentNicIndex].connected}
            <div title="Interface is connected to a network" style="height: 1rem;">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    style="color: var(--color-text-green);"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <circle cx="12.1" cy="12.1" r="4" />
                </svg>
            </div>
        {:else}
            <div title="Interface is NOT connected to a network" style="height: 1rem;">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    style="color: var(--color-text-red);"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <circle cx="12.1" cy="12.1" r="4" />
                </svg>
            </div>
        {/if}
        <div class="grow" title="Interface Name">{$nics[$currentNicIndex].interface_name}</div>
        <small class="color-dim" title="Interface Metric">[{$nics[$currentNicIndex].interface_metric}]</small>
    </div>
    <div class="pad-sm" style="padding-block: 0;">
        <hr />
    </div>
    <div class="flex wrap top">
        <div class="grid pad-sm grow">
            {#each $nics[$currentNicIndex].ips as ipMask, index}
                <div class="grid center-y gap-sm" style="grid-template-columns: 2.6rem 1fr 2rem;">
                    {#if $nics[$currentNicIndex].ip_is_dhcp}
                        <span class="color-dim grid right" style="font-size: 0.875rem;" title="DHCP is ON">
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
                                <span>ip:</span>
                            </div>
                        </span>
                    {:else}
                        <span class="color-dim grid right" style="font-size: 0.875rem;" title="Static">ip:</span>
                    {/if}
                    <span class="mono" style="color: var(--color-text-heading); font-weight: bold;">
                        {ipMask.ip_address}
                    </span>
                    {#if index === 0}
                        <div class="flex right"></div>
                    {/if}
                </div>
                <div class="grid center-y gap-sm" style="grid-template-columns: 2.6rem 1fr;">
                    <span class="color-dim grid right" style="font-size: 0.875rem;">mask:</span>
                    <span class="mono">{ipMask.subnet_mask}</span>
                </div>
            {/each}

            {#each $nics[$currentNicIndex].gateways as gateway}
                <div class="grid center-y gap-sm" style="grid-template-columns: 2.6rem 1fr;">
                    <span class="color-dim grid right" style="font-size: 0.875rem;">gate:</span>
                    <span class="mono">{gateway.gateway_address}</span>
                </div>
            {/each}

            {#each $nics[$currentNicIndex].dns_servers as dns_server}
                <div class="grid center-y gap-sm" style="grid-template-columns: 2.6rem 1fr;">
                    <span class="color-dim grid right" style="font-size: 0.875rem;">dns:</span>
                    <span class="mono">{dns_server}</span>
                </div>
            {/each}
        </div>
        <div class="pad-xs">
            <button class="ip-icon-button color-dim" title="Edit Interface" on:click={() => setPage("Edit Interface")}>
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
    </div>
</section>
