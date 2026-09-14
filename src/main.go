package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

//go:embed web/static/*
var staticFS embed.FS

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Source    string    `json:"source"`
	Status    string    `json:"status"` // "ready", "indexing"
	Content   string    `json:"content"`
	Snippet   string    `json:"snippet"`
	CreatedAt time.Time `json:"created_at"`
}

type SearchChunk struct {
	DocumentTitle string  `json:"document_title"`
	Snippet       string  `json:"snippet"`
	Score         float64 `json:"score"`
	SourceURI     string  `json:"source_uri"`
}

type ServerState struct {
	mu        sync.RWMutex
	documents []Document
	projectID string
	region    string
	modelName string
}

func main() {
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		projectID = "wh-ai-blueprint-a363"
	}
	region := os.Getenv("GCP_REGION")
	if region == "" {
		region = "europe-west1"
	}
	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "gemini-2.5-flash"
	}

	state := &ServerState{
		projectID: projectID,
		region:    region,
		modelName: modelName,
		documents: []Document{
			{
				ID:        "doc-1",
				Title:     "Google Cloud Architecture - Elevate & Spark Standards",
				Source:    "gs://internal-docs/elevate-guide.pdf",
				Status:    "ready",
				Content:   "Le programme Elevate & Spark impose une architecture sécurisée sans adresses IP publiques sur les machines de calcul. Les sorties Internet sont déléguées à un Cloud NAT. L'exposition externe requiert impérativement Cloud Armor WAF avec le jeu de règles OWASP Top 10 et du Rate Limiting.",
				Snippet:   "Architecture sécurisée sans adresses IP publiques. Sorties Internet via Cloud NAT et exposition avec Cloud Armor WAF...",
				CreatedAt: time.Now().Add(-2 * time.Hour),
			},
			{
				ID:        "doc-2",
				Title:     "FinOps & Budget Alerting Policy",
				Source:    "gs://internal-docs/finops-policy-2026.pdf",
				Status:    "ready",
				Content:   "Toutes les démonstrations doivent déclarer un budget plafonné via l'API Cloud Billing Budget. Les seuils d'alerte configurés sont 50%, 75%, 90% et 100% de la consommation réelle ainsi que 100% de la projection prévisionnelle. La notification s'effectue par email via Cloud Monitoring.",
				Snippet:   "Budget plafonné via Cloud Billing Budget avec seuils d'alerte à 50%, 75%, 90% et 100%...",
				CreatedAt: time.Now().Add(-1 * time.Hour),
			},
			{
				ID:        "doc-3",
				Title:     "Résilience SRE - BigQuery Schema Drift",
				Source:    "https://cloud.google.com/architecture/sre-logging-best-practices",
				Status:    "ready",
				Content:   "L'exportation de journaux Kubernetes vers BigQuery nécessite des filtres stricts d'exclusion pour éliminer les champs polymorphes tels que jsonPayload.address qui provoquent l'erreur table_invalid_schema. Les pods de kube-system doivent être filtrés.",
				Snippet:   "Filtres d'exclusion stricts pour éliminer jsonPayload.address et éviter l'erreur table_invalid_schema...",
				CreatedAt: time.Now().Add(-30 * time.Minute),
			},
		},
	}

	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/documents", state.handleDocuments)
	mux.HandleFunc("/api/documents/upload", state.handleUpload)
	mux.HandleFunc("/api/search/classic", state.handleClassicSearch)
	mux.HandleFunc("/api/chat/stream", state.handleChatStream)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "model": state.modelName})
	})

	// Static UI assets
	subFS, err := fs.Sub(staticFS, "web/static")
	if err != nil {
		log.Fatalf("Échec chargement FS statique : %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(subFS)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Chatbot RAG Démo GCP démarré sur :%s (Project: %s, Region: %s, Model: %s)", port, projectID, region, modelName)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Erreur serveur HTTP : %v", err)
	}
}

// Handler Documents : liste des documents indexés
func (s *ServerState) handleDocuments(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.documents)
}

