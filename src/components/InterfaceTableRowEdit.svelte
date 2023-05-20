<script>
    import { clone, validIPv4, validMask } from "../js/helper";

    // Components
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
    const edit = clone(nic);

    let contextMenu;
    $: contextMenuItems = [
        {
            text: "Confirm",
            hide: !isValidInterface(edit),
            class: "fa-solid fa-check green",
            onClick: () => confirm(),
        },
        {
            text: "Cancel",
            class: "fa-solid fa-xmark red",
            onClick: () => cancel(),
        },
        {
            text: "hr",
        },
        {
            text: "Remove Blanks",
            class: "fa-solid fa-minus purple",
            onClick: () => {
                removeBlankIPs();
                removeBlankDNS();
            },
        },
    ];

    function focus(element, bool) {
        if (bool) element.focus();
    }
    function removeBlankIPs() {
        edit.ip_and_masks = edit.ip_and_masks.filter((ip_mask) => ip_mask.ip_address !== "");
    }
    function removeLastIP() {
        edit.ip_and_masks.pop()
        edit.ip_and_masks = edit.ip_and_masks
    }
    function addBlankIP() {
        removeBlankIPs();
        if (edit.ip_and_masks.length < 5) {
            edit.ip_and_masks.push({
                ip_address: "",
                subnet_mask: "",
            });
        }
    }
    function removeBlankDNS() {
        edit.dns_servers = edit.dns_servers.filter((dns_server) => dns_server !== "");
    }
    function removeLastDNS() {
        edit.dns_servers.pop()
        edit.dns_servers = edit.dns_servers
    }
    function addBlankDNS() {
        removeBlankDNS();
        if (edit.dns_servers.length < 2) {
            edit.dns_servers.push("");
        }
    }
    function confirm() {
        removeBlankIPs();
        removeBlankDNS();
        dispatch("confirm", edit);
    }
    function cancel() {
        dispatch("cancel");
    }

    let validInterface = true;
    $: validInterface = isValidInterface(edit);
    function isValidInterface(edit) {
        let validInterface = true;
        edit.ip_and_masks.forEach((ip_and_mask) => {
            if (ip_and_mask.ip_address === "" && ip_and_mask.subnet_mask === "") {
            } else {
                if (validIPv4(ip_and_mask.ip_address) === false) {
                    validInterface = false;
                }
                if (validMask(ip_and_mask.subnet_mask) === false) {
                    validInterface = false;
                }
            }
        });
        edit.dns_servers.forEach((dns_server) => {
            if (dns_server === "") {
            } else {
                if (validIPv4(dns_server) === false) {
                    validInterface = false;
                }
            }
        });
        if (edit.gateway === "") {
        } else {
            if (validIPv4(edit.gateway) === false) {
                validInterface = false;
            }
        }
        return validInterface;
    }
</script>

<ContextMenu
    items={contextMenuItems}
    bind:menu={contextMenu}
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
        <!-- <span>{edit.interface_name} {nic.ip_is_dhcp ? "(DHCP)" : ""}</span> -->
        <input type="text" bind:value={edit.interface_name} />
    </td>
    <td>
        {#each edit.ip_and_masks as ip_and_mask, index}
            {#if index === 0}
                <input
                    use:focus={index === 0}
                    type="text"
                    class:red={!validIPv4(ip_and_mask.ip_address)}
                    bind:value={ip_and_mask.ip_address}
                />
            {:else}
                <input
                    type="text"
                    class:red={!validIPv4(ip_and_mask.ip_address)}
                    bind:value={ip_and_mask.ip_address}
                />
            {/if}
        {/each}
        <aside class="flex even">
            <button class="add" on:click={addBlankIP}>
                <i class="fa-solid fa-plus" />
            </button>
            <button class="add" on:click={removeLastIP}>
                <i class="fa-solid fa-minus" />
            </button>
        </aside>
    </td>
    <td>
        {#each edit.ip_and_masks as ip_and_mask}
            <input
                type="text"
                class:red={!validIPv4(ip_and_mask.subnet_mask)}
                bind:value={ip_and_mask.subnet_mask}
            />
        {/each}
    </td>
    <td>
        <input
            type="text"
            class:red={!validIPv4(edit.gateway)}
            bind:value={edit.gateway}
        />
    </td>
    <td>
        {#each edit.dns_servers as dns_server}
            <input
                type="text"
                class:red={!validIPv4(dns_server)}
                bind:value={dns_server}
            />
        {/each}
        <aside class="flex even">
            <button class="add" on:click={addBlankDNS}>
                <i class="fa-solid fa-plus" />
            </button>
            <button class="add" on:click={removeLastDNS}>
                <i class="fa-solid fa-minus" />
            </button>
        </aside>
    </td>
    <td style="padding: 0;">
        <!-- {#if validInterface}
                <button on:click={() => confirm()}>
                    <i class="fa-solid fa-check" />
                </button>
            {:else}
                <button on:click={() => cancel()}>
                    <i class="fa-solid fa-xmark" />
                </button>
            {/if} -->
        <button style="color: currentcolor;" on:click={(event) => contextMenu.showAtEvent(event)}>
            <i class="fa-solid fa-ellipsis-vertical" />
        </button>
    </td>
</tr>

<style>
    tr.selected td {
        /* background-color: var(--color-text-purple); */
        /* color: var(--color-bg-purple); */
        color: var(--color-text-bright);
    }
    tr.selected {
        /* background-color: var(--color-text-orange); */
        /* color: var(--color-bg-orange); */
        /* outline: var(--border); */
        /* outline-color: var(--color-bg-purple); */
        /* outline-width: 2px; */
        z-index: 2;
    }
    td {
        /* padding: 0 calc(var(--pad) / 2) calc(var(--pad) / 2) calc(var(--pad) / 2); */
        padding: calc(var(--pad) / 2);
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        gap: calc(var(--pad) / 2);
        /* width: 100%; */
    }
    aside {
        width: 100%;
        gap: calc(var(--pad) / 2);
    }
    .add {
        background-color: var(--color-bg-input);
        color: var(--color-text-input);
        padding: calc(var(--pad) / 2);
        width: 100%;

        /* border: var(--border); */
        /* border-color: var(--color-text); */
    }
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
        color: var(--color-text-input);
        padding: var(--pad);
    }
    tr {
        box-shadow: inset 0px 0px 0px var(--border-thickness)
            var(--color-bg-red);
    }
    button,
    input {
        padding: var(--pad);
        background-color: transparent;
        border-radius: 0;
    }
    input {
        width: auto;
    } */
</style>
