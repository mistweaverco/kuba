<script lang="ts">
	import HeadComponent from '$lib/HeadComponent.svelte';
	import ClickableHeadline from '$lib/ClickableHeadline.svelte';
	import Prism from 'prismjs';
	import 'prismjs/plugins/toolbar/prism-toolbar';
	import 'prismjs/plugins/copy-to-clipboard/prism-copy-to-clipboard';
	import 'prismjs/components/prism-bash';
	import 'prismjs/components/prism-yaml';
	import 'prismjs/components/prism-json';
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
		title: 'Providers Setup - Kuba',
		description:
			'Set up authentication and permissions for GCP, AWS, Azure, and OpenBao to use with Kuba.'
	}}
/>

<div class="container mx-auto px-4 py-8">
	<div class="max-w-4xl mx-auto">
		<div class="text-center mb-12">
			<ClickableHeadline level={1} id="providers-setup" className="text-4xl font-bold mb-4"
				>Providers Setup</ClickableHeadline
			>
			<p class="text-xl text-base-content/70">
				Configure authentication and permissions for your providers to start using Kuba securely.
			</p>
		</div>

		<div class="space-y-12">
			<section>
				<ClickableHeadline level={2} id="supported-providers" className="text-3xl font-bold mb-6"
					>Supported Providers</ClickableHeadline
				>

				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					<div class="card bg-base-200 text-center">
						<div class="card-body">
							<div class="text-4xl mb-2">☁️</div>
							<a class="hover:link" href="#gcp">
								<h3 class="card-title justify-center">Google Cloud Platform</h3>
							</a>
							<p class="text-sm">
								Secret Manager integration with service accounts and workload identity
							</p>
						</div>
					</div>

					<div class="card bg-base-200 text-center">
						<div class="card-body">
							<div class="text-4xl mb-2">☁️</div>
							<a class="hover:link" href="#aws">
								<h3 class="card-title justify-center">AWS</h3>
							</a>
							<p class="text-sm">Secrets Manager with IAM roles and access keys</p>
						</div>
					</div>

					<div class="card bg-base-200 text-center">
						<div class="card-body">
							<div class="text-4xl mb-2">☁️</div>
							<a class="hover:link" href="#azure">
								<h3 class="card-title justify-center">Azure</h3>
							</a>
							<p class="text-sm">Key Vault with service principals and managed identity</p>
						</div>
					</div>

					<div class="card bg-base-200 text-center">
						<div class="card-body">
							<div class="text-4xl mb-2">☁️</div>
							<a class="hover:link" href="#openbao">
								<h3 class="card-title justify-center">OpenBao</h3>
							</a>
							<p class="text-sm">Self-hosted secrets with tokens and namespaces</p>
						</div>
					</div>

					<div class="card bg-base-200 text-center">
						<div class="card-body">
							<div class="text-4xl mb-2">🛠️</div>
							<a class="hover:link" href="#local">
								<h3 class="card-title justify-center">Local</h3>
							</a>
							<p class="text-sm">Use for hard-coded non-sensitive values during development</p>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline level={2} id="gcp" className="text-3xl font-bold mb-6"
					>Google Cloud Platform (gcp)</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="gcp-authentication-methods" className="card-title"
								>1. Authentication Methods</ClickableHeadline
							>
							<div class="space-y-4">
								<div>
									<h4 class="font-bold">Application Default Credentials</h4>
									<p class="mb-2">Use gcloud for local development:</p>
									<pre><code
											class="language-bash"
											data-toolbar-order="copy-to-clipboard"
											data-prismjs-copy="📋">gcloud auth application-default login</code
										></pre>
								</div>
								<div class="divider">OR</div>
								<div>
									<h4 class="font-bold">Service Account Key</h4>
									<p class="mb-2">
										Set the <code>GOOGLE_APPLICATION_CREDENTIALS</code> environment variable:
									</p>
									<pre><code
											class="language-bash"
											data-toolbar-order="copy-to-clipboard"
											data-prismjs-copy="📋"
											>export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"</code
										></pre>
								</div>
								<div class="divider">OR</div>
								<div>
									<h4 class="font-bold">Workload Identity</h4>
									<p class="mb-2">For GKE or other GCP services, use workload identity.</p>
								</div>
								<div class="alert alert-info">
									<i class="fa-solid fa-info-circle"></i>
									<p>
										If running on Compute Engine, the default service account will be used
										automatically.
									</p>
									<div></div>
								</div>
							</div>
						</div>

						<div class="card bg-base-200">
							<div class="card-body">
								<h3 class="card-title">2. Enable Secret Manager API</h3>
								<p class="mb-4">Make sure the Secret Manager API is enabled in your GCP project:</p>
								<pre><code
										class="language-bash"
										data-toolbar-order="copy-to-clipboard"
										data-prismjs-copy="📋">gcloud services enable secretmanager.googleapis.com</code
									></pre>
							</div>
						</div>

						<div class="card bg-base-200">
							<div class="card-body">
								<h3 class="card-title">3. IAM Permissions</h3>
								<p class="mb-4">
									Ensure your service account has the <code>Secret Manager Secret Accessor</code> role:
								</p>
								<pre><code
										class="language-bash"
										data-toolbar-order="copy-to-clipboard"
										data-prismjs-copy="📋"
										>gcloud projects add-iam-policy-binding PROJECT_ID \
    --member="serviceAccount:YOUR_SERVICE_ACCOUNT@PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"</code
									></pre>
							</div>
						</div>
						<div class="card bg-base-200">
							<div class="card-body">
								<h3 class="card-title">4. Configuration Example</h3>
								<pre><code
										class="language-yaml"
										data-toolbar-order="copy-to-clipboard"
										data-prismjs-copy="📋"
										>default:
  provider: gcp
  project: 1337
  env:
    DATABASE_URL:
      secret-key: "database-connection-string"
    API_KEY:
      secret-key: "external-api-key"</code
									></pre>
							</div>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline level={2} id="aws" className="text-3xl font-bold mb-6"
					>AWS Secrets Manager (aws)</ClickableHeadline
				>
				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="aws-authentication-methods" className="card-title"
								>1. Authentication Methods</ClickableHeadline
							>
							<div class="space-y-4">
								<div>
									<h4 class="font-bold">Environment Variables</h4>
									<p class="mb-2">Set AWS credentials:</p>
									<pre><code
											class="language-bash"
											data-toolbar-order="copy-to-clipboard"
											data-prismjs-copy="📋"
											>export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_REGION="us-east-1"</code
										></pre>
								</div>
								<div class="divider">OR</div>
								<div>
									<h4 class="font-bold">AWS Profile</h4>
									<p class="mb-2">Use a specific profile from your AWS credentials file:</p>
									<pre><code
											class="language-bash"
											data-toolbar-order="copy-to-clipboard"
											data-prismjs-copy="📋"
											>export AWS_PROFILE="my-profile"
