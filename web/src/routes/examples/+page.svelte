<script lang="ts">
	import HeadComponent from '$lib/HeadComponent.svelte';
	import Prism from 'prismjs';
	import 'prismjs/plugins/toolbar/prism-toolbar';
	import 'prismjs/plugins/copy-to-clipboard/prism-copy-to-clipboard';
	import 'prismjs/components/prism-yaml';
	import 'prismjs/components/prism-bash';
	import 'prismjs/components/prism-javascript';
	import 'prismjs/components/prism-python';
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
		title: 'Examples - Kuba',
		description:
			'Practical examples and use cases for using Kuba with different applications and frameworks.'
	}}
/>

<div class="container mx-auto px-4 py-8">
	<div class="max-w-4xl mx-auto">
		<div class="text-center mb-12">
			<h1 class="text-4xl font-bold mb-4">Examples & Use Cases</h1>
			<p class="text-xl text-base-content/70">
				See practical examples of how to use Kuba with different applications, frameworks, and
				deployment scenarios.
			</p>
		</div>

		<div class="space-y-12">
			<section>
				<h2 class="text-3xl font-bold mb-6">Web Application Examples</h2>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Node.js Express Application</h3>
							<p class="mb-4">
								Run a Node.js Express application with database credentials and API keys:
							</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run --env production -- node app.js</code
								></pre>

							<h4 class="font-bold mt-4 mb-2">Configuration (kuba.yaml):</h4>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>production:
  provider: gcp
  project: my-project
  mappings:
    - environment-variable: "DATABASE_URL"
      secret-key: "prod-database-url"
    - environment-variable: "JWT_SECRET"
      secret-key: "jwt-secret"
    - environment-variable: "STRIPE_SECRET_KEY"
      secret-key: "stripe-secret-key"
    - environment-variable: "REDIS_URL"
      value: "redis://$&lbrace;REDIS_HOST:-localhost&rbrace;:6379"</code
								></pre>

							<h4 class="font-bold mt-4 mb-2">Application Code:</h4>
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

							<h4 class="font-bold mt-4 mb-2">Configuration (kuba.yaml):</h4>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>development:
  provider: aws
  mappings:
    - environment-variable: "FLASK_ENV"
      value: "development"
    - environment-variable: "DATABASE_URL"
      secret-key: "dev-database-url"
    - environment-variable: "SECRET_KEY"
      secret-key: "flask-secret-key"
    - environment-variable: "DEBUG"
      value: "true"</code
								></pre>

							<h4 class="font-bold mt-4 mb-2">Application Code:</h4>
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
				<h2 class="text-3xl font-bold mb-6">Database & API Examples</h2>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Database Migrations</h3>
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

							<h4 class="font-bold mt-4 mb-2">Configuration:</h4>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>production:
  provider: gcp
  project: my-project
  mappings:
    - environment-variable: "DATABASE_URL"
      secret-key: "prod-postgres-url"
    - environment-variable: "DB_PASSWORD"
      secret-key: "prod-db-password"

development:
  provider: aws
  mappings:
    - environment-variable: "DATABASE_URL"
      secret-key: "dev-postgres-url"
    - environment-variable: "DB_PASSWORD"
      secret-key: "dev-db-password"</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">External API Integration</h3>
							<p class="mb-4">Connect to external APIs with secure keys:</p>
							<pre><code
									class="language-bash"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋">kuba run --env staging -- python api_client.py</code
								></pre>

							<h4 class="font-bold mt-4 mb-2">Configuration:</h4>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>staging:
  provider: azure
  mappings:
    - environment-variable: "STRIPE_API_KEY"
      secret-key: "stripe-staging-key"
    - environment-variable: "SENDGRID_API_KEY"
      secret-key: "sendgrid-staging-key"
    - environment-variable: "TWILIO_ACCOUNT_SID"
      secret-key: "twilio-account-sid"
    - environment-variable: "TWILIO_AUTH_TOKEN"
      secret-key: "twilio-auth-token"</code
								></pre>

							<h4 class="font-bold mt-4 mb-2">API Client Code:</h4>
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
				<h2 class="text-3xl font-bold mb-6">Docker & Container Examples</h2>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Docker Container with Secrets</h3>
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
  myapp</code
								></pre>

							<h4 class="font-bold mt-4 mb-2">Dockerfile:</h4>
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
							<h3 class="card-title">Docker Compose Integration</h3>
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

							<h4 class="font-bold mt-4 mb-2">docker-compose.yml:</h4>
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
				<h2 class="text-3xl font-bold mb-6">CI/CD Pipeline Examples</h2>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">GitHub Actions</h3>
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
							<h3 class="card-title">GitLab CI</h3>
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
				<h2 class="text-3xl font-bold mb-6">Development Workflow Examples</h2>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Local Development</h3>
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

							<h4 class="font-bold mt-4 mb-2">Development Configuration:</h4>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>development:
  provider: gcp
  project: my-dev-project
  mappings:
    - environment-variable: "DATABASE_URL"
      secret-key: "dev-database-url"
    - environment-variable: "API_KEY"
      secret-key: "dev-api-key"
    - environment-variable: "DEBUG"
      value: "true"
    - environment-variable: "LOG_LEVEL"
      value: "debug"

