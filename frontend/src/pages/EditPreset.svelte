<script lang="ts">
    import { presets, presetTemp, presetSelectedIndex, editPreset, removePreset, selectPreset } from "../ts/presets";
    import { setPage } from "../ts/router";

    function setPresetName(event: any) {
        const val = event.target.value;
        $presetTemp.name = val;
        console.log("Temp Preset", $presetTemp);
    }
    function setIpAddress(index: number, event: any) {
        const val = event.target.value;
        $presetTemp.ips[index].ip_address = val;
        console.log("Temp Preset", $presetTemp);
    }
    function setSubnetMask(index: number, event: any) {
        const val = event.target.value;
        $presetTemp.ips[index].subnet_mask = val;
        console.log("Temp Preset", $presetTemp);
    }
    function setGateway(index: number, event: any) {
        const val = event.target.value;
        $presetTemp.gateways[index] = {
            gateway_address: val
        };;
        console.log("Temp Preset", $presetTemp);
    }
    function setDnsServer(index: number, event: any) {
        const val = event.target.value;
        $presetTemp.dns_servers[index] = val;
        console.log("Temp Preset", $presetTemp);
    }
</script>

<form
    class="grid gap-4 pad-1"
    on:submit|preventDefault={async () => {
        console.log("submit", $presetTemp);
        editPreset($presets[presetSelectedIndex].name, $presetTemp);
        setPage("IPv4 Presets");
    }}
