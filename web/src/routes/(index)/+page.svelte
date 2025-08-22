<script lang="ts">
	import Prism from 'prismjs';
	import 'prismjs/plugins/toolbar/prism-toolbar';
	import 'prismjs/plugins/copy-to-clipboard/prism-copy-to-clipboard';
	import 'prismjs/components/prism-json';
	import 'prismjs/components/prism-yaml';
	import 'prismjs/components/prism-bash';
	import 'prismjs/components/prism-powershell';
	import 'prismjs/themes/prism-okaidia.css';
	import { onMount } from 'svelte';
	import HeadComponent from '$lib/HeadComponent.svelte';

	let installUsing = 'manual';

	const handleAnchorClick = (evt: Event) => {
		evt.preventDefault();
		const link = evt.currentTarget as HTMLAnchorElement;
		const anchorId = new URL(link.href).hash.replace('#', '');
		const anchor = document.getElementById(anchorId);
		window.scrollTo({
			top: anchor?.offsetTop,
			behavior: 'smooth'
		});
	};

	const onInstallUsingChange = (evt: Event) => {
		const select = evt.currentTarget as HTMLSelectElement;
		installUsing = select.value;
	};

	onMount(() => {
		Prism.plugins.toolbar.registerButton('fullscreen-code', function (env: {element: HTMLElement}) {
			const button = document.createElement('button');
			button.innerHTML = '🔍';
			button.addEventListener('click', function () {
				const parent = env.element.parentNode as HTMLElement;
				parent.requestFullscreen();
			});

			return button;
		});

		Prism.highlightAll();
	});
</script>

<HeadComponent
	data={{
		title: 'Kuba - Securely and easily access your environment variables',
		description:
			'Pass env directly from GCP Secret Manager, AWS Secrets Manager, and Azure Key Vault to your application.',
	}}
/>

<div id="start" class="hero bg-base-200 min-h-screen">
	<div class="hero-content text-center">
		<div class="max-w-md">
			<img src="/logo.svg" alt="Kuba" class="m-5 mx-auto w-64" />
			<h1 class="text-5xl font-bold">Kuba</h1>
			<p class="py-6">Securely and easily access your environment variables</p>
			<a href="#install" on:click={handleAnchorClick}
				><button class="btn btn-primary">Get Started</button></a
			>
		</div>
	</div>
</div>
<div id="install" class="hero bg-base-200 min-h-screen">
	<div class="hero-content text-center">
		<div class="max-w-md">
			<h1 class="text-5xl font-bold">Install ⚡</h1>
			<p class="py-6">Install Kuba using ...</p>
			<select on:input={onInstallUsingChange} class="select select-bordered mb-5">
				<option value="manual" selected>manual</option>
				<option value="curl">curl (linux/mac)</option>
				<option value="wget">wget (linux/mac)</option>
				<option value="pwsh">pwsh (windows)</option>
			</select>
			<div class={installUsing === 'curl' ? '' : 'hidden'}>
				<pre><code
						class="language-bash"
						data-toolbar-order="copy-to-clipboard"
						data-prismjs-copy="📋">curl -sSL https://kuba.mwco.app/install.sh | bash</code
					></pre>
			</div>
			<div class={installUsing === 'wget' ? '' : 'hidden'}>
				<pre><code
						class="language-bash"
						data-toolbar-order="copy-to-clipboard"
						data-prismjs-copy="📋">wget -qO- https://kuba.mwco.app/install.sh | bash</code
					></pre>
			</div>
			<div class={installUsing === 'pwsh' ? '' : 'hidden'}>
				<pre><code
						class="language-powershell"
						data-toolbar-order="copy-to-clipboard"
						data-prismjs-copy="📋">iwr https://kuba.mwco.app/install.ps1 -useb | iex</code
					></pre>
			</div>
			<div class={installUsing === 'manual' ? '' : 'hidden'}>
				<p class="py-6">
					Download the latest release from the <a class="text-secondary" href="/download">download page</a>.
			</div>
			<p>
				<a href="#configure" on:click={handleAnchorClick}
					><button class="btn btn-primary mt-5">Configure</button></a
				>
			</p>
		</div>
	</div>
</div>
<div id="configure" class="hero bg-base-200 min-h-screen">
	<div class="hero-content text-center">
		<div class="max-w-md">
			<h1 class="text-5xl font-bold">Configure 🔧</h1>
			<p class="py-6">
				Configure Kuba using a simple configuration file <code>kuba.yaml</code>.
			</p>
			<div class="mb-5">
				<pre><code
						class="language-bash"
						data-toolbar-order="copy-to-clipboard"
						data-prismjs-copy="📋">kuba init</code
					></pre>
			</div>
			<div role="alert" class="alert alert-info">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					class="h-6 w-6 shrink-0 stroke-current"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
					></path>
				</svg>
				<span>
					This will generate a default configuration file for you (if it does not yet exist),
					which you can customize to your needs.
				</span>
			</div>
			<p>
				<a href="#usage" on:click={handleAnchorClick}
					><button class="btn btn-primary mt-5">Usage</button></a
				>
			</p>
		</div>
	</div>
</div>
<div id="usage" class="hero bg-base-200 min-h-screen">
	<div class="hero-content text-center">
		<div class="max-w-md">
			<h1 class="text-5xl font-bold">Usage 🔐</h1>
			<p class="py-6">
				Once you have your <code>kuba.yaml</code> file set up, you can run
				<code>kuba run -- &lt;your-command&gt;</code>
				to execute any command with the secrets
				from your cloud provider's secret management system.
			</p>
			<p>
				<a href="#why" on:click={handleAnchorClick}
					><button class="btn btn-secondary mt-5">Why?</button></a
				>
			</p>
		</div>
	</div>
</div>
<div id="why" class="hero bg-base-200 min-h-screen">
	<div class="hero-content text-center">
		<div class="max-w-md">
			<h1 class="text-5xl font-bold">Why? 🤔</h1>
			<p class="py-6">What are the benefits of using Kuba?</p>
			<p class="py-6">
				Environment variables are a common way to manage configuration in applications,
				especially when deploying to different environments like development,
				staging, and production.
			</p>
			<p class="py-6">
				However, managing these variables can become cumbersome,
				especially when dealing with multiple cloud providers and
				secret management systems.
			</p>
			<p class="py-6">
				This often leads to the use of <code>.env</code> files,
				which can be problematic for several reasons:
			</p>
			<ul class="list-disc list-inside">
				<li>Insecure sharing of secrets</li>
				<li>Manual management of secrets</li>
				<li>Security risks of committing secrets to version control</li>
				<li>Lack of standardization across cloud providers</li>
			</ul>
			<p class="py-6">
				Kuba is designed to address these issues by providing a unified way to
				manage environment variables across different cloud providers.
			</p>
			<p>
				<a href="#get-involved" on:click={handleAnchorClick}
					><button class="btn btn-secondary mt-5">Get involved</button></a
				>
			</p>
		</div>
	</div>
</div>
<div id="get-involved" class="hero bg-base-200 min-h-screen">
	<div class="hero-content text-center">
		<div class="max-w-md">
			<h1 class="text-5xl font-bold">Get involved 📦</h1>
			<p class="py-6">Kuba is open-source and we welcome contributions.</p>
			<p>
				View the <a class="text-secondary" href="https://github.com/mistweaverco/kuba">code</a>,
				and/or check out the <a class="text-secondary" href="https://kuba.mwco.app/docs">docs</a>.
			</p>
		</div>
	</div>
</div>
