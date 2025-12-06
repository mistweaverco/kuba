<script lang="ts">
	import '../app.css';
	import '@fortawesome/fontawesome-free/css/all.min.css';
	import Navigation from '$lib/Navigation.svelte';

	let { children } = $props();
</script>

<Navigation />

<main>
	{@render children()}
</main>
