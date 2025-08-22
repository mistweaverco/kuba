import { sveltePreprocess } from 'svelte-preprocess';
import { mdsvex } from 'mdsvex';
import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: [vitePreprocess(), sveltePreprocess(), mdsvex()],
	kit: {
		adapter: adapter()
	},
	extensions: ['.svelte', '.svx', '.mdx', '.md']
};

export default config;
