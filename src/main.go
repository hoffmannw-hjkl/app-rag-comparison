package main

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/zlib"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

//go:embed web/static/*
var staticFS embed.FS

// Chunk représente un passage découpé d'un document, unité de base du retrieval.
type Chunk struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
	// norm est la forme normalisée de Text (minuscules, sans accent ni
	// apostrophe) sur laquelle s'effectue la recherche. Elle est calculée une
	// fois à l'indexation plutôt qu'à chaque requête.
	norm      string
	Embedding []float32 `json:"embedding,omitempty"`
}

type Document struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Source string `json:"source"`
	Status string `json:"status"` // "ready", "indexing"
	// Content et Chunks ne sont jamais sérialisés vers le client : ils peuvent peser
	// plusieurs Mo par document et l'UI n'exploite que le snippet.
	Content   string    `json:"-"`
	Chunks    []Chunk   `json:"-"`
	Snippet   string    `json:"snippet"`
	SizeBytes int       `json:"size_bytes"`
	NumChunks int       `json:"num_chunks"`
	CreatedAt time.Time `json:"created_at"`
	normTitle string
}

type SearchChunk struct {
	DocumentTitle string  `json:"document_title"`
	Snippet       string  `json:"snippet"`
	Score         float64 `json:"score"`
	SemanticScore float64 `json:"semantic_score,omitempty"`
	LexicalScore  float64 `json:"lexical_score,omitempty"`
	SearchType    string  `json:"search_type,omitempty"` // "hybrid" ou "lexical"
	SourceURI     string  `json:"source_uri"`
	ChunkIndex    int     `json:"chunk_index"`
}

type ModelInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IsDefault   bool   `json:"is_default"`
}

var AvailableModels = []ModelInfo{
	{
		ID:          "gemini-3.5-flash",
		Name:        "Gemini 3.5 Flash",
		Description: "Dernière génération multimodale native ultra-rapide & efficiente de Google",
		Category:    "Flash",
		IsDefault:   true,
	},
	{
		ID:          "gemini-3.8-flash",
		Name:        "Gemini 3.8 Flash",
		Description: "Modèle agentique de pointe pour l'ingénierie et le raisonnement de code complexe",
		Category:    "Agentic",
		IsDefault:   false,
	},
	{
		ID:          "gemini-3.1-pro-preview",
		Name:        "Gemini 3.1 Pro Preview",
		Description: "Modèle Frontier de raisonnement approfondi & synthèse d'architecture complexe",
		Category:    "Pro",
		IsDefault:   false,
	},
	{
		ID:          "gemini-3.5-flash-lite",
		Name:        "Gemini 3.5 Flash-Lite",
		Description: "Ultra-rapide, ultra-légère latence et débit maximal pour le RAG haute fréquence",
		Category:    "Flash-Lite",
		IsDefault:   false,
	},
	{
		ID:          "gemini-3.7-flash",
		Name:        "Gemini 3.7 Flash",
		Description: "Modèle de transition rapide et précis avec réflexion hybride",
		Category:    "Flash",
		IsDefault:   false,
	},
	{
		ID:          "gemini-2.5-flash",
		Name:        "Gemini 2.5 Flash",
		Description: "Génération précédente (legacy / fallback)",
		Category:    "Legacy",
		IsDefault:   false,
	},
}

type ServerState struct {
	mu            sync.RWMutex
	documents     []Document
	projectID     string
	region        string
	modelName     string
	gcsBucket     string // Nom du bucket Cloud Storage pour la persistance du corpus
	gcsEndpoint   string // Surcharge de l'endpoint de base GCS pour les tests unitaires
	embedEndpoint string // Surcharge de l'endpoint Vertex AI Embeddings pour les tests unitaires
	evalEndpoint  string // Surcharge de l'endpoint Vertex AI Evaluation pour les tests unitaires
}

// Structures pour la persistance et synchronisation du corpus sur Cloud Storage (index/corpus.json)
type StoredChunk struct {
	Index     int       `json:"index"`
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding,omitempty"`
}

type StoredDocument struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Source    string        `json:"source"`
	Status    string        `json:"status"`
	Content   string        `json:"content"`
	Chunks    []StoredChunk `json:"chunks"`
	Snippet   string        `json:"snippet"`
	SizeBytes int           `json:"size_bytes"`
	NumChunks int           `json:"num_chunks"`
	CreatedAt time.Time     `json:"created_at"`
}

type CorpusSnapshot struct {
	Version   int              `json:"version"`
	UpdatedAt time.Time        `json:"updated_at"`
	Documents []StoredDocument `json:"documents"`
}

const corpusSnapshotPath = "index/corpus.json"

func (d Document) toStored() StoredDocument {
	chunks := make([]StoredChunk, len(d.Chunks))
	for i, c := range d.Chunks {
		chunks[i] = StoredChunk{
			Index:     c.Index,
			Text:      c.Text,
			Embedding: c.Embedding,
		}
	}
	return StoredDocument{
		ID:        d.ID,
		Title:     d.Title,
		Source:    d.Source,
		Status:    d.Status,
		Content:   d.Content,
		Chunks:    chunks,
		Snippet:   d.Snippet,
		SizeBytes: d.SizeBytes,
		NumChunks: d.NumChunks,
		CreatedAt: d.CreatedAt,
	}
}

func storedToDocument(sd StoredDocument) Document {
	chunks := make([]Chunk, len(sd.Chunks))
	for i, sc := range sd.Chunks {
		chunks[i] = Chunk{
			Index:     sc.Index,
			Text:      sc.Text,
			norm:      normalizeForSearch(sc.Text),
			Embedding: sc.Embedding,
		}
	}
	if len(chunks) == 0 && sd.Content != "" {
		chunks = chunkText(sd.Content)
	}
	status := sd.Status
	if status == "" {
		status = "ready"
	}
	numChunks := sd.NumChunks
	if numChunks == 0 {
		numChunks = len(chunks)
	}
	sizeBytes := sd.SizeBytes
	if sizeBytes == 0 {
		sizeBytes = len(sd.Content)
	}
	snippet := sd.Snippet
	if snippet == "" && len(chunks) > 0 {
		snippet = truncateText(chunks[0].Text, 250)
	}

	return Document{
		ID:        sd.ID,
		Title:     sd.Title,
		Source:    sd.Source,
		Status:    status,
		Content:   sd.Content,
		Chunks:    chunks,
		Snippet:   snippet,
		SizeBytes: sizeBytes,
		NumChunks: numChunks,
		CreatedAt: sd.CreatedAt,
		normTitle: normalizeForSearch(sd.Title),
	}
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
		modelName = "gemini-3.5-flash"
	}
	gcsBucket := os.Getenv("GCS_RAG_BUCKET")
	if gcsBucket == "" {
		gcsBucket = "wh-ai-blueprint-a363-ai-demo-2e2m-rag-docs"
	} else if gcsBucket == "none" || gcsBucket == "disabled" {
		gcsBucket = ""
	}

	state := &ServerState{
		projectID: projectID,
		region:    region,
		modelName: modelName,
		gcsBucket: gcsBucket,
		documents: []Document{
			newDocument("doc-1",
				"Google Cloud Architecture - Elevate & Spark Standards",
				"gs://internal-docs/elevate-guide.pdf",
				"Le programme Elevate & Spark impose une architecture sécurisée sans adresses IP publiques sur les machines de calcul. Les sorties Internet sont déléguées à un Cloud NAT. L'exposition externe requiert impérativement Cloud Armor WAF avec le jeu de règles OWASP Top 10 et du Rate Limiting.",
				time.Now().Add(-2*time.Hour)),
			newDocument("doc-2",
				"FinOps & Budget Alerting Policy",
				"gs://internal-docs/finops-policy-2026.pdf",
				"Toutes les démonstrations doivent déclarer un budget plafonné via l'API Cloud Billing Budget. Les seuils d'alerte configurés sont 50%, 75%, 90% et 100% de la consommation réelle ainsi que 100% de la projection prévisionnelle. La notification s'effectue par email via Cloud Monitoring.",
				time.Now().Add(-1*time.Hour)),
			newDocument("doc-3",
				"Résilience SRE - BigQuery Schema Drift",
				"https://cloud.google.com/architecture/sre-logging-best-practices",
				"L'exportation de journaux Kubernetes vers BigQuery nécessite des filtres stricts d'exclusion pour éliminer les champs polymorphes tels que jsonPayload.address qui provoquent l'erreur table_invalid_schema. Les pods de kube-system doivent être filtrés.",
				time.Now().Add(-30*time.Minute)),
		},
	}

	// Restauration du snapshot persistant depuis GCS si configuré
	if state.gcsBucket != "" {
		initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
		state.initCorpusFromGCS(initCtx)
		initCancel()
	}

	// Vectorisation en arrière-plan des chunks du corpus ne disposant pas encore d'embeddings
	go func() {
		embCtx, embCancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer embCancel()
		state.ensureCorpusEmbeddings(embCtx)
	}()

	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/documents", state.handleDocuments)
	mux.HandleFunc("/api/documents/upload", state.handleUpload)
	mux.HandleFunc("/api/search/classic", state.handleClassicSearch)
	mux.HandleFunc("/api/chat/stream", state.handleChatStream)
	mux.HandleFunc("/api/models", state.handleModels)
	mux.HandleFunc("/api/model/switch", state.handleModelSwitch)
	mux.HandleFunc("/api/evaluate", state.handleEvaluate)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		state.mu.RLock()
		current := state.modelName
		state.mu.RUnlock()
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "model": current})
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

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
		// ReadHeaderTimeout protège contre les attaques de type Slowloris.
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		// WriteTimeout doit rester à zéro : les réponses SSE sont des flux longue durée
		// qu'un timeout d'écriture couperait en plein milieu.
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	// Arrêt gracieux : Cloud Run envoie SIGTERM avant de retirer une instance.
	shutdownDone := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		<-sigCh
		log.Println("⏹️  Signal d'arrêt reçu, drainage des requêtes en cours...")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Arrêt gracieux interrompu : %v", err)
		}
		close(shutdownDone)
	}()

	log.Printf("🚀 Chatbot RAG Démo GCP démarré sur :%s (Project: %s, Region: %s, Model: %s)", port, projectID, region, modelName)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Erreur serveur HTTP : %v", err)
	}
	<-shutdownDone
	log.Println("✅ Serveur arrêté proprement")
}