export AWS_REGION="us-east-1"</code
										></pre>
								</div>
								<div class="divider">OR</div>
								<div>
									<h4 class="font-bold">IAM Roles</h4>
									<p class="mb-2">If running on EC2, ECS, or other AWS services, use IAM roles.</p>
								</div>
								<div class="divider">OR</div>
								<div>
									<h4 class="font-bold">AWS CLI</h4>
									<p class="mb-2">Use <code>aws configure</code> to set up your credentials.</p>
								</div>
							</div>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">2. IAM Permissions</h3>
							<p class="mb-4">
								Ensure your AWS credentials have the <code>secretsmanager:GetSecretValue</code> permission:
							</p>
							<pre><code
									class="language-json"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>&lbrace;
  "Version": "2012-10-17",
  "Statement": [
    &lbrace;
      "Effect": "Allow",
      "Action": "secretsmanager:GetSecretValue",
      "Resource": "arn:aws:secretsmanager:region:account:secret:secret-name-*"
    &rbrace;
  ]
&rbrace;</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">3. Configuration Example</h3>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>default:
  provider: aws
  env:
    DATABASE_URL:
      secret-key: "database-connection-string"
    API_KEY:
      secret-key: "external-api-key"</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline level={2} id="azure" className="text-3xl font-bold mb-6"
					>Azure Key Vault (azure)</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="azure-authentication-methods" className="card-title"
								>1. Authentication Methods</ClickableHeadline
							>
							<div class="space-y-4">
								<div>
									<h4 class="font-bold">Service Principal</h4>
									<p class="mb-2">Set the following environment variables:</p>
									<pre><code
											class="language-bash"
											data-toolbar-order="copy-to-clipboard"
											data-prismjs-copy="📋"
											>export AZURE_KEY_VAULT_URL="https://yourvault.vault.azure.net/"
export AZURE_TENANT_ID="your-tenant-id"
export AZURE_CLIENT_ID="your-client-id"
export AZURE_CLIENT_SECRET="your-client-secret"</code
										></pre>
								</div>
								<div class="divider">OR</div>
								<div>
									<h4 class="font-bold">Managed Identity</h4>
									<p class="mb-2">If running on Azure services with managed identity enabled.</p>
								</div>
								<div class="divider">OR</div>
								<div>
									<h4 class="font-bold">Default Azure Credential</h4>
									<p class="mb-2">Uses Azure CLI, Visual Studio Code, or other Azure tools.</p>
								</div>
							</div>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">2. Key Vault Permissions</h3>
							<p class="mb-4">
								Ensure your Azure credentials have the <code>Get</code> and <code>List</code> permissions
								for secrets in your Key Vault.
							</p>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">3. Configuration Example</h3>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>default:
  provider: azure
  env:
    DATABASE_URL:
      secret-key: "database-connection-string"
    API_KEY:
      secret-key: "external-api-key"</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline level={2} id="openbao" className="text-3xl font-bold mb-6"
					>OpenBao (openbao)</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">1. Setup</h3>
							<p class="mb-4">Make sure you have an OpenBao server running and accessible.</p>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline
								level={3}
								id="openbao-authentication-methods"
								className="card-title">2. Authentication Methods</ClickableHeadline
							>
							<p class="mb-4">Set up authentication using environment variables:</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>export OPENBAO_ADDR="http://localhost:8200"  # Required: OpenBao server address
