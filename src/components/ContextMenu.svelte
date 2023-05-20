<script>
    // Event Dispatcher
    import { createEventDispatcher } from "svelte";
    const dispatch = createEventDispatcher();

    // Variables
    export let menu = {
        show: false,
        position: { x: 0, y: 0 },
        size: { h: 0, w: 0 },
        hide: false,
        showAtEvent: (event) => {
            if (event.pointerId === -1) {
                console.log("showAtEvent", event);
                console.log("target", event.target.getBoundingClientRect());
                const buttonPosition = event.target.getBoundingClientRect();
                menu.position = { x: buttonPosition.x, y: buttonPosition.y + buttonPosition.height };
            }
            else {
                menu.position = { x: event.clientX, y: event.clientY };
            }
            checkBounds();
            setTimeout(() => menu.show = true, 10);
        },
        showAtXY: (x, y) => {
            menu.position = { x: x, y: y };
            checkBounds();
            setTimeout(() => menu.show = true, 10);
        },
        hide: () => {
            menu.show = false;
        },
    };
    export let items = [
        {
            text: "Test Item",
            class: "fa-solid fa-plus",
            onClick: () => {},
        },
        {
            text: "hr",
        },
        {
            text: "Test Item",
            class: "fa-solid fa-plus",
            onClick: () => {},
        },
    ]

    // Functions
    function checkBounds() {
        const browser = { w: window.innerWidth, h: window.innerHeight };
        if (browser.h - menu.position.y < menu.size.h) menu.position.y -= menu.size.h;
        if (browser.w - menu.position.x < menu.size.w) menu.position.x -= menu.size.w;
    }
    function getContextMenuDimensions(node) {
        menu.size.h = node.offsetHeight;
        menu.size.w = node.offsetWidth;
        checkBounds();
    }
    function focus(element, bool) {
        if (bool) element.focus();
    }
</script>

<svelte:window
    on:click={() => dispatch("any_click")}
    on:contextmenu={() => dispatch("any_contextmenu")}
/>

{#if menu?.show}
    <nav
        use:getContextMenuDimensions
        style="top:{menu.position.y}px; left:{menu.position.x}px; z-index: 99"
    >
        <div class="navbar">
            {#each items as item, index}
                {#if item.text == "hr"}
                    <hr/>
                {:else if item.hide === true}
                    <button disabled={item.hide} use:focus={index === 0} class="flex nowrap" on:click={item.onClick}>
                        <i class={item.class}/>
                        <span>{item.text}</span>
                    </button>
                {:else}
                    <button use:focus={index === 0} class="flex nowrap" on:click={item.onClick}>
                        <i class={item.class}/>
                        <span>{item.text}</span>
                    </button>
                {/if}
            {/each}
        </div>
    </nav>
{/if}

<style>
    nav {
        position: absolute;
        position: fixed;
    }
    .navbar {
        display: flex;
        gap: 0;
        min-width: 15rem;
        background-color: var(--color-bg-section);
        /* border-radius: var(--radius); */
        border: var(--border);
        flex-direction: column;
        padding: calc(var(--pad)/2);
        /* padding: var(--pad); */
        /* border-radius: var(--radius-lg); */
        /* overflow: hidden; */
        /* gap: var(--pad); */
    }
    button {
        text-align: left;
        background-color: var(--color-bg-section);
    }
    hr {
        background-color: var(--color-border);
        margin: 0 auto;
        margin: var(--pad) auto;
    }
</style>
