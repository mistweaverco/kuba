<script lang="ts">
	import '../app.css';
	import Navigation from '$lib/Navigation.svelte';

	let { children } = $props();
</script>

<Navigation />

<main>
	{@render children()}
</main>
