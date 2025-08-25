<script lang="ts">
	import HeadComponent from '$lib/HeadComponent.svelte';
	import Prism from 'prismjs';
	import 'prismjs/plugins/toolbar/prism-toolbar';
	import 'prismjs/plugins/copy-to-clipboard/prism-copy-to-clipboard';
	import 'prismjs/components/prism-yaml';
	import 'prismjs/components/prism-bash';
	import 'prismjs/themes/prism-okaidia.css';
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
		title: 'Configuration Guide - Kuba',
		description:
			'Learn how to configure Kuba with kuba.yaml, environment variable interpolation, and secret path mapping.'
	}}
/>

<div class="container mx-auto px-4 py-8">
	<div class="max-w-4xl mx-auto">
		<div class="text-center mb-12">
			<h1 class="text-4xl font-bold mb-4">Configuration Guide</h1>
			<p class="text-xl text-base-content/70">
				Learn how to configure Kuba with the <code>kuba.yaml</code> file and understand advanced features
				like variable interpolation and secret paths.
			</p>
		</div>

		<div class="space-y-12">
			<section>
				<h2 class="text-3xl font-bold mb-6">Getting Started</h2>

				<div class="card bg-base-200 mb-6">
					<div class="card-body">
						<h3 class="card-title">Initialize Configuration</h3>
						<p class="mb-4">
							Start by creating a configuration file using the <code>kuba init</code> command:
						</p>
						<pre><code
								class="language-bash"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋">kuba init</code
							></pre>
						<p class="mt-4">
							This will generate a default <code>kuba.yaml</code> file that you can customize for your
							needs.
						</p>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Configuration File Structure</h2>

				<div class="card bg-base-200 mb-6">
					<div class="card-body">
						<h3 class="card-title">Basic Structure</h3>
						<p class="mb-4">
							The <code>kuba.yaml</code> file is organized into environment sections, each with its own
							provider and mappings:
						</p>
						<pre><code
								class="language-yaml"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋"
								># yaml-language-server: $schema=https://kuba.mwco.app/kuba.schema.json
---
default:
  provider: gcp
  project: 1337
  mappings:
    - environment-variable: "DATABASE_URL"
      secret-key: "database-connection-string"
    - environment-variable: "API_KEY"
      secret-key: "external-api-key"

development:
  provider: gcp
  project: 1337
  mappings:
    - environment-variable: "DEV_DATABASE_URL"
      secret-key: "dev-database-connection-string"

production:
  provider: gcp
  project: 1337
  mappings:
    - environment-variable: "PROD_DATABASE_URL"
      secret-key: "prod-database-connection-string"</code
							></pre>
					</div>
				</div>

				<div class="grid md:grid-cols-2 gap-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Environment Sections</h3>
							<p>
								Each top-level section (like <code>default</code>, <code>development</code>,
								<code>production</code>) represents a different environment configuration.
							</p>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Provider Configuration</h3>
							<p>
								The <code>provider</code> field specifies which cloud provider to use (gcp, aws, azure,
								openbao).
							</p>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Project ID</h3>
							<p>
								The <code>project</code> field specifies the project ID for the cloud provider (required
								for GCP and Azure).
							</p>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Mappings</h3>
							<p>
								The <code>mappings</code> array defines how secrets are mapped to environment variables.
							</p>
						</div>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Mapping Types</h2>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Individual Secrets (secret-key)</h3>
							<p class="mb-4">Fetch a single secret from your cloud provider:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>mappings:
  - environment-variable: "DATABASE_URL"
    secret-key: "database-connection-string"
  - environment-variable: "API_KEY"
    secret-key: "external-api-key"</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Secret Paths (secret-path)</h3>
							<p class="mb-4">Fetch all secrets under a specific path prefix:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>mappings:
  - environment-variable: "DB"
    secret-path: "database"
  - environment-variable: "API"
    secret-path: "external-apis"</code
								></pre>
							<p class="mt-4 text-sm">
								This will create environment variables like <code>DB_CONNECTION_STRING</code>,
								<code>DB_USERNAME</code>, <code>API_STRIPE_KEY</code>, etc.
							</p>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Hard-coded Values (value)</h3>
							<p class="mb-4">Set static environment variables:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>mappings:
  - environment-variable: "APP_ENV"
    value: "production"
  - environment-variable: "DEBUG"
    value: "false"</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Environment Variable Interpolation</h2>

				<div class="card bg-base-200 mb-6">
					<div class="card-body">
						<h3 class="card-title">Basic Interpolation</h3>
						<p class="mb-4">
							Kuba supports environment variable interpolation using <code
								>$&lbrace;VAR_NAME&rbrace;</code
							> syntax:
						</p>
						<pre><code
								class="language-yaml"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋"
								>mappings:
  - environment-variable: "DB_PASSWORD"
    secret-key: "db-password"
  - environment-variable: "DB_HOST"
    value: "mydbhost"
  - environment-variable: "DB_CONNECTION_STRING"
    value: "postgresql://user:$&lbrace;DB_PASSWORD&rbrace;@$&lbrace;DB_HOST&rbrace;:5432/mydb"</code
							></pre>
					</div>
				</div>

				<div class="grid md:grid-cols-2 gap-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">System Environment Variables</h3>
							<p>Reference system environment variables:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>- environment-variable: "API_URL"
  value: "https://api.$&lbrace;DOMAIN&rbrace;/v1"</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Default Values</h3>
							<p>Provide fallback values with <code>$&lbrace;VAR:-default&rbrace;</code> syntax:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>- environment-variable: "REDIS_URL"
  value: "redis://$&lbrace;REDIS_HOST:-localhost&rbrace;:$&lbrace;REDIS_PORT:-6379&rbrace;/0"</code
								></pre>
						</div>
					</div>
				</div>

				<div class="alert alert-info mt-6">
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
						<strong>Important:</strong> Interpolation is processed in order, so you can reference variables
						defined earlier in the same configuration.
					</span>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Cross-Provider Mappings</h2>

				<div class="card bg-base-200 mb-6">
					<div class="card-body">
						<h3 class="card-title">Multiple Cloud Providers</h3>
						<p class="mb-4">
							You can fetch secrets from different cloud providers in the same configuration:
						</p>
						<pre><code
								class="language-yaml"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋"
								>default:
  provider: gcp
  project: 1337
  mappings:
    - environment-variable: "GCP_PROJECT_ID"
      secret-key: "gcp_project_secret"
    - environment-variable: "AWS_PROJECT_ID"
      secret-key: "aws_project_secret"
      provider: aws
    - environment-variable: "AZURE_PROJECT_ID"
      secret-key: "azure_project_secret"
      provider: azure
      project: "my-azure-project"</code
							></pre>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Complete Example</h2>

				<div class="card bg-base-200">
					<div class="card-body">
						<h3 class="card-title">Full Configuration Example</h3>
						<p class="mb-4">Here's a comprehensive example showing all features:</p>
						<pre><code
								class="language-yaml"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋"
								># yaml-language-server: $schema=https://kuba.mwco.app/kuba.schema.json
