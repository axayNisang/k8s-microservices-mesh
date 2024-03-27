# Cloud Microservices on Kubernetes with gRPC & Observability 🚀

[![K8s CI](https://github.com/axayNisang/k8s-microservices-mesh/actions/workflows/k8s-ci.yml/badge.svg)](https://github.com/axayNisang/k8s-microservices-mesh/actions)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.28+-326CE5.svg?logo=kubernetes)](https://kubernetes.io)
[![Helm](https://img.shields.io/badge/Helm-v3-0F1689.svg?logo=helm)](https://helm.sh)
[![Go](https://img.shields.io/badge/Go-1.21-00ADD8.svg?logo=go)](https://golang.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0-3178C6.svg?logo=typescript)](https://www.typescriptlang.org)

An enterprise-grade, polyglot microservices platform architected for Kubernetes. Demonstrates production patterns including gRPC inter-service communication, Horizontal Pod Autoscaling (HPA), zero-trust NetworkPolicies, non-root scratch containers, and automated deployment via Helm 3.

---

## 🏗 System Architecture

```
                    Internet Clients
                           │
                           ▼ [HTTPS :443]
              +--------------------------+
              |    NGINX Ingress Pods    |
              +--------------------------+
                │                      │
   /orders      ▼                      ▼      /auth
+────────────────────────+     +────────────────────────+
|      Order Service     |     |      Auth Service      |
|  (Node.js/TypeScript)  |     |         (Go)           |
|   Port 8083 | HPA 3-15 |     |       Port 8081        |
+────────────────────────+     +────────────────────────+
         │                                  ▲
         │ (Internal gRPC / Private Mesh)   │
         └──────────────────────────────────┘
```

---

## ⚡ Key Highlights & Production Engineering

1. **Polyglot Microservices**:
   - **Auth Service** (*Go*): High-performance JWT issuance and token cryptographic validation, compiling to a secure, zero-CVE scratch container running as an unprivileged UID.
   - **Catalog Service** (*Go*): Item catalog with sync primitives and fast in-memory query endpoints.
   - **Order Service** (*TypeScript / Node.js*): Business workflow engine handling state transitions, Prometheus `/metrics` exposition, and upstream auth calls.
2. **Kubernetes Resiliency & Autoscaling**:
   - Dynamic **Horizontal Pod Autoscaling (HPA)** automatically scales Order Service pods from 3 to 15 based on real-time CPU (75%) and Memory (80%) thresholds with stabilization windows to avoid flapping.
   - Strict `readinessProbe` and `livenessProbe` definitions to guarantee zero downtime during rolling upgrades.
3. **Zero-Trust Network Hardening**:
   - Default deny-all `NetworkPolicy` isolates pods within the namespace.
   - Explicit ingress whitelist rules ensure only the Order Service can communicate with the internal Auth Service.
4. **Declarative Helm 3 Infrastructure**:
   - Parameterized chart supporting multi-environment deployments (`staging` vs `production`).

---

## 📂 Repository Structure

```
├── proto/                     # Protocol Buffer interface contracts (auth.proto, catalog.proto)
├── services/
│   ├── auth-service/          # Go authentication service & scratch Dockerfile
│   ├── catalog-service/       # Go product catalog service
│   └── order-service/         # TypeScript order processing service & metrics
├── k8s/                       # Raw Kubernetes manifests (Deployments, Services, HPA, Ingress)
├── helm/                      # Production Helm 3 chart (Chart.yaml, values.yaml, templates)
└── .github/workflows/         # Automated Helm lint and container validation CI
```

---

## 🚀 Local Development & Cluster Deployment

### Prerequisites
- Docker & Kubernetes cluster (`kind`, `minikube`, or AWS EKS)
- [Helm v3](https://helm.sh/docs/intro/install/)

### Deploying via Helm

```bash
# Add or update local Helm dependencies
cd helm

# Dry-run render templates
helm template cloud-mesh . -f values.yaml

# Install or upgrade release onto cluster
helm upgrade --install cloud-mesh . --namespace microservices --create-namespace
```

### Checking Pod Health & Autoscaling

```bash
kubectl get pods -n microservices -o wide
kubectl get hpa -n microservices
```

---

## 📄 License
MIT License. Created by [Akshay Nisang](https://github.com/axayNisang).
