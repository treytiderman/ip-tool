<script>
    import { release_dhcp, renew_dhcp, flushdns } from "../js/tauri";

    // Components
    import ContextMenu from "./ContextMenu.svelte";

    // Event Dispatcher
    import { createEventDispatcher } from "svelte";
    const dispatch = createEventDispatcher();

    // Variables
    let contextMenu;
    const contextMenuItems = [
        {
            text: "Release & Renew DHCP",
            class: "fa-solid fa-shuffle purple",
            onClick: () => {
                release_dhcp();
                setTimeout(() => {
                    renew_dhcp();
                }, 1000);
            },
        },
        {
            text: "hr",
        },
        {
            text: "Refresh Interfaces",
            class: "fa-solid fa-arrows-rotate purple",
            onClick: () => dispatch("get_interfaces"),
        },
        {
            text: "Flush DNS",
            class: "fa-solid fa-toilet purple",
            onClick: () => flushdns(),
        },
    ];
</script>

<thead>
    <ContextMenu
        items={contextMenuItems}
        bind:menu={contextMenu}
        on:any_click={() => contextMenu.hide()}
        on:any_contextmenu={() => contextMenu.hide()}
    />
    <tr>
        <th>
            <div class="flex align-center nowrap">
                <i class="fa-solid fa-ethernet purples" />
                Interface
            </div>
        </th>
        <th><div>IP Address(s)</div></th>
        <th><div>Subnet Mask(s)</div></th>
        <th><div>Gateway</div></th>
        <th><div>DNS Server(s)</div></th>
        <th>
            <button style="color: currentcolor;" on:click={(e) => contextMenu.showAtEvent(e)}>
                <i class="fa-solid fa-ellipsis-vertical" />
            </button>
        </th>
    </tr>
</thead>

<style>
    /* th {
        min-width: fit-content;
        padding: 0;
    }
    span {
        padding: 0 var(--pad);
    }
    th:last-child {
        min-width: 2rem;
    }
    button {
        padding: var(--pad);
        background-color: transparent;
        color: inherit;
        border-radius: 0;
    } */
</style>
