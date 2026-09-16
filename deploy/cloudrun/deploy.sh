#!/usr/bin/env bash
set -euo pipefail

# Déploiement serverless sur Google Cloud Run
GCP_PROJECT="${GCP_PROJECT:-wh-ai-blueprint-a363}"
GCP_REGION="${GCP_REGION:-europe-west1}"
SERVICE_NAME="rag-comparison-demo"
IMAGE_TAG="${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT}/ai-demo-repo/rag-comparison:latest"

GCS_RAG_BUCKET="${GCS_RAG_BUCKET:-wh-ai-blueprint-a363-ai-demo-2e2m-rag-docs}"

# Service account dédié de moindre privilège issu de gcp-ai-foundation-blueprint
RUNTIME_SA="${RUNTIME_SA:-ai-demo-2e2m-gke-ai-sa@${GCP_PROJECT}.iam.gserviceaccount.com}"
VPC_NETWORK="${VPC_NETWORK:-ai-demo-2e2m-vpc}"
SUBNET_NAME="${SUBNET_NAME:-ai-demo-2e2m-subnet}"

echo "🔨 Construction de l'image de production multi-stage..."
gcloud builds submit --project="${GCP_PROJECT}" --tag "${IMAGE_TAG}" src/

echo "🚀 Déploiement sur Cloud Run (${GCP_REGION})..."

# Notes d'architecture :
#
# --ingress=internal-and-cloud-load-balancing
#     L'URL *.run.app ne doit pas être joignable directement : tout le trafic doit
#     transiter par le Load Balancer afin de passer par IAP et Cloud Armor. Sans
#     cette restriction, l'authentification et le WAF sont contournables.
#
# --no-allow-unauthenticated
#     L'authentification des utilisateurs est assurée par IAP en amont.
#
# --network / --subnet / --vpc-egress=all-traffic
#     Direct VPC Egress natif Cloud Run v2 vers le VPC privé de la fondation IA.
#
# --max-instances=5
#     Le corpus documentaire est persisté sur Cloud Storage (gs://${GCS_RAG_BUCKET}/index/corpus.json)
#     et synchronisé à chaque écriture. L'auto-scaling multi-instances est actif.
#
# --no-cpu-throttling
#     Garantit que les traitements engagés disposent de CPU jusqu'à leur terme,
#     y compris après l'envoi des premiers octets d'une réponse en streaming.
DEPLOY_ARGS=(
  --project="${GCP_PROJECT}"
  --image "${IMAGE_TAG}"
  --region "${GCP_REGION}"
  --platform managed
  --ingress=internal-and-cloud-load-balancing
  --no-allow-unauthenticated
  --network="${VPC_NETWORK}"
  --subnet="${SUBNET_NAME}"
  --vpc-egress=all-traffic
  --service-account="${RUNTIME_SA}"
  --max-instances=5
  --no-cpu-throttling
  --memory=1Gi
  --timeout=300s
  --set-env-vars "GCP_PROJECT=${GCP_PROJECT},GCP_REGION=${GCP_REGION},GEMINI_MODEL=gemini-3.5-flash,GCS_RAG_BUCKET=${GCS_RAG_BUCKET}"
)


gcloud run deploy "${SERVICE_NAME}" "${DEPLOY_ARGS[@]}"

URL=$(gcloud run services describe "${SERVICE_NAME}" --project="${GCP_PROJECT}" --region "${GCP_REGION}" --format='value(status.url)')
echo "✅ Déploiement Cloud Run réussi : ${URL}"
echo "ℹ️  Accès utilisateur via le Load Balancer : https://rag.hoffmannw.demo.altostrat.com"
