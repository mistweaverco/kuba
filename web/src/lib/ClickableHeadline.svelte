<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	export let level: 1 | 2 | 3 | 4 | 5 | 6 = 2;
	export let id: string;
	export let className: string = '';

	let element: HTMLElement;
	let showToast = false;
	let toastTimeout: ReturnType<typeof setTimeout>;

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

	onDestroy(() => {
		// Clean up timeout when component is destroyed
		if (toastTimeout) {
			clearTimeout(toastTimeout);
		}
	});

	function showToastNotification() {
		showToast = true;
		// Clear any existing timeout
		if (toastTimeout) {
			clearTimeout(toastTimeout);
		}
		// Hide toast after 3 seconds
		toastTimeout = setTimeout(() => {
			showToast = false;
		}, 3000);
	}

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
				// Show toast notification
				showToastNotification();
			})
			.catch((err) => {
				console.error('Failed to copy link:', err);
			});
	}
</script>

<!-- Toast notification -->
{#if showToast}
	<div class="toast toast-top toast-end z-50" role="alert" aria-live="polite">
		<div class="alert alert-success shadow-lg">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="stroke-current shrink-0 h-6 w-6"
				fill="none"
				viewBox="0 0 24 24"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
			<span>Link copied to clipboard!</span>
			<button 
				class="btn btn-sm btn-ghost" 
				on:click={() => showToast = false}
				aria-label="Close notification"
			>
				✕
			</button>
		</div>
	</div>
{/if}

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

	/* Toast animation styles */
	.toast {
		animation: slideInRight 0.3s ease-out;
	}

	@keyframes slideInRight {
		from {
			transform: translateX(100%);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	.alert {
		transition: all 0.2s ease-in-out;
	}

	.alert:hover {
		transform: translateY(-2px);
		box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
	}
</style>
