<script lang="ts">
</script>

<section class="grid gap">
    <form class="grid gap">
        <div class="grid gap-xs">
            <label for="preset-name">Preset Name</label>
            <input type="text" id="preset-name" name="preset-name" placeholder="eth0" required />
        </div>

        <div class="grid gap-xs">
            <label for="ip-address">IP Address</label>
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
            <label for="dns">DNS</label>
            <input type="text" id="dns" name="dns" placeholder="1.1.1.1" required />
        </div>

        <button type="submit">Save Changes</button>
    </form>
</section>

<style>

</style>