---
default:
  provider: gcp
  project: 1337
  mappings:
    # Individual secrets
    - environment-variable: "DATABASE_URL"
      secret-key: "database-connection-string"
    - environment-variable: "STRIPE_API_KEY"
      secret-key: "stripe-api-key"
    
    # Secret paths for bulk loading
    - environment-variable: "DB"
      secret-path: "database"
    - environment-variable: "API"
      secret-path: "external-apis"
    
    # Hard-coded values
    - environment-variable: "APP_ENV"
      value: "development"
    - environment-variable: "DEBUG"
      value: "true"
    
    # Interpolated values
    - environment-variable: "REDIS_URL"
      value: "redis://$&lbrace;REDIS_HOST:-localhost$&rbrace;:$&lbrace;REDIS_PORT:-6379&rbrace;/0"
    - environment-variable: "LOG_LEVEL"
      value: "$&lbrace;LOG_LEVEL:-info&rbrace;"

development:
  provider: gcp
  project: 1337
  mappings:
    - environment-variable: "DEV_DATABASE_URL"
      secret-key: "dev-database-connection-string"
    - environment-variable: "DEV_STRIPE_API_KEY"
      secret-key: "dev-stripe-api-key"

staging:
  provider: gcp
  project: 1337
  mappings:
    - environment-variable: "STAGING_DATABASE_URL"
      secret-key: "staging-database-connection-string"
    - environment-variable: "STAGING_STRIPE_API_KEY"
      secret-key: "staging-stripe-api-key"

production:
  provider: gcp
  project: 1337
  mappings:
    - environment-variable: "PROD_DATABASE_URL"
      secret-key: "prod-database-connection-string"
    - environment-variable: "PROD_STRIPE_API_KEY"
      secret-key: "prod-stripe-api-key"
    - environment-variable: "APP_ENV"
      value: "production"
    - environment-variable: "DEBUG"
      value: "false"</code
							></pre>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Best Practices</h2>

				<div class="grid md:grid-cols-2 gap-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Organization</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Use descriptive environment variable names</li>
								<li>Group related secrets with secret paths</li>
								<li>Keep environment-specific overrides minimal</li>
								<li>Document your configuration structure</li>
							</ul>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Security</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Never commit secrets to version control</li>
								<li>Use environment-specific configurations</li>
								<li>Limit access to production secrets</li>
								<li>Rotate secrets regularly</li>
							</ul>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Maintenance</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Keep configurations in sync across teams</li>
								<li>Use consistent naming conventions</li>
								<li>Test configurations in staging first</li>
								<li>Version control your configuration templates</li>
							</ul>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Performance</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Use secret paths for bulk operations</li>
								<li>Avoid unnecessary cross-provider calls</li>
								<li>Cache configurations when possible</li>
								<li>Monitor secret access patterns</li>
							</ul>
						</div>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Next Steps</h2>

				<div class="grid md:grid-cols-2 gap-6">
					<div class="card bg-primary text-primary-content">
						<div class="card-body">
							<h3 class="card-title">Cloud Providers Setup</h3>
							<p>Configure authentication and permissions for your cloud providers.</p>
							<a href="/providers" class="btn btn-secondary">Cloud Providers Guide</a>
						</div>
					</div>

					<div class="card bg-secondary text-secondary-content">
						<div class="card-body">
							<h3 class="card-title">Usage Examples</h3>
							<p>See practical examples of how to use your configuration.</p>
							<a href="/examples" class="btn btn-primary">Examples Guide</a>
						</div>
					</div>
				</div>
			</section>
		</div>
	</div>
</div>
