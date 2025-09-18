<script lang="ts">
	import HeadComponent from '$lib/HeadComponent.svelte';
	import ClickableHeadline from '$lib/ClickableHeadline.svelte';
	import Prism from 'prismjs';
	import 'prismjs/plugins/toolbar/prism-toolbar';
	import 'prismjs/plugins/copy-to-clipboard/prism-copy-to-clipboard';
	import 'prismjs/components/prism-yaml';
	import 'prismjs/components/prism-bash';
	import 'prismjs/components/prism-javascript';
	import 'prismjs/components/prism-python';
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
		title: 'Examples - Kuba',
		description:
			'Practical examples and use cases for using Kuba with different applications and frameworks.'
	}}
/>

<div class="container mx-auto px-4 py-8">
	<div class="max-w-4xl mx-auto">
		<div class="text-center mb-12">
			<ClickableHeadline level={1} id="examples-and-use-cases" className="text-4xl font-bold mb-4"
				>Examples &amp; Use Cases</ClickableHeadline
			>
			<p class="text-xl text-base-content/70">
				See practical examples of how to use Kuba with different applications, frameworks, and
				deployment scenarios.
			</p>
		</div>

		<div class="space-y-12">
			<section>
				<ClickableHeadline
					level={2}
					id="web-application-examples"
					className="text-3xl font-bold mb-6">Web Application Examples</ClickableHeadline
				>
				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="nodejs-express-application" className="card-title"
								>Node.js Express Application</ClickableHeadline
							>
							<p class="mb-4">
								Run a Node.js Express application with database credentials and API keys:
							</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run --env production -- node app.js</code
								></pre>

							<ClickableHeadline
								level={4}
								id="nodejs-kuba-configuration"
								className="font-bold mt-4 mb-2 text-left"
								>Configuration (kuba.yaml):</ClickableHeadline
							>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>production:
  provider: gcp
  project: 1337
  env:
    DATABASE_URL:
      secret-key: "prod-database-url"
    JWT_SECRET:
      secret-key: "jwt-secret"
    STRIPE_SECRET_KEY:
      secret-key: "stripe-secret-key"
    REDIS_URL:
      value: "redis://$&lbrace;REDIS_HOST:-localhost&rbrace;:6379"</code
								></pre>

							<ClickableHeadline
								level={4}
								id="nodejs-application-code"
								className="font-bold mt-4 mb-2 text-left">Application Code:</ClickableHeadline
							>
							<pre><code
									class="language-javascript"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>const express = require('express');
const app = express();

// Environment variables are automatically available
const dbUrl = process.env.DATABASE_URL;
const jwtSecret = process.env.JWT_SECRET;
const stripeKey = process.env.STRIPE_SECRET_KEY;
const redisUrl = process.env.REDIS_URL;

console.log('Database URL:', dbUrl);
console.log('Redis URL:', redisUrl);

app.listen(3000, () => &lbrace;
  console.log('Server running on port 3000');
&rbrace;);</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Python Flask Application</h3>
							<p class="mb-4">
								Run a Python Flask application with environment-specific configurations:
							</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run --env development -- python app.py</code
								></pre>

							<ClickableHeadline
								level={4}
								id="python-kuba-configuration"
								className="font-bold mt-4 mb-2 text-left"
								>Configuration (kuba.yaml):</ClickableHeadline
							>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>development:
  provider: aws
  env:
    FLASK_ENV:
      value: "development"
    DATABASE_URL:
      secret-key: "dev-database-url"
    SECRET_KEY:
      secret-key: "flask-secret-key"
    DEBUG:
      value: "true"</code
								></pre>

							<ClickableHeadline
								level={4}
								id="python-application-code"
								className="font-bold mt-4 mb-2 text-left">Application Code:</ClickableHeadline
							>
							<pre><code
									class="language-python"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>from flask import Flask
import os

app = Flask(__name__)

# Environment variables are automatically available
app.config['DATABASE_URL'] = os.environ.get('DATABASE_URL')
app.config['SECRET_KEY'] = os.environ.get('SECRET_KEY')
app.config['DEBUG'] = os.environ.get('DEBUG', 'false').lower() == 'true'