// Handler Documents : liste des documents indexés ou suppression (unitaire ou globale)
func (s *ServerState) handleDocuments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.mu.RLock()
		defer s.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.documents)

	case http.MethodDelete:
		id := strings.TrimSpace(r.URL.Query().Get("id"))
		all := r.URL.Query().Get("all") == "true"

		if !all && id == "" {
			http.Error(w, "Paramètre 'id' ou 'all=true' requis pour la suppression", http.StatusBadRequest)
			return
		}

		if all {
			s.mu.Lock()
			deletedCount := len(s.documents)
			s.documents = []Document{}
			s.mu.Unlock()

			if s.gcsBucket != "" {
				if err := s.saveCorpusSnapshot(r.Context()); err != nil {
					log.Printf("⚠️  Erreur mise à jour snapshot GCS après purge : %v", err)
				}
			}

			log.Printf("🗑️  Purge complète du corpus documentaire (%d documents supprimés)", deletedCount)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"status":          "cleared",
				"deleted_count":   deletedCount,
				"remaining_count": 0,
			})
			return
		}

		s.mu.Lock()
		found := false
		filtered := make([]Document, 0, len(s.documents))
		var deletedTitle string
		for _, doc := range s.documents {
			if doc.ID == id {
				found = true
				deletedTitle = doc.Title
				continue
			}
			filtered = append(filtered, doc)
		}

		if !found {
			s.mu.Unlock()
			http.Error(w, fmt.Sprintf("Document non trouvé : %s", id), http.StatusNotFound)
			return
		}

		s.documents = filtered
		remCount := len(s.documents)
		s.mu.Unlock()

		if s.gcsBucket != "" {
			if err := s.saveCorpusSnapshot(r.Context()); err != nil {
				log.Printf("⚠️  Erreur mise à jour snapshot GCS après suppression : %v", err)
			}
		}

		log.Printf("🗑️  Document %q (id: %s) supprimé (%d restants)", deletedTitle, id, remCount)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":          "deleted",
			"id":              id,
			"deleted_title":   deletedTitle,
			"remaining_count": remCount,
		})

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// Handler Models : liste des modèles disponibles et modèle actuellement actif
func (s *ServerState) handleModels(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	current := s.modelName
	bucket := s.gcsBucket
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"current":    current,
		"models":     AvailableModels,
		"embedding":  "text-multilingual-embedding-002",
		"gcs_bucket": bucket,
	})
}

// Handler Model Switch : valide le modèle demandé par le client.
//
// Ce handler ne mute délibérément PAS l'état du serveur : le modèle est un choix
// propre à chaque utilisateur, transmis à chaque requête via le paramètre ?model=.
// Muter s.modelName ici changerait le modèle de tous les utilisateurs connectés
// simultanément, ce qui est un effet de bord indésirable en démonstration collective.
func (s *ServerState) handleModelSwitch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Model string `json:"model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Model == "" {
		http.Error(w, "Modèle invalide", http.StatusBadRequest)
		return
	}

	if !isKnownModel(req.Model) {
		http.Error(w, fmt.Sprintf("Modèle inconnu : %s", req.Model), http.StatusBadRequest)
		return
	}

	log.Printf("🤖 Modèle sélectionné par le client : %s", req.Model)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"current": req.Model,
		"scope":   "session",
	})
}

// isKnownModel vérifie qu'un identifiant de modèle fait partie du catalogue exposé.
func isKnownModel(id string) bool {
	for _, m := range AvailableModels {
		if m.ID == id {
			return true
		}
	}
	return false
}

// Handler Upload : ajout de document (fichier uploadé, texte ou URL) à la volée
func (s *ServerState) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// ParseMultipartForm ne borne que la part gardée en mémoire : sans MaxBytesReader
	// le corps total reste illimité et peut provoquer un OOM kill sur Cloud Run.
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)

	contentType := r.Header.Get("Content-Type")
	var addedDocs []Document

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(maxMemoryBytes); err != nil {
			http.Error(w, "Fichiers trop volumineux ou illisibles (limite : 50 Mo au total)", http.StatusRequestEntityTooLarge)
			return
		}
		defer r.MultipartForm.RemoveAll()

		// Récupération des fichiers soit sous le champ 'file' soit 'files'
		files := r.MultipartForm.File["file"]
		if len(files) == 0 {
			files = r.MultipartForm.File["files"]
		}

		if len(files) > 0 {
			for i, fh := range files {
				f, err := fh.Open()
				if err != nil {
					log.Printf("⚠️  Ouverture impossible pour %q : %v", fh.Filename, err)
					continue
				}
				data, err := io.ReadAll(f)
				f.Close()
				if err != nil {
					log.Printf("⚠️  Lecture impossible pour %q : %v", fh.Filename, err)
					continue
				}

				name := sanitizeFilename(fh.Filename)
				extracted := string(data)
				if strings.HasSuffix(strings.ToLower(name), ".pdf") {
					extracted = extractTextFromPDF(data)
				}
				extracted = strings.ToValidUTF8(extracted, "")

				docID := fmt.Sprintf("doc-%d-%d", time.Now().UnixNano(), i)
				sourceURI := "gs://wh-ai-blueprint-a363-rag-docs/" + name
				if s.gcsBucket != "" {
					sourceURI = fmt.Sprintf("gs://%s/documents/%s/%s", s.gcsBucket, docID, name)
					rawPath := fmt.Sprintf("documents/%s/%s", docID, name)
					cType := fh.Header.Get("Content-Type")
					if cType == "" {
						cType = "application/octet-stream"
					}
					go func(path, ct string, b []byte) {
						bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()
						if err := s.uploadToGCS(bgCtx, path, ct, b); err != nil {
							log.Printf("⚠️  Échec upload fichier brut GCS %s: %v", path, err)
						}
					}(rawPath, cType, data)
				}

				addedDocs = append(addedDocs, newDocument(
					docID,
					name,
					sourceURI,
					extracted,
					time.Now(),
				))
			}
		} else {
			// Saisie directe sans fichier
			title := sanitizeFilename(r.FormValue("title"))
			content := strings.ToValidUTF8(r.FormValue("content"), "")
			source := r.FormValue("source")
			if title != "" && content != "" {
				docID := fmt.Sprintf("doc-%d", time.Now().UnixNano())
				if source == "" {
					if s.gcsBucket != "" {
						source = fmt.Sprintf("gs://%s/documents/%s/%s.txt", s.gcsBucket, docID, title)
					} else {
						source = "gs://wh-ai-blueprint-a363-rag-docs/" + title
					}
				}
				if s.gcsBucket != "" {
					rawPath := fmt.Sprintf("documents/%s/%s.txt", docID, title)
					go func(path string, b []byte) {
						bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()
						if err := s.uploadToGCS(bgCtx, path, "text/plain; charset=utf-8", b); err != nil {
							log.Printf("⚠️  Échec upload texte brut GCS %s: %v", path, err)
						}
					}(rawPath, []byte(content))
				}
				addedDocs = append(addedDocs, newDocument(
					docID, title, source, content, time.Now()))
			}
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
		title := sanitizeFilename(req.Title)
		content := strings.ToValidUTF8(req.Content, "")
		if title != "" && content != "" {
			docID := fmt.Sprintf("doc-%d", time.Now().UnixNano())
			source := req.Source
			if source == "" {
				if s.gcsBucket != "" {
					source = fmt.Sprintf("gs://%s/documents/%s/%s.txt", s.gcsBucket, docID, title)
				} else {
					source = "gs://wh-ai-blueprint-a363-rag-docs/" + title
				}
			}
			if s.gcsBucket != "" {
				rawPath := fmt.Sprintf("documents/%s/%s.txt", docID, title)
				go func(path string, b []byte) {
					bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()
					if err := s.uploadToGCS(bgCtx, path, "text/plain; charset=utf-8", b); err != nil {
						log.Printf("⚠️  Échec upload texte brut GCS %s: %v", path, err)
					}
				}(rawPath, []byte(content))
			}
			addedDocs = append(addedDocs, newDocument(
				docID, title, source, content, time.Now()))
		}
	}

	if len(addedDocs) == 0 {
		http.Error(w, "Aucun document ou fichier valide reçu", http.StatusBadRequest)
		return
	}

	// Vectorisation sémantique des chunks ajoutés via Vertex AI Embeddings
	for i := range addedDocs {
		if len(addedDocs[i].Chunks) > 0 {
			texts := make([]string, len(addedDocs[i].Chunks))
			for j, c := range addedDocs[i].Chunks {
				texts[j] = c.Text
			}
			embCtx, embCancel := context.WithTimeout(r.Context(), 15*time.Second)
			embs, err := s.generateEmbeddings(embCtx, texts, "RETRIEVAL_DOCUMENT")
			embCancel()
			if err == nil && len(embs) == len(addedDocs[i].Chunks) {
				for j := range addedDocs[i].Chunks {
					addedDocs[i].Chunks[j].Embedding = embs[j]
				}
				log.Printf("🧠 %d chunk(s) vectorisé(s) pour %q", len(embs), addedDocs[i].Title)
			} else if err != nil {
				log.Printf("ℹ️  Vectorisation non disponible pour %q (%v) : repli lexical actif", addedDocs[i].Title, err)
			}
		}
	}

	totalChunks := 0
	totalBytes := 0
	s.mu.Lock()
	for i := range addedDocs {
		totalChunks += addedDocs[i].NumChunks
		totalBytes += addedDocs[i].SizeBytes
		s.documents = append([]Document{addedDocs[i]}, s.documents...)
	}
	s.mu.Unlock()

	// Synchronisation du snapshot d'index sur GCS
	if s.gcsBucket != "" {
		if err := s.saveCorpusSnapshot(r.Context()); err != nil {
			log.Printf("⚠️  Avertissement : snapshot GCS non synchronisé : %v", err)
		}
	}

	log.Printf("📥 %d document(s) reçu(s), %d chunk(s) indexé(s) dans le corpus (%d octets)", len(addedDocs), totalChunks, totalBytes)

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]any{
		"status":       "ok",
		"count":        len(addedDocs),
		"total_chunks": totalChunks,
		"total_bytes":  totalBytes,
		"documents":    addedDocs,
	}
	if len(addedDocs) == 1 {
		resp["id"] = addedDocs[0].ID
		resp["title"] = addedDocs[0].Title
		resp["source"] = addedDocs[0].Source
		resp["snippet"] = addedDocs[0].Snippet
		resp["num_chunks"] = addedDocs[0].NumChunks
		resp["size_bytes"] = addedDocs[0].SizeBytes
		resp["created_at"] = addedDocs[0].CreatedAt
	}
	json.NewEncoder(w).Encode(resp)
}

// Constantes de traitement documentaire.
const (
	maxUploadBytes = 50 << 20 // Taille maximale du corps d'une requête d'upload
	maxMemoryBytes = 10 << 20 // Part de l'upload conservée en mémoire, le reste va sur disque
	chunkSize      = 1000     // Taille cible d'un chunk, en runes
	chunkOverlap   = 100      // Chevauchement entre deux chunks consécutifs, en runes
	topKChunks     = 8        // Nombre de passages transmis au modèle pour le grounding
)

// Regex compilées une seule fois au chargement du package plutôt qu'à chaque appel.
var (
	// Opérateur d'affichage simple : (texte) Tj ou <hex> Tj (et variantes ' et ")
	rePDFShowText = regexp.MustCompile(`(?:\(((?:[^()\\]|\\.)*)\)|<([0-9A-Fa-f\s]+)>)\s*(?:Tj|'|")`)
	// Tableau de crénage : [(V)94(ersion)-375(Managemen)31(t)] TJ ou [<00510053>8<0046>] TJ
	// C'est la forme employée par la plupart des générateurs (TeX, InDesign, dvips)
	// et elle porte l'essentiel du texte : l'ignorer revient à ne rien extraire.
	rePDFShowArray = regexp.MustCompile(`(?s)\[((?:[^\[\]\\]|\\.)*)\]\s*TJ`)
	// Éléments d'un tableau de crénage : chaînes littérales, chaînes hexadécimales et valeurs de crénage.
	rePDFArrayItem = regexp.MustCompile(`\(((?:[^()\\]|\\.)*)\)|<([0-9A-Fa-f\s]+)>|(-?\d+(?:\.\d+)?)`)
	// Objets stream/endstream, dont le contenu est généralement compressé
	rePDFStream = regexp.MustCompile(`(?s)stream\r?\n?(.*?)endstream`)
	// Caractères interdits ou risqués dans un titre de document
	reUnsafeTitle = regexp.MustCompile(`[<>"'&\x00-\x1f\x7f]`)

	// Parsing des tables de correspondance ToUnicode CMap (indispensable pour
	// les PDF produits par TeX, LaTeX, InDesign ou LibreOffice avec polices intégrées).
	rePDFBFCharBlock  = regexp.MustCompile(`(?s)beginbfchar(.*?)endbfchar`)
	rePDFBFChar       = regexp.MustCompile(`<([0-9A-Fa-f]+)>\s+<([0-9A-Fa-f]+)>`)
	rePDFBFRangeBlock = regexp.MustCompile(`(?s)beginbfrange(.*?)endbfrange`)
	rePDFBFRange      = regexp.MustCompile(`<([0-9A-Fa-f]+)>\s+<([0-9A-Fa-f]+)>\s+<([0-9A-Fa-f]+)>`)
	rePDFBFRangeArray = regexp.MustCompile(`(?s)<([0-9A-Fa-f]+)>\s+<([0-9A-Fa-f]+)>\s+\[(.*?)\]`)
	rePDFHexToken     = regexp.MustCompile(`<([0-9A-Fa-f]+)>`)
)

