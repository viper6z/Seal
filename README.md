# Seal

Seal is a small, single-host application platform built as a hands-on cloud and platform-engineering project.

It provides a narrow developer interface for deploying containerised applications to one AWS EC2 instance running Docker Compose and Nginx.

```text
application image in GHCR
→ Seal YAML manifest
→ Go CLI validation and config generation
→ pull request
→ CI validation
→ merge
→ host reconciliation
→ Docker Compose
→ Nginx
→ application
```

## Architecture

```text
GitHub
├── GitHub Actions
├── Terraform
└── application manifests

AWS
└── EC2
    ├── Docker
    ├── Docker Compose
    ├── Nginx
    └── application containers
```

Terraform provisions the AWS infrastructure. Docker Compose defines the services running on the host. Nginx is the public HTTP entry point, while application containers remain private on the backend Docker network.

## Application interface

Applications are declared through a small YAML manifest:

```yaml
name: example-api
image: ghcr.io/example/example-api:v1
internal_port: 8080
exposure_type: public
allowed_public_routes:
  - /
  - /health
```

Seal validates the manifest and generates the platform configuration needed to run the application.

Public applications receive explicitly declared Nginx routes. Internal applications join the backend network without being exposed through Nginx.

## Current status

The project currently includes:

* reproducible AWS infrastructure with Terraform
* GitHub Actions using AWS OIDC
* EC2 management through AWS Systems Manager
* a root Docker Compose stack with separate edge and backend networks
* Nginx as the only public entry point
* a Go CLI for manifest validation and Compose generation
* safe Compose updates using temporary-file validation and atomic replacement

Per-application Nginx generation, Git automation and the host-side reconciliation agent are still under development.

## Scope

Seal is intentionally not Kubernetes or a general-purpose production platform. It is a focused learning project for understanding how infrastructure provisioning, configuration generation, CI/CD, container networking, ingress and reconciliation fit together on a single host.

The project build-up and engineering decisions are documented in [`docs/logbook.md`](docs/logbook.md).