print(f"Database URL: &lbrace;app.config['DATABASE_URL']&rbrace;")
print(f"Debug mode: &lbrace;app.config['DEBUG']&rbrace;")

if __name__ == '__main__':
    app.run(debug=app.config['DEBUG'])</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline
					level={2}
					id="database-and-api-examples"
					className="text-3xl font-bold mb-6">Database Migrations</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="database-migrations" className="card-title"
								>Database Migrations</ClickableHeadline
							>
							<p class="mb-4">Run database migrations with production credentials:</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									># Run migrations with production database credentials
kuba run --env production -- npm run migrate

# Run seed data with development database
kuba run --env development -- npm run seed</code
								></pre>

							<ClickableHeadline
								level={4}
								id="database-migrations-kuba-configuraton"
								className="font-bold mt-4 mb-2 text-left"
								>Configuration (kuba.yaml):</ClickableHeadline
							>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>production:
  provider: gcp
  project: 1337
  env:
    DATABASE_URL:
      secret-key: "prod-postgres-url"
    DB_PASSWORD:
      secret-key: "prod-db-password"

development:
  provider: aws
  env:
    DATABASE_URL:
      secret-key: "dev-postgres-url"
    DB_PASSWORD:
      secret-key: "dev-db-password"</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="external-api-integration" className="card-title"
								>External API Integration</ClickableHeadline
							>
							<p class="mb-4">Connect to external APIs with secure keys:</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run --env staging -- python api_client.py</code
								></pre>

							<ClickableHeadline
								level={4}
								id="external-api-integration-kuba-configuration"
								className="font-bold mt-4 mb-2 text-left"
								>Configuration (kuba.yaml):</ClickableHeadline
							>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>staging:
  provider: azure
  env:
    STRIPE_API_KEY:
      secret-key: "stripe-staging-key"
    SENDGRID_API_KEY:
      secret-key: "sendgrid-staging-key"
    TWILIO_ACCOUNT_SID:
      secret-key: "twilio-account-sid"
    TWILIO_AUTH_TOKEN:
      secret-key: "twilio-auth-token"</code
								></pre>

							<ClickableHeadline
								level={4}
								id="external-api-integration-api-client-code"
								className="font-bold mt-4 mb-2 text-left">API Client Code:</ClickableHeadline
							>
							<pre><code
									class="language-python"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>import os
import stripe
import sendgrid
from twilio.rest import Client

# API keys are automatically available
stripe.api_key = os.environ.get('STRIPE_API_KEY')
sendgrid_client = sendgrid.SendGridAPIClient(
    api_key=os.environ.get('SENDGRID_API_KEY')
)
twilio_client = Client(
    os.environ.get('TWILIO_ACCOUNT_SID'),
    os.environ.get('TWILIO_AUTH_TOKEN')
)

print("Stripe API key configured:", bool(stripe.api_key))
print("SendGrid API key configured:", bool(os.environ.get('SENDGRID_API_KEY')))
print("Twilio credentials configured:", bool(os.environ.get('TWILIO_ACCOUNT_SID')))</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline
					level={2}
					id="docker-and-container-examples"
					className="text-3xl font-bold mb-6">Docker &amp; Container Examples</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="docker-container-with-secrets" className="card-title"
								>Docker Container with Secrets</ClickableHeadline
							>
							<p class="mb-4">Run Docker containers with environment variables from Kuba:</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									># Build image with secrets available during build
kuba run --env production -- docker build \
  --build-arg DATABASE_URL \
  --build-arg API_KEY \
  -t myapp .

# Run container with secrets as environment variables
kuba run --env production -- docker run \
  -e DATABASE_URL \
  -e API_KEY \
  -e REDIS_URL \
  -p 3000:3000 \
  myapp

# Use --contain to avoid inheriting host environment
docker run --env-file=&lt;(kuba run --env production --contain -- env) myapp

