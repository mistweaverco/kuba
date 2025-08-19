<div align="center">

![kuba logo](assets/logo.svg)

# Kuba

[![Made with love](assets/badge-made-with-love.svg)](https://github.com/mistweaverco/kuba/graphs/contributors)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/mistweaverco/kuba?style=for-the-badge)](https://github.com/mistweaverco/kuba/releases/latest)
[![License](https://img.shields.io/github/license/mistweaverco/kuba?style=for-the-badge)](./LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/mistweaverco/kuba?style=for-the-badge)](https//:github.com/mistweaverco/kuba/issues)
[![Discord](assets/badge-discord.svg)](https://mistweaverco.com/discord)

[Why?](#why) • [Installation](#installation) • [Usage](#usage)

<p></p>

Kuba is [Swahili](https://en.wikipedia.org/wiki/Swahili_language) for "vault."

Kuba helps you to get rid of `.env` files.

Pass env directly from GCP, AWS and Azure Secrets to your application

<p></p>

</div>

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
and AWS, as well as Azure Key Vault.

This eliminates the need for `.env` files and provides a more secure,
consistent, and scalable way to manage environment variables across different environments.

## Installation

Kuba is a single binary, so you can install it easily.

### Manual Installation

Download the latest release from [GitHub Releases](https://github.com/mistweaverco/kuba/releases/latest).

### Automatic Linux and macOS Installation

You can install it using `curl`:

```sh
curl -sSL https://getkuba.net/install | sh
```

### Automatic Windows Installation

Run the following command in PowerShell:

```powershell
iwr https://getkuba.net/install -useb | iex
```

## Usage

```sh
kuba run -- <your-command>
```

This will fetch all secrets definded in
`kuba.yaml` and pass them as
environment variables to any arbitrary command.

Let's say you want to pass
some secrets from GCP to your node application.

```sh
kuba run -- node dist/server.js
```

and your `kuba.yaml` would look something like this:

```yaml
---
# Top-level sections for different environments.
default:
  provider: "gcp"
  project: "my-gcp-project-default"

  # Mapping of cloud projects to environment variables and secret keys.
  mappings:
    - environment-variable: "GCP_PROJECT_ID"
      secret-key: "gcp_project_secret"
    - environment-variable: "AWS_PROJECT_ID"
      secret-key: "aws_project_secret"
      provider: "aws"
      project: "my-aws-project-default"
    - environment-variable: "AZURE_PROJECT_ID"
      secret-key: "azure_project_secret"
      provider: "azure"
      project: "my-azure-project-default"

---

# Settings for the development environment.
development:
  provider: "gcp"
  project: "my-gcp-project-development"
  
  # You can override specific mappings here or add new ones.
  mappings:
    - environment-variable: "DEV_GCP_PROJECT_ID"
      secret-key: "dev_gcp_project_secret"
    - environment-variable: "DEV_AWS_PROJECT_ID"
      secret-key: "dev_aws_project_secret"
      provider: "aws"
      project: "my-aws-project-development"

---

# Settings for the staging environment.
staging:
  provider: "gcp"
  project: "my-gcp-project-staging"
  
  mappings:
    - environment-variable: "STAGING_GCP_PROJECT_ID"
      secret-key: "staging_gcp_project_secret"
    - environment-variable: "STAGING_AWS_PROJECT_ID"
      secret-key: "staging_aws_project_secret"
      provider: "aws"
      project: "my-aws-project-staging"

---

# Settings for the production environment.
production:
  provider: "gcp"
  project: "my-gcp-project-production"
  
  mappings:
    - environment-variable: "PROD_GCP_PROJECT_ID"
      secret-key: "prod_gcp_project_secret"
    - environment-variable: "PROD_AWS_PROJECT_ID"
      secret-key: "prod_aws_project_secret"
      provider: "aws"
      project: "my-aws-project-production"
```