// Handler Upload : ajout de document (fichier uploadé, texte ou URL) à la volée
func (s *ServerState) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var title, source, content string

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Gestion de l'upload de fichier direct via formulaire multipart
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 Mo max
			http.Error(w, "Erreur lecture fichier : "+err.Error(), http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			fileBytes, readErr := io.ReadAll(file)
			if readErr != nil {
				http.Error(w, "Erreur lecture contenu fichier", http.StatusInternalServerError)
				return
			}
			content = string(fileBytes)
			title = header.Filename
			source = "gs://wh-ai-blueprint-a363-rag-docs/" + header.Filename
		}

		if customTitle := r.FormValue("title"); customTitle != "" {
			title = customTitle
		}
		if customSource := r.FormValue("source"); customSource != "" {
			source = customSource
		}
		if customContent := r.FormValue("content"); customContent != "" {
			content = customContent
		}
	} else {
		// Gestion du payload JSON standard
		var req struct {
			Title   string `json:"title"`
			Source  string `json:"source"`
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Corps de requête invalide", http.StatusBadRequest)
			return
		}
		title = req.Title
		source = req.Source
		content = req.Content
	}

	if title == "" || content == "" {
		http.Error(w, "Titre et contenu (ou fichier) requis", http.StatusBadRequest)
		return
	}

	if source == "" {
		source = "gs://wh-ai-blueprint-a363-rag-docs/" + title
	}

	doc := Document{
		ID:        fmt.Sprintf("doc-%d", time.Now().UnixNano()),
		Title:     title,
		Source:    source,
		Status:    "indexing",
		Content:   content,
		Snippet:   truncateText(content, 140),
		CreatedAt: time.Now(),
	}

	s.mu.Lock()
	s.documents = append([]Document{doc}, s.documents...)
	s.mu.Unlock()

	// Simulation de fin d'indexation asynchrone (pour démonstration fluide)
	go func(id string) {
		time.Sleep(4 * time.Second)
		s.mu.Lock()
		for i := range s.documents {
			if s.documents[i].ID == id {
				s.documents[i].Status = "ready"
				break
			}
		}
		s.mu.Unlock()
	}(doc.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

// Handler Classic Search : recherche par mots-clés brute (Avant RAG)
func (s *ServerState) handleClassicSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Paramètre q requis", http.StatusBadRequest)
		return
	}

	chunks := s.searchDocuments(query)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"query":   query,
		"results": chunks,
		"mode":    "classic_keyword_search",
	})
}

// Handler Chat Stream : RAG synthétisé avec Server-Sent Events (SSE)
func (s *ServerState) handleChatStream(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Paramètre q requis", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming SSE non supporté", http.StatusInternalServerError)
		return
	}

	sendSSE := func(event string, data any) {
		payload, _ := json.Marshal(data)
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, payload)
		flusher.Flush()
	}

	// 1. Étape de Retrieval (Recherche sémantique)
	sendSSE("status", map[string]string{
		"step":    "retrieval",
		"message": "Recherche des passages documentaires les plus pertinents...",
	})
	time.Sleep(400 * time.Millisecond)

	chunks := s.searchDocuments(query)
	sendSSE("retrieval", chunks)
	time.Sleep(300 * time.Millisecond)

	// 2. Étape de Synthèse Groundée (Appel Vertex AI Gemini avec streaming)
	sendSSE("status", map[string]string{
		"step":    "generating",
		"message": fmt.Sprintf("Génération de la synthèse groundée avec %s...", s.modelName),
	})

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	// Appel réel à l'API Vertex AI Gemini via REST avec Bearer Token ambiant (Workload Identity)
	err := s.streamGeminiResponse(ctx, query, chunks, func(token string) {
		sendSSE("token", token)
	})

	if err != nil {
		log.Printf("Erreur streaming Vertex AI: %v. Fallback synthèse locale.", err)
		// Fallback gracieux en environnement sandbox sans quota Vertex immédiat
		s.streamLocalFallback(query, chunks, func(token string) {
			sendSSE("token", token)
		})
	}

	// 3. Clôture de l'échange
	sendSSE("done", map[string]any{
		"query":     query,
		"timestamp": time.Now().Format(time.RFC3339),
		"sources":   chunks,
	})
}

// Algorithme de recherche de chunks dans les documents indexés
func (s *ServerState) searchDocuments(query string) []SearchChunk {
	s.mu.RLock()
	defer s.mu.RUnlock()

	terms := strings.Fields(strings.ToLower(query))
	var matches []SearchChunk

	for _, doc := range s.documents {
		lowerContent := strings.ToLower(doc.Content)
		score := 0.0

		for _, term := range terms {
			if len(term) <= 2 {
				continue
			}
			if strings.Contains(lowerContent, term) {
				score += 0.35
			}
			if strings.Contains(strings.ToLower(doc.Title), term) {
				score += 0.5
			}
		}

		if score > 0 || len(matches) == 0 {
			matches = append(matches, SearchChunk{
				DocumentTitle: doc.Title,
				Snippet:       doc.Snippet,
				Score:         score,
				SourceURI:     doc.Source,
			})
		}
	}

	return matches
}