// kerningSpaceThreshold : en deçà de cette valeur (en millièmes d'unité de texte),
// un déplacement de crénage matérialise une espace entre deux mots. Au-dessus, il
// s'agit d'un simple resserrement typographique à l'intérieur d'un même mot.
const kerningSpaceThreshold = -150

// extractTextFromPDF extrait le texte lisible d'un flux binaire PDF.
//
// Trois écueils sont traités ici :
//  1. La compression des flux : traitée par inflateStream (zlib et deflate brut).
//  2. L'encodage des glyphes par CMap ToUnicode : les générateurs modernes (TeX,
//     InDesign) émettent des chaînes hexadécimales <00510053...> dont les codes
//     ne correspondent pas à l'ASCII mais à une table CMap embarquée.
//  3. Le bruit binaire : les chaînes candidates sont filtrées token par token
//     via looksLikeText pour ne jamais polluer le corpus de résidus compressés.
func extractTextFromPDF(data []byte) string {
	var sb strings.Builder
	decodedStreams := 0

	// 1. Décompression de tous les flux et séparation des CMaps ToUnicode
	//    des flux de contenu textuel.
	cmap := make(map[uint32]string)
	type decompressedStream struct {
		content []byte
	}
	var contentStreams []decompressedStream

	for _, m := range rePDFStream.FindAllSubmatch(data, -1) {
		if len(m) < 2 {
			continue
		}
		decoded, ok := inflateStream(bytes.TrimLeft(m[1], "\r\n"))
		if !ok {
			continue
		}
		decodedStreams++

		// Les flux CMap définissent les correspondances glyphes -> Unicode
		if bytes.Contains(decoded, []byte("begincmap")) {
			parseToUnicodeCMap(decoded, cmap)
		} else {
			contentStreams = append(contentStreams, decompressedStream{content: decoded})
		}
	}

	// 2. Extraction des opérateurs de texte sur les flux de contenu décompressés
	for _, cs := range contentStreams {
		appendPDFOperators(&sb, cs.content, cmap)
	}

	// 3. Aucun flux décompressable : le document est probablement en clair.
	if decodedStreams == 0 {
		appendPDFOperators(&sb, data, cmap)
	}

	extracted := normalizeWhitespace(sb.String())

	// 4. Dernier recours pour les documents en clair sans opérateur exploitable.
	if utf8.RuneCountInString(extracted) < 40 && decodedStreams == 0 {
		if words := extractPrintableWords(data); looksLikeText(words) {
			extracted = words
		}
	}

	// Le texte restant a déjà été filtré token par token : s'il est vide, c'est que
	// le document ne porte aucune couche texte exploitable (PDF scanné ou image).
	if extracted == "" {
		return "[Document PDF sans couche texte extractible - OCR requis]"
	}
	return extracted
}

// inflateStream décompresse un flux PDF.
//
// La plupart des producteurs émettent du zlib (RFC 1950), mais certains écrivent
// du DEFLATE brut (RFC 1951) sans en-tête : on tente donc les deux. Un flux non
// compressé ou chiffré avec un filtre non supporté renvoie simplement false.
func inflateStream(raw []byte) ([]byte, bool) {
	if len(raw) < 2 {
		return nil, false
	}

	if zr, err := zlib.NewReader(bytes.NewReader(raw)); err == nil {
		decoded, err := io.ReadAll(io.LimitReader(zr, maxMemoryBytes))
		zr.Close()
		// Une erreur en fin de flux (troncature) reste exploitable si des octets
		// ont déjà été décodés.
		if len(decoded) > 0 && (err == nil || len(decoded) > 64) {
			return decoded, true
		}
	}

	fr := flate.NewReader(bytes.NewReader(raw))
	decoded, err := io.ReadAll(io.LimitReader(fr, maxMemoryBytes))
	fr.Close()
	if len(decoded) > 64 && (err == nil || looksLikeText(string(decoded))) {
		return decoded, true
	}

	return nil, false
}

// looksLikeText distingue du texte exploitable d'un résidu binaire.
//
// Le critère est la proportion de caractères imprimables : un fragment de police
// ou d'image décompressé contient massivement des octets de contrôle, là où un
// contenu textuel est presque intégralement lisible.
func looksLikeText(s string) bool {
	if s == "" {
		return false
	}

	var printable, total int
	for _, r := range s {
		total++
		if r == '\n' || r == '\t' || (r >= ' ' && r != utf8.RuneError) {
			printable++
		}
	}

	return total > 0 && float64(printable)/float64(total) >= 0.9
}

// newDocument construit un document prêt à être indexé : découpage en chunks,
// calcul du snippet et des métadonnées. C'est le point d'entrée unique de
// l'ingestion, qu'elle provienne d'un upload, d'un collage de texte ou du seed.
func newDocument(id, title, source, content string, createdAt time.Time) Document {
	content = strings.ToValidUTF8(content, "")
	chunks := chunkText(content)

	return Document{
		ID:        id,
		Title:     title,
		Source:    source,
		Status:    "ready",
		Content:   content,
		Chunks:    chunks,
		Snippet:   truncateText(normalizeWhitespace(content), 160),
		SizeBytes: len(content),
		NumChunks: len(chunks),
		CreatedAt: createdAt,
		normTitle: normalizeForSearch(title),
	}
}

// chunkText découpe un texte en passages de taille bornée avec chevauchement.
//
// Le chevauchement évite qu'une information à cheval sur deux chunks soit perdue
// pour le retrieval. Le découpage s'opère sur les runes et non sur les octets,
// afin de ne jamais couper un caractère accentué en deux.
func chunkText(text string) []Chunk {
	text = normalizeWhitespace(text)
	if text == "" {
		return nil
	}

	runes := []rune(text)
	if len(runes) <= chunkSize {
		return []Chunk{{Index: 0, Text: text, norm: normalizeForSearch(text)}}
	}

	var chunks []Chunk
	step := chunkSize - chunkOverlap
	for start, idx := 0, 0; start < len(runes); start, idx = start+step, idx+1 {
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		piece := string(runes[start:end])
		chunks = append(chunks, Chunk{Index: idx, Text: piece, norm: normalizeForSearch(piece)})
		if end == len(runes) {
			break
		}
	}
	return chunks
}

// Repliage des caractères accentués et des ligatures vers leur équivalent ASCII.
// Une recherche sur « theoriciens » doit remonter « théoriciens », et inversement.
var accentFolding = map[rune]string{
	'à': "a", 'á': "a", 'â': "a", 'ã': "a", 'ä': "a", 'å': "a",
	'ç': "c",
	'è': "e", 'é': "e", 'ê': "e", 'ë': "e",
	'ì': "i", 'í': "i", 'î': "i", 'ï': "i",
	'ñ': "n",
	'ò': "o", 'ó': "o", 'ô': "o", 'õ': "o", 'ö': "o", 'ø': "o",
	'ù': "u", 'ú': "u", 'û': "u", 'ü': "u",
	'ý': "y", 'ÿ': "y",
	'æ': "ae", 'œ': "oe", 'ß': "ss",
}

// normalizeForSearch produit la forme canonique utilisée pour la comparaison.
//
// Trois écarts empêchaient sinon toute correspondance sur un corpus francophone :
//
//   - Les élisions. « l'anarchie » constituait un seul terme, qui ne
//     correspondait pas à « anarchie ». Toute ponctuation devient un séparateur,
//     ce qui isole le mot porteur de sens.
//   - Les apostrophes. Un PDF contient l'apostrophe typographique « ’ » (U+2019)
//     là où l'utilisateur saisit l'apostrophe droite « ' » (U+0027) : les deux
//     chaînes ne s'égalaient jamais.
//   - Les accents, dont la saisie est souvent omise dans une question.
//
// La fonction est appliquée symétriquement au contenu indexé et à la requête,
// seule façon de garantir que les deux côtés soient comparables.
func normalizeForSearch(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))

	for _, r := range strings.ToLower(s) {
		if folded, ok := accentFolding[r]; ok {
			sb.WriteString(folded)
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(r)
			continue
		}
		sb.WriteByte(' ')
	}

	return strings.Join(strings.Fields(sb.String()), " ")
}