testing:
  provider: gcp
  project: my-dev-project
  mappings:
    - environment-variable: "DATABASE_URL"
      secret-key: "test-database-url"
    - environment-variable: "API_KEY"
      secret-key: "test-api-key"
    - environment-variable: "NODE_ENV"
      value: "test"</code
								></pre>
						</div>
					</div>

					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Team Collaboration</h3>
							<p class="mb-4">Share configuration templates with your team:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									># kuba.yaml.template (commit this to version control)
default:
  provider: gcp
  project: YOUR_PROJECT_ID
  mappings:
    - environment-variable: "DATABASE_URL"
      secret-key: "database-url"
    - environment-variable: "API_KEY"
      secret-key: "api-key"
    - environment-variable: "REDIS_URL"
      value: "redis://$&lbrace;REDIS_HOST:-localhost&rbrace;:6379"

development:
  provider: gcp
  project: YOUR_PROJECT_ID
  mappings:
    - environment-variable: "DATABASE_URL"
      secret-key: "dev-database-url"
    - environment-variable: "DEBUG"
      value: "true"</code
								></pre>

							<p class="mt-4">
								<strong>Instructions for team members:</strong>
							</p>
							<ol class="list-decimal list-inside space-y-2">
								<li>Copy <code>kuba.yaml.template</code> to <code>kuba.yaml</code></li>
								<li>Replace <code>YOUR_PROJECT_ID</code> with your actual project ID</li>
								<li>Set up authentication for your cloud provider</li>
								<li>Create the necessary secrets in your cloud provider</li>
								<li>Run <code>kuba run --env development -- npm run dev</code></li>
							</ol>
						</div>
					</div>
				</div>
			</section>

			<section>
				<h2 class="text-3xl font-bold mb-6">Advanced Configuration Examples</h2>

				<div class="space-y-6">
					<div class="card bg-base-200">
						<div class="card-body">
							<h3 class="card-title">Multi-Environment with Secret Paths</h3>
							<p class="mb-4">Use secret paths to bulk-load related secrets:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>production:
  provider: gcp
  project: my-project
  mappings:
    # Individual secrets
    - environment-variable: "APP_ENV"
      value: "production"
    
    # Database secrets (bulk load)
    - environment-variable: "DB"
      secret-path: "database"
    
    # API keys (bulk load)
    - environment-variable: "API"
      secret-path: "external-apis"
    
    # Service secrets (bulk load)
    - environment-variable: "SERVICE"
      secret-path: "microservices"
    
    # Interpolated connection strings
    - environment-variable: "DATABASE_URL"
      value: "postgresql://$&lbrace;DB_USERNAME&rbrace;:$&lbrace;DB_PASSWORD&rbrace;@$&lbrace;DB_HOST&rbrace;:$&lbrace;DB_PORT&rbrace;/$&lbrace;DB_NAME&rbrace;"
    
    - environment-variable: "REDIS_URL"
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
							<h3 class="card-title">Cross-Provider Configuration</h3>
							<p class="mb-4">Use different cloud providers for different types of secrets:</p>
							<pre><code
									class="language-yaml"
									data-toolbar-order="copy-to-clipboard"
									data-prismjs-copy="📋"
									>production:
  provider: gcp
  project: my-gcp-project
  mappings:
    # GCP secrets
    - environment-variable: "GCP_PROJECT_ID"
      secret-key: "project-id"
    
    # AWS secrets
    - environment-variable: "AWS_ACCESS_KEY"
      secret-key: "aws-access-key"
      provider: aws
    
    # Azure secrets
    - environment-variable: "AZURE_TENANT_ID"
      secret-key: "tenant-id"
      provider: azure
      project: "my-azure-project"
    
    # OpenBao secrets
    - environment-variable: "INTERNAL_API_KEY"
      secret-key: "internal-api-key"
      provider: openbao
    
    # Hard-coded values
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
				<h2 class="text-3xl font-bold mb-6">Next Steps</h2>

				<div class="grid md:grid-cols-2 gap-6">
					<div class="card bg-primary text-primary-content">
						<div class="card-body">
							<h3 class="card-title">Configuration Guide</h3>
							<p>Learn more about advanced configuration options and best practices.</p>
							<a href="/configuration" class="btn btn-secondary">Configuration Guide</a>
						</div>
					</div>

					<div class="card bg-secondary text-secondary-content">
						<div class="card-body">
							<h3 class="card-title">Cloud Providers</h3>
							<p>Set up authentication and permissions for your cloud providers.</p>
							<a href="/providers" class="btn btn-primary">Cloud Providers Guide</a>
						</div>
					</div>
				</div>
			</section>
		</div>
	</div>
</div>
