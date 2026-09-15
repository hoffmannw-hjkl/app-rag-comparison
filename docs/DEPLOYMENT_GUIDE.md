# Guide de Déploiement - Application RAG Comparison Demo 🤖⚖️

Ce guide explique pas-à-pas comment déployer et exploiter l'application **RAG Comparison Demo** sur l'infrastructure Google Cloud, soit en mode **Serverless (Cloud Run)** soit sur le **GCP AI Foundation Blueprint (GKE Autopilot)**.

---

## 🏛️ Architecture de Déploiement

```
                          ┌──────────────────────────┐
                          │   Navigateur Client      │
                          │   (Interface Web SSE)    │
                          └─────────────┬────────────┘
                                        │ HTTPS
                     ┌──────────────────┴──────────────────┐
                     │                                     │
          [Option 1: Serverless]               [Option 2: Enterprise]
          ┌─────────────────────┐             ┌─────────────────────┐
          │   Google Cloud Run  │             │   Google Cloud Armor│
          │   (Service Géré)    │             │   + Cloud IAP / WAF │
          └──────────┬──────────┘             └──────────┬──────────┘
                     │                                   │ Ingress
                     │                        ┌──────────┴──────────┐
                     │                        │   GKE Autopilot     │
                     │                        │   Pod (rag-chatbot) │
                     │                        └──────────┬──────────┘
                     │                                   │ Workload Identity
                     └─────────────────┬─────────────────┘
                                       │ Vertex AI API
                          ┌────────────┴────────────┐
                          │   Gemini 3 Suite        │
                          │   (3.5/3.8 Flash, 3.1)  │
                          │   + Grounding Search    │
                          └─────────────────────────┘
```

---

## 🚀 Option 1 : Déploiement Serverless Cloud Run (Recommandé Avant-Vente)

Pour les démos clients rapides et les ateliers techniques :

```bash
# 1. Configurer vos variables d'environnement
export GCP_PROJECT="votre-projet-gcp"
export GCP_REGION="europe-west1"

# 2. Exécuter le script automatisé
./deploy/cloudrun/deploy.sh
```

Le script :
1. Compile le conteneur Go statique via Google Cloud Build (< 30s).
2. Déploie le service sur Cloud Run avec l'autoscaling de 0 à 5 instances.
3. Configure l'accès direct en streaming SSE.

---

## ☸️ Option 2 : Déploiement sur le GCP AI Foundation Blueprint (GKE Autopilot)

Pour les déploiements d'entreprise haute sécurité et souverains :

### Prérequis
- Avoir déployé le blueprint Terraform : [`cloud-gtm/gcp-ai-foundation-blueprint`](https://github.com/cloud-gtm/gcp-ai-foundation-blueprint).

### Déploiement One-Click
```bash
./scripts/deploy-to-blueprint.sh --blueprint-dir=/chemin/vers/gcp-ai-foundation-blueprint
```

Ce script extrait automatiquement les outputs Terraform (`project_id`, `gke_cluster_name`, `gke_app_service_account_email`), associe les rôles Workload Identity (`roles/iam.workloadIdentityUser`), build l'image conteneur et applique les manifests Kubernetes.

---

## 🔒 Sécurité & IAM

- **Workload Identity** : Aucun secret ou clé de service account JSON n'est injecté dans le conteneur.
- **Rôles GCP requis** :
  - `roles/aiplatform.user` : Pour interroger Gemini sur Vertex AI.
  - `roles/artifactregistry.reader` : Pour récupérer les images de conteneurs.