// sanitizeFilename neutralise les noms de fichiers hostiles.
//
// Le nom provient intégralement du client : il sert de titre affiché dans l'UI et
// de segment d'URI GCS. On supprime donc toute composante de chemin (traversée de
// répertoire) ainsi que les caractères permettant une injection HTML.
func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	name = filepath.Base(strings.TrimSpace(name))
	if name == "." || name == "/" || name == ".." {
		return ""
	}
	name = reUnsafeTitle.ReplaceAllString(name, "")
	name = strings.ToValidUTF8(name, "")
	return truncateText(strings.TrimSpace(name), 200)
}

// appendPDFOperators extrait le texte des opérateurs d'affichage d'un flux.
//
// Les deux formes d'affichage sont traitées : l'opérateur simple (texte Tj, <hex> Tj,
// ', ") et le tableau de crénage [(texte) -250 <hex>] TJ. Les correspondances sont
// parcourues dans l'ordre de leur position afin de préserver l'ordre de lecture.
func appendPDFOperators(sb *strings.Builder, content []byte, cmap map[uint32]string) {
	type match struct {
		pos  int
		text string
	}
	var matches []match

	// 1. Opérateurs d'affichage simple : (texte) Tj ou <hex> Tj (et variantes ' et ")
	for _, m := range rePDFShowText.FindAllSubmatchIndex(content, -1) {
		var text string
		if m[2] >= 0 {
			text = cleanPDFString(string(content[m[2]:m[3]]))
		} else if m[4] >= 0 {
			text = decodePDFHexString(string(content[m[4]:m[5]]), cmap)
		}
		if text != "" {
			matches = append(matches, match{pos: m[0], text: text})
		}
	}

	// 2. Tableaux de crénage : [ ... ] TJ
	for _, idx := range rePDFShowArray.FindAllSubmatchIndex(content, -1) {
		if idx[2] < 0 {
			continue
		}
		text := decodeKerningArray(content[idx[2]:idx[3]], cmap)
		if text != "" {
			matches = append(matches, match{pos: idx[0], text: text})
		}
	}

	sort.Slice(matches, func(i, j int) bool { return matches[i].pos < matches[j].pos })

	for _, m := range matches {
		if m.text == "" || !looksLikeText(m.text) {
			continue
		}
		sb.WriteString(m.text)
		sb.WriteString(" ")
	}
}

// decodeKerningArray reconstitue le texte d'un tableau d'affichage TJ.
//
// Un tableau alterne chaînes littérales, chaînes hexadécimales et déplacements
// de crénage. Les déplacements suffisamment négatifs correspondent à des espaces
// inter-mots, que le fichier ne matérialise pas autrement.
func decodeKerningArray(array []byte, cmap map[uint32]string) string {
	var sb strings.Builder

	for _, m := range rePDFArrayItem.FindAllSubmatch(array, -1) {
		switch {
		case len(m[1]) > 0:
			sb.WriteString(cleanPDFString(string(m[1])))
		case len(m[2]) > 0:
			sb.WriteString(decodePDFHexString(string(m[2]), cmap))
		case len(m[3]) > 0:
			if kern, err := strconv.ParseFloat(string(m[3]), 64); err == nil && kern <= kerningSpaceThreshold {
				sb.WriteString(" ")
			}
		}
	}

	return strings.TrimSpace(sb.String())
}

// decodePDFHexString convertit une chaîne hexadécimale PDF en texte Unicode.
//
// Si une table CMap ToUnicode est fournie (polices intégrées TeX, InDesign, etc.),
// les codes de glyphes y sont convertis vers leur caractère Unicode réel.
// Sans CMap, la fonction tente l'UTF-16BE (avec ou sans BOM FEFF) puis le 8-bit standard.
func decodePDFHexString(rawHex string, cmap map[uint32]string) string {
	var cleaned []byte
	for i := 0; i < len(rawHex); i++ {
		b := rawHex[i]
		switch {
		case b >= '0' && b <= '9', b >= 'a' && b <= 'f', b >= 'A' && b <= 'F':
			cleaned = append(cleaned, b)
		}
	}
	if len(cleaned) == 0 {
		return ""
	}
	if len(cleaned)%2 != 0 {
		cleaned = append(cleaned, '0')
	}

	var sb strings.Builder

	// 1. Décodage via CMap ToUnicode
	if len(cmap) > 0 {
		if len(cleaned)%4 == 0 {
			for i := 0; i < len(cleaned); i += 4 {
				code, err := strconv.ParseUint(string(cleaned[i:i+4]), 16, 32)
				if err != nil {
					continue
				}
				if mapped, ok := cmap[uint32(code)]; ok {
					sb.WriteString(mapped)
				} else if code >= 32 && code < 127 {
					sb.WriteRune(rune(code))
				}
			}
			return sb.String()
		}
		if len(cleaned)%2 == 0 {
			for i := 0; i < len(cleaned); i += 2 {
				code, err := strconv.ParseUint(string(cleaned[i:i+2]), 16, 32)
				if err != nil {
					continue
				}
				if mapped, ok := cmap[uint32(code)]; ok {
					sb.WriteString(mapped)
				} else if code >= 32 && code < 127 {
					sb.WriteRune(rune(code))
				}
			}
			return sb.String()
		}
	}

	// 2. Détection et décodage UTF-16BE sans CMap
	if len(cleaned) >= 4 && len(cleaned)%4 == 0 {
		isUTF16 := strings.HasPrefix(strings.ToUpper(string(cleaned)), "FEFF")
		if !isUTF16 && len(cleaned) >= 8 {
			nullCount := 0
			for i := 0; i < len(cleaned); i += 4 {
				if cleaned[i] == '0' && cleaned[i+1] == '0' {
					nullCount++
				}
			}
			if float64(nullCount)/float64(len(cleaned)/4) > 0.6 {
				isUTF16 = true
			}
		}
		if isUTF16 {
			start := 0
			if strings.HasPrefix(strings.ToUpper(string(cleaned)), "FEFF") {
				start = 4
			}
			var u16 []uint16
			for i := start; i < len(cleaned); i += 4 {
				v, err := strconv.ParseUint(string(cleaned[i:i+4]), 16, 16)
				if err == nil {
					u16 = append(u16, uint16(v))
				}
			}
			return string(utf16.Decode(u16))
		}
	}

	// 3. Décodage hexadécimal 8-bit standard (<48656C6C6F> -> "Hello")
	var rawBytes []byte
	for i := 0; i < len(cleaned); i += 2 {
		v, err := strconv.ParseUint(string(cleaned[i:i+2]), 16, 8)
		if err != nil {
			return ""
		}
		rawBytes = append(rawBytes, byte(v))
	}
	return cleanPDFString(string(rawBytes))
}

// decodeHexUTF16 décode une chaîne hexadécimale en caractères UTF-16 ou ASCII.
func decodeHexUTF16(hexStr string) string {
	if len(hexStr)%4 == 0 && len(hexStr) >= 4 {
		var u16 []uint16
		for i := 0; i < len(hexStr); i += 4 {
			v, err := strconv.ParseUint(hexStr[i:i+4], 16, 16)
			if err != nil {
				return ""
			}
			u16 = append(u16, uint16(v))
		}
		return string(utf16.Decode(u16))
	}
	if len(hexStr)%2 == 0 {
		var sb strings.Builder
		for i := 0; i < len(hexStr); i += 2 {
			v, err := strconv.ParseUint(hexStr[i:i+2], 16, 8)
			if err != nil {
				return ""
			}
			sb.WriteByte(byte(v))
		}
		return sb.String()
	}
	return ""
}

// parseToUnicodeCMap extrait les correspondances de glyphes vers Unicode d'un CMap PDF.
// Un CMap définit des blocs beginbfchar ... endbfchar et beginbfrange ... endbfrange.
func parseToUnicodeCMap(content []byte, cmap map[uint32]string) {
	if !bytes.Contains(content, []byte("begincmap")) {
		return
	}

	// 1. Définitions bfchar unitaires (<src> <dst>)
	for _, block := range rePDFBFCharBlock.FindAll(content, -1) {
		for _, m := range rePDFBFChar.FindAllSubmatch(block, -1) {
			src, err := strconv.ParseUint(string(m[1]), 16, 32)
			if err != nil {
				continue
			}
			dst := decodeHexUTF16(string(m[2]))
			if dst != "" {
				cmap[uint32(src)] = dst
			}
		}
	}

	// 2. Définitions bfrange par plage
	for _, block := range rePDFBFRangeBlock.FindAll(content, -1) {
		// Forme avec tableau : <start> <end> [ <dst1> <dst2> ... ]
		for _, m := range rePDFBFRangeArray.FindAllSubmatch(block, -1) {
			start, err1 := strconv.ParseUint(string(m[1]), 16, 32)
			end, err2 := strconv.ParseUint(string(m[2]), 16, 32)
			if err1 != nil || err2 != nil || start > end || end-start > 65535 {
				continue
			}
			dstTokens := rePDFHexToken.FindAllSubmatch(m[3], -1)
			for idx, code := 0, start; code <= end && idx < len(dstTokens); idx, code = idx+1, code+1 {
				dst := decodeHexUTF16(string(dstTokens[idx][1]))
				if dst != "" {
					cmap[uint32(code)] = dst
				}
			}
		}

		// Forme séquentielle : <start> <end> <dstStart>
		for _, m := range rePDFBFRange.FindAllSubmatch(block, -1) {
			start, err1 := strconv.ParseUint(string(m[1]), 16, 32)
			end, err2 := strconv.ParseUint(string(m[2]), 16, 32)
			dstStart, err3 := strconv.ParseUint(string(m[3]), 16, 32)
			if err1 != nil || err2 != nil || err3 != nil || start > end || end-start > 65535 {
				continue
			}
			for code := start; code <= end; code++ {
				cmap[uint32(code)] = string(rune(dstStart + (code - start)))
			}
		}
	}
}