>
    <div class="grid gap-1">
        <div class="flex bottom gap-1">
            <label class="text-dark thin grow" for="preset-name">Preset Name</label>
            <button
                class="transparent text-dark error-bg-hover"
                type="button"
                title="Delete Preset"
                on:click={() => {
                    removePreset($presets[presetSelectedIndex].name);
                    setPage("IPv4 Presets");
                }}
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
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6" />
                    <path d="M3 6h18" />
                    <path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                </svg>
            </button>
        </div>
        <input
            type="text"
            class="mono border shadow-inset"
            id="preset-name"
            name="preset-name"
            placeholder={$presets[presetSelectedIndex].name}
            value={$presetTemp.name}
            on:input={setPresetName}
            required
        />
    </div>

    <div class="grid gap-1">
        <div class="flex bottom gap-1">
            <label class="text-dark thin grow" for="ip-address-1">IP Address</label>
            {#if $presetTemp.ips.length > 1}
                <button
                    class="transparent text-dark"
                    type="button"
                    title="Remove IP Address"
                    on:click={() => {
                        $presetTemp.ips.pop();
                        $presetTemp = $presetTemp;
                    }}
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
                        <path d="M5 12h14" />
                    </svg>
                </button>
            {/if}
            <button
                class="transparent text-dark"
                type="button"
                title="Add IP Address"
                on:click={() => {
                    $presetTemp.ips.push({ ip_address: "", subnet_mask: "", cidr: 24 });
                    $presetTemp = $presetTemp;
                }}
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
                    <path d="M5 12h14" />
                    <path d="M12 5v14" />
                </svg>
            </button>
        </div>
        <input
            type="text"
            class="mono border shadow-inset"
            id="ip-address-1"
            name="ip-address-1"
            placeholder={$presets[presetSelectedIndex].ips[0]?.ip_address ?? ""}
            value={$presetTemp.ips[0]?.ip_address ?? ""}
            on:input={(ev) => setIpAddress(0, ev)}
            required
            autofocus
        />
        {#if $presetTemp.ips.length > 1}
            {#each $presetTemp.ips as ip, index}
                {#if index > 0}
                    <div class="grid gap-1">
                        <label class="text-dark thin" for="ip-address-{index}" hidden>IP Address</label>
                        <input
                            type="text"
                            id="ip-address-{index}"
                            class="mono border shadow-inset"
                            name="ip-address-{index}"
                            placeholder={$presets[presetSelectedIndex].ips[0]?.ip_address ?? ""}
                            value={$presetTemp.ips[index]?.ip_address ?? ""}
                            on:input={(ev) => setIpAddress(index, ev)}
                            required
                        />
                    </div>
                {/if}
            {/each}
        {/if}
    </div>

    <div class="grid gap-1">
        <label class="text-dark thin" for="subnet-mask-1">Subnet Mask</label>
        <input
            type="text"
            id="subnet-mask-1"
            class="mono border shadow-inset"
            name="subnet-mask-1"
            placeholder={$presets[presetSelectedIndex].ips[0]?.subnet_mask ?? ""}
            value={$presetTemp.ips[0]?.subnet_mask ?? ""}
            on:input={(ev) => setSubnetMask(0, ev)}
            required
        />
        {#if $presetTemp.ips.length > 1}
            {#each $presetTemp.ips as ip, index}
                {#if index > 0}
                    <div class="grid gap-1">
                        <label class="text-dark thin" for="subnet-mask-{index}" hidden>Subnet Mask</label>
                        <input
                            type="text"
                            id="subnet-mask-{index}"
                            class="mono border shadow-inset"
                            name="subnet-mask-{index}"
                            placeholder={$presets[presetSelectedIndex].ips[0]?.subnet_mask ?? ""}
                            value={$presetTemp.ips[index]?.subnet_mask ?? ""}
                            on:input={(ev) => setSubnetMask(index, ev)}
                            required
                        />
                    </div>
                {/if}
            {/each}
        {/if}
    </div>

    <div class="grid gap-1">
        <label class="text-dark thin" for="gateway">Gateway</label>
        <input
            type="text"
            id="gateway"
            class="mono border shadow-inset"
            name="gateway"
            placeholder={$presets[presetSelectedIndex].gateways[0]?.gateway_address ?? ""}
            value={$presetTemp.gateways[0]?.gateway_address ?? ""}
            on:input={(ev) => setGateway(0, ev)}
        />
    </div>

    <div class="grid gap-1">
        <div class="flex bottom gap-1">
            <label class="text-dark thin grow" for="dns">DNS Servers</label>
            {#if $presetTemp.dns_servers.length > 1}
                <button
                    class="transparent text-dark"
                    type="button"
                    title="Remove DNS Server"
                    on:click={() => {
                        $presetTemp.dns_servers.pop();
                        $presetTemp = $presetTemp;
                    }}
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
                        <path d="M5 12h14" />
                    </svg>
                </button>
            {/if}
            <button
                class="transparent text-dark"
                type="button"
                title="Add DNS Server"
                on:click={() => {
                    $presetTemp.dns_servers.push("");
                    $presetTemp = $presetTemp;
                }}
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
                    <path d="M5 12h14" />
                    <path d="M12 5v14" />
                </svg>
            </button>
        </div>
        <input
            type="text"
            id="dns"
            class="mono border shadow-inset"
            name="dns"
            placeholder={$presets[presetSelectedIndex].dns_servers[0] ?? ""}
            value={$presetTemp.dns_servers[0] ?? ""}
            on:input={(ev) => setDnsServer(0, ev)}
        />
        {#if $presetTemp.dns_servers.length > 1}
            {#each $presetTemp.dns_servers as dns, index}
                {#if index > 0}
                    <div class="grid gap-1">
                        <label class="text-dark thin" for="dns-{index}" hidden>IP Address</label>
                        <input
                            type="text"
                            id="dns-{index}"
                            class="mono border shadow-inset"
                            name="dns-{index}"
                            placeholder={$presets[presetSelectedIndex].dns_servers[0]}
                            value={$presetTemp.dns_servers[index]}
                            on:input={(ev) => setDnsServer(index, ev)}
                        />
                    </div>
                {/if}
            {/each}
        {/if}
    </div>

    <div class="flex even gap-2">
        <button type="submit" class="warning warning-bg warning-border shadow">Update</button>
        <button
            type="button"
            class="transparent text-dark"
            on:click={() => setPage("IPv4 Presets")}
            title="Back to IPv4 Presets"
        >
            Cancel
        </button>
    </div>
</form>

<style>
</style>
