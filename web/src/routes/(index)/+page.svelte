<script lang="ts">
	import { onMount } from 'svelte';
	import HeadComponent from '$lib/HeadComponent.svelte';
</script>

<HeadComponent
	data={{
		title: 'Kuba - Securely and easily access your environment variables',
		description:
			'Pass env directly from GCP Secret Manager, AWS Secrets Manager, Azure Key Vault and OpenBao to your application.'
	}}
/>

<div class="hero bg-base-200 min-h-screen">
	<div class="hero-content text-center">
		<div class="max-w-4xl">
			<img src="/logo.svg" alt="Kuba" class="m-5 mx-auto w-64" />
			<h1 class="text-5xl font-bold">Kuba</h1>
			<p class="py-6 text-xl">Securely and easily access your environment variables</p>
			<p class="py-4 text-lg">
				Kuba is <a href="https://en.wikipedia.org/wiki/Swahili_language" class="link link-primary"
					>Swahili</a
				> for "vault."
			</p>
			<p class="py-4">
				Kuba helps you to get rid of <code>.env</code> files. Pass env directly from GCP Secret Manager,
				AWS Secrets Manager, Azure Key Vault, and OpenBao to your application.
			</p>
			<div class="flex gap-4 justify-center mt-8">
				<a href="/installation" class="btn btn-primary">Get Started</a>
				<a href="/usage" class="btn btn-secondary">Learn More</a>
			</div>
		</div>
	</div>
</div>

<div class="py-16 bg-base-100">
	<div class="container mx-auto px-4">
		<h2 class="text-3xl font-bold text-center mb-12">Why Choose Kuba?</h2>
		<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h3 class="card-title">🔒 Secure</h3>
					<p>
						Eliminate insecure <code>.env</code> files and leverage your existing cloud provider security
					</p>
				</div>
			</div>
			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h3 class="card-title">🚀 Simple</h3>
					<p>Single binary with a simple YAML configuration file</p>
				</div>
			</div>
			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h3 class="card-title">☁️ Multi-Cloud</h3>
					<p>Support for GCP, AWS, Azure, and OpenBao all in one tool</p>
				</div>
			</div>
			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h3 class="card-title">💰 Free</h3>
					<p>No subscription fees - leverage your existing cloud infrastructure</p>
				</div>
			</div>
			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h3 class="card-title">⚡ Fast</h3>
					<p>Single binary written in Go for optimal performance</p>
				</div>
			</div>
			<div class="card bg-base-200 shadow-xl">
				<div class="card-body">
					<h3 class="card-title">🔧 Flexible</h3>
					<p>Environment-specific configurations with variable interpolation</p>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="py-16 bg-base-200">
	<div class="container mx-auto px-4">
		<h2 class="text-3xl font-bold text-center mb-12">Quick Start</h2>
		<div class="grid md:grid-cols-3 gap-8">
			<div class="text-center">
				<div class="text-4xl mb-4">1️⃣</div>
				<h3 class="text-xl font-bold mb-4">Install</h3>
				<p class="mb-4">Download and install Kuba on your system</p>
				<a href="/installation" class="btn btn-primary">Installation Guide</a>
			</div>
			<div class="text-center">
				<div class="text-4xl mb-4">2️⃣</div>
				<h3 class="text-xl font-bold mb-4">Configure</h3>
				<p class="mb-4">Set up your <code>kuba.yaml</code> configuration file</p>
				<a href="/configuration" class="btn btn-primary">Configuration Guide</a>
			</div>
			<div class="text-center">
				<div class="text-4xl mb-4">3️⃣</div>
				<h3 class="text-xl font-bold mb-4">Run</h3>
				<p class="mb-4">Execute your applications with secure environment variables</p>
				<a href="/usage" class="btn btn-primary">Usage Guide</a>
			</div>
		</div>
	</div>
</div>

<div class="py-16 bg-base-100">
	<div class="container mx-auto px-4 text-center">
		<h2 class="text-3xl font-bold mb-8">Ready to Get Started?</h2>
		<p class="text-lg mb-8">
			Join the community and start managing your secrets securely with Kuba.
		</p>
		<div class="flex gap-4 justify-center">
			<a href="/installation" class="btn btn-primary btn-lg">Install Kuba</a>
			<a href="https://github.com/mistweaverco/kuba" class="btn btn-outline btn-lg"
				>View on GitHub</a
			>
		</div>
	</div>
</div>