export OPENBAO_TOKEN="your-openbao-token"    # Optional: Authentication token
export OPENBAO_NAMESPACE="your-namespace"     # Optional: Namespace (if using enterprise features)</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">3. Permissions</h3>
							<p class="mb-4">
								Ensure your OpenBao token has read permissions for the secrets you want to access.
							</p>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">4. Configuration Example</h3>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>default:
  provider: openbao
  env:
    DATABASE_URL:
      secret-key: "secret/database-url"
    API_KEY:
      secret-key: "secret/api-key"</code
								></pre>
							<p class="mt-4 text-sm">
								<strong>Note:</strong> OpenBao secrets are stored as key-value pairs. If a secret contains
								multiple keys, Kuba will return the first string value it finds.
							</p>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline level={2} id="local" className="text-3xl font-bold mb-6"
					>Local (local)</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="alert alert-error text-xl">
						<i class="fa-solid fa-triangle-exclamation"></i>
						<span>
							Make sure to use the <code>local</code> provider only for non-sensitive values.
						</span>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">2. Configuration Example</h3>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>default:
  provider: local
  env:
    LOG_ENV:
      value: "local"
    LOG_LEVEL:
      value: "debug"</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Multi-Provider Configuration</h2>

				<div class="card bg-base-200">
					<div class="card-body">
						<h3 class="card-title">Using Multiple Providers</h3>
						<p class="mb-4">You can use different cloud providers in the same configuration:</p>
						<pre><code
								class="language-yaml"
								data-toolbar-order="copy-to-clipboard"
								data-prismjs-copy="📋"
								>default:
  provider: gcp
  project: 1337
  env:
    GCP_SECRETS:
      secret-path: "app-config"
      provider: gcp
    AWS_SECRETS:
      secret-path: "app-config"
      provider: aws
    AZURE_SECRETS:
      secret-path: "app-config"
      provider: azure
      project: "my-azure-project"
    OPENBAO_SECRETS:
      secret-path: "app-config"
      provider: openbao</code
							></pre>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Security Best Practices</h2>

				<div class="grid md:grid-cols-2 gap-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Authentication</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Use service accounts with minimal permissions</li>
								<li>Rotate credentials regularly</li>
								<li>Use managed identities when possible</li>
								<li>Avoid hardcoding credentials</li>
							</ul>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Permissions</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Follow principle of least privilege</li>
								<li>Use role-based access control</li>
								<li>Limit access to production secrets</li>
								<li>Monitor access patterns</li>
							</ul>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Network Security</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Use private networks when possible</li>
								<li>Enable VPC endpoints for AWS</li>
								<li>Use private service connect for GCP</li>
								<li>Restrict access by IP when applicable</li>
							</ul>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Monitoring</h3>
							<ul class="list-disc list-inside space-y-2">
								<li>Enable audit logging</li>
								<li>Set up alerts for unusual access</li>
								<li>Monitor secret rotation</li>
								<li>Track usage patterns</li>
							</ul>
						</div>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Troubleshooting</h2>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Common Issues</h3>
							<div class="space-y-4">
								<div>
									<h4 class="font-bold">Authentication Errors</h4>
									<p class="text-sm">
										Check your credentials and ensure they haven't expired. Verify the
										authentication method you're using.
									</p>
								</div>
								<div>
									<h4 class="font-bold">Permission Errors</h4>
									<p class="text-sm">
										Ensure your credentials have the necessary permissions to access the secrets
										specified in your configuration.
									</p>
								</div>
								<div>
									<h4 class="font-bold">Network Issues</h4>
									<p class="text-sm">
										Check your network connectivity and firewall settings. Ensure you can reach the
										cloud provider APIs.
									</p>
								</div>
							</div>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Debug Mode</h3>
							<p class="mb-4">
								Enable debug mode to see detailed information about authentication and API calls:
							</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run --debug -- echo "Testing connection"</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Next Steps</h2>

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
							<h3 class="card-title">Usage Examples</h3>
							<p>See practical examples of how to use your configured providers.</p>
							<a href="/examples" class="btn btn-outline bg-lg">Examples Guide</a>
						</div>
					</div>
				</div>
			</section>
		</div>
	</div>
</div>
