<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { initNics, pollNics } from "../ts/nic";
    import { settingsStore } from "../ts/settings";

    // Components
    import Presets_Presets from "./Presets_Presets.svelte";
    import Presets_Interface from "./Presets_Interface.svelte";

    // Page start
    let interval: number;
    onMount(async () => {
        await initNics();
        interval = pollNics($settingsStore.interfacePollRateMs);
    });
    
    // Page close
    onDestroy(async () => {
        clearInterval(interval);
    });
</script>

<div class="grow height-100 flex column gap-1 pad-1">
    <Presets_Interface />
    <div></div>
    <Presets_Presets />
</div>
