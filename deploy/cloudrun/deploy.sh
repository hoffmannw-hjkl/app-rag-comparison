#!/usr/bin/env bash
set -euo pipefail

# Déploiement serverless ultra-rapide sur Google Cloud Run
GCP_PROJECT="${GCP_PROJECT:-$(gcloud config get-value project 2>/dev/null)}"
GCP_REGION="${GCP_REGION:-europe-west1}"
SERVICE_NAME="rag-comparison-demo"
IMAGE_TAG="${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT}/ai-demo-repo/rag-comparison:latest"

echo "🔨 Construction de l'image de production multi-stage..."
gcloud builds submit --project="${GCP_PROJECT}" --tag "${IMAGE_TAG}" src/

echo "🚀 Déploiement sur Cloud Run (${GCP_REGION})..."
gcloud run deploy "${SERVICE_NAME}" \
  --project="${GCP_PROJECT}" \
  --image "${IMAGE_TAG}" \
  --region "${GCP_REGION}" \
  --platform managed \
  --allow-unauthenticated \
  --set-env-vars "GCP_PROJECT=${GCP_PROJECT},GCP_REGION=${GCP_REGION},GEMINI_MODEL=gemini-3.5-flash"

URL=$(gcloud run services describe "${SERVICE_NAME}" --project="${GCP_PROJECT}" --region "${GCP_REGION}" --format='value(status.url)')
echo "✅ Déploiement Cloud Run réussi : ${URL}"