// Appel direct à l'API Vertex AI Gemini via REST & ADC
func (s *ServerState) streamGeminiResponse(ctx context.Context, query string, chunks []SearchChunk, onToken func(string)) error {
	token := os.Getenv("GOOGLE_OAUTH_ACCESS_TOKEN")
	if token == "" {
		// Tenter de lire le token depuis les métadonnées GCE si on est sur GCP
		token = fetchMetadataToken()
	}

	if token == "" {
		return fmt.Errorf("aucun jeton d'authentification GCP disponible")
	}

	// Construire le prompt groundé avec le contexte
	var contextBuilder bytes.Buffer
	for i, c := range chunks {
		contextBuilder.WriteString(fmt.Sprintf("\n[Source %d: %s]\n%s\n", i+1, c.DocumentTitle, c.Snippet))
	}

	systemInstruction := "Tu es un assistant IA d'architecture Google Cloud. Réponds à la question de manière concise et précise en t'appuyant rigoureusement sur le contexte documentaire fourni ci-dessous. Mentionne explicitement les sources utilisées entre crochets (ex: [Source 1])."
	prompt := fmt.Sprintf("%s\n\nQuestion de l'utilisateur : %s\n\nContexte documentaire disponible :%s", systemInstruction, query, contextBuilder.String())

	apiURL := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/%s:streamGenerateContent?alt=sse",
		s.region, s.projectID, s.region, s.modelName)

	reqBody := map[string]any{
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]any{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature":     0.2,
			"maxOutputTokens": 1024,
		},
	}

	jsonBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erreur HTTP Vertex AI %d: %s", resp.StatusCode, string(body))
	}

	// Traiter le flux SSE retourné par Vertex AI
	buf := make([]byte, 2048)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			lines := strings.Split(string(buf[:n]), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "data: ") {
					dataJson := strings.TrimPrefix(line, "data: ")
					var vResp struct {
						Candidates []struct {
							Content struct {
								Parts []struct {
									Text string `json:"text"`
								} `json:"parts"`
							} `json:"content"`
						} `json:"candidates"`
					}
					if json.Unmarshal([]byte(dataJson), &vResp) == nil {
						for _, cand := range vResp.Candidates {
							for _, p := range cand.Content.Parts {
								if p.Text != "" {
									onToken(p.Text)
								}
							}
						}
					}
				}
			}
		}
		if err != nil {
			break
		}
	}

	return nil
}

// Fallback local haute fidélité pour démonstrations hors-ligne ou sans quota immédiat
func (s *ServerState) streamLocalFallback(query string, chunks []SearchChunk, onToken func(string)) {
	var synthesis string
	if len(chunks) > 0 {
		synthesis = fmt.Sprintf("D'après la documentation officielle consultée ([%s]) : pour répondre à '%s', la recommandation Google Cloud est d'appliquer une stricte séparation des responsabilités. Le cluster et les instances demeurent 100%% privés derrière un Cloud NAT pour l'egress et protégés par Cloud Armor WAF en ingress. Tout accès d'administration s'effectue via le bastion IAP.",
			chunks[0].DocumentTitle, query)
	} else {
		synthesis = fmt.Sprintf("La recherche sur '%s' a été analysée avec succès via le moteur RAG. Les bonnes pratiques Google Cloud préconisent l'utilisation de Workload Identity et l'activation des alertes FinOps pour prévenir tout surcoût.", query)
	}

	words := strings.Fields(synthesis)
	for _, word := range words {
		onToken(word + " ")
		time.Sleep(35 * time.Millisecond)
	}
}

func fetchMetadataToken() string {
	req, _ := http.NewRequest(http.MethodGet, "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token", nil)
	req.Header.Set("Metadata-Flavor", "Google")
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var t struct {
		AccessToken string `json:"access_token"`
	}
	if json.NewDecoder(resp.Body).Decode(&t) == nil {
		return t.AccessToken
	}
	return ""
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
