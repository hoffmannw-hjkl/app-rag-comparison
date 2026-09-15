# RAG Comparison Demo - Application de Comparaison Recherche Classique vs RAG 🤖⚖️

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Google Cloud](https://img.shields.io/badge/Google_Cloud-Vertex_AI-4285F4?style=flat&logo=google-cloud)](https://cloud.google.com/vertex-ai)
[![Gemini 3](https://img.shields.io/badge/Gemini-3.5_Flash_%7C_3.8_Flash_%7C_3.1_Pro-8E75B2?style=flat&logo=google-gemini)](https://ai.google.dev/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

> **Application de démonstration avant-vente prête à l'emploi**, conçue pour illustrer en clientèle l'apport immédiat du **RAG (Retrieval-Augmented Generation) et du Grounding avec la suite Gemini 3** face à un moteur de recherche classique par mots-clés.

Cette application est découplée de l'infrastructure et conçue pour se déployer instantanément sur [**GCP AI Foundation Blueprint**](https://github.com/cloud-gtm/gcp-ai-foundation-blueprint) (GKE Autopilot) ou sur **Google Cloud Run**.

---

## 🌟 Fonctionnalités Clés

1. **Suite de Modèles Gemini 3 (Dernière Génération)** :
   - **Gemini 3.5 Flash** (Défaut) : Rapidité d'exécution, ultra-faible latence et efficience multimodale native.
   - **Gemini 3.8 Flash** : Modèle agentique de pointe pour l'ingénierie et les raisonnements de code complexes.
   - **Gemini 3.1 Pro Preview** : Modèle Frontier pour les analyses d'architecture denses et le raisonnement en profondeur.
   - **Gemini 3.5 Flash-Lite** : Modèle ultra-léger haute cadence pour cas d'usage à fort débit.
   - **Sélecteur Dynamique** : Bascule instantanée de modèle directement depuis la barre d'en-tête, sans redémarrage du conteneur.
2. **Comparateur de Modèles en Direct (Arena LLM)** :
   - **Bataille de Modèles Côte à Côte** : Évaluez simultanément 2 modèles différents (ex: `Gemini 3.5 Flash` vs `Gemini 3.8 Flash` ou `Gemini 3.1 Pro`) sur la même requête documentaire.
   - **Télémétrie en Temps Réel** : Benchmarks affichés sous chaque réponse : temps de réponse au premier token (TTFT en ms), durée totale de génération et nombre total de tokens consommés.
3. **Effet Avant / Après (Recherche Mots-Clés vs RAG Groundé)** :
   - **Mode Recherche Classique** : Affiche les extraits bruts sans synthèse ni compréhension (obligeant l'utilisateur à tout lire).
   - **Mode RAG GenAI** : Génère une synthèse rédigée et sourcée avec citations interactives.
4. **Transparence du Grounding** :
   - L'interface affiche en temps réel les documents et passages consultés avant de streamer la réponse finale.
5. **Gestion Documentaire à la Volée** :
   - Panneau latéral avec statut d'indexation (`Prêt` vs `En cours`).
   - Modal d'ajout immédiat de nouveaux textes, fichiers ou URLs dans le corpus.
6. **Architecture Ultra-Sobre & Rapide** :
   - Backend Go compilé en binaire statique unique (< 25 Mo).
   - Frontend intégré via `embed.FS` (zéro serveur Node.js séparé, démarrage en < 1 seconde sur Cloud Run).
   - Streaming temps réel par **Server-Sent Events (SSE)**.

---

## 📐 Architecture & Schéma GCP Draw

> 📐 **Schéma interactif GCP Draw** : Le schéma officiel au format GCP Draw (`go/gcpdraw`) est disponible dans [docs/architecture-gcpdraw.md](docs/architecture-gcpdraw.md).

## 🏗️ Structure du Dépôt

```
app-rag-comparison/
├── src/                           # 🧠 Code Source Applicatif
│   ├── main.go                    # Serveur Go (API REST, Vertex AI Gemini, SSE)
│   ├── go.mod                     # Dépendances Go
│   ├── Dockerfile                 # Conteneur multi-stage ultra-léger (<25MB)
│   └── web/                       # Interface Web (HTML5/CSS3/JS natif, Lucide icons)
│       ├── index.html
│       ├── app.js
│       └── style.css
│
├── deploy/                        # 📦 Manifests de Déploiement
│   ├── cloudrun/                  # Déploiement Serverless Cloud Run
│   │   └── deploy.sh
│   └── k8s/                       # Manifests GKE Autopilot (Workload Identity)
│       └── deployment.yaml
│
├── scripts/                       # ⚡ Scripts d'Automatisation
│   ├── deploy-to-blueprint.sh     # Déploiement "One-Click" sur le Blueprint GCP
│   ├── sync-gtm.sh                # Synchronisation Git vers le dépôt officiel cloud-gtm
│   └── proxy.py                   # Proxy local authentifié IAP / Cloud Run
│
└── docs/                          # 📚 Documentation Technique
    └── DEPLOYMENT_GUIDE.md        # Guide complet de déploiement Cloud Run & GKE
```

---

## 🚀 Démarrage Rapide en Local

```bash
# 1. Se placer dans les sources
cd src

# 2. Configurer vos identifiants GCP
export GCP_PROJECT=$(gcloud config get-value project)
export GCP_REGION="europe-west1"
export GEMINI_MODEL="gemini-2.5-flash"
export GOOGLE_OAUTH_ACCESS_TOKEN=$(gcloud auth print-access-token)

# 3. Lancer le serveur Go
go run main.go
# Accédez à http://localhost:8080
```

---

## ☁️ Déploiements

Consultez le [**Guide de Déploiement Complet**](docs/DEPLOYMENT_GUIDE.md) pour les instructions pas-à-pas :
- **Serverless Cloud Run** : `./deploy/cloudrun/deploy.sh`
- **GKE Autopilot sur Blueprint** : `./scripts/deploy-to-blueprint.sh --blueprint-dir=../gcp-ai-foundation-blueprint`

---

## 🔒 Sécurité

- **Workload Identity** natif : Zéro clé statique ou secrète stockée dans le conteneur.
- **Accès IAP / WAF** : Compatible avec Cloud Armor et Google Cloud Identity-Aware Proxy.

---

## 📄 Licence

Apache License 2.0. Voir [LICENSE](LICENSE) pour plus de détails.