# or pass full host environment including Kuba-managed vars
docker run --env-file=&lt;(kuba run --env production -- env) myapp
							</code></pre>

							<ClickableHeadline level={4} id="dockerfile" className="font-bold mt-4 mb-2 text-left"
								>Dockerfile:</ClickableHeadline
							>
							<pre><code
									class="language-dockerfile"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY . .

# Build arguments for secrets
ARG DATABASE_URL
ARG API_KEY

# Set environment variables
ENV DATABASE_URL=$DATABASE_URL
ENV API_KEY=$API_KEY

EXPOSE 3000

CMD ["npm", "start"]</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="docker-compose-integration" className="card-title"
								>Docker Compose Integration</ClickableHeadline
							>
							<p class="mb-4">Use Kuba with Docker Compose for multi-service applications:</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									># Start all services with production secrets
kuba run --env production -- docker-compose up -d

# Start specific service with development secrets
kuba run --env development -- docker-compose up web</code
								></pre>

							<ClickableHeadline
								level={4}
								id="docker-compose-yaml"
								className="font-bold mt-4 mb-2 text-left">docker-compose.yml:</ClickableHeadline
							>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>version: '3.8'
services:
  web:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL
      - API_KEY
      - REDIS_URL
    depends_on:
      - db
      - redis

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=myapp
      - POSTGRES_USER=myapp
      - POSTGRES_PASSWORD
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline
					level={2}
					id="ci-cd-pipeline-examples"
					className="text-3xl font-bold mb-6">CI/CD Pipeline Examples</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="github-actions" className="card-title"
								>GitHub Actions</ClickableHeadline
							>
							<p class="mb-4">Integrate Kuba into GitHub Actions workflows:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install Kuba
        run: |
          curl -sSL https://kuba.mwco.app/install.sh | bash

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: $&lbrace;&lbrace; secrets.AWS_ACCESS_KEY_ID &rbrace;&rbrace;
          aws-secret-access-key: $&lbrace;&lbrace; secrets.AWS_SECRET_ACCESS_KEY &rbrace;&rbrace;
          aws-region: us-east-1

      - name: Build and deploy
        run: |
          kuba run --env production -- npm run build
          kuba run --env production -- npm run deploy</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="gitlab-ci" className="card-title"
								>GitLab CI</ClickableHeadline
							>
							<p class="mb-4">Use Kuba in GitLab CI/CD pipelines:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>stages:
  - test
  - deploy

variables:
  KUBE_CONFIG_FILE: $CI_PROJECT_DIR/kuba.yaml

test:
  stage: test
  image: node:18
  before_script:
    - curl -sSL https://kuba.mwco.app/install.sh | bash
  script:
    - kuba run --env testing -- npm test
  only:
    - merge_requests

deploy:
  stage: deploy
  image: node:18
  before_script:
    - curl -sSL https://kuba.mwco.app/install.sh | bash
  script:
    - kuba run --env production -- npm run deploy
  only:
    - main</code
								></pre>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline
					level={3}
					id="development-workflow-examples"
					className="text-3xl font-bold mb-6">Development Workflow Examples</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="local-development" className="card-title"
								>Local Development</ClickableHeadline
							>
							<p class="mb-4">
								Use Kuba for local development without managing <code>.env</code> files:
							</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									># Start development server
kuba run --env development -- npm run dev

# Run tests
kuba run --env testing -- npm test

# Run database migrations
kuba run --env development -- npm run migrate

# Start background services
kuba run --env development -- npm run start:services</code
								></pre>

							<ClickableHeadline
								level={4}
								id="local-development-kuba-configuration"
								className="font-bold mt-4 mb-2 text-left"
								>Configuration (kuba.yaml):</ClickableHeadline
							>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>development:
  provider: gcp
  project: 1337
  env:
    DATABASE_URL:
      secret-key: "dev-database-url"
    API_KEY:
      secret-key: "dev-api-key"
    DEBUG:
      value: "true"
    LOG_LEVEL:
      value: "debug"

