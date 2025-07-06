<script lang="ts">
    import { setPage } from "./router";
    let selectedInfterface: string = "Ethernet";

    const interfaces = [
        {
            interfaceName: "Ethernet",
            interfaceMetric: 25,
            ipAndMasks: [
                {
                    ipAddress: "192.168.1.9",
                    subnetMask: "255.255.255.0",
                },
            ],
            gateway: "192.168.1.254",
            gatewayMetric: 1,
            dnsServers: ["192.168.1.1"],
        },
        {
            interfaceName: "WiFi",
            interfaceMetric: 50,
            ipAndMasks: [
                {
                    ipAddress: "10.10.1.9",
                    subnetMask: "255.255.0.0",
                },
                {
                    ipAddress: "10.10.2.9",
                    subnetMask: "255.255.0.0",
                },
                {
                    ipAddress: "10.10.3.9",
                    subnetMask: "255.255.0.0",
                },
            ],
            gateway: "10.10.1.1",
            gatewayMetric: 1,
            dnsServers: ["1.1.1.1", "8.8.8.8"],
        },
        {
            interfaceName: "Hyper-V Virtual Ethernet Adapter",
            interfaceMetric: 125,
            ipAndMasks: [
                {
                    ipAddress: "10.10.1.9",
                    subnetMask: "255.255.0.0",
                },
            ],
            gateway: "10.10.1.1",
            gatewayMetric: 1,
            dnsServers: ["10.10.1.1"],
        },
    ];
</script>

<div class="flex column gap-sm h-full">
    <!-- Network Interface -->
    <div class="ip-nic flex column gap-sm">
        <select bind:value={selectedInfterface} id="presets" name="presets" title="Select Network Interface">
            <button>
                <selectedcontent></selectedcontent>
            </button>
            <option value="" hidden>
                <span>Select an Interface</span>
            </option>

            {#each interfaces as nic}
                <option value={nic.interfaceName} title={`${nic.interfaceName} (${nic.interfaceMetric})`}>
                    <div>
                        <span>{nic.interfaceName} <span>({nic.interfaceMetric})</span></span>
                        <span class="mono" style="font-size: 0.875rem;">192.168.1.99</span>
                    </div>
                </option>
            {/each}
        </select>

        <div class="grid">
            <div class="grid center-y gap-sm" style="grid-template-columns: 3rem 1fr;">
                <div class="flex right">
                    <button class="ip-icon-button" on:click={() => setPage("Edit Interface")} title="Edit Interface">
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
                <span class="mono" style="color: var(--color-text-heading); font-weight: bold;">192.168.1.99</span>
            </div>
            <div class="grid center-y gap-sm" style="grid-template-columns: 3rem 1fr;">
                <span class="color-dim grid right" style="font-size: 0.875rem;">Mask:</span>
                <span class="mono">255.255.255.0</span>
            </div>
            <div class="grid center-y gap-sm" style="grid-template-columns: 3rem 1fr;">
                <span class="color-dim grid right" style="font-size: 0.875rem;">Gate:</span>
                <span class="mono">192.168.1.1</span>
            </div>
            <div class="grid center-y" style="grid-template-columns: 3rem 1fr; column-gap: var(--gap-sm);">
                <span class="color-dim grid right" style="font-size: 0.875rem;">DNS:</span>
                <span class="mono">1.1.1.1</span>
                <!-- <span class="color-dim grid right" style="font-size: 0.875rem;">DNS:</span>
                <span class="mono">8.8.8.8</span> -->
            </div>
        </div>
    </div>

    <div class="grid gap-sm overflow grow" style="padding-right: 0.25rem;">
        <div
            class="flex center-y gap-sm"
            style="background-color: var(--color-bg-1); padding: var(--gap-sm); padding-left: var(--pad-x); border-radius: var(--border-radius); box-shadow: 2px 2px 2px -1px #00000077;"
        >
            <button
                class="ip-icon-button"
                style="background-color: var(--color-bg-3);"
                on:click={() => {}}
                title="Set Preset"
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

        <div
            class="flex center-y gap-sm"
            style="background-color: var(--color-bg-1); padding: var(--gap-sm); padding-left: var(--pad-x); border-radius: var(--border-radius); box-shadow: 2px 2px 2px -1px #00000077;"
        >
            <button
                class="ip-icon-button"
                style="background-color: var(--color-bg-3);"
                on:click={() => {}}
                title="Set Preset"
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
                <span>192.168.1.99/24</span>
            </div>
            <button class="ip-icon-button" on:click={() => setPage("Edit Preset")} title="Edit Preset">
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

        <div
            class="flex center-y gap-sm"
            style="background-color: var(--color-bg-1); padding: var(--gap-sm); padding-left: var(--pad-x); border-radius: var(--border-radius);"
        >
            <button
                class="ip-icon-button"
                style="background-color: var(--color-bg-3);"
                on:click={() => {}}
                title="Set Preset"
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
                <span>169.254.0.9/16</span>
            </div>
            <button class="ip-icon-button" on:click={() => {}} title="Edit Preset">
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

        <div
            class="flex center-y gap-sm"
            style="background-color: var(--color-bg-1); padding: var(--gap-sm); padding-left: var(--pad-x); border-radius: var(--border-radius);"
        >
            <button
                class="ip-icon-button"
                style="background-color: var(--color-bg-3);"
                on:click={() => {}}
                title="Set Preset"
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
                <span>10.0.0.7/8</span>
            </div>
            <button class="ip-icon-button" on:click={() => {}} title="Edit Preset">
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
        padding: 0 0 var(--gap-sm) 0;
        border-radius: var(--border-radius);
        background-color: var(--color-bg-1);
        box-shadow: 2px 2px 2px -1px #00000077;
    }

    /* Select & Option */

    .ip-nic select,
    .ip-nic ::picker(select) {
        appearance: base-select;
    }

    .ip-nic select {
        width: 100%;
        max-width: 100%;
        border-radius: var(--border-radius) var(--border-radius) 0 0;
    }

    @media (hover: hover) {
        .ip-nic select:hover:not(:disabled) {
            color: var(--color-text-button, #fff);
            background-color: var(--color-bg-hover, #333);
        }
    }

    .ip-nic ::picker(select) {
        border: none;
        box-shadow: 4px 8px 1px 1000px #000000cc;
    }

    .ip-nic select::picker-icon {
        color: var(--color-text, #ccc);
    }

    .ip-nic select:open::picker-icon {
        rotate: 180deg;
    }

    .ip-nic option {
        border: none;
        outline: none;
        padding: var(--pad-y, 0.5rem) var(--pad-x, 0.8rem);
        color: var(--color-text, #fff);
        background-color: var(--color-bg-1, #444);
        display: flex;
        flex-wrap: wrap;
        gap: var(--gap-sm);
    }

    .ip-nic option div {
        display: flex;
        flex-direction: column;
    }

    .ip-nic span:nth-child(1) {
        max-width: 65vw;
        text-overflow: ellipsis;
        overflow: hidden;

        display: flex;
        align-items: center;
        gap: var(--gap-sm);
    }

    .ip-nic span:nth-child(1) span {
        color: var(--color-text-placeholder);
        font-size: small;
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
        filter: brightness(1.1);
    }

    .ip-nic option::checkmark {
        width: 1.25rem;
        height: 1.25rem;
        display: grid;
        place-items: center;
    }

    selectedcontent div span:nth-child(2) {
        display: none !important;
        background-color: red;
    }

    selectedcontent div span:nth-child(1) {
        white-space: nowrap;

        max-width: 50vw !important;
        overflow: hidden;
        text-overflow: ellipsis;
    }
</style>
