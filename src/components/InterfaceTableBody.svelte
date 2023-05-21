<script>
    import { set_preset, set_interface_name, set_metric } from "../js/tauri"
    import { ipv4 } from "../js/store_ipv4"
    import InterfaceTableRow from "./InterfaceTableRow.svelte"
    import InterfaceTableRowEdit from "./InterfaceTableRowEdit.svelte"
</script>

<tbody>
    {#each $ipv4.interfaces as nic}
        {#if $ipv4.interface_active.interface_name === nic.interface_name && $ipv4.interface_editing}
            <InterfaceTableRowEdit
                {nic}
                on:confirm={(edit) => {
                    $ipv4.interface_editing = false
                    console.log("edit confimed", edit.detail)
                    set_preset(nic.interface_name, edit.detail)
                    if (nic.interface_name !== edit.detail.interface_name) {
                        console.log("set_interface_name", nic.interface_name, edit.detail.interface_name)
                        set_interface_name(nic.interface_name, edit.detail.interface_name)
                    }
                    if (nic.interface_metric !== edit.detail.interface_metric) {
                        console.log("set_interface_metric", nic.interface_metric, edit.detail.interface_metric)
                        set_metric(nic.interface_name, edit.detail.interface_metric)
                    }
                }}
                on:cancel={() => {
                    $ipv4.interface_editing = false
                }}
            />
        {:else}
            <InterfaceTableRow
                {nic}
                selected={$ipv4.interface_active.interface_name === nic.interface_name}
                on:edit={() => {
                    $ipv4.interface_editing = true
                }}
                on:select={() => {
                    $ipv4.interface_editing = false
                    $ipv4.interface_active = nic
                }}
            />
        {/if}
    {/each}
</tbody>

<style>
</style>
