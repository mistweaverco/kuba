<script lang="ts">
	import '../app.css';
	import '@fortawesome/fontawesome-free/css/all.min.css';
	import '@mistweaverco/mdsvex-shiki/styles.css';
	import { copyAction } from '@mistweaverco/mdsvex-shiki/copyAction';
	import Navigation from '$lib/Navigation.svelte';

	let { children } = $props();
</script>

<Navigation />

<main use:copyAction>
	{@render children()}
</main>
