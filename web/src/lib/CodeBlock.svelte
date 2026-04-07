<script lang="ts">
	/**
	 * Fallback when the Vite `inline-shiki-codeblocks` plugin has not transformed the page yet.
	 * Production and dev builds normally replace `<CodeBlock>` with inlined Shiki HTML first.
	 */
	let { code, lang, meta }: { code: string; lang: string; meta?: string } = $props();
	void lang;
	void meta;
</script>

<div class="mb-0">
	<pre class="rounded border border-base-300 bg-base-200 p-4 text-sm"><code>{code}</code></pre>
</div>