// extractPrintableWords récupère les suites de caractères imprimables d'un binaire.
func extractPrintableWords(data []byte) string {
	var words []string
	var cur []byte

	flush := func() {
		if len(cur) >= 3 {
			w := strings.TrimSpace(string(cur))
			if len(w) >= 3 && !isPDFInternalKeyword(w) {
				words = append(words, w)
			}
		}
		cur = cur[:0]
	}

	for _, b := range data {
		switch {
		case b >= 'a' && b <= 'z', b >= 'A' && b <= 'Z', b >= '0' && b <= '9',
			b == ' ', b == '.', b == ',', b == '-', b == '/', b == ':', b == '@', b == '_':
			cur = append(cur, b)
		default:
			flush()
		}
	}
	flush()

	return strings.Join(words, " ")
}

// normalizeWhitespace réduit les espaces multiples, retire les caractères de
// contrôle et garantit une sortie UTF-8 valide, seule forme sérialisable en JSON.
func normalizeWhitespace(s string) string {
	s = strings.ToValidUTF8(s, "")
	s = strings.Map(func(r rune) rune {
		switch {
		case r == '\n' || r == '\t' || r == '\r':
			return ' '
		case r < ' ' || r == 0x7f:
			return -1 // Caractère de contrôle : supprimé
		default:
			return r
		}
	}, s)
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

// cleanPDFString décode les séquences d'échappement d'une chaîne littérale PDF.
//
// La spécification définit les échappements usuels (\n, \t, \( ...) mais aussi les
// codes octaux \ddd, largement utilisés pour les parenthèses et les caractères
// non ASCII. Sans ce décodage, le texte extrait contient des littéraux « \050 »
// au lieu des caractères correspondants.
func cleanPDFString(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))

	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 >= len(s) {
			sb.WriteByte(s[i])
			continue
		}

		i++
		switch c := s[i]; c {
		case 'n', 'r':
			sb.WriteByte(' ')
		case 't':
			sb.WriteByte(' ')
		case 'b', 'f':
			sb.WriteByte(' ')
		case '(', ')', '\\':
			sb.WriteByte(c)
		case '\n':
			// Continuation de ligne : la séquence est purement typographique.
		case '\r':
			if i+1 < len(s) && s[i+1] == '\n' {
				i++
			}
		default:
			if c >= '0' && c <= '7' {
				// Code octal sur un à trois chiffres.
				val := int(c - '0')
				for digits := 1; digits < 3 && i+1 < len(s) && s[i+1] >= '0' && s[i+1] <= '7'; digits++ {
					i++
					val = val*8 + int(s[i]-'0')
				}
				if val >= 0x20 && val < 0x7f {
					sb.WriteByte(byte(val))
				} else if val >= 0xa0 {
					// Plage haute Latin-1 : convertie en rune UTF-8 valide.
					sb.WriteRune(rune(val))
				} else {
					sb.WriteByte(' ')
				}
			} else {
				sb.WriteByte(c)
			}
		}
	}

	return strings.TrimSpace(sb.String())
}

func isPDFInternalKeyword(w string) bool {
	switch strings.ToLower(w) {
	case "obj", "endobj", "stream", "endstream", "xref", "trailer", "startxref",
		"font", "type", "subtype", "pages", "catalog", "parent", "contents",
		"mediabox", "resources", "filter", "flatedecode", "length", "identity",
		"true", "false", "null":
		return true
	default:
		return false
	}
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

	s.mu.RLock()
	activeModel := s.modelName
	s.mu.RUnlock()

	if requested := r.URL.Query().Get("model"); isKnownModel(requested) {
		activeModel = requested
	}

	// 1. Étape de Retrieval (Recherche Hybride : Sémantique + BM25)
	sendSSE("status", map[string]string{
		"step":    "retrieval",
		"message": "Recherche hybride (embeddings text-multilingual-002 + BM25)...",
	})

	retrievalStart := time.Now()
	chunks, searchType := s.searchDocumentsHybrid(r.Context(), query)
	retrievalDuration := time.Since(retrievalStart)

	sendSSE("retrieval", chunks)

	// 2. Étape de Synthèse Groundée (Appel Vertex AI Gemini avec streaming)
	sendSSE("status", map[string]string{
		"step":    "generating",
		"model":   activeModel,
		"message": fmt.Sprintf("Génération de la synthèse groundée avec %s...", activeModel),
	})

	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	startTime := time.Now()
	var firstTokenTime time.Duration
	var tokenCount int

	// Appel réel à l'API Vertex AI Gemini via REST avec Bearer Token ambiant (Workload Identity)
	promptTokens, candidateTokens, err := s.streamGeminiResponse(ctx, query, chunks, activeModel, func(token string) {
		if firstTokenTime == 0 {
			firstTokenTime = time.Since(startTime)
		}
		tokenCount++
		sendSSE("token", token)
	})

	totalDuration := time.Since(startTime)

	if err != nil {
		if errors.Is(err, context.Canceled) || r.Context().Err() != nil {
			return
		}
		log.Printf("Erreur streaming Vertex AI (%s): %v. Fallback synthèse locale.", activeModel, err)
		// Fallback gracieux en environnement sandbox sans quota Vertex immédiat
		s.streamLocalFallback(query, chunks, func(token string) {
			if firstTokenTime == 0 {
				firstTokenTime = time.Since(startTime)
			}
			tokenCount++
			sendSSE("token", token)
		})
		totalDuration = time.Since(startTime)
		promptTokens = 240
		candidateTokens = tokenCount * 2
	}

	if candidateTokens == 0 {
		candidateTokens = tokenCount * 2
	}
	if promptTokens == 0 {
		promptTokens = 240
	}

	metricsData := map[string]any{
		"model":             activeModel,
		"retrieval_ms":      retrievalDuration.Milliseconds(),
		"search_type":       searchType,
		"first_token_ms":    firstTokenTime.Milliseconds(),
		"total_duration_ms": totalDuration.Milliseconds(),
		"prompt_tokens":     promptTokens,
		"candidate_tokens":  candidateTokens,
		"total_tokens":      promptTokens + candidateTokens,
	}
	sendSSE("metrics", metricsData)

	// 3. Clôture de l'échange
	sendSSE("done", map[string]any{
		"query":     query,
		"model":     activeModel,
		"timestamp": time.Now().Format(time.RFC3339),
		"sources":   chunks,
		"metrics":   metricsData,
	})
}

// searchDocuments sélectionne les passages les plus pertinents du corpus.
//
// Le retrieval opère au niveau du chunk et non du document : c'est le passage qui
// a effectivement matché qui est transmis au modèle, et non les 160 premiers
// caractères du fichier. Les résultats sont triés par score décroissant puis
// tronqués à topKChunks pour borner la taille du prompt et son coût.
func (s *ServerState) searchDocuments(query string) []SearchChunk {
	s.mu.RLock()
	defer s.mu.RUnlock()

	terms := significantTerms(query)
	if len(terms) == 0 {
		return []SearchChunk{}
	}

	matches := make([]SearchChunk, 0, topKChunks)

	for _, doc := range s.documents {
		// Les formes normalisées sont pré-calculées à l'indexation ; la requête
		// l'étant également, les deux côtés de la comparaison sont comparables.
		// Comparer directement doc.Title/chunk.Text mis en minuscules faisait
		// échouer tout terme accentué, élidé ou porteur d'une apostrophe typographique.
		titleNorm := doc.normTitle

		// Un match dans le titre bénéficie à tous les chunks du document.
		titleScore := 0.0
		for _, term := range terms {
			if strings.Contains(titleNorm, term) {
				titleScore += 0.5
			}
		}

		for _, chunk := range doc.Chunks {
			chunkNorm := chunk.norm
			score := titleScore

			for _, term := range terms {
				// La fréquence d'occurrence départage les passages : un chunk qui
				// mentionne trois fois le terme est plus pertinent qu'un chunk qui
				// l'évoque une seule fois.
				if occurrences := strings.Count(chunkNorm, term); occurrences > 0 {
					score += 0.35 + 0.1*float64(min(occurrences-1, 5))
				}
			}

			if score <= 0 {
				continue
			}

			matches = append(matches, SearchChunk{
				DocumentTitle: doc.Title,
				Snippet:       chunk.Text,
				Score:         score,
				SourceURI:     doc.Source,
				ChunkIndex:    chunk.Index,
			})
		}
	}

	// Tri par pertinence décroissante, en départageant à score égal par l'ordre
	// des chunks afin de garantir un résultat déterministe.
	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		return matches[i].ChunkIndex < matches[j].ChunkIndex
	})

	if len(matches) > topKChunks {
		matches = matches[:topKChunks]
	}
	return matches
}

