<script>
    import { get_interfaces } from "../js/tauri";
    import { ipv4 } from "../js/store_ipv4";
    import { settings } from "../js/store_settings";
    import InterfaceTableHeader from "../components/InterfaceTableHeader.svelte";
    import InterfaceTableBody from "../components/InterfaceTableBody.svelte";

    async function poll_interfaces() {
        const interfaces = await get_interfaces();
        if (JSON.stringify($ipv4.interfaces) !== JSON.stringify(interfaces)) {
            $ipv4.interfaces = interfaces;
        }
    }

    let interval;
    import { onDestroy, onMount } from "svelte";
    onMount(async () => {
        // Initial
        await poll_interfaces();
        $ipv4.interface_active = $ipv4.interfaces[0];

        // Poll
        interval = setInterval(async () => {
            await poll_interfaces();
        }, $settings.ipPollRate_ms);
    });
    onDestroy(async () => {
        clearInterval(interval);
    });
</script>

<!-- Interface Table -->
<table id="interface_table">
    <InterfaceTableHeader
        on:get_interfaces={async () => {
            $ipv4.interfaces = []
            await poll_interfaces();
        }}
    />
    <InterfaceTableBody
        on:get_interfaces={async () => {
            $ipv4.interfaces = []
            await poll_interfaces();
        }}
    />
</table>

<style>
</style>
