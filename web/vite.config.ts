import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite-plus";
import { sveltekit } from "@sveltejs/kit/vite";
import { inlineShikiCodeblocks } from "./vite-plugins/inline-shiki-codeblocks.js";

export default defineConfig({
	plugins: [inlineShikiCodeblocks(), sveltekit(), tailwindcss()],
});
