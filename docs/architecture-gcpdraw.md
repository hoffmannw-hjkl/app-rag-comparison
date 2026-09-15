# Schéma d'Architecture GCP Draw : RAG Comparison Chatbot

Ce schéma représente l'architecture de l'application **RAG Comparison Chatbot** formalisée en syntaxe **GCP Draw** (`go/gcpdraw`).

## 📋 Instructions d'utilisation
1. Rendez-vous sur l'outil officiel Google Cloud : **[GCP Draw (go/gcpdraw)](https://gcpdraw.corp.google.com)**.
2. Cliquez sur **Import / Code**.
3. Copiez-collez l'intégralité du bloc ci-dessous pour visualiser, éditer et exporter le schéma.

---

```text
meta {
  title "RAG Comparison Chatbot Architecture"
}

elements {
  card users as users {
    display_name "Web Users & Evaluators"
  }

  gcp {
    card armor as waf {
      name "Cloud Armor"
      description "Security Policy & Rate Limiting"
    }

    card load_balancer as lb {
      name "Cloud Load Balancing"
      description "IAP Routing & Frontend Ingress"
    }

    card iap as iap {
      name "Identity-Aware Proxy"
      description "Authenticated User Identity"
    }

    card run as rag_service {
      name "Cloud Run RAG Service"
      description "Go Microservice & In-Memory Index"
    }

    group knowledge_data {
      name "Persistent Corpus Storage"

      card storage as gcs_corpus {
        name "Cloud Storage"
        description "Raw Documents & corpus.json Index Snapshot"
      }
    }

    group vertex_ai_suite {
      name "Vertex AI Services"

      card vertex_ai as embeddings {
        name "Vertex AI Embeddings"
        description "text-multilingual-embedding-002"
      }

      card vertex_ai as gemini_arena {
        name "Gemini Model Arena"
        description "Gemini 2.5 Flash vs Gemini 3.8 Flash"
      }

      card vertex_ai as rapid_eval {
        name "Rapid Evaluation API"
        description "Groundedness & Relevance Assessment"
      }
    }
  }
}

paths {
  users -down-> waf
  waf --> lb
  lb --> iap
  iap -down-> rag_service

  rag_service -right-> gcs_corpus : "Corpus Snapshot & Sync"
  rag_service --> embeddings : "Vectorization"
  rag_service --> gemini_arena : "Dual Streaming SSE"
  rag_service --> rapid_eval : "On-Demand Quality Evaluation"
}
```
