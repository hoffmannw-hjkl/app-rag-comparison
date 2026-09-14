# RAG Comparison Demo - Classic Search vs RAG Comparative Application 🤖⚖️

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Google Cloud](https://img.shields.io/badge/Google_Cloud-Vertex_AI-4285F4?style=flat&logo=google-cloud)](https://cloud.google.com/vertex-ai)
[![Gemini 3](https://img.shields.io/badge/Gemini-3.5_Flash_%7C_3.8_Flash_%7C_3.1_Pro-8E75B2?style=flat&logo=google-gemini)](https://ai.google.dev/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

> **Production-ready pre-sales demonstration application**, designed to illustrate the immediate value of **RAG (Retrieval-Augmented Generation) and Gemini 3 Grounding** compared to standard keyword search engines.

This application is decoupled from infrastructure and ready for instant deployment on [**GCP AI Foundation Blueprint**](https://github.com/cloud-gtm/gcp-ai-foundation-blueprint) (GKE Autopilot) or **Google Cloud Run**.

---

## 🌟 Key Features

1. **Gemini 3 Generation Models Suite**:
   - **Gemini 3.5 Flash** (Default): Native multimodal reasoning, ultra-low latency, and high efficiency.
   - **Gemini 3.8 Flash**: Frontier agentic model optimized for code understanding and complex reasoning.
   - **Gemini 3.1 Pro Preview**: Deep reasoning model for exhaustive technical synthesis and architecture reviews.
   - **Gemini 3.5 Flash-Lite**: High-throughput, ultra-fast model for demanding low-latency RAG workloads.
   - **Dynamic Switching**: Change models on-the-fly directly from the header dropdown without restarting the app.
2. **Live Multi-Model Arena (Side-by-Side LLM Comparison)**:
   - **Parallel Model Evaluation**: Run and compare two different Gemini models simultaneously (e.g. `Gemini 3.5 Flash` vs `Gemini 3.8 Flash` or `Gemini 3.1 Pro`) on the exact same grounded query.
   - **Real-time Telemetry & Benchmarks**: Granular latency and consumption cards displaying time-to-first-token (TTFT in ms), total generation duration, and total token count.
3. **Before / After Effect (Classic Search vs Grounded RAG)**:
   - **Classic Search Mode**: Displays raw excerpts without synthesis (forcing the user to read through all passages).
   - **GenAI RAG Mode**: Generates an authoritative, cited synthesis with interactive references.
4. **Grounding Transparency**:
   - Live display of documents and retrieved passages before streaming the final response.
5. **On-the-fly Document Management**:
   - Sidebar panel with indexing status (`Ready` vs `Indexing`).
   - Modal to add new text, files, or URLs directly into the corpus.
6. **Ultra-Lean & Fast Architecture**:
   - Go backend compiled into a single static binary (< 25 MB).
   - Frontend embedded via `embed.FS` (zero separate Node.js server, cold start in < 1 second on Cloud Run).
   - Real-time streaming via **Server-Sent Events (SSE)**.

---

## 🏗️ Repository Structure

```
app-rag-comparison/
├── src/                           # 🧠 Application Source Code
│   ├── main.go                    # Go HTTP server (REST API, Vertex AI Gemini, SSE)
│   ├── go.mod                     # Go dependencies
│   ├── Dockerfile                 # Lightweight multi-stage container (<25MB)
│   └── web/                       # Web UI (HTML5/CSS3/Vanilla JS, Lucide icons)
│
├── deploy/                        # 📦 Deployment Manifests
│   ├── cloudrun/                  # Serverless Cloud Run deployment
│   │   └── deploy.sh
│   └── k8s/                       # GKE Autopilot manifests (Workload Identity)
│       └── deployment.yaml
│
├── scripts/                       # ⚡ Automation Scripts
│   ├── deploy-to-blueprint.sh     # One-Click deployment to GCP Blueprint
│   ├── sync-gtm.sh                # Git synchronization to official cloud-gtm repo
│   └── proxy.py                   # Local authenticated IAP / Cloud Run proxy
│
└── docs/                          # 📚 Technical Documentation
    └── DEPLOYMENT_GUIDE.md        # Complete Cloud Run & GKE deployment guide
```

---

## 🚀 Local Quickstart

```bash
cd src
export GCP_PROJECT=$(gcloud config get-value project)
export GCP_REGION="europe-west1"
export GEMINI_MODEL="gemini-2.5-flash"
export GOOGLE_OAUTH_ACCESS_TOKEN=$(gcloud auth print-access-token)

go run main.go
# Open http://localhost:8080
```

---

## 📄 License

Apache License 2.0. See [LICENSE](LICENSE) for details.
