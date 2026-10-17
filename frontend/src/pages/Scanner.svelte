<script lang="ts">
    import { onDestroy, onMount, tick } from "svelte";
    import { main } from "../../wailsjs/go/models";
    import * as app from "../../wailsjs/go/main/App.js";
    import { EventsOn } from "../../wailsjs/runtime/runtime.js";
    import { setPage } from "../ts/router";
    import { currentNicIndex, initNics, nics, nicsUpdateTrigger } from "../ts/nic";
    import { settingsStore } from "../ts/settings";

    export let scanMode: "ping" | "mdns" = "ping";

    let scanning = false;
    let scanComplete = false;
    let scanCancelled = false;
    let cancelRequested = false;
    let loadError = "";
    let scanError = "";
    let results: main.ScanResult[] = [];
    let mdnsResults: main.MDNSResult[] = [];
    let scanIPAddress = "";
    let scanSubnetMask = "";
    let scannedHosts = 0;
    let totalHosts = 0;
    let mdnsDevicesFound = 0;
    let unsubscribeScanProgress: (() => void) | undefined;
    let unsubscribeMDNSProgress: (() => void) | undefined;

    $: selectedInterface = $nics[$currentNicIndex];
    $: scanAddress = selectedInterface?.ips?.find(
        (ip) => isIPv4Address(ip.ip_address) && getMaskPrefix(ip.subnet_mask) !== null,
    );
    $: maskPrefix = getMaskPrefix(scanSubnetMask);
    $: subnet = maskPrefix !== null && isIPv4Address(scanIPAddress)
        ? `${getNetworkAddress(scanIPAddress, scanSubnetMask)}/${maskPrefix}`
        : "";
    $: canScan = !!selectedInterface && (
        scanMode === "mdns" ||
        (isIPv4Address(scanIPAddress) && maskPrefix !== null && maskPrefix >= 16)
    );
    onDestroy(() => {
        unsubscribeScanProgress?.();
        unsubscribeMDNSProgress?.();
    });

    function isIPv4Address(address: string) {
        const octets = address.split(".");
        return octets.length === 4 && octets.every(
            (octet) =>
                /^\d{1,3}$/.test(octet) &&
                Number(octet) <= 255 &&
                String(Number(octet)) === octet,
        );
    }

    function getMaskPrefix(mask: string): number | null {
        if (!isIPv4Address(mask)) return null;

        const bits = mask
            .split(".")
            .map((octet) => Number(octet).toString(2).padStart(8, "0"))
            .join("");
        if (!/^1*0*$/.test(bits)) return null;

        const prefix = bits.indexOf("0") === -1 ? 32 : bits.indexOf("0");
        return prefix > 0 ? prefix : null;
    }

    function getNetworkAddress(address: string, mask: string) {
        const maskOctets = mask.split(".");
        return address
            .split(".")
            .map((octet, index) => String(Number(octet) & Number(maskOctets[index])))
            .join(".");
    }

    onMount(async () => {
        unsubscribeScanProgress = EventsOn("lan-scan-progress", (progress: { scanned: number; total: number }) => {
            scannedHosts = progress.scanned;
            totalHosts = progress.total;
        });
        unsubscribeMDNSProgress = EventsOn("lan-mdns-progress", (found: number) => {
            mdnsDevicesFound = found;
        });
        try {
            await initNics();
            await tick();
            scanIPAddress = scanAddress?.ip_address || "";
            scanSubnetMask = scanAddress?.subnet_mask || "";
        } catch (error) {
            loadError = error instanceof Error ? error.message : String(error);
        }
    });

    async function scan() {
        if (!selectedInterface || !canScan || scanning) return;

        scanning = true;
        scanComplete = false;
        scanCancelled = false;
        cancelRequested = false;
        scanError = "";
        results = [];
        mdnsResults = [];
        scannedHosts = 0;
        totalHosts = 0;
        mdnsDevicesFound = 0;
        try {
            if (scanMode === "ping") {
                results = await app.ScanSubnet(
                    selectedInterface.interface_name,
                    scanIPAddress,
                    scanSubnetMask,
                );
            } else {
                mdnsResults = await app.ScanMDNS(selectedInterface.interface_name);
            }
            scanComplete = true;
        } catch (error) {
            if (cancelRequested) {
                scanCancelled = true;
            } else {
                scanError = error instanceof Error ? error.message : String(error);
            }
        } finally {
            scanning = false;
            cancelRequested = false;
        }
    }

    async function cancelScan() {
        if (!scanning || cancelRequested) return;
        cancelRequested = true;
        try {
            const cancelled = await app.CancelScan();
            if (!cancelled) {
                cancelRequested = false;
                scanError = "The scan could not be cancelled.";
            }
        } catch (error) {
            cancelRequested = false;
            scanError = error instanceof Error ? error.message : String(error);
        }
    }

