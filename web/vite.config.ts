import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite-plus';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
	plugins: [sveltekit(), tailwindcss()]
});
