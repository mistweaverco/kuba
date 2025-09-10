<div align="center">

![kuba logo](assets/logo.svg)

# Kuba

[![Made with love](assets/badge-made-with-love.svg)](https://github.com/mistweaverco/kuba/graphs/contributors)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/mistweaverco/kuba?style=for-the-badge)](https://github.com/mistweaverco/kuba/releases/latest)
[![License](https://img.shields.io/github/license/mistweaverco/kuba?style=for-the-badge)](./LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/mistweaverco/kuba?style=for-the-badge)](https://github.com/mistweaverco/kuba/issues)
[![Discord](assets/badge-discord.svg)](https://mistweaverco.com/discord)

[Why?](#why) • [Installation](#installation) • [Usage](#usage)

<p></p>

Kuba is [Swahili](https://en.wikipedia.org/wiki/Swahili_language) for "vault."

Kuba helps you to get rid of `.env` files.

Pass env directly from GCP Secret Manager,
AWS Secrets Manager,
Azure Key Vault, and OpenBao to your application

<p></p>

</div>

## Table of Contents

- [Why?](#why)
  - [Advantages over other services](#advantages-over-other-services)
- [Installation](#installation)
  - [Manual Installation](#manual-installation)
  - [Automatic Linux and macOS Installation](#automatic-linux-and-macos-installation)
  - [Automatic Windows Installation](#automatic-windows-installation)
- [Usage](#usage)
  - [Configuration File Structure](#configuration-file-structure)
  - [Environment Variable Interpolation](#environment-variable-interpolation)
  - [Secret Path Mapping](#secret-path-mapping)
  - [Running with a specific environment](#running-with-a-specific-environment)
  - [Testing configuration and access](#testing-configuration-and-access)
- [Cloud Provider Setup](#cloud-provider-setup)
  - [Google Cloud Platform (GCP)](#google-cloud-platform-gcp)
  - [AWS Secrets Manager](#aws-secrets-manager)
  - [Azure Key Vault](#azure-key-vault)
  - [OpenBao](#openbao)

---

## Why?

Environment variables are a common way to manage configuration in applications,
especially when deploying to different environments like development,
staging, and production.

However, managing these variables can become cumbersome,
especially when dealing with multiple cloud providers and
secret management systems.

This often leads to the use of `.env` files,
which can be problematic for several reasons:

- Onboarding new developers, often involves sharing `.env` files.
  This often leads to `.env` files being shared insecurely,
  such as through email or chat applications,
  which can expose sensitive information.
- **Manual Management**: Keeping `.env` files up-to-date with the latest secrets
  from cloud providers can be tedious and error-prone.
- **Security Risks**: `.env` files can accidentally be committed to version control,
  exposing sensitive information.
- **Lack of Standardization**: Each cloud provider has its own way of managing secrets,
  leading to a fragmented approach that can complicate development and deployment.

Kuba addresses these issues by allowing you to define your environment variables
in a single `kuba.yaml` file and fetch them directly from cloud providers like GCP
Secret Manager, AWS Secrets Manager, Azure Key Vault, and OpenBao.

This eliminates the need for `.env` files and provides a more secure,
consistent, and scalable way to manage environment variables across
different environments.

### Advantages over other services

To be clear, there are many other tools that can help you manage secrets:

- [Doppler](https://www.doppler.com/)
- [Vault](https://www.vaultproject.io/)
- [1Password Secrets Automation](https://developer.1password.com/docs/secrets-automation/)
- [Infisical](https://infisical.com/)

… and many more.

> [!CAUTION]
> Most of them require a whopping subscription fee,
> or setting up and maintaining a separate service yourself,
> which can be a barrier for small teams or individual developers.

However, Kuba is designed to be straightforward and easy to use,
by leveraging the existing secret management systems of cloud providers,
that you might already be using.

## Installation

Kuba is a single binary, so you can install it easily.

### Manual Installation

Download the latest release from [GitHub Releases](https://github.com/mistweaverco/kuba/releases/latest).

### Automatic Linux and macOS Installation

You can install it using `curl`:

```sh
curl -sSL https://kuba.mwco.app/install.sh | sh
```

### Automatic Windows Installation

Run the following command in PowerShell:

```powershell
iwr https://kuba.mwco.app/install.ps1 -useb | iex
```

## Usage

```sh
kuba run -- <your-command>
```

This will fetch all secrets defined in
`kuba.yaml` and pass them as
environment variables to any arbitrary command.

### Debug Mode

For troubleshooting configuration issues and seeing detailed execution steps, you can enable debug mode:

```sh
kuba --debug run -- <your-command>
# or use the short form
kuba -d run -- <your-command>
```

Debug mode provides verbose logging that shows:
- Configuration file discovery and loading
- Environment selection and validation
- Secret provider initialization
- Secret retrieval attempts and results
- Environment variable mapping
- Command execution details

This is particularly useful for:
- Diagnosing cloud provider authentication issues
- Troubleshooting configuration file syntax errors
- Understanding why certain secrets aren't being loaded
- Verifying environment variable interpolation
- Debugging provider-specific errors

### Available Commands and Flags

Kuba provides several commands to help you manage your configuration:

```sh
# Initialize a new configuration file
kuba init

# Run a command with secrets
kuba run -- <command>

# Test secret retrieval without running a command
kuba test --env <environment>

# Show version information
kuba version

# Get help
kuba --help
```

**Global Flags:**
- `--debug, -d`: Enable debug mode for verbose logging
- `--version`: Show version information
- `--help, -h`: Show help information

**Run Command Flags:**
- `--env, -e`: Specify environment (default: "default")
- `--config, -c`: Path to configuration file

**Test Command Flags:**
- `--env, -e`: Specify environment (default: "default")
- `--config, -c`: Path to configuration file

Let's say you want to pass
some secrets from GCP to your node application.

```sh
kuba run -- node dist/server.js
```

and your `kuba.yaml` would look something like this:

```yaml
# yaml-language-server: $schema=https://kuba.mwco.app/kuba.schema.json
---
# Top-level sections for different environments.
default:
  provider: gcp
  project: 1337

  # Mapping of cloud projects to environment variables and secret keys.
  mappings:
    - environment-variable: "GCP_PROJECT_ID"
      secret-key: "gcp_project_secret"
    - environment-variable: "AWS_PROJECT_ID"
      secret-key: "aws_project_secret"
      provider: aws
    - environment-variable: "AZURE_PROJECT_ID"
      secret-key: "azure_project_secret"
      provider: azure
      project: "my-azure-project-default"
    - environment-variable: "OPENBAO_SECRET"
      secret-key: "secret/openbao-secret"
      provider: openbao
    - environment-variable: "DATABASE_CONFIG"
      secret-path: "database"
    - environment-variable: "API_KEYS"
      secret-path: "external-apis"
      provider: aws
    - environment-variable: "SOME_HARD_CODED_ENV"
      value: "hard-coded-value"

---

# Settings for the development environment.
development:
  provider: gcp
  project: 1337

  # You can override specific mappings here or add new ones.
  mappings:
    - environment-variable: "DEV_GCP_PROJECT_ID"
      secret-key: "dev_gcp_project_secret"
    - environment-variable: "DEV_AWS_PROJECT_ID"
      secret-key: "dev_aws_project_secret"
      provider: aws

---

# Settings for the staging environment.
staging:
  provider: gcp
  project: 1337

  mappings:
    - environment-variable: "STAGING_GCP_PROJECT_ID"
      secret-key: "staging_gcp_project_secret"
    - environment-variable: "STAGING_AWS_PROJECT_ID"
      secret-key: "staging_aws_project_secret"
      provider: aws
---

# Settings for the production environment.
production:
  provider: gcp
  project: 1337

  mappings:
    - environment-variable: "PROD_GCP_PROJECT_ID"
      secret-key: "prod_gcp_project_secret"
    - environment-variable: "PROD_AWS_PROJECT_ID"
      secret-key: "prod_aws_project_secret"
      provider: aws
```

This `kuba.yaml` file defines the secrets for different environments
and maps them to environment variables. The example includes:

- **Individual secrets** using `secret-key`
  (e.g., GCP_PROJECT_ID, AWS_PROJECT_ID)
- **Secret paths** using `secret-path` to
  fetch all secrets under a prefix (e.g., DATABASE_CONFIG, API_KEYS)
- **Hard-coded values** using `value` for static configuration
- **Cross-provider mappings** where different secrets come
  from different cloud providers

### Confguration File Structure

Each top-level section corresponds to a different environment,
such as `default`, `development`, `staging`, and `production`.
They're completely arbitrary and can be named as you wish.

Each section specifies the cloud provider, the project ID,
and a list of mappings between environment variables and secret keys.

You can also specify the provider and project ID for each mapping,
allowing you to fetch secrets from different cloud providers
or projects as needed. Kuba currently supports GCP Secret Manager,
AWS Secrets Manager, Azure Key Vault, and OpenBao.

### Environment Variable Interpolation

Kuba supports environment variable interpolation
in the `value` field using `${VAR_NAME}` syntax.

This allows you to:

- Reference previously defined environment variables from the same configuration
- Use system environment variables
- Build complex connection strings and URLs dynamically

**Example with interpolation:**

```yaml
default:
  provider: gcp
  project: 1337
  mappings:
    - environment-variable: "DB_PASSWORD"
      secret-key: "db-password"
    - environment-variable: "DB_HOST"
      value: "mydbhost"
    - environment-variable: "DB_CONNECTION_STRING"
      value: "postgresql://user:${DB_PASSWORD}@${DB_HOST}:5432/mydb"
    - environment-variable: "API_URL"
      value: "https://api.${DOMAIN}/v1"
    - environment-variable: "APP_ENV"
      value: "${NODE_ENV:-development}"
    - environment-variable: "REDIS_URL"
      value: "redis://${REDIS_HOST:-localhost}:${REDIS_PORT:-6379}/0"
```

In this example:

- `${DB_PASSWORD}` will be replaced with the value from the secret
- `${DB_HOST}` will be replaced with the literal value "mydbhost"
- `${DOMAIN}` will be replaced with the system environment variable if it exists
- `${NODE_ENV:-development}` will use the `NODE_ENV` environment variable if set, otherwise default to "development"
- `${REDIS_HOST:-localhost}` will use the `REDIS_HOST` environment variable if set, otherwise default to "localhost"

**Note**: Interpolation is processed in order,
so you can reference variables defined earlier in the same configuration.
Unresolved variables will remain unchanged in the output.

**Shell-style default values**: You can use `${VAR_NAME:-default}` syntax to
provide fallback values when environment variables aren't set.

This is particularly useful for providing sensible defaults
while allowing overrides through environment variables.

**Environment variable naming**:
All environment variable names (including those from secrets) are
automatically sanitized to be valid POSIX environment variable names.

This means:

- Names are converted to uppercase
- Non-alphanumeric characters are replaced with underscores
- Names that don't start with a letter or underscore get a leading underscore
- This ensures compatibility across different operating systems and shells

### Secret Path Mapping

In addition to individual secret keys, Kuba supports **secret path mapping** using the `secret-path` field.
This feature allows you to fetch all secrets that start with a given path prefix,
which is particularly useful for:

- **Bulk secret retrieval**: Fetch all secrets under a specific namespace or directory
- **Organized secret management**: Group related secrets under common path prefixes
- **Environment-specific configurations**: Load all secrets for a specific environment or service

**How it works:**
- When you specify a `secret-path`, Kuba will fetch all secrets that start with that path
- Each secret found will be converted to an environment variable using the pattern: `{ENVIRONMENT_VARIABLE}_{SECRET_NAME}`
- Secret names are automatically sanitized to be valid POSIX environment variable names (uppercase, underscores only)

**Example with secret paths:**

```yaml
default:
  provider: gcp
  project: 1337
  mappings:
    - environment-variable: "DB"
      secret-path: "database"
    - environment-variable: "API"
      secret-path: "external-apis"
    - environment-variable: "SERVICE"
      secret-path: "microservices"
    - environment-variable: "HARD_CODED"
      value: "static-value"
```

If your GCP Secret Manager contains secrets like:

- `database-connection-string`
- `database-username`
- `database-password`
- `external-apis-stripe-key`
- `external-apis-sendgrid-key`
- `microservices-auth-service-token`
- `microservices-user-service-token`

Kuba will create these environment variables:

- `DB_CONNECTION_STRING` = value of `database-connection-string`
- `DB_USERNAME` = value of `database-username`
- `DB_PASSWORD` = value of `database-password`
- `API_STRIPE_KEY` = value of `external-apis-stripe-key`
- `API_SENDGRID_KEY` = value of `external-apis-sendgrid-key`
- `SERVICE_AUTH_SERVICE_TOKEN` = value of `microservices-auth-service-token`
- `SERVICE_USER_SERVICE_TOKEN` = value of `microservices-user-service-token`

**Cross-provider secret paths:**
You can also use secret paths with different providers:

```yaml
default:
  provider: gcp
  project: 1337
  mappings:
    - environment-variable: "GCP_SECRETS"
      secret-path: "app-config"
      provider: gcp
    - environment-variable: "AWS_SECRETS"
      secret-path: "app-config"
      provider: aws
    - environment-variable: "AZURE_SECRETS"
      secret-path: "app-config"
      provider: azure
      project: "my-azure-project"
```

**Important notes:**
- Secret paths work with all supported providers (GCP, AWS, Azure, OpenBao)
- The resulting environment variable names are automatically sanitized and uppercased
- You can mix `secret-key`, `secret-path`, and `value` mappings in the same configuration
- Secret paths are processed after individual secret keys, so you can reference path-based variables in value interpolations

### Running with a specific environment

You can also specify the environment you want to use:

```sh
kuba run --env development -- node dist/server.js
```

### Testing configuration and access

Use the `test` subcommand to verify that Kuba can load your configuration and
retrieve all mapped values for an environment without executing any program:

```sh
# Use default environment
kuba test

# Specify an environment
kuba test --env staging

# Point to a specific configuration file
kuba test --config ./config/kuba.yaml --env production
```

This is useful for validating credentials, permissions, and
configuration mappings during setup or CI.

## Cloud Provider Setup

### Google Cloud Platform (GCP)

Kuba supports GCP Secret Manager for fetching secrets. To use GCP:

1. **Enable Secret Manager API**: Make sure the Secret Manager API is enabled in your GCP project.

2. **Authentication**: Set up authentication using one of these methods:
   - **Service Account Key**: Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to your service account JSON key file:
     ```sh
     export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
     ```
   - **Application Default Credentials**: Use `gcloud auth application-default login` to set up local development credentials
   - **Workload Identity**: If running on GKE or other GCP services, use workload identity
   - **Compute Engine**: If running on Compute Engine, the default service account will be used automatically

3. **IAM Permissions**: Ensure your service account has the `Secret Manager Secret Accessor` role for the secrets you want to access.

4. **Example Configuration**:

   ```yaml
   default:
     provider: gcp
     project: 1337
     mappings:
       - environment-variable: "DATABASE_URL"
         secret-key: "database-connection-string"
       - environment-variable: "API_KEY"
         secret-key: "external-api-key"
       - environment-variable: "SOME_HARD_CODED_ENV"
         value: "hard-coded-value"
   ```

### AWS Secrets Manager

Kuba supports AWS Secrets Manager for fetching secrets. To use AWS:

1. **Authentication**: Set up authentication using one of these methods:
   - **Environment Variables**: Set `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`:
     ```sh
     export AWS_ACCESS_KEY_ID="your-access-key"
     export AWS_SECRET_ACCESS_KEY="your-secret-key"
     export AWS_REGION="us-east-1"
     ```
   - **AWS Profile**: Set `AWS_PROFILE` to use a specific profile from your AWS credentials file:
     ```sh
     export AWS_PROFILE="my-profile"
     export AWS_REGION="us-east-1"
     ```
   - **IAM Roles**: If running on EC2, ECS, or other AWS services, use IAM roles
   - **AWS CLI**: Use `aws configure` to set up your credentials

2. **IAM Permissions**: Ensure your AWS credentials have the `secretsmanager:GetSecretValue` permission for the secrets you want to access.

3. **Example Configuration**:

   ```yaml
   default:
     provider: aws
     mappings:
       - environment-variable: "DATABASE_URL"
         secret-key: "database-connection-string"
       - environment-variable: "API_KEY"
         secret-key: "external-api-key"
       - environment-variable: "SOME_HARD_CODED_ENV"
         value: "hard-coded-value"
   ```

### Azure Key Vault

Kuba supports Azure Key Vault for fetching secrets. To use Azure Key Vault:

1. **Authentication**: Kuba supports multiple authentication methods:
   - **Service Principal**: Set the following environment variables:
     ```bash
     export AZURE_KEY_VAULT_URL="https://yourvault.vault.azure.net/"
     export AZURE_TENANT_ID="your-tenant-id"
     export AZURE_CLIENT_ID="your-client-id"
     export AZURE_CLIENT_SECRET="your-client-secret"
     ```
   - **Managed Identity**: If running on Azure services with managed identity enabled
   - **Default Azure Credential**: Uses Azure CLI, Visual Studio Code, or other Azure tools

2. **Key Vault Permissions**: Ensure your Azure credentials have the `Get` and `List` permissions for secrets in your Key Vault.

3. **Configuration**: In your `kuba.yaml`, specify the Azure provider:
   ```yaml
   default:
     provider: azure
     mappings:
       - environment-variable: "DATABASE_URL"
         secret-key: "database-connection-string"
       - environment-variable: "SOME_HARD_CODED_ENV"
         value: "hard-coded-value"
   ```

### OpenBao

Kuba supports OpenBao for fetching secrets.
OpenBao is a fork of HashiCorp Vault that provides secure secret storage and access.

To use OpenBao:

1. **Setup**: Make sure you have an OpenBao server running and accessible.

2. **Authentication**: Set up authentication using environment variables:
   ```bash
   export OPENBAO_ADDR="http://localhost:8200"  # Required: OpenBao server address
   export OPENBAO_TOKEN="your-openbao-token"    # Optional: Authentication token
   export OPENBAO_NAMESPACE="your-namespace"     # Optional: Namespace (if using enterprise features)
   ```

3. **Permissions**: Ensure your OpenBao token has read permissions for the secrets you want to access.

4. **Configuration**: In your `kuba.yaml`, specify the OpenBao provider:
   ```yaml
   default:
     provider: openbao
     mappings:
       - environment-variable: "DATABASE_URL"
         secret-key: "secret/database-url"
       - environment-variable: "API_KEY"
         secret-key: "secret/api-key"
       - environment-variable: "SOME_HARD_CODED_ENV"
         value: "hard-coded-value"
   ```

**Note**: OpenBao secrets are stored as key-value pairs. If a secret contains multiple keys, Kuba will return the first string value it finds. For more precise control, structure your secrets with single values or use the project field to namespace your secrets:

```yaml
default:
  provider: openbao
  mappings:
    - environment-variable: "DATABASE_URL"
      secret-key: "database-url"
      project: "secret"  # This will look for secret/database-url
```
