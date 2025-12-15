# Inventory Service

A lightweight Golang microservice built as part of the **Scenario-Based Assignment: Multi-Platform Inventory Service** for the BlueSky Logistics Platform Engineering Team.

The service exposes two simple HTTP endpoints and demonstrates a **production-grade CI/CD pipeline** using **GitHub reusable workflows**, **matrix testing**, **Docker image automation**, and **Kubernetes staging deployments**.

---

## Service Endpoints

* `GET /status` — returns `OK` for health checks and CI validation
* `GET /items` — returns a JSON array of sample inventory items

---

## Repository Structure

```
inventory-service/
├── .github/
│   └── workflows/
│       └── pipeline.yml        # Service-level pipeline calling reusable workflows
├── k8s/
│   └── deployment.yaml         # Example Kubernetes deployment (staging)
├── Dockerfile
├── go.mod
├── main.go
├── main_test.go                # Unit tests for HTTP handlers
├── .dockerignore
└── README.md
```

---

## CI/CD Pipeline Overview

The CI/CD pipeline is defined in `.github/workflows/pipeline.yml`. Instead of duplicating logic, it delegates work to **reusable workflows** hosted in a separate platform repository:

**Reusable workflows repository:**
[https://github.com/johnafariogun/workflows](https://github.com/johnafariogun/workflows)

[https://github.com/johnafariogun/workflows](https://github.com/johnafariogun/workflows)

The pipeline consists of three mandatory stages:

---

### Stage 1 — Matrix Tests (Reusable Workflow)

**Workflow used:**
`.github/workflows/go-matrix-test.yml`

* Executes unit tests across a matrix defined *inside the reusable workflow*
* The service repository does **not** define the matrix directly
* Matrix dimensions:

  * **Go versions:** `1.22`, `1.23`, `1.24`, `1.25`
  * **Operating systems:**

    * `ubuntu-latest`
    * `macos-latest`
    * `windows-latest`
* Ensures cross-platform and multi-version compatibility

**Secrets explicitly passed:**

* `GITHUB_TOKEN`

---

### Stage 2 — Docker Build & Push (Reusable Workflow)

**Workflow used:**
`.github/workflows/build-and-push-image.yml`

* Builds a Docker image from the provided `Dockerfile`
* Tags the image using the commit SHA
* Pushes the image to **GitHub Container Registry (GHCR)**

**Image format:**

```
ghcr.io/johnafariogun/inventory-service:<commit-sha>
```

**Secrets explicitly passed:**

* `REGISTRY_USERNAME`
* `REGISTRY_PASSWORD`
* `GITHUB_TOKEN` 'couldn't work with this due to github not allowing secrets starting with *GITHUB_*'

---

### Stage 3 — Deploy to Kubernetes (Reusable Workflow)

**Workflow used:**
`.github/workflows/deploy-k8s.yml`

* Runs **only on the `main` branch**
* Executes after a successful Docker build and push
* Deploys the service to the staging environment (as no staging environment credentials, it for now just skips the deployment)

**Deployment target:**

* Namespace: `inventory-staging`
* Deployment name: `inventory-service`
* Image: `ghcr.io/johnafariogun/inventory-service:<commit-sha>`

**Secrets explicitly passed:**

* `KUBECONFIG_STAGING` (forwarded as `KUBECONFIG_CONTENT`)
* `GITHUB_TOKEN`

---

## Secrets Management

All secrets are stored securely in the repository under:

**Settings → Secrets → Actions**

| Secret Name          | Purpose                                                              |
| -------------------- | -------------------------------------------------------------------- |
| `REGISTRY_USERNAME`  | Docker registry username / identity                                  |
| `REGISTRY_PASSWORD`  | Docker registry PAT or password                                      |
| `KUBECONFIG_STAGING` | kubeconfig contents for the staging cluster                          |
| `GITHUB_TOKEN`       | Explicitly passed to all reusable workflows                          |

Secrets are forwarded explicitly using the `secrets:` block when calling reusable workflows, ensuring clear ownership and secure propagation.

---

## Branching Strategy & PR Workflow

This repository follows a **production-style Git workflow**:

* Feature development occurs on branches such as:

  ```
  feature/pipeline-setup
  ```
* Changes are merged into `main` via Pull Requests only
* Branch protection rules are enabled:

  * Pull Request required before merge
  * All CI checks must pass
  * Direct pushes to `main` are disallowed

This guarantees that all code is validated by matrix tests, Docker builds, and deployment checks before release.

---

## Local Development

### Build and run locally

```bash
go build -o inventory-service .
./inventory-service
```

Test endpoints:

```bash
curl http://localhost:8080/status
curl http://localhost:8080/items
```

### Run unit tests

```bash
go test ./...
```

---

## Kubernetes (Staging)

An example deployment manifest is provided in `k8s/deployment.yaml`. During CI/CD execution, the reusable deployment workflow updates the `inventory-service` deployment in the `inventory-staging` namespace with the newly built Docker image.

---

## Optional Enhancements / Future Work

* Add code coverage reporting and badges
* Implement semantic versioning and release automation
* Enable multi-architecture Docker builds (`amd64`, `arm64`)
* Extend pipeline for per-branch or preview environments

---

## Summary

This project demonstrates **real-world CI/CD engineering practices**, including:

* Reusable GitHub Actions workflows
* Matrix testing across platforms and Go versions
* Secure secrets management
* Automated Docker builds and registry publishing
* Kubernetes deployment orchestration
* Enforced PR-based development and branch protection
* Automated notification on workflow deployment
