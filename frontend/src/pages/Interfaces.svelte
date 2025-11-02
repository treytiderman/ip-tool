<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { setPage } from "../ts/router";
    import { nics, setNic, nicsUpdateTrigger, initNics, pollNics } from "../ts/nic";
    import { settingsStore } from "../ts/settings";

    // Page start
    let interval: number;
    onMount(async () => {
        await initNics();
        interval = pollNics($settingsStore.interfacePollRateMs);
    });
    
    // Page close
    onDestroy(async () => {
        clearInterval(interval);
    });
</script>

<div class="grid gap-1 pad-1">
    <div class="flex bottom gap-1">
        <div class="grow text-dark thin">Select Interface</div>
        <button class="transparent text-dark" style="margin-left: auto; visibility: hidden;" title="Change Interface">
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
                <circle cx="12.1" cy="12.1" r="4" />
            </svg>
        </button>
    </div>

    {#each $nics as nic}
        <section class="bg border radius shadow">
            <div class="flex center-y gap-2 pad-2">
                {#if nic.connected}
                    <div
                        class:pulse-out={$nicsUpdateTrigger && $settingsStore.showPollAnimation}
                        title="Interface is connected to a network
Every pulse indicates the interfaces were polled"
                        style="height: 1rem;"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="var(--success-border)"
                            stroke="var(--success)"
                            stroke-width="var(--border-width)"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <circle cx="12.1" cy="12.1" r="6" />
                        </svg>
                    </div>
                {:else}
                    <div
                        class:pulse-out={$nicsUpdateTrigger && $settingsStore.showPollAnimation}
                        class="pulse-out-error"
                        title="Interface is NOT connected to a network
Every pulse indicates the interfaces were polled"
                        style="height: 1rem;"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            style="color: var(--error);"
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="var(--error-border)"
                            stroke="var(--error)"
                            stroke-width="var(--border-width)"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <circle cx="12.1" cy="12.1" r="6" />
                        </svg>
                    </div>
                {/if}
                <button
                    on:click={() => {
                        setNic(nic.interface_name);
                        setPage("IPv4 Presets");
                    }}
                    class="border grow shadow"
                    title="Select Interface"
                    disabled={nic.disabled}
                >
                    {nic.interface_name}</button
                >
                <small class="text-dark thin pad-inline-1">[{nic.interface_metric}]</small>
            </div>
            <div class="pad-inline-2">
                <hr />
            </div>
            <div class="grid pad-2">
                {#each nic.ips as ipMask, index}
                    <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr 2rem;">
                        {#if nic.ip_is_dhcp}
                            <span class="text-dark thin grid right small" title="DHCP is ON">
                                <div class="flex center-y gap-1">
                                    <svg
                                        style="height: 0.75rem;"
                                        xmlns="http://www.w3.org/2000/svg"
                                        width="16"
                                        height="16"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="var(--text-light)"
                                        stroke-width="var(--border-width)"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    >
                                        <path
                                            fill="var(--text-light)"
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
                            <span class="text-dark thin grid right small" title="Static">ip:</span>
                        {/if}
                        <span class="mono" style="color: var(--text-light); font-weight: bold;">
                            {ipMask.ip_address}
                        </span>
                        {#if index === 0}
                            <div class="flex right"></div>
                        {/if}
                    </div>
                    <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr;">
                        <span class="text-dark thin grid right small">mask:</span>
                        <span class="mono">{ipMask.subnet_mask}</span>
                    </div>
                {/each}

                {#each nic.gateways as gateway}
                    <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr;">
                        <span class="text-dark thin grid right small">gate:</span>
                        <span class="mono">{gateway.gateway_address}</span>
                    </div>
                {/each}

                {#each nic.dns_servers as dns_server}
                    <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr;">
                        <span class="text-dark thin grid right small">dns:</span>
                        <span class="mono">{dns_server}</span>
                    </div>
                {/each}
            </div>
        </section>
        <div></div>
    {/each}
</div>

<style>
    .pulse-out {
        z-index: 0;
        position: relative;
    }
    .pulse-out::before {
        content: "";
        z-index: -1;
        position: absolute;
        background-color: var(--success);
        border-radius: 50%;
        width: 1rem;
        height: 1rem;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%) scale(0);
        animation: pulse 450ms ease-out forwards;
    }
    .pulse-out-error::before {
        background-color: var(--error);
    }
    @keyframes pulse {
        0% {
            opacity: 0.8;
            transform: translate(-50%, -50%) scale(0.4);
        }
        100% {
            opacity: 0;
            transform: translate(-50%, -50%) scale(1.4);
        }
    }
</style>
