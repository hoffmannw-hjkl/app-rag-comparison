---
name: rag-benchmark-and-ci
description: >-
  Pre-commit verification and CI validation runbook for app-rag-comparison.
  Enforces gofmt formatting, executes the full Go unit test suite, and validates Cloud Run Direct VPC Egress deployment configuration.
---

# Go Format, Test & RAG Benchmark Runbook

Always execute this runbook in `src/` before committing any change to `app-rag-comparison` to guarantee a 100% green GitHub Actions pipeline.

## Step 1: Enforce `gofmt` Formatting

GitHub Actions fails immediately if `gofmt -l .` returns any file. Run:

```bash
cd src
gofmt -w .
if [ -n "$(gofmt -l .)" ]; then
  echo "❌ Unformatted Go files remaining!"
  exit 1
fi
```

## Step 2: Run Unit & Mock Server Tests

Run the entire test suite (PDF extraction, CMap decoding, French elisions, GCS snapshot mock server, hybrid cosine similarity, and GenAI evaluation):

```bash
cd src
go test -v ./...
```

## Step 3: Cloud Run Production Verification (`deploy/cloudrun/deploy.sh`)

When deploying to `wh-ai-blueprint-a363`, verify that the service uses:
- Least-privilege Service Account: `ai-demo-2e2m-gke-ai-sa@wh-ai-blueprint-a363.iam.gserviceaccount.com`
- Direct VPC Egress: `--network=ai-demo-2e2m-vpc --subnet=ai-demo-2e2m-subnet --vpc-egress=all-traffic`
- Internal + Load Balancing Ingress: `--ingress=internal-and-cloud-load-balancing`
