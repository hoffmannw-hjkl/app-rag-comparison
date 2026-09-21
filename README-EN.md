# RAG Comparison Demo

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Google Cloud](https://img.shields.io/badge/Google_Cloud-Vertex_AI-4285F4?style=flat&logo=google-cloud)](https://cloud.google.com/vertex-ai)
[![Gemini Models](https://img.shields.io/badge/Gemini-3.5_Flash_%7C_3.8_Flash_%7C_3.1_Pro-8E75B2?style=flat&logo=google-gemini)](https://ai.google.dev/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Technical demonstration web application comparing keyword-based lexical search with grounded Retrieval-Augmented Generation (RAG) using the Google Gemini 3 model family on Vertex AI.

The application runs as a standalone Go binary compiled with zero external dependencies (`0 dependencies`, `0 CVEs`). It is deployed on **Google Cloud Run** or **Google Kubernetes Engine (GKE) Autopilot** as part of the [GCP AI Foundation Blueprint](https://github.com/cloud-gtm/gcp-ai-foundation-blueprint) architecture.

---

## Features

### 1. Hybrid Search (Semantic & Lexical)
- **Dense embeddings**: Vectorization using the Vertex AI `text-multilingual-embedding-002` model (768-dimensional vectors).
- **Native Go cosine similarity**: High-performance in-memory vector comparison without external vector databases or CGO dependencies.
- **Combined ranking**: Weighted scoring formula $0.65 \times \text{CosineSimilarity} + 0.35 \times \text{LexicalScore}$ (normalized BM25).
- **Automatic fallback**: If the Vertex AI Embedding API is unreachable or rate-limited, the engine automatically falls back to full lexical search with linguistic normalization (accent folding, French elision handling, stop words filtering).

### 2. Gemini 3 Model Suite & Multi-Model Arena
- **Supported models**:
  - `gemini-3.5-flash` (default): Ultra-low latency, high throughput, and native multimodal reasoning.
  - `gemini-3.8-flash`: Frontier agentic capabilities and code-level technical reasoning.
  - `gemini-3.1-pro-preview`: Deep reasoning model for complex architectural analysis and dense synthesis.
  - `gemini-3.5-flash-lite`: High-frequency, cost-sensitive operational workloads.
- **Dynamic model switching**: Hot-swap active models via `/api/model/switch` without container restarts.
- **Display modes**:
  - **Standard RAG**: Unified grounded response with retrieved source cards and telemetry.
  - **Split View**: Direct side-by-side comparison between raw lexical excerpts and grounded RAG synthesis.
  - **Model Arena (Dual-Model)**: Concurrent side-by-side inference with two distinct Gemini models on the exact same query.
  - **Triple View**: Three-column comparison (Model A, Model B, Classical Lexical Search).
- **Real-time telemetry**: In-stream tracking of time-to-first-token (TTFT in ms), total generation latency, and exact token counts parsed from Vertex AI responses.

### 3. On-Demand GenAI Evaluation (Vertex AI Autorater)
- **Unit-level evaluation**: Interactive evaluation button under every generated model response.
- **Calculated metrics**:
  - **Groundedness**: Measures factual alignment between the generated response and the retrieved document chunks (rating from 1 to 5).
  - **QA Relevance**: Evaluates how well the response directly addresses the user's query (rating from 1 to 5).
- **Inspection reports**: Interactive modal displaying numerical ratings, textual rationale produced by the Vertex AI Autorater, and exact contextual excerpts evaluated.
- **Heuristic fallback**: In offline or sandbox environments without Rapid Evaluation quotas, a deterministic n-gram overlap heuristic provides continuous feedback.

### 4. Cloud Storage Persistence & Horizontal Scaling
- **Serialized index**: Uploaded documents and their embeddings are persisted as JSON snapshots in `gs://{GCS_RAG_BUCKET}/index/corpus.json`.
- **Multi-instance support**: Cloud Run instances scale horizontally (`--max-instances=5`) with consistent corpus synchronization.
- **Instant cold starts**: Snapshots are restored from GCS in under one second during instance startup.

---

## Architecture

Review the complete architectural diagram in GCP Draw format in [docs/architecture-gcpdraw.md](docs/architecture-gcpdraw.md).

```
                            [ Client Browser ]
                                     │
                   HTTPS / Identity-Aware Proxy (IAP)
                                     │
                                     ▼
                     [ Google Cloud Load Balancer ]
                                     │
                                     ▼
                     [ Cloud Run : rag-comparison ]
                     ┌────────────────────────────┐
                     │ • Native Go HTTP Server    │
                     │ • Hybrid Search Engine     │
                     │ • Static Assets (embed.FS) │
                     └──────┬──────────────┬──────┘
                            │              │
           Vectorization &  │              │  Corpus Storage
           LLM Inference    │              │  & Snapshots
                            ▼              ▼
                    [ Vertex AI ]    [ Cloud Storage ]
                    • Gemini 3.5/3.8 • gs://{BUCKET}/index/
                    • Embeddings 002
                    • Rapid Eval
```

---

## REST API Reference

All endpoints are served by the Go HTTP server on the configured port (`PORT`, default `8080`).

| Method | Endpoint | Parameters | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/documents` | None | Returns the list of all indexed documents and metadata (title, size, chunks). |
| `POST` | `/api/documents/upload` | Multipart form (`files[]`) | Uploads and indexes new documents (PDF, text, markdown) with embeddings. |
| `DELETE` | `/api/documents?id={id}` | `id` (document ID) | Deletes a document from the corpus and updates the Cloud Storage snapshot. |
| `DELETE` | `/api/documents?all=true` | `all=true` | Purges all documents from memory and removes the Cloud Storage snapshot. |
| `GET` | `/api/search/classic` | `q={query}` | Performs keyword-based lexical search and returns raw snippet excerpts. |
| `GET` | `/api/chat/stream` | `q={query}`, `model={model_id}` | Streams the grounded RAG response via Server-Sent Events (SSE). |
| `GET` | `/api/models` | None | Lists available Gemini models, specifications, and the active default model. |
| `POST` | `/api/model/switch` | JSON body `{"model": "id"}` | Dynamically changes the active default Gemini model. |
| `POST` | `/api/evaluate` | JSON body `{"query", "prediction", "context"}` | Runs an evaluation for Groundedness and QA Relevance via Vertex AI Rapid Evaluation. |
| `GET` | `/healthz` | None | Liveness and readiness health check probe for Cloud Run and Kubernetes. |

---

## Environment Variables

Configure application behavior using the following environment variables:

| Variable | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `GCP_PROJECT` (or `PROJECT_ID`) | String | Auto-detected via GCP metadata | Google Cloud project ID hosting Vertex AI. |
| `GCP_REGION` (or `REGION`) | String | `europe-west1` | Target Google Cloud region for regional Vertex AI calls. |
| `GEMINI_MODEL` | String | `gemini-3.5-flash` | Default Gemini model configured at startup. |
| `GCS_RAG_BUCKET` | String | None (optional) | Cloud Storage bucket name for corpus snapshot persistence (`corpus.json`). |
| `PORT` | Integer | `8080` | TCP port on which the HTTP server listens. |

---

## Local Development

### Prerequisites
- Go 1.25 or higher installed.
- Google Cloud SDK (`gcloud`) authenticated with access to a Google Cloud project.
- Required IAM roles: `roles/aiplatform.user` and `roles/storage.objectAdmin`.

### Run Locally

```bash
# 1. Navigate to the source directory
cd src

# 2. Set environment variables (GCP_PROJECT or PROJECT_ID)
export GCP_PROJECT=$(gcloud config get-value project)
export GCP_REGION="europe-west1"
export GEMINI_MODEL="gemini-3.5-flash"
export GCS_RAG_BUCKET="${GCP_PROJECT}-rag-docs"

# 3. Run unit tests
go test -v ./...

# 4. Start the server
go run main.go
```

Access the web interface at `http://localhost:8080`.

---

## Deployment

Refer to the [Deployment Guide](docs/DEPLOYMENT_GUIDE.md) for step-by-step production deployment procedures.

### Option 1: Google Cloud Run (Serverless)

```bash
./deploy/cloudrun/deploy.sh
```

### Option 2: GKE Autopilot (GCP AI Foundation Blueprint)

```bash
./scripts/deploy-to-blueprint.sh --blueprint-dir=../gcp-ai-foundation-blueprint
```

---

## Security & Governance

- **Zero hardcoded credentials**: The container image contains no service account keys or static secrets. Authentication is managed exclusively through **Workload Identity** (GKE) and Cloud Run execution service accounts via instance metadata (`http://metadata.google.internal`).
- **Network perimeter protection**: Recommended deployment with `--ingress=internal-and-cloud-load-balancing`, protected by Cloud Armor WAF and Identity-Aware Proxy (IAP).
- **Zero third-party runtime dependencies**: The Go backend imports only the Go standard library (`0 direct dependencies`, `0 CVEs` in container vulnerability scans).

---

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