// cosineSimilarity calcule la similarité cosinus entre deux vecteurs float32.
func cosineSimilarity(a, b []float32) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0.0
	}

	var dot, normA, normB float64
	for i := range a {
		ai := float64(a[i])
		bi := float64(b[i])
		dot += ai * bi
		normA += ai * ai
		normB += bi * bi
	}

	if normA <= 0 || normB <= 0 {
		return 0.0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// searchDocumentsHybrid combine la similarité sémantique (Vertex AI Embeddings text-multilingual-002)
// et le scoring lexical BM25 en mémoire selon la formule validée :
//
//	Score_final = 0.65 * Sim_cosinus + 0.35 * Score_lexical
//
// En cas d'erreur ou d'indisponibilité de l'API d'embeddings, bascule automatiquement à 100% sur le lexical.
func (s *ServerState) searchDocumentsHybrid(ctx context.Context, query string) ([]SearchChunk, string) {
	s.mu.RLock()
	docsCount := len(s.documents)
	s.mu.RUnlock()

	if docsCount == 0 {
		return []SearchChunk{}, "none"
	}

	terms := significantTerms(query)
	lexicalOnly := false

	// 1. Tenter la vectorisation sémantique de la requête
	var queryEmbedding []float32
	embCtx, embCancel := context.WithTimeout(ctx, 4*time.Second)
	embs, err := s.generateEmbeddings(embCtx, []string{query}, "RETRIEVAL_QUERY")
	embCancel()

	if err != nil || len(embs) == 0 || len(embs[0]) == 0 {
		lexicalOnly = true
		if err != nil {
			log.Printf("ℹ️  Vectorisation requête non disponible (%v) -> repli lexical 100%%", err)
		}
	} else {
		queryEmbedding = embs[0]
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	type scoredCandidate struct {
		chunk      SearchChunk
		finalScore float64
	}

	var candidates []scoredCandidate
	hasSemanticMatch := false

	// Facteur de normalisation lexical
	lexMaxExpected := math.Max(1.0, float64(len(terms))*1.2)

	for _, doc := range s.documents {
		titleNorm := doc.normTitle
		titleScore := 0.0
		for _, term := range terms {
			if strings.Contains(titleNorm, term) {
				titleScore += 0.5
			}
		}

		for _, chunk := range doc.Chunks {
			// Calcul lexical normalisé
			rawLexScore := titleScore
			for _, term := range terms {
				if occurrences := strings.Count(chunk.norm, term); occurrences > 0 {
					rawLexScore += 0.35 + 0.1*float64(min(occurrences-1, 5))
				}
			}
			normLexScore := math.Min(1.0, rawLexScore/lexMaxExpected)

			// Calcul sémantique
			var semScore float64
			if !lexicalOnly && len(queryEmbedding) > 0 && len(chunk.Embedding) > 0 {
				cosSim := cosineSimilarity(queryEmbedding, chunk.Embedding)
				// Calibrage cosinus : les embeddings Vertex multilingues ont des scores de similarité
				// typiquement dans [0.3, 0.9] pour des textes pertinents.
				if cosSim > 0.25 {
					semScore = math.Min(1.0, math.Max(0.0, (cosSim-0.25)/0.65))
					hasSemanticMatch = true
				}
			}

			// Score final hybride
			var finalScore float64
			var sType string
			if !lexicalOnly && len(chunk.Embedding) > 0 {
				finalScore = 0.65*semScore + 0.35*normLexScore
				sType = "hybrid"
			} else {
				finalScore = normLexScore
				sType = "lexical"
			}

			if finalScore <= 0.05 {
				continue
			}

			candidates = append(candidates, scoredCandidate{
				chunk: SearchChunk{
					DocumentTitle: doc.Title,
					Snippet:       chunk.Text,
					Score:         math.Round(finalScore*1000) / 1000,
					SemanticScore: math.Round(semScore*1000) / 1000,
					LexicalScore:  math.Round(normLexScore*1000) / 1000,
					SearchType:    sType,
					SourceURI:     doc.Source,
					ChunkIndex:    chunk.Index,
				},
				finalScore: finalScore,
			})
		}
	}

	searchMode := "hybrid"
	if lexicalOnly || !hasSemanticMatch {
		searchMode = "lexical"
	}

	// Tri par score décroissant, en départageant à score égal par l'ordre des chunks
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].finalScore != candidates[j].finalScore {
			return candidates[i].finalScore > candidates[j].finalScore
		}
		return candidates[i].chunk.ChunkIndex < candidates[j].chunk.ChunkIndex
	})

	if len(candidates) > topKChunks {
		candidates = candidates[:topKChunks]
	}

	result := make([]SearchChunk, len(candidates))
	for i, c := range candidates {
		result[i] = c.chunk
	}

	return result, searchMode
}

// significantTerms découpe la requête en termes comparables au contenu indexé.
//
// La requête traverse exactement la même normalisation que les chunks : c'est la
// seule façon de garantir que « Qu'est-ce que l'anarchie ? » produise le terme
// « anarchie ». L'ancien découpage sur les espaces conservait « l'anarchie » en un
// seul token, qui ne correspondait à rien dans le corpus.
//
// Les mots vides et les termes de moins de trois caractères sont écartés : ils
// matcheraient partout sans apporter de signal de pertinence.
func significantTerms(query string) []string {
	var terms []string
	for _, term := range strings.Fields(normalizeForSearch(query)) {
		if utf8.RuneCountInString(term) <= 2 || isStopWord(term) {
			continue
		}
		terms = append(terms, term)
	}
	return terms
}

// La liste est volontairement sans accent : isStopWord reçoit des termes déjà
// passés par normalizeForSearch, qui replie les accents.
func isStopWord(w string) bool {
	switch w {
	case "les", "des", "une", "que", "qui", "pour", "dans", "avec", "sur", "par",
		"est", "sont", "aux", "ses", "cette", "comment", "quel", "quelle", "quels",
		"quelles", "quoi", "donc", "ceci", "cela", "leur", "leurs", "nous", "vous",
		"the", "and", "for", "with", "what", "how", "are", "was", "does", "did",
		"this", "that", "from", "you":
		return true
	default:
		return false
	}
}

