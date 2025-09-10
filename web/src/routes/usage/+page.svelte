<script lang="ts">
	import HeadComponent from '$lib/HeadComponent.svelte';
	import ClickableHeadline from '$lib/ClickableHeadline.svelte';
	import Prism from 'prismjs';
	import 'prismjs/plugins/toolbar/prism-toolbar';
	import 'prismjs/plugins/copy-to-clipboard/prism-copy-to-clipboard';
	import 'prismjs/components/prism-bash';
	import 'prismjs/components/prism-yaml';
	import 'dracula-prism/dist/css/dracula-prism.css';
	import { onMount } from 'svelte';

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
		title: 'Usage Guide - Kuba',
		description:
			'Learn how to use Kuba to run applications with secure environment variables from cloud providers.'
	}}
/>

<div class="container mx-auto px-4 py-8">
	<div class="max-w-4xl mx-auto">
		<div class="text-center mb-12">
			<ClickableHeadline level={1} id="usage-guide" className="text-4xl font-bold mb-4"
				>Usage Guide</ClickableHeadline
			>
			<p class="text-xl text-base-content/70">
				Learn how to use Kuba to securely run your applications with environment variables from
				cloud providers.
			</p>
		</div>

		<div class="space-y-12">
			<section>
				<ClickableHeadline level={2} id="basic-usage" className="text-3xl font-bold mb-6"
					>Basic Usage</ClickableHeadline
				>

				<div class="card bg-base-200 mb-6">
					<div class="card-body">
						<h3 class="card-title">Running Commands with Kuba</h3>
						<p class="mb-4">The basic syntax for using Kuba is:</p>
						<pre><code
								class="language-bash"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋">kuba run -- &lt;your-command&gt;</code
							></pre>
						<p class="mt-4">
							This will fetch all secrets defined in your <code>kuba.yaml</code> file and pass them as
							environment variables to your command.
						</p>
					</div>
				</div>

				<div class="grid md:grid-cols-2 gap-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Node.js Application</h3>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run -- node dist/server.js</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Python Application</h3>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run -- python app.py</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Docker Container</h3>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run -- docker run myapp</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Shell Script</h3>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run -- ./deploy.sh</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline level={2} id="testing-configuration" className="text-3xl font-bold mb-6"
					>Testing Configuration</ClickableHeadline
				>

				<div class="card bg-base-200 mb-6">
					<div class="card-body">
						<h3 class="card-title">Validate Access and Mappings</h3>
						<p class="mb-4">
							Use the <code>test</code> subcommand to verify that Kuba can load your configuration and
							retrieve all mapped values for an environment without executing a program.
						</p>
						<pre><code
								class="language-bash"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋"
								># Use default environment
kuba test

# Specify an environment
kuba test --env staging

# Point to a specific configuration file
kuba test --config ./config/kuba.yaml --env production</code
							></pre>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline
					level={2}
					id="environment-specific-usage"
					className="text-3xl font-bold mb-6">Environment-Specific Usage</ClickableHeadline
				>

				<div class="card bg-base-200 mb-6">
					<div class="card-body">
						<h3 class="card-title">Specifying Environments</h3>
						<p class="mb-4">
							You can specify which environment configuration to use with the <code>--env</code> flag:
						</p>
						<pre><code
								class="language-bash"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋">kuba run --env development -- node app.js</code
							></pre>
						<pre><code
								class="language-bash"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋">kuba run --env staging -- python app.py</code
							></pre>
						<pre><code
								class="language-bash"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋">kuba run --env production -- docker run myapp</code
							></pre>
					</div>
				</div>

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
						If no environment is specified, Kuba will use the <code>default</code> environment from your
						configuration.
					</span>
				</div>
			</section>

			<section>
				<ClickableHeadline
					level={2}
					id="advanced-usage-patterns"
					className="text-3xl font-bold mb-6">Advanced Usage Patterns</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Development Workflow</h3>
							<p class="mb-4">
								Use Kuba during development to avoid managing local <code>.env</code> files:
							</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									># Start development server with secrets
kuba run --env development -- npm run dev

# Run tests with test environment secrets
kuba run --env testing -- npm test

# Run database migrations
kuba run --env development -- npm run migrate</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">CI/CD Integration</h3>
							<p class="mb-4">Integrate Kuba into your CI/CD pipelines:</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									># Build and test with staging secrets
kuba run --env staging -- npm run build
kuba run --env staging -- npm test

# Deploy with production secrets
kuba run --env production -- docker build -t myapp .
kuba run --env production -- docker push myapp</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Docker Integration</h3>
							<p class="mb-4">Use Kuba with Docker containers:</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									># Run container with secrets as environment variables
kuba run -- docker run -e DATABASE_URL -e API_KEY myapp

# Build container with secrets available during build
kuba run -- docker build --build-arg DATABASE_URL --build-arg API_KEY .</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline level={2} id="troubleshooting" className="text-3xl font-bold mb-6"
					>Troubleshooting</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Common Issues</h3>
							<div class="space-y-4">
								<div>
									<h4 class="font-bold">Authentication Errors</h4>
									<p class="text-sm">
										Ensure your cloud provider credentials are properly configured. Check the <a
											href="/providers"
											class="link link-primary">Cloud Providers</a
										> guide for setup instructions.
									</p>
								</div>
								<div>
									<h4 class="font-bold">Configuration Errors</h4>
									<p class="text-sm">
										Validate your <code>kuba.yaml</code> file. Use <code>kuba init</code> to generate
										a valid template.
									</p>
								</div>
								<div>
									<h4 class="font-bold">Permission Errors</h4>
									<p class="text-sm">
										Ensure your credentials have the necessary permissions to access the secrets
										specified in your configuration.
									</p>
								</div>
							</div>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Debug Mode</h3>
							<p class="mb-4">
								Enable debug mode to see detailed information about what Kuba is doing:
							</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run --debug -- node app.js</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline level={2} id="best-practices" className="text-3xl font-bold mb-6"
					>Best Practices</ClickableHeadline
				>

				<div class="grid md:grid-cols-2 gap-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Security</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Never commit secrets to version control</li>
								<li>Use environment-specific configurations</li>
								<li>Rotate secrets regularly</li>
								<li>Limit access to production secrets</li>
							</ul>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Configuration</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Use descriptive environment variable names</li>
								<li>Group related secrets with secret paths</li>
								<li>Leverage variable interpolation</li>
								<li>Document your configuration structure</li>
							</ul>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Deployment</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Test configurations in staging first</li>
								<li>Use CI/CD for consistent deployments</li>
								<li>Monitor secret access and usage</li>
								<li>Have a rollback strategy</li>
							</ul>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Development</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Use local development environments</li>
								<li>Share configuration templates, not secrets</li>
								<li>Test with different cloud providers</li>
								<li>Keep configurations in sync across teams</li>
							</ul>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline level={2} id="next-steps" className="text-3xl font-bold mb-6"
					>Next Steps</ClickableHeadline
				>

				<div class="grid md:grid-cols-2 gap-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Configuration Guide</h3>
							<p>Learn how to set up your <code>kuba.yaml</code> configuration file.</p>
							<a href="/configuration" class="btn btn-outline bg-lg">Configuration Guide</a>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Cloud Providers</h3>
							<p>Set up authentication and permissions for your cloud providers.</p>
							<a href="/providers" class="btn btn-outline btn-lg">Cloud Providers Guide</a>
						</div>
					</div>
				</div>
			</section>
		</div>
	</div>
</div>
