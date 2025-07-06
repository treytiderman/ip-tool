<script lang="ts">
    import { Greet } from "../wailsjs/go/main/App.js";
    let name: string;
    let resultText: string = "Please enter your name below 👇";
    function greet(): void {
        Greet(name).then((result) => (resultText = result));
    }
</script>