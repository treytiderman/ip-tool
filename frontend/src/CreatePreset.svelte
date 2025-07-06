<script lang="ts">
</script>

<section class="grid gap">
    <form class="grid gap">
        <div class="grid gap-xs">
            <label for="preset-name">Preset Name</label>
            <input type="text" id="preset-name" name="preset-name" placeholder="eth0" required autofocus/>
        </div>

        <div class="grid gap-xs">
            <div class="flex center-y gap-xs">
                <label for="ip-address" class="grow">IP Address</label>
                <button class="ip-icon-button" type="button" title="Add IP Address">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M5 12h14"/>
                        <path d="M12 5v14"/>
                    </svg>
                </button>
            </div>
            <input type="text" id="ip-address" name="ip-address" placeholder="192.168.1.99" required />
        </div>

        <div class="grid gap-xs">
            <label for="subnet-mask">Subnet Mask</label>
            <input type="text" id="subnet-mask" name="subnet-mask" placeholder="255.255.255.0" required />
        </div>

        <div class="grid gap-xs">
            <label for="gateway">Gateway</label>
            <input type="text" id="gateway" name="gateway" placeholder="192.168.1.1" required />
        </div>

        <div class="grid gap-xs">
            <div class="flex center-y gap-xs">
                <label for="dns" class="grow">DNS</label>
                <button class="ip-icon-button" type="button" title="Add DNS Server">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M5 12h14"/>
                        <path d="M12 5v14"/>
                    </svg>
                </button>
            </div>
            <input type="text" id="dns" name="dns" placeholder="1.1.1.1" required />
        </div>

        <button type="submit">
            <div>Save Changes</div>
        </button>
    </form>
</section>

<style>

</style>
