<script>
    import ContextMenu from "./ContextMenu.svelte";
    import { ipv4 } from "../js/store_ipv4";

    const contextMenuItems = [
        {
            text: "Confirm",
            hide: validInterface,
            class: "fa-solid fa-check",
            onClick: () => confirm(),
        },
        {
            text: "Cancel",
            class: "fa-solid fa-xmark",
            onClick: () => cancel(),
        },
        {
            text: "hr",
        },
        {
            text: "Add IP",
            class: "fa-solid fa-plus",
            onClick: () => addBlankIP(),
        },
        {
            text: "Add DNS server",
            class: "fa-solid fa-plus",
            onClick: () => addBlankDNS(),
        },
        {
            text: "hr",
        },
        {
            text: "Remove Blanks",
            class: "fa-solid fa-trash",
            onClick: () => {
                removeBlankIPs();
                removeBlankDNS();
            },
        },
    ];
</script>