</script>

<div class="grow flex column gap-1 pad-1 overflow-hidden">

    <div class="flex bottom">
        <div class="grow text-dark thin">Interface</div>
        <button
            class="transparent text-dark thin"
            on:click={() => setPage("Ping Scanner")}
            title="Ping Scanner"
        >
            Ping
        </button>
        <button
            class="transparent text-dark thin "
            on:click={() => setPage("mDNS Scanner")}
            title="mDNS Scanner"
        >
            mDNS
        </button>
        <button class="transparent text-dark thin" title="Change Interface" on:click={() => setPage("Interfaces")}>
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
                <path d="M12 12h.01" />
                <path d="M16 12h.01" />
                <path d="m17 7 5 5-5 5" />
                <path d="m7 7-5 5 5 5" />
                <path d="M8 12h.01" />
            </svg>
        </button>
    </div>

    <!-- <div></div> -->

    <section class="bg border radius shadow">
        <div class="flex center-y gap-2 pad-2">
            {#if $nics[$currentNicIndex].connected}
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
            <div class="grow" title="Interface Name">{$nics[$currentNicIndex].interface_name}</div>
            <small class="text-dark thin pad-inline-1" title="Interface Metric">
                [{$nics[$currentNicIndex].interface_metric}]
            </small>
        </div>
        <div class="pad-inline-2">
            <hr />
        </div>
        <div class="flex wrap top pad-2">
            <div class="grid grow">
                {#each $nics[$currentNicIndex].ips as ipMask, index}
                    <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem minmax(0, 1fr);">
                        {#if $nics[$currentNicIndex].ip_is_dhcp}
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
                        {#if scanMode === "ping" && index === 0 && scanAddress}
                            <div style="padding-bottom: var(--gap);">
                                <input
                                    class="mono"
                                    style="width: 100%; min-width: 0;"
                                    aria-label="Scan IPv4 address"
                                    inputmode="decimal"
                                    placeholder="IPv4 address"
                                    bind:value={scanIPAddress}
                                />
                            </div>
                        {:else}
                            <span class="mono" style="color: var(--text-light); font-weight: bold;">
                                {ipMask.ip_address}
                            </span>
                        {/if}
                    </div>
                    <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr;">
                        <span class="text-dark thin grid right small">mask:</span>
                        {#if scanMode === "ping" && index === 0 && scanAddress}
                            <div>
                                <input
                                    class="mono"
                                    style="width: 100%; min-width: 0;"
                                    aria-label="Scan IPv4 subnet mask"
                                    inputmode="decimal"
                                    placeholder="Subnet mask"
                                    bind:value={scanSubnetMask}
                                />
                            </div>
                        {:else}
                            <span class="mono">{ipMask.subnet_mask}</span>
                        {/if}
                    </div>
                {/each}

                <!-- {#each $nics[$currentNicIndex].gateways as gateway}
                    <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr;">
                        <span class="text-dark thin grid right small">gate:</span>
                        <span class="mono">{gateway.gateway_address}</span>
                    </div>
                {/each}

                {#each $nics[$currentNicIndex].dns_servers as dns_server}
                    <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr;">
                        <span class="text-dark thin grid right small">dns:</span>
                        <span class="mono">{dns_server}</span>
                    </div>
                {/each} -->

                {#if loadError}
                    <div class="error error-bg error-border radius pad-1" role="alert">
                        Could not load interfaces: {loadError}
                    </div>
                {/if}
                {#if scanError}
                    <div class="error error-bg error-border radius pad-1" role="alert">
                        Scan failed: {scanError}
                    </div>
                {/if}
            </div>
        </div>
    </section>
    
    <div></div>
    
    <button
        class:error={scanning}
        class:error-bg={scanning}
        class:error-border={scanning}
        class="accent accent-bg accent-border"
        on:click={scanning ? cancelScan : scan}
        disabled={(!scanning && !canScan) || cancelRequested}
    >
        {#if scanning}
            {#if cancelRequested}
                Cancelling...
            {:else if scanMode === "ping"}
                Cancel scan... {scannedHosts.toLocaleString()} / {totalHosts.toLocaleString()} hosts
            {:else}
                Cancel mDNS listen... {mdnsDevicesFound.toLocaleString()} found
            {/if}
        {:else if scanComplete}
            {#if scanCancelled}
                Scan cancelled
            {:else if scanMode === "ping"}
                Scanned {results.length} {results.length === 1 ? "device" : "devices"}
            {:else}
                Found {mdnsResults.length} mDNS {mdnsResults.length === 1 ? "device" : "devices"}
            {/if}
        {:else if scanCancelled}
            Scan cancelled
        {:else if scanMode === "ping"}
            Scan {subnet}
        {:else}
            Listen for mDNS
        {/if}
    </button>
    {#if scanMode === "ping" && scanIPAddress && (!isIPv4Address(scanIPAddress) || maskPrefix === null || maskPrefix < 16)}
        <div class="warning small">
            Enter a valid IPv4 address and contiguous subnet mask. The scan is limited to 65,536 addresses (/16 or narrower).
        </div>
    {/if}
    
    <div></div>

    <section class="grow overflow-y bg border radius shadow pad-2 flex column gap-1" aria-label="Scan results">
        {#if scanMode === "ping" && results.length > 0}

            {#each results as result, index}
                {#if index > 0}
                    <div></div>
                    <hr>
                    <div></div>
                {/if}
                <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr auto;">
                    <div class="flex">
                        <span class="grow small">{index+1}</span>
                        <span class="text-dark thin small">ip:</span>
                    </div>
                    {#if result.is_local}
                        <div class="mono success grow" title="This device">{result.ip_address}</div>
                    {:else}
                        <div class="mono text-light grow">{result.ip_address}</div>
                    {/if}
                    <span class="text-dark small">{result.latency_ms} ms</span>
                </div>
                <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr auto;">
                    <span class="text-dark thin grid right small">mac:</span>
                    <div class="mono grow">{result.mac_address || "Not available"}</div>
                </div>
                <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr auto;">
                    <span class="text-dark thin grid right small">make:</span>
                    <div
                        class="mono grow"
                        title="{result.vendor || "Unknown"}"
                        style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis;"
                    >
                        {result.vendor || "Unknown"}
                    </div>
                </div>
            {/each}

        {:else if scanMode === "mdns" && mdnsResults.length > 0}
            {#each mdnsResults as result, index}
                {#if index > 0}
                    <div></div>
                    <hr>
                    <div></div>
                {/if}
                <details class="mdns-details">
                    <summary class="mdns-summary">
                        <div>
                            <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr auto;">
                                <span class="text-dark thin grid right small">host:</span>
                                <div class="mono text-light grow">{result.hostname}</div>
                            </div>
                            <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr auto;">
                                <span class="text-dark thin grid right small">ip:</span>
                                <div class="mono grow">
                                    {(result.ip_addresses || []).join(", ") || "Not available"}
                                </div>
                            </div>
                        </div>
                    </summary>
                    {#each (result.services || []) as service}
                        <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr auto;">
                            <span class="text-dark thin grid right small">serv:</span>
                            <div class="mono grow">{service}</div>
                        </div>
                    {/each}
                    {#each (result.details || []) as detail}
                        <div class="grid center-y gap-2" style="grid-template-columns: 2.6rem 1fr auto;">
                            <span class="text-dark thin grid right small">info:</span>
                            <div
                                class="mono grow"
                                title={detail}
                                style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis;"
                            >
                                {detail}
                            </div>
                        </div>
                    {/each}
                </details>
            {/each}
        {:else if scanComplete}
            <div class="pad-1 text-dark small">
                {#if scanMode === "ping"}
                    No hosts responded to ping.
                {:else}
                    No mDNS services were discovered.
                {/if}
            </div>
        {:else if scanCancelled}
            <div class="pad-1 text-dark small">Scan cancelled.</div>
        {:else if !scanning}
            <div class="pad-1 text-dark small">
                {#if scanMode === "ping"}
                    Run a scan to discover devices...
                {:else}
                    Listen for mDNS service announcements on the selected interface...
                {/if}
            </div>
        {/if}

    </section>
</div>

<style>
    .mdns-summary {
        position: relative;
        padding-right: 1.25rem;
    }

    .mdns-summary:before {
        content: "";
        position: absolute;
        top: 0.45rem;
        right: 0.25rem;
        width: 0.45rem;
        height: 0.45rem;
        padding: 0;
        border-right: 2px solid var(--text-dark);
        border-bottom: 2px solid var(--text-dark);
        transform: rotate(45deg);
    }

    .mdns-details[open] .mdns-summary:before {
        content: "";
        top: 0.65rem;
        transform: rotate(225deg);
    }
</style>
