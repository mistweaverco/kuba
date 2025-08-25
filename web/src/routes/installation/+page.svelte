<script lang="ts">
	import HeadComponent from '$lib/HeadComponent.svelte';
	import Prism from 'prismjs';
	import 'prismjs/plugins/toolbar/prism-toolbar';
	import 'prismjs/plugins/copy-to-clipboard/prism-copy-to-clipboard';
	import 'prismjs/components/prism-bash';
	import 'prismjs/components/prism-powershell';
	import 'prismjs/themes/prism-okaidia.css';
	import { onMount } from 'svelte';

	let installUsing = 'manual';

	const handleInstallUsingChange = (evt: Event) => {
		const select = evt.currentTarget as HTMLSelectElement;
		installUsing = select.value;
	};

	onMount(() => {
		Prism.plugins.toolbar.registerButton(
			'fullscreen-code',
			function (env: { element: HTMLElement }) {
				const button = document.createElement('button');
				button.innerHTML = '🔍';
				button.addEventListener('click', function () {
					const parent = env.element.parentNode as HTMLElement;
					parent.requestFullscreen();
				});

				return button;
			}
		);

		Prism.highlightAll();
	});
</script>

<HeadComponent
	data={{
		title: 'Installation - Kuba',
		description: 'Install Kuba on Linux, macOS, and Windows using various methods.'
	}}
/>

<div class="container mx-auto px-4 py-8">
	<div class="max-w-4xl mx-auto">
		<div class="text-center mb-12">
			<h1 class="text-4xl font-bold mb-4">Installation Guide</h1>
			<p class="text-xl text-base-content/70">
				Get Kuba up and running on your system with these simple installation methods.
			</p>
		</div>

		<div class="grid lg:grid-cols-2 gap-8">
			<div>
				<h2 class="text-2xl font-bold mb-6">Automatic Installation</h2>

				<div class="mb-6">
					<label for="install-method" class="block text-sm font-medium mb-2">
						Choose your installation method:
					</label>
					<select
						id="install-method"
						bind:value={installUsing}
						on:change={handleInstallUsingChange}
						class="select select-bordered w-full"
					>
						<option value="manual">Manual Installation</option>
						<option value="curl-zsh">curl & zsh (Linux/macOS)</option>
						<option value="curl-bash">curl & bash (Linux/macOS)</option>
						<option value="wget-zsh">wget & zsh (Linux/macOS)</option>
						<option value="wget-bash">wget & bash (Linux/macOS)</option>
						<option value="pwsh">PowerShell (Windows)</option>
					</select>
				</div>

				<div class="space-y-4">
					<div class="card bg-base-200 {installUsing === 'curl-zsh' ? '' : 'hidden'}">
						<div class="card-body">
							<h3 class="card-title">curl & zsh (Linux/macOS)</h3>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">curl -sSL https://kuba.mwco.app/install.sh | zsh</code
								></pre>
						</div>
					</div>
					<div class="card bg-base-200 {installUsing === 'curl-bash' ? '' : 'hidden'}">
						<div class="card-body">
							<h3 class="card-title">curl & bash (Linux/macOS)</h3>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">curl -sSL https://kuba.mwco.app/install.sh | bash</code
								></pre>
						</div>
					</div>
					<div class="card bg-base-200 {installUsing === 'wget-zsh' ? '' : 'hidden'}">
						<div class="card-body">
							<h3 class="card-title">wget & zsh (Linux/macOS)</h3>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">wget -qO- https://kuba.mwco.app/install.sh | zsh</code
								></pre>
						</div>
					</div>
					<div class="card bg-base-200 {installUsing === 'wget-bash' ? '' : 'hidden'}">
						<div class="card-body">
							<h3 class="card-title">wget & bash (Linux/macOS)</h3>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">wget -qO- https://kuba.mwco.app/install.sh | bash</code
								></pre>
						</div>
					</div>
					<div class="card bg-base-200 {installUsing === 'pwsh' ? '' : 'hidden'}">
						<div class="card-body">
							<h3 class="card-title">PowerShell (Windows)</h3>
							<pre><code
									class="language-powershell"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">iwr https://kuba.mwco.app/install.ps1 -useb | iex</code
								></pre>
						</div>
					</div>
					<div class="card bg-base-200 {installUsing === 'manual' ? '' : 'hidden'}">
						<div class="card-body">
							<h3 class="card-title">Manual Installation</h3>
							<p class="mb-4">
								Download the latest release from the <a href="/download" class="link link-primary"
									>download page</a
								>.
							</p>
							<p>
								Or visit our <a
									href="https://github.com/mistweaverco/kuba/releases/latest"
									class="link link-primary">GitHub releases</a
								> page.
							</p>
						</div>
					</div>
				</div>
			</div>

			<div>
				<h2 class="text-2xl font-bold mb-6">System Requirements</h2>

				<div class="space-y-4">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Operating Systems</h3>
							<ul class="list-disc list-inside space-y-2">
								<li><strong>Linux:</strong> x86_64, ARM64, ARMv7</li>
								<li><strong>macOS:</strong> Intel, Apple Silicon (M1/M2)</li>
								<li><strong>Windows:</strong> x86_64</li>
							</ul>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Dependencies</h3>
							<p>Kuba is a single binary with no external dependencies required.</p>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Network Access</h3>
							<p>Required for downloading releases and accessing cloud provider APIs.</p>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div class="mt-12">
			<h2 class="text-2xl font-bold mb-6">Verification</h2>

			<div class="card bg-base-200">
				<div class="card-body">
					<p class="mb-4">After installation, verify that Kuba is working correctly:</p>
					<pre><code
							class="language-bash"
							data-toolbar-order="copy-to-clipboard"
							data-prismjs-copy="📋">kuba --version</code
						></pre>
					<p class="mt-4 text-sm text-base-content/70">
						You should see the current version of Kuba displayed.
					</p>
				</div>
			</div>
		</div>

		<div class="mt-12">
			<h2 class="text-2xl font-bold mb-6">Next Steps</h2>

			<div class="grid md:grid-cols-2 gap-6">
				<div class="card bg-base-200 text-primary-content">
					<div class="card-body">
						<h3 class="card-title">Configure Kuba</h3>
						<p>Set up your configuration file to start using Kuba with your cloud providers.</p>
						<a href="/configuration" class="btn btn-outline bg-lg">Configuration Guide</a>
					</div>
				</div>

				<div class="card bg-base-200 text-primary-content">
					<div class="card-body">
						<h3 class="card-title">Learn Usage</h3>
						<p>
							Discover how to use Kuba to run your applications with secure environment variables.
						</p>
						<a href="/usage" class="btn btn-outline bg-lg">Usage Guide</a>
					</div>
				</div>
			</div>
		</div>

		<div class="mt-12 text-center">
			<div class="alert alert-info">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					class="stroke-current shrink-0 w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
					></path>
				</svg>
				<span>
					<strong>Need help?</strong> Check out our
					<a href="https://github.com/mistweaverco/kuba/issues" class="link">GitHub issues</a>
					or join our <a href="https://mistweaverco.com/discord" class="link">Discord community</a>.
				</span>
			</div>
		</div>
	</div>
</div>
