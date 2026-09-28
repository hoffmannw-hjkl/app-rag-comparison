.PHONY: help verify test run demo-agentic demo-evaluate deploy sync-gtm

PORT     ?= 8080
BASE_URL ?= https://rag.hoffmannw.demo.altostrat.com

help: ## Affiche l'aide interactive des commandes de démo RAG & Agentic CRAG
	@echo "================================================================================"
	@echo " 🤖 GCP AI Foundation — RAG & Agentic CRAG Comparator (Commandes de Démo)"
	@echo "================================================================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

verify: ## Exécute le Gatekeeper M1L1 (gofmt + suite complète de tests unitaires Go)
	@./.agents/skills/rag-benchmark-and-ci/scripts/verify.sh

test: ## Lance les tests unitaires Go avec détection des conditions de concurrence (-race)
	@cd src && go test -v -race ./...

run: ## Démarre le serveur Go RAG localement sur http://localhost:8080
	@cd src && PORT=$(PORT) go run .

demo-agentic: ## Teste en direct le flux SSE du Swarm Agentic CRAG (4 sous-agents) en CLI
	@curl -sN "$(BASE_URL)/api/chat/stream?mode=agentic&model=gemini-3.5-flash&q=Compare+Cloud+Armor+WAF+et+Backup+DR+WORM" | head -n 30

demo-evaluate: ## Déclenche l'Autorater LLM-as-a-Judge (Groundedness & Relevance /5) en CLI
	@curl -s -X POST "$(BASE_URL)/api/evaluate" \
		-H "Content-Type: application/json" \
		-d '{"query":"Comment Cloud Armor protège-t-il Vertex AI ?","model":"gemini-3.5-flash","prediction":"Cloud Armor filtre le trafic L7 OWASP Top 10 en amont du Load Balancer.","context":"Cloud Armor WAF protège les endpoints exposés contre les attaques OWASP Top 10."}' | jq .

deploy: ## Déploie l'application sur Google Cloud Run (europe-west1)
	@./deploy/cloudrun/deploy.sh

sync-gtm: ## Synchronise la branche courante vers cloud-gtm/app-rag-comparison via Pull Request
	@./scripts/sync-gtm.sh --force
