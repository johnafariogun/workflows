# Inventory Service

A lightweight Go microservice for BlueSky Logistics that exposes two endpoints:

- `GET /status` — returns `OK` so health checks and CI can validate the service
- `GET /items` — returns a small JSON list of inventory items

This repository demonstrates a professional CI/CD layout using reusable GitHub Actions workflows, matrix testing across Go versions and OSes, automated Docker image builds and pushes to GHCR, and automatic staging Kubernetes deployments.

---

## Repository Structure

```
inventory-service/
├── .github/
│   └── workflows/
│       └── pipeline.yml        # service-level pipeline that calls reusable workflows
├── k8s/
│   └── deployment.yaml         # example Kubernetes deployment manifest for staging
├── Dockerfile
├── go.mod
├── main.go
├── main_test.go                # unit tests for handlers
├── README.md
└── .dockerignore
```

---

## CI/CD Pipeline Overview

The pipeline is defined in `.github/workflows/pipeline.yml` and delegates the heavy lifting to reusable workflows hosted in `your-org/platform-ci`:

- Stage 1 — Matrix Tests (reusable): `your-org/platform-ci/.github/workflows/go-matrix-test.yml@v1`
  - The matrix (Go versions × OSes) is defined in the reusable workflow.
  - This service passes `working-directory: "."` and `GITHUB_TOKEN` as a secret.
  - Go versions targeted by the platform workflow: `1.22`, `1.23`, `1.24`, `1.25`.
  - Operating systems targeted by the platform workflow: `ubuntu-latest`, `macos-latest`, `windows-latest`.

- Stage 2 — Docker Build & Push (reusable): `your-org/platform-ci/.github/workflows/build-and-push-image.yml@v1`
  - Builds the Docker image and pushes to GHCR with the tag: `ghcr.io/<your-username>/inventory-service:<sha>`
  - The service provides `REGISTRY_USERNAME`, `REGISTRY_PASSWORD` and `GITHUB_TOKEN` via repository secrets.

- Stage 3 — Deploy to Kubernetes (reusable): `your-org/platform-ci/.github/workflows/deploy-k8s.yml@v1`
  - Runs only on `main` branch after successful Docker push.
  - Deploys to namespace `inventory-staging`, deployment `inventory-service` with the image built above.
  - Uses `KUBECONFIG_STAGING` repository secret (passed as `KUBECONFIG_CONTENT` to the reusable workflow).

---

## Secrets

Store the following secrets in the repository (`Settings → Secrets → Actions`):

- `REGISTRY_USERNAME` — username or registry identity used to push images
- `REGISTRY_PASSWORD` — registry personal access token / password
- `KUBECONFIG_STAGING` — kubeconfig file contents for the staging cluster
- `GITHUB_TOKEN` — usually available automatically; also passed explicitly when calling reusable workflows

When the pipeline calls the reusable workflows, secrets are forwarded explicitly in the `secrets:` mapping inside `pipeline.yml`.

---

## Branching & PR Workflow

- Create feature branches (example): `feature/pipeline-setup`
- Open a Pull Request into `main` for review
- Branch protection recommendations (GitHub repo settings):
  - Require PR reviews before merge
  - Require all CI checks to pass before merge
  - Disallow direct pushes to `main`

This ensures changes are validated by the matrix tests and Docker builds before deployment.

---

## Local development and testing

Build and run locally:

```bash
go build -o inventory-service .
./inventory-service
# then
curl http://localhost:8080/status
curl http://localhost:8080/items
```

Run unit tests:

```bash
go test ./...
```

---

## Kubernetes (Staging) example

An example `k8s/deployment.yaml` is included. The platform's reusable deploy workflow will apply the image to the `inventory-staging` namespace and the `inventory-service` deployment.

---

## Submitting the assignment

What to provide once you push the repository:

1. Repository URL
2. Pull Request URL (feature branch → main)
3. Screenshots of the matrix test jobs, Docker build job, and the staging deployment job in GitHub Actions
4. Final Docker image URL: `ghcr.io/<your-username>/inventory-service:<sha>`
5. README.md (this file) — explanation of pipeline and secrets

---

## Next steps / Optional bonuses

- Add Go unit tests and code coverage reporting
- Add semantic-release workflow + GitHub Releases
- Extend Docker build to multi-arch images (amd64 + arm64)
- Add notifications (Slack / Teams) on deploy

---

If you'd like, I can also generate the feature branch, commit these files, and give the exact commands to create the repository on GitHub and enable branch protection rules. Let me know how you want to proceed.
