<script>
    import { route_print_v4 } from "../js/tauri";
    import { settings } from "../js/store_settings";

    let routes = [];
    async function poll_interfaces() {
        routes = await route_print_v4();
        console.log(routes);
    }

    let interval;
    import { onDestroy, onMount } from "svelte";
    onMount(async () => {
        // Initial
        await poll_interfaces();

        // Poll
        interval = setInterval(async () => {
            await poll_interfaces();
        }, $settings.ipPollRate_ms);
    });
    onDestroy(async () => {
        clearInterval(interval);
    });
</script>

<article class="mono" class:coloredHeaders={$settings.coloredHeaders}>
    <table>
        <thead>
            <tr>
                <th>
                    <div class="flex align-center nowrap">
                        <i class="fa-solid fa-route purples" />
                        Destination
                    </div>
                </th>
                <th>Netmask</th>
                <th>Gateway</th>
                <th>Interface</th>
                <th>Metric</th>
                <th>Persistent</th>
            </tr>
        </thead>
        <tbody>
            {#each routes as route}
                {#if !($settings.hideDefaultRoutes && route.gateway === "On-link")}                    
                    <tr>
                        <td><div>{route.destination}</div></td>
                        <td><div>{route.netmask}</div></td>
                        <td><div>{route.gateway}</div></td>
                        <td><div>{route.interface}</div></td>
                        <td><div>{route.metric}</div></td>
                        <td><div>{route.persistent}</div></td>
                    </tr>
                {/if}
            {/each}
        </tbody>
    </table>
</article>

<style>
    article {
        display: flex;
        flex-wrap: wrap;
        gap: var(--gap);
        column-gap: calc(var(--gap) * 4);
        justify-content: flex-start;
        margin: var(--pad);
        align-items: flex-start;
    }
</style>
