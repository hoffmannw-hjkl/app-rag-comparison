# Guide de Déploiement - RAG Comparison Demo

Ce guide décrit les procédures de déploiement et d'exploitation de l'application **RAG Comparison Demo** sur Google Cloud Platform, en mode **Serverless (Cloud Run)** ou sur cluster **GKE Autopilot** au sein du blueprint [GCP AI Foundation Blueprint](https://github.com/cloud-gtm/gcp-ai-foundation-blueprint).

---

## Architecture de Déploiement

Consultez le schéma interactif officiel au format GCP Draw dans [docs/architecture-gcpdraw.md](architecture-gcpdraw.md).

```
                          ┌──────────────────────────┐
                          │   Navigateur Client      │
                          │   (Interface Web SSE)    │
                          └─────────────┬────────────┘
                                        │ HTTPS
                     ┌──────────────────┴──────────────────┐
                     │                                     │
          [Option 1: Serverless]               [Option 2: Enterprise Blueprint]
          ┌─────────────────────┐             ┌─────────────────────┐
          │   Google Cloud Run  │             │   Cloud Armor WAF   │
          │   (Service Géré)    │             │   + Cloud IAP       │
          └──────────┬──────────┘             └──────────┬──────────┘
                     │                                   │ Ingress
                     │                        ┌──────────┴──────────┐
                     │                        │   GKE Autopilot     │
                     │                        │   Pod (rag-chatbot) │
                     │                        └──────────┬──────────┘
                     │                                   │ Workload Identity
                     └─────────────────┬─────────────────┘
                                       │
                      ┌────────────────┴────────────────┐
                      │                                 │
                      ▼                                 ▼
          ┌────────────────────────┐        ┌────────────────────────┐
          │   Vertex AI API        │        │   Cloud Storage        │
          │   • Gemini 3.5 / 3.8   │        │   • gs://{BUCKET}/     │
          │   • text-embedding-002 │        │     index/corpus.json  │
          │   • Rapid Evaluation   │        └────────────────────────┘
          └────────────────────────┘
```

---

## Prérequis Généraux

1. Projet Google Cloud actif avec facturation activée.
2. APIs Google Cloud activées :
   ```bash
   gcloud services enable \
     run.googleapis.com \
     aiplatform.googleapis.com \
     storage.googleapis.com \
     cloudbuild.googleapis.com \
     artifactregistry.googleapis.com
   ```
3. Permissions IAM minimales sur le compte de service d'exécution :
   - `roles/aiplatform.user` : Exécution des inférences Gemini, vectorisation et Rapid Evaluation.
   - `roles/storage.objectAdmin` : Lecture et écriture du snapshot de corpus dans le bucket GCS.

---

## Option 1 : Déploiement sur Google Cloud Run

Ce mode est recommandé pour les environnements de démonstration, les validations de concept et les ateliers clients.

### 1. Variables d'Environnement

Configurez les variables requises avant l'exécution du déploiement :

```bash
export GCP_PROJECT="$(gcloud config get-value project)"
export GCP_REGION="europe-west1"
export GCS_RAG_BUCKET="${GCP_PROJECT}-rag-docs"
```

### 2. Déploiement Automatisé

Exécutez le script de déploiement Cloud Run :

```bash
./deploy/cloudrun/deploy.sh
```

Le script exécute les opérations suivantes :
1. Crée le bucket Cloud Storage `gs://${GCS_RAG_BUCKET}` si nécessaire.
2. Déclenche la compilation du conteneur multi-stage via Google Cloud Build.
3. Déploie le service `rag-comparison-demo` sur Cloud Run avec les paramètres :
   - `--memory=1Gi` et `--cpu=1`
   - `--min-instances=0` et `--max-instances=5`
   - `--timeout=300s` (indispensable pour les flux de streaming longs et l'indexation)
   - `--ingress=internal-and-cloud-load-balancing` (ou `--ingress=all` selon la politique réseau du projet)

### 3. Accès Authentifié IAP en Local

Lorsque le service Cloud Run est protégé par Cloud Load Balancing et Identity-Aware Proxy (IAP), lancez le proxy local d'authentification :

```bash
python3 scripts/proxy.py
# Accédez ensuite à l'application sur http://localhost:8085
```

---

## Option 2 : Déploiement sur GKE Autopilot (GCP AI Foundation Blueprint)

Ce mode est recommandé pour les architectures d'entreprise nécessitant un périmètre de sécurité renforcé (VPC Service Controls, Private Service Connect).

### 1. Prérequis Blueprint
Déployez au préalable l'infrastructure de base via Terraform :
```bash
git clone https://github.com/cloud-gtm/gcp-ai-foundation-blueprint.git
cd gcp-ai-foundation-blueprint
terraform init && terraform apply
```

### 2. Déploiement Applicatif One-Click
Revenez dans le répertoire du projet `app-rag-comparison` et exécutez le script d'automatisation :

```bash
./scripts/deploy-to-blueprint.sh --blueprint-dir=/chemin/vers/gcp-ai-foundation-blueprint
```

Ce script réalise automatiquement les étapes suivantes :
1. Lecture des outputs Terraform (`project_id`, `gke_cluster_name`, `gke_app_service_account_email`).
2. Configuration de l'accès kubectl au cluster GKE Autopilot régional.
3. Liaison IAM Workload Identity entre le Kubernetes Service Account (`rag-comparison-sa`) et le Google Service Account.
4. Construction et publication de l'image de conteneur sur Google Artifact Registry.
5. Déploiement du manifest Kubernetes (`deploy/k8s/deployment.yaml`).

---

## Vérification et Tests Post-Déploiement

### 1. Sonde de Santé
Interrogez le point de terminaison `/healthz` pour valider la disponibilité :
```bash
curl -f https://<VOTRE-DOMAINE-OU-URL-CLOUDRUN>/healthz
```

### 2. Test du Streaming RAG
Effectuez un appel de streaming SSE pour vérifier la chaîne complète (retrieval hybride + inférence Vertex AI) :
```bash
curl -N -H "Authorization: Bearer $(gcloud auth print-access-token)" \
  "https://<VOTRE-DOMAINE>/api/chat/stream?q=Quel+est+le+r%C3%B4le+du+module+security-waf+%3F&model=gemini-3.5-flash"
```

### 3. Test de l'Évaluation GenAI
Lancez une requête POST sur `/api/evaluate` :
```bash
curl -X POST -H "Authorization: Bearer $(gcloud auth print-access-token)" \
  -H "Content-Type: application/json" \
  -d '{"query":"Quelle est la politique de rétention ?","prediction":"La rétention est de 30 jours.","context":"Politique de rétention : 30 jours."}' \
  "https://<VOTRE-DOMAINE>/api/evaluate"
```

---

## Dépannage

| Problème | Cause potentielle | Solution |
| :--- | :--- | :--- |
| Erreur HTTP 403 sur `/api/chat/stream` | Compte de service sans rôle Vertex AI | Attribuer `roles/aiplatform.user` au Service Account Cloud Run ou GKE. |
| Erreur HTTP 403 sur snapshot GCS | Compte de service sans accès au bucket | Attribuer `roles/storage.objectAdmin` sur `gs://${GCS_RAG_BUCKET}`. |
| Coupure de streaming SSE | Timeout de proxy ou de Cloud Run trop court | Vérifier `--timeout=300s` sur Cloud Run et `proxy_read_timeout 300s` sur le Load Balancer. |
| Réponses tronquées | Plafond `maxOutputTokens` ou taille de snippet trop faible | La configuration par défaut alloue 8 192 tokens de sortie et 8 passages complets de grounding. |

