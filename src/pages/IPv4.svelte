<script>
    import { ipv4 } from "../js/store_ipv4";
    import { settings } from "../js/store_settings";
    import InterfaceTable from "../components/InterfaceTable.svelte";
    import PresetTable from "../components/PresetTable.svelte";

    // Debug
    // $: console.log("$ipv4", $ipv4)
    // $: console.log("$ipv4.interface_active", $ipv4.interface_active)
    // $: console.log("$ipv4", $ipv4)
    // $: console.log("$ipv4", $ipv4)
</script>

<article class="mono tableStyle">
    <InterfaceTable />
    <!-- <br /> -->
    <PresetTable />
</article>

<style>
    article {
        display: grid;
        gap: var(--gap);
        margin: var(--pad);
    }

    /* tableStyle */
    :global(.tableStyle table),
    :global(.tableStyle tbody),
    :global(.tableStyle thead),
    :global(.tableStyle tr),
    :global(.tableStyle th),
    :global(.tableStyle td) {
        display: block;
    }
    :global(.tableStyle table) {
        font-family: var(--font-mono);
        display: grid;
        width: auto;
    }
    :global(.tableStyle tbody),
    :global(.tableStyle thead) {
        display: grid;
    }
    :global(.tableStyle thead) {
        color: var(--color-text-bright);
    }
    :global(.tableStyle tr) {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr 1fr 1fr 2.6rem;
    }
    :global(.tableStyle th) {
        display: flex;
        align-items: center;
    }
    :global(.tableStyle th),
    :global(.tableStyle td) {
        padding: 0;
    }
    :global(.tableStyle th div),
    :global(.tableStyle td div) {
        padding: var(--pad);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        width: 100%;
    }
    :global(.tableStyle table button) {
        padding: var(--pad);
        border-radius: 0;
        background: transparent;
        width: 100%;
    }
    :global(.tableStyle table input) {
        border-radius: calc(var(--radius)/2);
        border-radius: 0;
        width: 100%;
        padding: calc(var(--pad)/2);
    }

    /* Alt row colors */
    :global(tr:nth-child(even) td) {
        background-color: var(--color-bg-section);
    }
    :global(table tr th) {
        background-color: var(--color-bg-header);
    }
    :global(table tr th) {
        background-color: var(--color-bg-header);
    }

    /* All Borders */
    :global(table th) {
        border-top: var(--border);
    }
    :global(table th),
    :global(table td) {
        border-bottom: var(--border);
        border-right: var(--border);
    }
    :global(table th:first-child),
    :global(table td:first-child) {
        border-left: var(--border);
    }

</style>
