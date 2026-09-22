---
name: rag-benchmark-and-ci
description: >-
  Formats Go source files, executes unit & mock-server tests, and verifies Cloud Run Direct VPC Egress deployments for app-rag-comparison.
  Use when the user asks to "run go tests", "fix gofmt CI failure", "test Agentic RAG",
  "evaluate RAG groundedness", or "deploy rag-comparison-demo to Cloud Run".
---

# Go Format, Test & RAG Benchmark Runbook (`M1L1 Skills Framework` Compliant)

**Skill Patterns Combined**: *Tool Wrapper Skill* (Slide 13) • *Generator / Experience Before Theory* (Slide 14) • *Reviewer Checklist* (Slide 15) • *Workflow Skill* (Slide 16).

---

## 1. Tool Wrapper & Usage Recipes (Slide 13)

Use the standard Go toolchain inside the `src/` directory and filter verbose test logs when running full suite checks:

```bash
GO_BIN="go"
GOFMT_BIN="gofmt"
```

### Usage
`go test -v ./...` emits verbose logs for all 30+ PDF, GCS mock, and SSE tests. To inspect only failures or target a specific test, use `-run`:

```bash
cd src && $GO_BIN test -v -run "TestDecomposeQueryAgent|TestHandleChatStreamAgenticMode" ./...
```

---

## 2. Experience Before Theory — Known Gotchas (Slide 14)

The following real-world gotchas were discovered in production and must always be prevented:

1. **GitHub Actions `gofmt -l .` Hard Gate (`src/main_test.go`)**:
   - *Gotcha*: The GitHub Actions workflow runs `unformatted=$(gofmt -l .)` and exits with code `1` if even a single blank line or indentation tab is unformatted.
   - *Rule*: **Never commit `.go` files without running `gofmt -w .` first.**
2. **Direct `*.run.app` URL Returning 404 (`--ingress=internal-and-cloud-load-balancing`)**:
   - *Gotcha*: Because `rag-comparison-demo` enforces `--ingress=internal-and-cloud-load-balancing` (so users cannot bypass Cloud Armor WAF and IAP), calling `https://rag-comparison-demo-*.run.app` directly from outside the VPC returns a Google Frontend `404 Not Found`.
   - *Rule*: Always access the public service via `https://rag.hoffmannw.demo.altostrat.com` or inspect startup readiness via `gcloud logging read`.
3. **Answer Truncation on Multi-Document RAG**:
   - *Gotcha*: Default `maxOutputTokens` (`1024`) truncates rich markdown comparisons.
   - *Rule*: Always keep `maxOutputTokens: 8192` and `topKChunks = 8` with a 120s context timeout.

---

## 3. Multi-Step Pre-Commit & Deployment Workflow (Slide 16)

Run the bundled verification script before every commit:

```bash
./.agents/skills/rag-benchmark-and-ci/scripts/verify.sh
```

### Step-by-Step Procedure
1. **Format & Check**:
   ```bash
   cd src && gofmt -w . && test -z "$(gofmt -l .)"
   ```
2. **Execute Deterministic Unit Tests**:
   ```bash
   cd src && go test ./...
   ```
3. **Cloud Run Deployment (with Least-Privilege SA & Direct VPC Egress)**:
   ```bash
   ./deploy/cloudrun/deploy.sh
   ```

---

## 4. Reviewer Assessment Checklist (Slide 15)

Before merging any change to `app-rag-comparison`, verify:
- [ ] `gofmt -l src/` returns zero files.
- [ ] Every access to `s.documents` acquires `s.mu.RLock()` or `s.mu.Lock()`.
- [ ] All user-supplied strings in `src/web/static/index.html` are wrapped with `esc(...)` to prevent DOM XSS.
- [ ] Cloud Run uses `ai-demo-2e2m-gke-ai-sa@wh-ai-blueprint-a363.iam.gserviceaccount.com` and Direct VPC Egress (`--network=ai-demo-2e2m-vpc --subnet=ai-demo-2e2m-subnet`).
