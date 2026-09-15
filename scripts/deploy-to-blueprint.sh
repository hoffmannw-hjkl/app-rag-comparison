#!/usr/bin/env bash
set -euo pipefail

# Déploiement "One-Click" sur le GCP AI Foundation Blueprint (GKE Autopilot)
BLUEPRINT_DIR=""

while [[ $# -gt 0 ]]; do
  case $1 in
    --blueprint-dir=*)
      BLUEPRINT_DIR="${1#*=}"
      shift
      ;;
    --blueprint-dir)
      BLUEPRINT_DIR="$2"
      shift 2
      ;;
    *)
      echo "Argument inconnu: $1"
      exit 1
      ;;
  esac
done

if [ -z "$BLUEPRINT_DIR" ]; then
  # Recherche automatique si cloné à côté
  if [ -d "../gcp-ai-foundation-blueprint" ]; then
    BLUEPRINT_DIR="../gcp-ai-foundation-blueprint"
  else
    echo "❌ Erreur: Répertoire du blueprint non spécifié."
    echo "Usage: ./scripts/deploy-to-blueprint.sh --blueprint-dir=/chemin/vers/gcp-ai-foundation-blueprint"
    exit 1
  fi
fi

echo "🔍 Lecture des outputs Terraform depuis : $BLUEPRINT_DIR"
cd "$BLUEPRINT_DIR"
PROJECT_ID=$(terraform output -raw project_id 2>/dev/null || gcloud config get-value project)
GKE_CLUSTER=$(terraform output -raw gke_cluster_name 2>/dev/null || echo "wh-djvagl-cluster")
GKE_REGION=$(terraform output -raw region 2>/dev/null || echo "europe-west1")
GKE_KSA="rag-comparison-ksa"
GKE_GSA=$(terraform output -raw gke_app_service_account_email 2>/dev/null || echo "")
GCS_RAG_BUCKET=$(terraform output -raw rag_bucket_name 2>/dev/null || echo "${PROJECT_ID}-rag-docs")

cd - >/dev/null

echo "📦 Projet GCP : $PROJECT_ID | Cluster GKE : $GKE_CLUSTER ($GKE_REGION) | Bucket RAG : $GCS_RAG_BUCKET"

echo "🔨 Construction de l'image de production Go..."
IMAGE_URI="${GKE_REGION}-docker.pkg.dev/${PROJECT_ID}/ai-demo-repo/rag-comparison:latest"
gcloud builds submit --tag "${IMAGE_URI}" src/

echo "🔐 Configuration des identités Workload Identity..."
if [ -n "$GKE_GSA" ]; then
  gcloud iam service-accounts add-iam-policy-binding "${GKE_GSA}" \
    --role="roles/iam.workloadIdentityUser" \
    --member="serviceAccount:${PROJECT_ID}.svc.id.goog[default/${GKE_KSA}]" \
    --project="${PROJECT_ID}" --quiet || true
fi

echo "🚀 Déploiement des manifests sur le cluster GKE Autopilot..."
# Connexion sécurisée au cluster
gcloud container clusters get-credentials "${GKE_CLUSTER}" --region "${GKE_REGION}" --project "${PROJECT_ID}"

# Préparation du manifest avec variables injectées
sed -e "s|\${GCP_PROJECT}|${PROJECT_ID}|g" -e "s|\${GCS_RAG_BUCKET}|${GCS_RAG_BUCKET}|g" deploy/k8s/deployment.yaml | kubectl apply -f -

echo "⏳ Attente de disponibilité du service..."
kubectl rollout status deployment/rag-comparison-demo -n default --timeout=120s

echo "✨ Déploiement terminé avec succès sur le GCP AI Foundation Blueprint !"
