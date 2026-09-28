# Specialized AI Agents & Skills — `app-rag-comparison`

This repository implements the **Google Cloud Elevate 2026 & EMEA SPARK** AI-Native Software Engineering architecture. It embeds specialized subagents (`.agents/agents/`) and procedural skills (`.agents/skills/`) discovered automatically by **Jetski**, **Antigravity**, and **Gemini CLI**.

---

## 🤖 Specialized Repository Subagents (`.agents/agents/`)

| Subagent Name | Role & Specialization | When to Invoke (`invoke_subagent`) |
| :--- | :--- | :--- |
| **[`rag-eval-scientist`](.agents/agents/rag-eval-scientist.md)** | **RAG Evaluation & Hybrid Retrieval Scientist** | When tuning hybrid search (768d Vertex AI embeddings + BM25 lexical scoring), configuring Vertex AI GenAI Evaluation (`Groundedness` & `Relevance`), or generating synthetic test datasets. |
| **[`go-concurrency-reviewer`](.agents/agents/go-concurrency-reviewer.md)** | **Senior Go Systems & Concurrency Reviewer** | Before committing changes to `src/main.go` or `src/main_test.go`: audits `sync.RWMutex` thread-safety, SSE streaming lifecycle, and `gofmt` compliance. |

---

## 🛠️ Repository Skills (`.agents/skills/`)

| Skill Name | Path | Description |
| :--- | :--- | :--- |
| **`rag-benchmark-and-ci`** | [`.agents/skills/rag-benchmark-and-ci/SKILL.md`](.agents/skills/rag-benchmark-and-ci/SKILL.md) | Pre-commit validation runbook (`gofmt -w .`, `go test -v -race ./...`, Zero-Dependency check, and Cloud Run Direct VPC Egress checks). |

---

## 🎬 Quick Demo Prompts (How to Trigger Each Agent Live)

- **Trigger Runtime 4-Subagent CRAG Swarm (UI / SSE)**:
  Click **`🤖 Agentic RAG`** in the web UI (`https://rag.hoffmannw.demo.altostrat.com`) or stream via `curl`:
  ```bash
  curl -N "https://rag.hoffmannw.demo.altostrat.com/api/chat/stream?mode=agentic&q=Architecture+Cloud+Armor+et+RRF"
  ```
- **Trigger `rag-eval-scientist` (IDE / CLI)**:
  > `"Invoke rag-eval-scientist to audit the Reciprocal Rank Fusion (k=60) and GraderCriticAgent self-healing loop in src/main.go."`
- **Trigger `go-concurrency-reviewer` (IDE / CLI)**:
  > `"Invoke go-concurrency-reviewer to audit sync.RWMutex usage and SSE http.Flusher context cancellation in src/main.go."`
- **Run the M1L1 Gatekeeper Script**:
  ```bash
  ./.agents/skills/rag-benchmark-and-ci/scripts/verify.sh
  ```