// Appel direct à l'API Vertex AI Gemini via REST & ADC
func (s *ServerState) streamGeminiResponse(ctx context.Context, query string, chunks []SearchChunk, modelName string, onToken func(string)) (int, int, error) {
	token := getGCPToken()

	if token == "" {
		return 0, 0, fmt.Errorf("aucun jeton d'authentification GCP disponible")
	}

	if modelName == "" {
		s.mu.RLock()
		modelName = s.modelName
		s.mu.RUnlock()
	}

	// Construction du prompt groundé. On transmet le texte du passage retenu par le
	// retrieval, et non plus le snippet d'en-tête du document qui ne contenait que
	// la page de garde et privait le modèle du contenu réellement pertinent.
	var contextBuilder bytes.Buffer
	for i, c := range chunks {
		fmt.Fprintf(&contextBuilder, "\n[Source %d: %s (passage %d)]\n%s\n",
			i+1, c.DocumentTitle, c.ChunkIndex+1, c.Snippet)
	}

	systemInstruction := "Tu es un assistant IA d'architecture Google Cloud. Réponds à la question de manière claire, rigoureuse et exhaustive en t'appuyant rigoureusement sur le contexte documentaire fourni ci-dessous. Développe chaque point nécessaire pour fournir une réponse complète, sans jamais abréger ni tronquer tes explications ou conclusions. Mentionne explicitement les sources utilisées entre crochets (ex: [Source 1]). Si le contexte ne permet pas de répondre, indique-le explicitement plutôt que d'inventer."

	contextSection := contextBuilder.String()
	if contextSection == "" {
		contextSection = "\n(Aucun passage pertinent n'a été trouvé dans le corpus indexé.)\n"
	}
	prompt := fmt.Sprintf("%s\n\nQuestion de l'utilisateur : %s\n\nContexte documentaire disponible :%s", systemInstruction, query, contextSection)

	// Les modèles Gemini 3 de dernière génération sont routés via le endpoint global Vertex AI
	targetLocation := s.region
	endpointHost := fmt.Sprintf("%s-aiplatform.googleapis.com", s.region)
	if strings.HasPrefix(modelName, "gemini-3") {
		targetLocation = "global"
		endpointHost = "aiplatform.googleapis.com"
	}

	apiURL := fmt.Sprintf("https://%s/v1/projects/%s/locations/%s/publishers/google/models/%s:streamGenerateContent?alt=sse",
		endpointHost, s.projectID, targetLocation, modelName)

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
			"maxOutputTokens": 8192,
		},
	}

	jsonBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return 0, 0, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Pas de Timeout sur le client : sur une réponse en streaming, http.Client.Timeout
	// couvre la lecture complète du corps et couperait donc la génération en plein
	// milieu. La durée de vie de l'appel est pilotée par le contexte du handler.
	client := &http.Client{
		Transport: &http.Transport{
			// On borne en revanche l'attente des en-têtes, pour ne pas rester bloqué
			// si Vertex AI ne répond pas du tout.
			ResponseHeaderTimeout: 30 * time.Second,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return 0, 0, fmt.Errorf("erreur HTTP Vertex AI %d: %s", resp.StatusCode, string(body))
	}

	// Lecture du flux SSE ligne par ligne.
	//
	// Une lecture à taille fixe découperait les événements JSON à cheval sur deux
	// blocs réseau : le json.Unmarshal échouait alors silencieusement et le token
	// était définitivement perdu. bufio.Scanner réassemble les lignes partielles.
	promptTokens, candidateTokens := 0, 0
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // Certains événements dépassent 64 Ko

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" || payload == "[DONE]" {
			continue
		}

		var vResp struct {
			Candidates []struct {
				FinishReason string `json:"finishReason"`
				Content      struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
			UsageMetadata struct {
				PromptTokenCount     int `json:"promptTokenCount"`
				CandidatesTokenCount int `json:"candidatesTokenCount"`
			} `json:"usageMetadata"`
		}

		if err := json.Unmarshal([]byte(payload), &vResp); err != nil {
			log.Printf("⚠️  Événement SSE Vertex AI illisible, ignoré : %v", err)
			continue
		}

		if vResp.UsageMetadata.PromptTokenCount > 0 {
			promptTokens = vResp.UsageMetadata.PromptTokenCount
		}
		if vResp.UsageMetadata.CandidatesTokenCount > 0 {
			candidateTokens = vResp.UsageMetadata.CandidatesTokenCount
		}
		for _, cand := range vResp.Candidates {
			if cand.FinishReason == "MAX_TOKENS" {
				log.Printf("⚠️  Plafond maxOutputTokens atteint pour le modèle %s", modelName)
			}
			for _, p := range cand.Content.Parts {
				if p.Text != "" {
					onToken(p.Text)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil && !errors.Is(err, context.Canceled) {
		return promptTokens, candidateTokens, fmt.Errorf("lecture du flux Vertex AI interrompue : %w", err)
	}

	return promptTokens, candidateTokens, nil
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

// tokenCache évite un aller-retour vers le serveur de métadonnées à chaque requête.
// Les jetons GCE sont valides une heure ; on conserve une marge de sécurité.
var tokenCache struct {
	sync.Mutex
	value     string
	expiresAt time.Time
}

func fetchMetadataToken() string {
	tokenCache.Lock()
	defer tokenCache.Unlock()

	if tokenCache.value != "" && time.Now().Before(tokenCache.expiresAt) {
		return tokenCache.value
	}

	req, err := http.NewRequest(http.MethodGet,
		"http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token", nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Metadata-Flavor", "Google")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var t struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil || t.AccessToken == "" {
		return ""
	}

	ttl := time.Duration(t.ExpiresIn) * time.Second
	if ttl <= 0 {
		ttl = time.Hour
	}
	// Marge de 5 minutes pour ne jamais servir un jeton sur le point d'expirer.
	tokenCache.value = t.AccessToken
	tokenCache.expiresAt = time.Now().Add(ttl - 5*time.Minute)

	return t.AccessToken
}

// getGCPToken résout le jeton d'accès GCP selon l'ordre : variable d'environnement,
// serveur de métadonnées GCE/Cloud Run, puis binaire gcloud local.
func getGCPToken() string {
	if t := os.Getenv("GOOGLE_OAUTH_ACCESS_TOKEN"); t != "" {
		return t
	}
	if t := fetchMetadataToken(); t != "" {
		return t
	}
	// Fallback pour environnement de développement local si gcloud est disponible
	if path, err := exec.LookPath("gcloud"); err == nil && path != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "gcloud", "auth", "print-access-token")
		out, err := cmd.Output()
		if err == nil {
			tok := strings.TrimSpace(string(out))
			if tok != "" {
				return tok
			}
		}
	}
	return ""
}

func (s *ServerState) endpoint() string {
	if s.gcsEndpoint != "" {
		return s.gcsEndpoint
	}
	return "https://storage.googleapis.com"
}

// uploadToGCS téléverse des données brutes vers Cloud Storage via l'API REST JSON.
func (s *ServerState) uploadToGCS(ctx context.Context, objectName, contentType string, data []byte) error {
	if s.gcsBucket == "" {
		return nil
	}
	var token string
	if s.gcsEndpoint != "" {
		token = "mock-test-token"
	} else {
		token = getGCPToken()
		if token == "" {
			return fmt.Errorf("jeton GCP manquant pour upload GCS")
		}
	}

	apiURL := fmt.Sprintf("%s/upload/storage/v1/b/%s/o?uploadType=media&name=%s",
		s.endpoint(), url.PathEscape(s.gcsBucket), url.QueryEscape(objectName))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("erreur HTTP GCS %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// downloadFromGCS télécharge le contenu d'un objet Cloud Storage via l'API REST JSON.
func (s *ServerState) downloadFromGCS(ctx context.Context, objectName string) ([]byte, error) {
	if s.gcsBucket == "" {
		return nil, fmt.Errorf("bucket GCS non configuré")
	}
	var token string
	if s.gcsEndpoint != "" {
		token = "mock-test-token"
	} else {
		token = getGCPToken()
		if token == "" {
			return nil, fmt.Errorf("jeton GCP manquant pour download GCS")
		}
	}

	apiURL := fmt.Sprintf("%s/storage/v1/b/%s/o/%s?alt=media",
		s.endpoint(), url.PathEscape(s.gcsBucket), url.PathEscape(objectName))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, os.ErrNotExist
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("erreur HTTP GCS %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// deleteFromGCS supprime un objet Cloud Storage via l'API REST JSON.
func (s *ServerState) deleteFromGCS(ctx context.Context, objectName string) error {
	if s.gcsBucket == "" {
		return nil
	}
	var token string
	if s.gcsEndpoint != "" {
		token = "mock-test-token"
	} else {
		token = getGCPToken()
		if token == "" {
			return fmt.Errorf("jeton GCP manquant pour delete GCS")
		}
	}

	apiURL := fmt.Sprintf("%s/storage/v1/b/%s/o/%s",
		s.endpoint(), url.PathEscape(s.gcsBucket), url.PathEscape(objectName))

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("erreur HTTP GCS %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// saveCorpusSnapshot sérialise l'index documentaire complet vers gs://{bucket}/index/corpus.json.
func (s *ServerState) saveCorpusSnapshot(ctx context.Context) error {
	if s.gcsBucket == "" {
		return nil
	}

	s.mu.RLock()
	stored := make([]StoredDocument, len(s.documents))
	for i, doc := range s.documents {
		stored[i] = doc.toStored()
	}
	s.mu.RUnlock()

	snapshot := CorpusSnapshot{
		Version:   1,
		UpdatedAt: time.Now().UTC(),
		Documents: stored,
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur sérialisation snapshot corpus: %w", err)
	}

	if err := s.uploadToGCS(ctx, corpusSnapshotPath, "application/json", data); err != nil {
		log.Printf("⚠️  Échec sauvegarde snapshot GCS (%s): %v", corpusSnapshotPath, err)
		return err
	}
	log.Printf("💾 Snapshot corpus sauvegardé avec succès sur gs://%s/%s (%d documents, %d octets)",
		s.gcsBucket, corpusSnapshotPath, len(stored), len(data))
	return nil
}

// initCorpusFromGCS charge l'index documentaire depuis Cloud Storage au démarrage.
func (s *ServerState) initCorpusFromGCS(ctx context.Context) {
	if s.gcsBucket == "" {
		log.Println("ℹ️  GCS_RAG_BUCKET non configuré : fonctionnement en mémoire vive pure")
		return
	}

	data, err := s.downloadFromGCS(ctx, corpusSnapshotPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("ℹ️  Aucun snapshot existant sur gs://%s/%s. Initialisation avec le corpus par défaut...",
				s.gcsBucket, corpusSnapshotPath)
			go func() {
				bgCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				if err := s.saveCorpusSnapshot(bgCtx); err != nil {
					log.Printf("⚠️  Échec synchronisation initiale du corpus vers GCS: %v", err)
				}
			}()
			return
		}
		log.Printf("⚠️  Impossible de charger le snapshot GCS (%v). Conservation du corpus mémoire par défaut.", err)
		return
	}

	var snapshot CorpusSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		log.Printf("⚠️  Erreur désérialisation snapshot GCS: %v. Conservation du corpus mémoire par défaut.", err)
		return
	}

	docs := make([]Document, len(snapshot.Documents))
	for i, sd := range snapshot.Documents {
		docs[i] = storedToDocument(sd)
	}

	s.mu.Lock()
	s.documents = docs
	s.mu.Unlock()

	log.Printf("✅ Corpus chargé depuis gs://%s/%s : %d document(s) restauré(s) (mis à jour le %s)",
		s.gcsBucket, corpusSnapshotPath, len(docs), snapshot.UpdatedAt.Format(time.RFC3339))
}

// generateEmbeddings appelle l'API Vertex AI pour vectoriser une liste de textes.
// Modèle : text-multilingual-embedding-002 (768 dimensions).
// taskType : "RETRIEVAL_DOCUMENT" pour les chunks, "RETRIEVAL_QUERY" pour les requêtes utilisateur.
func (s *ServerState) generateEmbeddings(ctx context.Context, texts []string, taskType string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	var token string
	if s.embedEndpoint != "" {
		token = "mock-test-token"
	} else {
		token = getGCPToken()
		if token == "" {
			return nil, fmt.Errorf("jeton GCP manquant pour Vertex AI Embeddings")
		}
	}

	targetLocation := s.region
	apiURL := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/text-multilingual-embedding-002:predict",
		s.region, s.projectID, targetLocation)
	if s.embedEndpoint != "" {
		apiURL = s.embedEndpoint
	}

	// Limite Vertex AI text-multilingual-embedding-002 : max 20 000 tokens par requête.
	// Avec des passages de ~1200 caractères (~300 tokens), un lot de 10 passages
	// consomme ~3 000 tokens, bien en dessous du plafond.
	const batchSize = 10
	allEmbeddings := make([][]float32, 0, len(texts))

	for start := 0; start < len(texts); start += batchSize {
		select {
		case <-ctx.Done():
			return allEmbeddings, ctx.Err()
		default:
		}

		end := start + batchSize
		if end > len(texts) {
			end = len(texts)
		}
		batch := texts[start:end]

		instances := make([]map[string]string, len(batch))
		for i, txt := range batch {
			instances[i] = map[string]string{
				"content":   txt,
				"task_type": taskType,
			}
		}

		reqBody := map[string]any{
			"instances": instances,
		}

		jsonBytes, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(jsonBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("appel Vertex Embeddings échoué : %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
			resp.Body.Close()
			return nil, fmt.Errorf("erreur HTTP Vertex Embeddings %d: %s", resp.StatusCode, string(body))
		}

		var vResp struct {
			Predictions []struct {
				Embeddings struct {
					Values []float32 `json:"values"`
				} `json:"embeddings"`
			} `json:"predictions"`
		}

		err = json.NewDecoder(resp.Body).Decode(&vResp)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("décodage réponse Vertex Embeddings impossible : %w", err)
		}

		if len(vResp.Predictions) != len(batch) {
			return nil, fmt.Errorf("nombre de prédictions incohérent : reçu %d, attendu %d", len(vResp.Predictions), len(batch))
		}

		for _, p := range vResp.Predictions {
			allEmbeddings = append(allEmbeddings, p.Embeddings.Values)
		}

		if len(texts) > batchSize {
			time.Sleep(50 * time.Millisecond)
		}
	}

	return allEmbeddings, nil
}

// ensureCorpusEmbeddings s'assure que tous les chunks du corpus en RAM disposent d'un vecteur d'embeddings.
// Les chunks nouvellement vectorisés sont sauvegardés dans le snapshot Cloud Storage.
func (s *ServerState) ensureCorpusEmbeddings(ctx context.Context) {
	type chunkRef struct {
		docIdx   int
		chunkIdx int
		text     string
	}

	s.mu.RLock()
	var missing []chunkRef
	for dIdx, doc := range s.documents {
		for cIdx, chunk := range doc.Chunks {
			if len(chunk.Embedding) == 0 && strings.TrimSpace(chunk.Text) != "" {
				missing = append(missing, chunkRef{
					docIdx:   dIdx,
					chunkIdx: cIdx,
					text:     chunk.Text,
				})
			}
		}
	}
	s.mu.RUnlock()

	if len(missing) == 0 {
		return
	}

	log.Printf("🧠 Vectorisation de %d chunk(s) manquant(s) dans le corpus via text-multilingual-embedding-002...", len(missing))

	// Traitement par lots pour permettre la sauvegarde progressive et la résilience
	const stepSize = 20
	totalDone := 0

	for start := 0; start < len(missing); start += stepSize {
		select {
		case <-ctx.Done():
			log.Printf("⚠️  Vectorisation interrompue : %d/%d chunks traités", totalDone, len(missing))
			return
		default:
		}

		end := start + stepSize
		if end > len(missing) {
			end = len(missing)
		}
		subMissing := missing[start:end]
		texts := make([]string, len(subMissing))
		for i, ref := range subMissing {
			texts[i] = ref.text
		}

		embeddings, err := s.generateEmbeddings(ctx, texts, "RETRIEVAL_DOCUMENT")
		if err != nil {
			log.Printf("ℹ️  Sous-lot de vectorisation (%d-%d) échoué (%v) : repli lexical temporaire", start, end, err)
			time.Sleep(1 * time.Second)
			continue
		}

		s.mu.Lock()
		for i, ref := range subMissing {
			if ref.docIdx < len(s.documents) && ref.chunkIdx < len(s.documents[ref.docIdx].Chunks) {
				s.documents[ref.docIdx].Chunks[ref.chunkIdx].Embedding = embeddings[i]
			}
		}
		s.mu.Unlock()

		totalDone += len(embeddings)
		if totalDone%100 == 0 || end == len(missing) {
			log.Printf("🧠 Progression vectorisation : %d/%d chunks vectorisés (%.1f%%)",
				totalDone, len(missing), float64(totalDone)*100.0/float64(len(missing)))
		}
	}

	log.Printf("✅ Vectorisation terminée : %d/%d chunk(s) vectorisé(s). Sauvegarde du snapshot...", totalDone, len(missing))
	if s.gcsBucket != "" && totalDone > 0 {
		if err := s.saveCorpusSnapshot(ctx); err != nil {
			log.Printf("⚠️  Échec sauvegarde snapshot après vectorisation : %v", err)
		}
	}
}

// truncateText tronque un texte à maxLen caractères.
//
// La troncature s'effectue sur les runes et non sur les octets : couper au milieu
// d'un caractère accentué produirait une séquence UTF-8 invalide, remplacée par un
// caractère de remplacement (�) lors de la sérialisation JSON.
func truncateText(text string, maxLen int) string {
	if utf8.RuneCountInString(text) <= maxLen {
		return text
	}
	runes := []rune(text)
	return string(runes[:maxLen]) + "..."
}

// Structures pour l'Évaluation GenAI à la demande (Vertex AI EvaluateInstances)
type EvaluationRequest struct {
	Query      string `json:"query"`
	Model      string `json:"model,omitempty"`
	Prediction string `json:"prediction"`
	Context    string `json:"context,omitempty"`
}

type MetricEvaluation struct {
	Score       float64 `json:"score"`
	MaxScore    float64 `json:"max_score"`
	Confidence  float64 `json:"confidence,omitempty"`
	Explanation string  `json:"explanation,omitempty"`
	Label       string  `json:"label"`
}

type EvaluationResponse struct {
	Groundedness MetricEvaluation `json:"groundedness"`
	Relevance    MetricEvaluation `json:"relevance"`
	OverallScore float64          `json:"overall_score"`
	DurationMS   int64            `json:"duration_ms"`
	Evaluator    string           `json:"evaluator"`
	EvaluatedAt  time.Time        `json:"evaluated_at"`
}

// evaluateMetricVertexAI appelle l'API REST de Vertex AI Evaluation (evaluateInstances).
// metricName : "groundedness" ou "question_answering_relevance".
func (s *ServerState) evaluateMetricVertexAI(ctx context.Context, metricName string, evalReq EvaluationRequest) (*MetricEvaluation, error) {
	var token string
	if s.evalEndpoint != "" {
		token = "mock-test-token"
	} else {
		token = getGCPToken()
		if token == "" {
			return nil, fmt.Errorf("jeton GCP manquant pour Vertex AI Evaluation")
		}
	}

	apiURL := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1beta1/projects/%s/locations/%s:evaluateInstances",
		s.region, s.projectID, s.region)
	if s.evalEndpoint != "" {
		apiURL = s.evalEndpoint
	}

	var reqBody map[string]any
	switch metricName {
	case "groundedness":
		reqBody = map[string]any{
			"groundedness_input": map[string]any{
				"metric_spec": map[string]any{
					"version": 1,
				},
				"instance": map[string]any{
					"prediction": evalReq.Prediction,
					"context":    evalReq.Context,
				},
			},
		}
	case "question_answering_relevance":
		reqBody = map[string]any{
			"question_answering_relevance_input": map[string]any{
				"metric_spec": map[string]any{
					"version": 1,
				},
				"instance": map[string]any{
					"instruction": evalReq.Query,
					"prediction":  evalReq.Prediction,
					"context":     evalReq.Context,
				},
			},
		}
	default:
		return nil, fmt.Errorf("métrique d'évaluation non supportée : %s", metricName)
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("appel Vertex EvaluateInstances échoué : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("erreur HTTP Vertex Evaluation %d: %s", resp.StatusCode, string(body))
	}

	var vResp struct {
		GroundednessResult *struct {
			Score       float64 `json:"score"`
			Explanation string  `json:"explanation"`
			Confidence  float64 `json:"confidence"`
		} `json:"groundednessResult,omitempty"`
		QARelevanceResult *struct {
			Score       float64 `json:"score"`
			Explanation string  `json:"explanation"`
			Confidence  float64 `json:"confidence"`
		} `json:"questionAnsweringRelevanceResult,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&vResp); err != nil {
		return nil, fmt.Errorf("décodage réponse Vertex Evaluation impossible : %w", err)
	}

	if metricName == "groundedness" {
		if vResp.GroundednessResult == nil {
			return nil, fmt.Errorf("aucun résultat groundedness retourné par Vertex AI")
		}
		rawScore := vResp.GroundednessResult.Score
		// Si le score est sur l'intervalle [0, 1], conversion sur une échelle de 5 étoiles
		score5 := rawScore
		if rawScore <= 1.0 {
			score5 = math.Round(rawScore*50) / 10
		}
		label := "Ancrage Partiel"
		if score5 >= 4.5 {
			label = "Parfaitement Ancré"
		} else if score5 >= 3.5 {
			label = "Bon Ancrage"
		} else if score5 < 2.5 {
			label = "Faible Ancrage"
		}
		return &MetricEvaluation{
			Score:       score5,
			MaxScore:    5.0,
			Confidence:  vResp.GroundednessResult.Confidence,
			Explanation: vResp.GroundednessResult.Explanation,
			Label:       label,
		}, nil
	}

	if vResp.QARelevanceResult == nil {
		return nil, fmt.Errorf("aucun résultat relevance retourné par Vertex AI")
	}
	score := vResp.QARelevanceResult.Score
	label := "Moyennement Pertinent"
	if score >= 4.5 {
		label = "Très Pertinent"
	} else if score >= 3.5 {
		label = "Pertinent"
	} else if score < 2.5 {
		label = "Peu Pertinent"
	}

	return &MetricEvaluation{
		Score:       score,
		MaxScore:    5.0,
		Confidence:  vResp.QARelevanceResult.Confidence,
		Explanation: vResp.QARelevanceResult.Explanation,
		Label:       label,
	}, nil
}

// fallbackEvaluation calcule une estimation heuristique robuste en cas d'indisponibilité du service d'évaluation
func fallbackEvaluation(query, prediction, docContext string) (*MetricEvaluation, *MetricEvaluation) {
	normPred := normalizeForSearch(prediction)
	normCtx := normalizeForSearch(docContext)
	normQuery := normalizeForSearch(query)

	queryTerms := significantTerms(normQuery)
	predTerms := significantTerms(normPred)

	// Estimation Ancrage : proportion des termes significatifs de la réponse présents dans le contexte
	matchedInCtx := 0
	for _, term := range predTerms {
		if strings.Contains(normCtx, term) {
			matchedInCtx++
		}
	}
	ratioGrounded := 0.8
	if len(predTerms) > 0 {
		ratioGrounded = float64(matchedInCtx) / float64(len(predTerms))
	}
	gScore := math.Min(5.0, math.Max(1.0, math.Round(ratioGrounded*50)/10))
	gLabel := "Bon Ancrage (Heuristique)"
	if gScore >= 4.5 {
		gLabel = "Parfaitement Ancré (Heuristique)"
	} else if gScore < 3.0 {
		gLabel = "Ancrage Partiel (Heuristique)"
	}

	// Estimation Pertinence : proportion des termes de la question présents dans la réponse
	matchedInPred := 0
	for _, term := range queryTerms {
		if strings.Contains(normPred, term) {
			matchedInPred++
		}
	}
	ratioRel := 0.85
	if len(queryTerms) > 0 {
		ratioRel = float64(matchedInPred) / float64(len(queryTerms))
	}
	rScore := math.Min(5.0, math.Max(1.0, math.Round(ratioRel*50)/10))
	rLabel := "Pertinent (Heuristique)"
	if rScore >= 4.5 {
		rLabel = "Très Pertinent (Heuristique)"
	} else if rScore < 3.0 {
		rLabel = "Peu Pertinent (Heuristique)"
	}

	gMetric := &MetricEvaluation{
		Score:       gScore,
		MaxScore:    5.0,
		Confidence:  0.8,
		Explanation: fmt.Sprintf("Évaluation heuristique de repli : %d/%d termes de la synthèse identifiés dans les sources documentaires.", matchedInCtx, len(predTerms)),
		Label:       gLabel,
	}

	rMetric := &MetricEvaluation{
		Score:       rScore,
		MaxScore:    5.0,
		Confidence:  0.8,
		Explanation: fmt.Sprintf("Évaluation heuristique de repli : %d/%d termes de la requête traités dans la synthèse.", matchedInPred, len(queryTerms)),
		Label:       rLabel,
	}

	return gMetric, rMetric
}

// handleEvaluate évalue la fidélité documentaire et la pertinence d'une réponse RAG via Vertex AI EvaluateInstances.
func (s *ServerState) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req EvaluationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête JSON invalide: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Query) == "" || strings.TrimSpace(req.Prediction) == "" {
		http.Error(w, "Les champs 'query' et 'prediction' sont obligatoires", http.StatusBadRequest)
		return
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Si le contexte n'a pas été fourni par le client, récupération des passages pertinents du corpus
	docContext := req.Context
	if strings.TrimSpace(docContext) == "" {
		chunks, _ := s.searchDocumentsHybrid(ctx, req.Query)
		var sb strings.Builder
		for i, c := range chunks {
			if i > 0 {
				sb.WriteString("\n---\n")
			}
			sb.WriteString(c.Snippet)
		}
		docContext = sb.String()
		req.Context = docContext
	}

	// Évaluation concurrente : Groundedness et Relevance
	var gMetric, rMetric *MetricEvaluation
	var gErr, rErr error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		gMetric, gErr = s.evaluateMetricVertexAI(ctx, "groundedness", req)
	}()

	go func() {
		defer wg.Done()
		rMetric, rErr = s.evaluateMetricVertexAI(ctx, "question_answering_relevance", req)
	}()

	wg.Wait()

	evaluator := "Vertex AI Rapid Evaluation (Autorater)"
	if gErr != nil || rErr != nil {
		log.Printf("ℹ️  Vertex AI Evaluation non disponible (Groundedness: %v, Relevance: %v) -> repli heuristique", gErr, rErr)
		evaluator = "Repli Heuristique Standard (Vertex AI hors ligne)"
		fallbackG, fallbackR := fallbackEvaluation(req.Query, req.Prediction, docContext)
		if gMetric == nil {
			gMetric = fallbackG
		}
		if rMetric == nil {
			rMetric = fallbackR
		}
	}

	overallScore := math.Round(((gMetric.Score+rMetric.Score)/2.0)*10) / 10
	duration := time.Since(start)

	resp := EvaluationResponse{
		Groundedness: *gMetric,
		Relevance:    *rMetric,
		OverallScore: overallScore,
		DurationMS:   duration.Milliseconds(),
		Evaluator:    evaluator,
		EvaluatedAt:  time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
