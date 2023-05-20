<script>
    // Components
    import ContextMenu from "./ContextMenu.svelte";

    // Event Dispatcher
    import { createEventDispatcher } from "svelte";
    const dispatch = createEventDispatcher();

    // Variables
    let contextMenu;
    const contextMenuItems = [
        {
            text: "Add Preset",
            class: "fa-solid fa-add orange",
            onClick: () => dispatch("new"),
        },
        {
            text: "hr",
        },
        {
            text: "Set ALL Presets Default",
            class: "fa-solid fa-rotate-right red",
            onClick: () => dispatch("reset"),
        },
        {
            text: "Delete ALL Presets",
            class: "fa-solid fa-trash red",
            onClick: () => dispatch("clear"),
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
                <i class="fa-solid fa-floppy-disk oranges" />
                Preset
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
</style>
