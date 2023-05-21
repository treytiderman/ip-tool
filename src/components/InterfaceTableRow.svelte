<script>
    import { set_dhcp, set_preset, set_metric_auto } from "../js/tauri";
    import { ipv4 } from "../js/store_ipv4";
    import ContextMenu from "../components/ContextMenu.svelte";

    // Event Dispatcher
    import { createEventDispatcher } from "svelte";
    const dispatch = createEventDispatcher();

    export let selected = true;
    export let nic = {
        interface_name: "Ethernet",
        interface_metric: 25,
        gateway: "192.168.1.254",
        gateway_metric: 1,
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "192.168.1.9",
                subnet_mask: "255.255.255.0",
            },
            {
                ip_address: "192.168.3.9",
                subnet_mask: "255.255.255.0",
            },
            {
                ip_address: "192.168.4.9",
                subnet_mask: "255.255.255.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["192.168.1.1"],
    };

    let contextMenu;
    const contextMenuItems = [
        {
            text: "Edit Interface",
            class: "fa-solid fa-pen-to-square purple",
            onClick: () => dispatch("edit"),
        },
        {
            text: "Set DHCP",
            class: "fa-solid fa-wand-magic-sparkles purple",
            onClick: () => set_dhcp(nic.interface_name),
        },
        {
            text: "Set Metric Auto",
            class: "fa-solid fa-route purple",
            onClick: () => set_metric_auto(nic.interface_name),
        },
        {
            text: "hr",
        },
        {
            text: "Set Selected Preset",
            class: "fa-solid fa-floppy-disk orange",
            onClick: () => set_preset($ipv4.interface_active.interface_name, $ipv4.preset_active),
        },
    ];
</script>

<ContextMenu
    bind:menu={contextMenu}
    items={contextMenuItems}
    on:any_click={() => contextMenu.hide()}
    on:any_contextmenu={() => contextMenu.hide()}
/>
{#if contextMenu?.show}<tr />{/if}
<tr
    class:selected
    on:click={() => dispatch("select")}
    on:contextmenu|preventDefault={(event) => {
        dispatch("select");
        contextMenu.showAtEvent(event);
    }}
>
    <td>
        <div>
            <span
                >{nic.interface_name}
                <small>
                    {@html nic.ip_is_dhcp ? `<i class="fa-solid fa-wand-magic-sparkles"/>` : ""} ({nic.interface_metric})
                </small>
            </span>
        </div>
    </td>
    <td>
        {#each nic.ip_and_masks as ip_and_mask}
            <div>{ip_and_mask.ip_address}</div>
        {/each}
    </td>
    <td>
        {#each nic.ip_and_masks as ip_and_mask}
            <div>{ip_and_mask.subnet_mask}</div>
        {/each}
    </td>
    <td>
        <div>{nic.gateway || ""}</div>
    </td>
    <td>
        {#each nic.dns_servers as dns_server}
            <div>{dns_server}</div>
        {/each}
    </td>
    <td>
        <button style="color: currentcolor;" on:click={(event) => contextMenu.showAtEvent(event)}>
            <i class="fa-solid fa-ellipsis-vertical" />
        </button>
    </td>
</tr>

<style>
    /* td {
        padding: 0;
        vertical-align: top;
    }
    td {
        min-width: fit-content;
        max-width: 10rem;
        overflow: scroll;
    }
    td::-webkit-scrollbar {
        display: none;
    }
    td:last-child {
        min-width: 2rem;
    }
    td > div {
        display: grid;
        gap: 0;
    }
    td:last-child > div {
        padding: 0;
        gap: 0;
    }

    td > div > span {
        padding: var(--pad);
    } */

    tr.selected td {
        background-color: var(--color-text-purple);
        color: var(--color-bg-purple);
    }
    tr.selected {
        /* background-color: var(--color-text-orange); */
        /* color: var(--color-bg-orange); */
        /* outline: var(--border); */
        /* outline-color: var(--color-bg-purple); */
        /* outline-width: 2px; */
        z-index: 2;
    }
    html.light tr.selected td {
        color: var(--color-text-bright);
    }
    td {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
    }
    td div {
        padding: calc(var(--pad) / 2) var(--pad);
    }
    td div:first-child {
        padding-top: var(--pad);
    }
    td div:last-child {
        padding-bottom: var(--pad);
    }
    /* tr.selected > td:first-child > div > span {
        padding: calc(var(--pad)/2) calc(3*var(--pad)/4);
        margin: calc(var(--pad)/2) calc(var(--pad)/4);
        width: fit-content;
        border-radius: var(--radius-lg);

        background-color: var(--color-bg-purple);
        color: var(--color-text-purple);
    } */
    /* 
    button {
        padding: var(--pad);
        background-color: transparent;
        border-radius: 0;
    } */
</style>
