<script lang="ts">
	import { onMount } from 'svelte';

	export let level: 1 | 2 | 3 | 4 | 5 | 6 = 2;
	export let id: string;
	export let className: string = '';

	let element: HTMLElement;

	onMount(() => {
		// Scroll to element if hash matches on page load
		if (typeof window !== 'undefined') {
			const hash = window.location.hash.slice(1);
			if (hash === id) {
				setTimeout(() => {
					element?.scrollIntoView({ behavior: 'smooth', block: 'start' });
				}, 100);
			}
		}
	});

	function handleClick() {
		// Update URL hash
		if (typeof window !== 'undefined') {
			window.history.replaceState(null, '', `#${id}`);
		}

		// Scroll to element
		element?.scrollIntoView({ behavior: 'smooth', block: 'start' });

		// Copy link to clipboard
		const url = `${window.location.origin}${window.location.pathname}#${id}`;
		navigator.clipboard
			.writeText(url)
			.then(() => {
				// Show a brief tooltip or notification (optional)
				console.log('Link copied to clipboard:', url);
			})
			.catch((err) => {
				console.error('Failed to copy link:', err);
			});
	}
</script>

{#if level === 1}
	<button
		bind:this={element}
		{id}
		on:click={handleClick}
		class="clickable-headline cursor-pointer hover:text-primary transition-colors duration-200 group relative bg-transparent border-none p-0 {className}"
		title="Click to copy link to this section"
		aria-label="Copy link to this section"
	>
		<span
			class="absolute right-[-2rem] opacity-0 group-hover:opacity-100 transition-opacity duration-200 text-primary text-sm"
		>
			🔗
		</span>
		<h1 class="m-0">
			<slot />
		</h1>
	</button>
{:else if level === 2}
	<button
		bind:this={element}
		{id}
		on:click={handleClick}
		class="clickable-headline cursor-pointer hover:text-primary transition-colors duration-200 group relative bg-transparent border-none p-0 {className}"
		title="Click to copy link to this section"
		aria-label="Copy link to this section"
	>
		<span
			class="absolute right-[-2rem] opacity-0 group-hover:opacity-100 transition-opacity duration-200 text-primary text-sm"
		>
			🔗
		</span>
		<h2 class="m-0">
			<slot />
		</h2>
	</button>
{:else if level === 3}
	<button
		bind:this={element}
		{id}
		on:click={handleClick}
		class="clickable-headline cursor-pointer hover:text-primary transition-colors duration-200 group relative bg-transparent border-none p-0 {className}"
		title="Click to copy link to this section"
		aria-label="Copy link to this section"
	>
		<span
			class="absolute right-[-2rem] opacity-0 group-hover:opacity-100 transition-opacity duration-200 text-primary text-sm"
		>
			🔗
		</span>
		<h3 class="m-0">
			<slot />
		</h3>
	</button>
{:else if level === 4}
	<button
		bind:this={element}
		{id}
		on:click={handleClick}
		class="clickable-headline cursor-pointer hover:text-primary transition-colors duration-200 group relative bg-transparent border-none p-0 {className}"
		title="Click to copy link to this section"
		aria-label="Copy link to this section"
	>
		<span
			class="absolute right-[-2rem] opacity-0 group-hover:opacity-100 transition-opacity duration-200 text-primary text-sm"
		>
			🔗
		</span>
		<h4 class="m-0">
			<slot />
		</h4>
	</button>
{:else if level === 5}
	<button
		bind:this={element}
		{id}
		on:click={handleClick}
		class="clickable-headline cursor-pointer hover:text-primary transition-colors duration-200 group relative bg-transparent border-none p-0 {className}"
		title="Click to copy link to this section"
		aria-label="Copy link to this section"
	>
		<span
			class="absolute right-[-2rem] opacity-0 group-hover:opacity-100 transition-opacity duration-200 text-primary text-sm"
		>
			🔗
		</span>
		<h5 class="m-0">
			<slot />
		</h5>
	</button>
{:else}
	<button
		bind:this={element}
		{id}
		on:click={handleClick}
		class="clickable-headline cursor-pointer hover:text-primary transition-colors duration-200 group relative bg-transparent border-none p-0 {className}"
		title="Click to copy link to this section"
		aria-label="Copy link to this section"
	>
		<span
			class="absolute right-[-2rem] opacity-0 group-hover:opacity-100 transition-opacity duration-200 text-primary text-sm"
		>
			🔗
		</span>
		<h6 class="m-0">
			<slot />
		</h6>
	</button>
{/if}

<style>
	.clickable-headline:hover {
		text-decoration: underline;
		text-decoration-color: hsl(var(--p));
		text-underline-offset: 4px;
	}
</style>
