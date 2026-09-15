#!/usr/bin/env bash
set -euo pipefail

# Script de synchronisation vers le dépôt officiel cloud-gtm
# Permet de regrouper les commits, créer une branche sync/batch-*,
# ouvrir une PR et l'auto-fusionner.

GTM_REMOTE="gtm"
GITHUB_REMOTE="origin"
MAIN_BRANCH="main"

echo "🔄 Vérification des remotes Git..."
if ! git remote | grep -q "^${GTM_REMOTE}$"; then
  echo "ℹ️ Remote '${GTM_REMOTE}' introuvable. Tentative de configuration..."
  git remote add "${GTM_REMOTE}" "https://github.com/cloud-gtm/app-rag-comparison.git" || true
fi

echo "🚀 Récupération des dernières modifications..."
git fetch "${GITHUB_REMOTE}" "${MAIN_BRANCH}" || true

BATCH_ID=$(date +%Y%m%d-%H%M%S)
SYNC_BRANCH="sync/batch-${BATCH_ID}"

echo "🌿 Création de la branche de synchronisation: ${SYNC_BRANCH}..."
git checkout -b "${SYNC_BRANCH}"

echo "📤 Push vers ${GTM_REMOTE}..."
if git push "${GTM_REMOTE}" "${SYNC_BRANCH}"; then
  echo "📬 Création de la Pull Request..."
  PR_URL=$(gh pr create \
    --repo cloud-gtm/app-rag-comparison \
    --head "${SYNC_BRANCH}" \
    --base "${MAIN_BRANCH}" \
    --title "sync: update app-rag-comparison from local workspace (${BATCH_ID})" \
    --body "Automated batch sync of application code, deployment manifests, and documentation.")
  
  echo "🔀 Auto-fusion de la Pull Request..."
  gh pr merge "${PR_URL}" --merge --delete-branch
  
  git checkout "${MAIN_BRANCH}"
  git pull "${GTM_REMOTE}" "${MAIN_BRANCH}"
  git push "${GITHUB_REMOTE}" "${MAIN_BRANCH}"
  echo "✅ Synchronisation réussie !"
else
  echo "⚠️ Le dépôt officiel cloud-gtm/app-rag-comparison n'est pas encore accessible ou configuré."
  echo "Poussé uniquement sur votre dépôt personnel GitHub."
  git checkout "${MAIN_BRANCH}"
fi
