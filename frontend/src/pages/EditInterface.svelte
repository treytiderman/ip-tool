<script lang="ts">
    import { nics, nicTemp, currentNicIndex, setNicToInterface } from "../ts/nic";
    import { setPage } from "../ts/router";

    function setInterfaceName(event: any) {
        const val = event.target.value;
        $nicTemp.interface_name = val;
        console.log("Temp Interface", $nicTemp);
    }
    function setIpAddress(index: number, event: any) {
        const val = event.target.value;
        $nicTemp.ips[index].ip_address = val;
        console.log("Temp Interface", $nicTemp);
    }
    function setSubnetMask(index: number, event: any) {
        const val = event.target.value;
        $nicTemp.ips[index].subnet_mask = val;
        console.log("Temp Interface", $nicTemp);
    }
    function setGateway(index: number, event: any) {
        const val = event.target.value;
        $nicTemp.gateways[index].gateway_address = val;
        console.log("Temp Interface", $nicTemp);
    }
    function setDnsServer(index: number, event: any) {
        const val = event.target.value;
        $nicTemp.dns_servers[index] = val;
        console.log("Temp Interface", $nicTemp);
    }
</script>

<form
    class="grid gap-4 pad-1"
    on:submit|preventDefault={async () => {
        console.log("submit", $nicTemp);
        setPage("IPv4 Presets");
        await setNicToInterface($nics[$currentNicIndex].interface_name, $nicTemp);
    }}
>
    <div class="grid gap-1">
        <div class="flex bottom gap-1">
            <label class="text-dark thin grow" for="preset-name">Interface Name</label>
            <button class="transparent text-dark invisible" type="button">
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
                    <path d="M5 12h14" />
                </svg>
            </button>
        </div>
        <input
            type="text"
            id="interface-name"
            class="mono border shadow-inset"
            name="interface-name"
            placeholder={$nics[$currentNicIndex].interface_name}
            value={$nicTemp.interface_name}
            on:input={setInterfaceName}
            required
        />
    </div>

    <div class="grid gap-1">
        <div class="flex bottom gap-1">
            <label class="text-dark thin grow" for="ip-address-1">IP Address</label>
            {#if $nicTemp.ips.length > 1}
                <button
                    class="transparent text-dark"
                    type="button"
                    title="Remove IP Address"
                    on:click={() => {
                        $nicTemp.ips.pop();
                        $nicTemp = $nicTemp;
                    }}
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
                        <path d="M5 12h14" />
                    </svg>
                </button>
            {/if}
            <button
                class="transparent text-dark"
                type="button"
                title="Add IP Address"
                on:click={() => {
                    $nicTemp.ips.push({ ip_address: "", subnet_mask: "", cidr: 24 });
                    $nicTemp = $nicTemp;
                }}
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
            placeholder={$nics[$currentNicIndex].ips[0]?.ip_address ?? ""}
            value={$nicTemp.ips[0]?.ip_address ?? ""}
            on:input={(ev) => setIpAddress(0, ev)}
            required
            autofocus
        />
        {#if $nicTemp.ips.length > 1}
            {#each $nicTemp.ips as ip, index}
                {#if index > 0}
                    <div class="grid gap-1">
                        <label class="text-dark thin" for="ip-address-{index}" hidden>IP Address</label>
                        <input
                            type="text"
                            id="ip-address-{index}"
                            class="mono border shadow-inset"
                            name="ip-address-{index}"
                            placeholder={$nics[$currentNicIndex].ips[0]?.ip_address ?? ""}
                            value={$nicTemp.ips[index]?.ip_address ?? ""}
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
            placeholder={$nics[$currentNicIndex].ips[0]?.subnet_mask ?? ""}
            value={$nicTemp.ips[0]?.subnet_mask ?? ""}
            on:input={(ev) => setSubnetMask(0, ev)}
            required
        />
        {#if $nicTemp.ips.length > 1}
            {#each $nicTemp.ips as ip, index}
                {#if index > 0}
                    <div class="grid gap-1">
                        <label class="text-dark thin" for="subnet-mask-{index}" hidden>Subnet Mask</label>
                        <input
                            type="text"
                            id="subnet-mask-{index}"
                            class="mono border shadow-inset"
                            name="subnet-mask-{index}"
                            placeholder={$nics[$currentNicIndex].ips[0]?.subnet_mask ?? ""}
                            value={$nicTemp.ips[index]?.subnet_mask ?? ""}
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
            placeholder={$nics[$currentNicIndex].gateways[0]?.gateway_address ?? ""}
            value={$nicTemp.gateways[0]?.gateway_address ?? ""}
            on:input={(ev) => setGateway(0, ev)}
        />
    </div>

    <div class="grid gap-1">
        <div class="flex bottom gap-1">
            <label class="text-dark thin grow" for="dns">DNS Servers</label>
            {#if $nicTemp.dns_servers.length > 1}
                <button
                    class="transparent text-dark"
                    type="button"
                    title="Remove DNS Server"
                    on:click={() => {
                        $nicTemp.dns_servers.pop();
                        $nicTemp = $nicTemp;
                    }}
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
                        <path d="M5 12h14" />
                    </svg>
                </button>
            {/if}
            <button
                class="transparent text-dark"
                type="button"
                title="Add DNS Server"
                on:click={() => {
                    $nicTemp.dns_servers.push("");
                    $nicTemp = $nicTemp;
                }}
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
            placeholder={$nics[$currentNicIndex].dns_servers[0] ?? ""}
            value={$nicTemp.dns_servers[0] ?? ""}
            on:input={(ev) => setDnsServer(0, ev)}
        />
        {#if $nicTemp.dns_servers.length > 1}
            {#each $nicTemp.dns_servers as dns, index}
                {#if index > 0}
                    <div class="grid gap-1">
                        <label class="text-dark thin" for="dns-{index}" hidden>IP Address</label>
                        <input
                            type="text"
                            id="dns-{index}"
                            class="mono border shadow-inset"
                            name="dns-{index}"
                            placeholder={$nics[$currentNicIndex].dns_servers[0]}
                            value={$nicTemp.dns_servers[index]}
                            on:input={(ev) => setDnsServer(index, ev)}
                        />
                    </div>
                {/if}
            {/each}
        {/if}
    </div>

    <div class="flex even gap-2">
        <button type="submit" class="info info-bg info-border shadow">Set</button>
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
