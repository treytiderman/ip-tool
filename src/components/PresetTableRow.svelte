<script>
    import { set_dhcp, set_preset, route_add, route_add_persistent, route_delete } from "../js/tauri";
    import { ipv4 } from "../js/store_ipv4";
    import ContextMenu from "./ContextMenu.svelte";

    // Event Dispatcher
    import { createEventDispatcher } from "svelte";
    const dispatch = createEventDispatcher();

    export let selected = true;
    export let preset = {
        name: "default_1",
        gateway: "192.168.1.1",
        ip_is_dhcp: false,
        ip_and_masks: [
            {
                ip_address: "192.168.1.7",
                subnet_mask: "255.255.255.0",
            },
        ],
        dns_is_dhcp: false,
        dns_servers: ["192.168.1.1"],
    };

    let contextMenu;
    const contextMenuItems = [
        {
            text: "Edit Preset",
            class: "fa-solid fa-pen-to-square orange",
            onClick: () => dispatch("edit"),
        },
        {
            text: "hr",
        },
        {
            text: "Delete Preset",
            class: "fa-solid fa-trash red",
            onClick: () => dispatch("delete"),
        },
        {
            text: "hr",
        },
        {
            text: "Set Selected Inteface",
            class: "fa-solid fa-ethernet purple",
            onClick: () => set_preset($ipv4.interface_active.interface_name, $ipv4.preset_active),
        },
        {
            text: "hr",
        },
        {
            text: "Add Route",
            class: "fa-solid fa-route purple",
            onClick: () => {
                $ipv4.preset_active.ip_and_masks.forEach(ip_and_mask => {                 
                    let destination = zeroLastOctet(ip_and_mask.ip_address)
                    let mask = ip_and_mask.subnet_mask
                    let gateway = $ipv4.interface_active.gateway
                    let metric = "1"
                    route_add(destination, mask, gateway, metric)
                });
            },
        },
        {
            text: "Add Persistent Route",
            class: "fa-solid fa-route purple",
            onClick: () => {
                $ipv4.preset_active.ip_and_masks.forEach(ip_and_mask => {                 
                    let destination = zeroLastOctet(ip_and_mask.ip_address)
                    let mask = ip_and_mask.subnet_mask
                    let gateway = $ipv4.interface_active.gateway
                    let metric = "1"
                    route_add_persistent(destination, mask, gateway, metric)
                });
            },
        },
        {
            text: "Delete Route",
            class: "fa-solid fa-trash red",
            onClick: () => {
                $ipv4.preset_active.ip_and_masks.forEach(ip_and_mask => {                 
                    let destination = zeroLastOctet(ip_and_mask.ip_address)
                    route_delete(destination)
                });
            },
        },
    ];

    function zeroLastOctet(ip) {
        return ip.slice(0, ip.lastIndexOf(".")) + ".0"
    }
</script>

<ContextMenu
    bind:menu={contextMenu}
    items={contextMenuItems}
    on:any_click={() => contextMenu.hide()}
    on:any_contextmenu={() => contextMenu.hide()}
/>
{#if contextMenu?.show}<tr></tr>{/if}
<tr
    class:selected
    on:click={() => dispatch("select")}
    on:contextmenu|preventDefault={(event) => {
        dispatch("select");
        contextMenu.showAtEvent(event);
    }}
>
    <td>
        <div>{preset.name} {preset.ip_is_dhcp ? "(DHCP)" : ""}</div>
    </td>
    <td>
        {#each preset.ip_and_masks as ip_and_mask}
            <div>{ip_and_mask.ip_address}</div>
        {/each}
    </td>
    <td>
        {#each preset.ip_and_masks as ip_and_mask}
            <div>{ip_and_mask.subnet_mask}</div>
        {/each}
    </td>
    <td>
        <div>{preset.gateway || "-"}</div>
    </td>
    <td>
        {#each preset.dns_servers as dns_server}
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
    td {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
    }
    td div {
        padding: calc(var(--pad)/2) var(--pad);
    }
    td div:first-child {
        padding-top: var(--pad);
    }
    td div:last-child {
        padding-bottom: var(--pad);
    }
    /* td {
        padding: 0;
        vertical-align: top;
    }
    td {
        min-width: fit-content;
    }

    
    td > div {
        display: grid;
        gap: 0;
    }
    td:last-child > div {
        padding: 0;
        gap: 0;
    }

    td > div > div {
        padding: var(--pad);
    } */
    tr.selected {
        /* background-color: var(--color-text-orange); */
        /* color: var(--color-bg-orange); */
        color: var(--color-text-bright);
        /* outline: var(--border); */
        /* outline-color: var(--color-bg-orange); */
        /* outline-width: 2px; */
        z-index: 2;
    }
    tr.selected td {
        background-color: var(--color-text-orange);
        color: var(--color-bg-orange);
    }
    html.light tr.selected td {
        color: var(--color-text-bright);
    }
    /* tr.selected > td:first-child > div > div {
        padding: calc(var(--pad)/2) calc(3*var(--pad)/4);
        margin: calc(var(--pad)/2) calc(var(--pad)/4);
        width: fit-content;
        border-radius: var(--radius-lg);

        background-color: var(--color-bg-orange);
        color: var(--color-text-orange);
    } */

    /* button {
        padding: var(--pad);
        background-color: transparent;
        color: inherit;
        border-radius: 0;
    } */
</style>
