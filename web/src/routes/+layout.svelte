<script lang="ts">
	import '../app.css';
	import '@fortawesome/fontawesome-free/css/all.min.css';
	import Navigation from '$lib/Navigation.svelte';
	import Prism from 'prismjs';
	import 'prismjs/plugins/toolbar/prism-toolbar';
	import 'prismjs/plugins/copy-to-clipboard/prism-copy-to-clipboard';
	import 'prismjs/components/prism-bash';
	import 'prismjs/components/prism-yaml';
	import 'prismjs/components/prism-json';
	import 'prismjs/components/prism-javascript';
	import 'prismjs/components/prism-python';
	import 'prismjs/components/prism-powershell';
	import 'prismjs/components/prism-docker';
	import 'dracula-prism/dist/css/dracula-prism.css';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { onMount, tick } from 'svelte';

	let { children } = $props();
	let mainEl: HTMLElement | null = null;
	let highlightTimer: ReturnType<typeof setTimeout> | undefined;

	function registerFullscreenButtonOnce() {
		const anyPrism = Prism as unknown as { __kubaFullscreenRegistered?: boolean; plugins?: any };
		if (anyPrism.__kubaFullscreenRegistered) return;
		anyPrism.__kubaFullscreenRegistered = true;

		Prism.plugins?.toolbar?.registerButton?.(
			'fullscreen-code',
			function (env: { element: HTMLElement }) {
				const button = document.createElement('button');
				button.textContent = '🔍';
				button.addEventListener('click', function () {
					const parent = env.element.parentNode as HTMLElement | null;
					parent?.requestFullscreen?.();
				});
				return button;
			}
		);
	}

	async function highlight() {
		registerFullscreenButtonOnce();
		// Wait until DOM settles, then highlight.
		await tick();
		requestAnimationFrame(() => {
			setTimeout(() => {
				if (mainEl) {
					Prism.highlightAllUnder(mainEl);
				} else {
					Prism.highlightAll();
				}
			}, 0);
		});
	}

	onMount(() => {
		highlight();
		return () => {
			if (highlightTimer) clearTimeout(highlightTimer);
		};
	});

	// SvelteKit client-side navigation updates route content without remounting the layout.
	// Run highlighting after each update; debounce to avoid repeated work in dev/HMR.
	$effect(() => {
		if (!browser) return;
		// Touch $page so this effect runs on client-side navigations.
		$page.url.pathname;
		$page.url.hash;
		if (highlightTimer) clearTimeout(highlightTimer);
		highlightTimer = setTimeout(() => {
			highlight();
		}, 0);
	});
</script>

<Navigation />

<main bind:this={mainEl}>
	{@render children()}
</main>
