# kuba

Kuba helps you to get rid of `.env` files.

Pass env directly from GCP, AWS and Azure Secrets to your application

## Example

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
