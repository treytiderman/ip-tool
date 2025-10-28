<script lang="ts">
    import { setPage } from "../ts/router";
    import { setTheme, themeStore } from "../ts/window";
    function changeTheme(event: any) {
        const themeSelected = event.target.value
        setTheme(themeSelected)
    }
</script>

<div class="grid gap-4 pad-1">

    <div class="grid gap-1">
        <label for="theme" class="text-dark thin">Theme</label>
        <select id="theme" class="border shadow" on:input={changeTheme} bind:value={$themeStore}>
            <option value="auto">Auto</option>
            <option value="black">Black</option>
            <option value="dark">Dark</option>
            <option value="white">White</option>
        </select>
    </div>

</div>