testing:
  provider: gcp
  project: 1337
  env:
    DATABASE_URL:
      secret-key: "test-database-url"
    API_KEY:
      secret-key: "test-api-key"
    NODE_ENV:
      value: "test"</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="team-collaboration" className="card-title"
								>Team Collaboration</ClickableHeadline
							>
							<p class="mb-4">Share configuration templates with your team:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									># kuba.yaml (commit this to version control)
default:
  provider: gcp
  project: 1337
  env:
    DATABASE_URL:
      secret-key: "database-url"
    API_KEY:
      secret-key: "api-key"
    REDIS_URL:
      value: "redis://$&lbrace;REDIS_HOST:-localhost&rbrace;:6379"

development:
  provider: gcp
  project: 1337
  env:
    DATABASE_URL:
      secret-key: "dev-database-url"
    DEBUG:
      value: "true"</code
								></pre>

							<p class="mt-4">
								<strong>Instructions for team members:</strong>
							</p>
							<ol class="list-decimal list-inside space-y-2">
								<li>Set up authentication for your cloud provider</li>
								<li>Create the necessary secrets in your cloud provider</li>
								<li>Run <code>kuba run --env development -- npm run dev</code></li>
							</ol>
						</div>
					</div>
				</div>
			</section>

			<section>
				<ClickableHeadline level={2} id="advanced-configuration" className="text-3xl font-bold mb-6"
					>Advanced Configuration</ClickableHeadline
				>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline
								level={3}
								id="multi-environment-with-secret-paths"
								className="card-title">Multi-Environment with Secret Paths</ClickableHeadline
							>
							<p class="mb-4">Use secret paths to bulk-load related secrets:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>production:
  provider: gcp
  project: 1337
  env:
    # Individual secrets
    APP_ENV:
      value: "production"

    # Database secrets (bulk load)
    DB:
      secret-path: "database"

    # API keys (bulk load)
    API:
      secret-path: "external-apis"

    # Service secrets (bulk load)
    SERVICE:
      secret-path: "microservices"

    # Interpolated connection strings
    DATABASE_URL:
      value: "postgresql://$&lbrace;DB_USERNAME&rbrace;:$&lbrace;DB_PASSWORD&rbrace;@$&lbrace;DB_HOST&rbrace;:$&lbrace;DB_PORT&rbrace;/$&lbrace;DB_NAME&rbrace;"

    REDIS_URL:
      value: "redis://$&lbrace;REDIS_HOST:-localhost&rbrace;:$&lbrace;REDIS_PORT:-6379&rbrace;/0"</code
								></pre>

							<p class="mt-4 text-sm">
								<strong>Note:</strong> This configuration will create environment variables like
								<code>DB_USERNAME</code>, <code>DB_PASSWORD</code>, <code>API_STRIPE_KEY</code>,
								<code>SERVICE_AUTH_TOKEN</code>, etc.
							</p>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<ClickableHeadline level={3} id="cross-provider-configuration" className="card-title"
								>Cross-Provider Configuration</ClickableHeadline
							>
							<p class="mb-4">Use different cloud providers for different types of secrets:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>production:
  provider: gcp
  project: 1337
  env:
    # GCP secrets
    GCP_PROJECT_ID:
      secret-key: "project-id"

    # AWS secrets
    AWS_ACCESS_KEY:
      secret-key: "aws-access-key"
      provider: aws

    # Azure secrets
    AZURE_TENANT_ID:
      secret-key: "tenant-id"
      provider: azure
      project: "my-azure-project"

    # OpenBao secrets
    INTERNAL_API_KEY:
      secret-key: "internal-api-key"
      provider: openbao

    # Hard-coded values
    APP_ENV:
      value: "production"
    DEBUG:
      value: "false"</code
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
							<p>Learn more about advanced configuration options and best practices.</p>
							<a href="/configuration" class="btn btn-outline bg-lg">Configuration Guide</a>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Cloud Providers</h3>
							<p>Set up authentication and permissions for your cloud providers.</p>
							<a href="/providers" class="btn btn-outline bg-lg">Cloud Providers Guide</a>
						</div>
					</div>
				</div>
			</section>
		</div>
	</div>
</div>
