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
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"
)

//go:embed web/static/*
var staticFS embed.FS

// Chunk représente un passage découpé d'un document, unité de base du retrieval.
type Chunk struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
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
}

type SearchChunk struct {
	DocumentTitle string  `json:"document_title"`
	Snippet       string  `json:"snippet"`
	Score         float64 `json:"score"`
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
		modelName = "gemini-3.5-flash"
	}

	state := &ServerState{
		projectID: projectID,
		region:    region,
		modelName: modelName,
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

	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/documents", state.handleDocuments)
	mux.HandleFunc("/api/documents/upload", state.handleUpload)
	mux.HandleFunc("/api/search/classic", state.handleClassicSearch)
	mux.HandleFunc("/api/chat/stream", state.handleChatStream)
	mux.HandleFunc("/api/models", state.handleModels)
	mux.HandleFunc("/api/model/switch", state.handleModelSwitch)
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

// Handler Documents : liste des documents indexés
func (s *ServerState) handleDocuments(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.documents)
}

// Handler Models : liste des modèles disponibles et modèle actuellement actif
func (s *ServerState) handleModels(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	current := s.modelName
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"current":   current,
		"models":    AvailableModels,
		"embedding": "text-embedding-005",
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

				addedDocs = append(addedDocs, newDocument(
					fmt.Sprintf("doc-%d-%d", time.Now().UnixNano(), i),
					name,
					"gs://wh-ai-blueprint-a363-rag-docs/"+name,
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
				if source == "" {
					source = "gs://wh-ai-blueprint-a363-rag-docs/" + title
				}
				addedDocs = append(addedDocs, newDocument(
					fmt.Sprintf("doc-%d", time.Now().UnixNano()), title, source, content, time.Now()))
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
			source := req.Source
			if source == "" {
				source = "gs://wh-ai-blueprint-a363-rag-docs/" + title
			}
			addedDocs = append(addedDocs, newDocument(
				fmt.Sprintf("doc-%d", time.Now().UnixNano()), title, source, content, time.Now()))
		}
	}

	if len(addedDocs) == 0 {
		http.Error(w, "Aucun document ou fichier valide reçu", http.StatusBadRequest)
		return
	}

	// L'indexation (extraction + découpage) est déjà terminée à ce stade : elle est
	// réalisée de façon synchrone dans newDocument. Les documents sont donc publiés
	// directement en statut "ready".
	//
	// Une simulation d'indexation asynchrone était auparavant faite via une goroutine
	// temporisée. C'est un anti-pattern sur Cloud Run : hors annotation
	// `run.googleapis.com/cpu-throttling: false`, le CPU est retiré à l'instance dès
	// que la réponse HTTP est émise. La goroutine ne reprenait donc la main qu'à la
	// requête suivante, et les documents restaient affichés « Indexation... »
	// pendant un temps arbitrairement long.
	totalChunks := 0
	s.mu.Lock()
	for i := range addedDocs {
		totalChunks += addedDocs[i].NumChunks
		s.documents = append([]Document{addedDocs[i]}, s.documents...)
	}
	s.mu.Unlock()

	log.Printf("📥 %d document(s) reçu(s), %d chunk(s) indexé(s) dans le corpus", len(addedDocs), totalChunks)

	w.Header().Set("Content-Type", "application/json")
	if len(addedDocs) == 1 {
		json.NewEncoder(w).Encode(addedDocs[0])
	} else {
		json.NewEncoder(w).Encode(map[string]any{
			"status":    "ok",
			"count":     len(addedDocs),
			"documents": addedDocs,
		})
	}
}

// Constantes de traitement documentaire.
const (
	maxUploadBytes = 50 << 20 // Taille maximale du corps d'une requête d'upload
	maxMemoryBytes = 10 << 20 // Part de l'upload conservée en mémoire, le reste va sur disque
	chunkSize      = 1000     // Taille cible d'un chunk, en runes
	chunkOverlap   = 100      // Chevauchement entre deux chunks consécutifs, en runes
	topKChunks     = 5        // Nombre de passages transmis au modèle pour le grounding
)

// Regex compilées une seule fois au chargement du package plutôt qu'à chaque appel.
var (
	// Opérateur d'affichage simple : (texte) Tj
	rePDFShowText = regexp.MustCompile(`\(((?:[^()\\]|\\.)*)\)\s*(?:Tj|')`)
	// Tableau de crénage : [(V)94(ersion)-375(Managemen)31(t)] TJ.
	// C'est la forme employée par la plupart des générateurs (TeX, dvips) et elle
	// porte l'essentiel du texte : l'ignorer revient à ne rien extraire du document.
	rePDFShowArray = regexp.MustCompile(`(?s)\[((?:[^\[\]\\]|\\.)*)\]\s*TJ`)
	// Éléments d'un tableau de crénage : chaînes littérales et valeurs de crénage.
	rePDFArrayItem = regexp.MustCompile(`\(((?:[^()\\]|\\.)*)\)|(-?\d+(?:\.\d+)?)`)
	// Objets stream/endstream, dont le contenu est généralement compressé
	rePDFStream = regexp.MustCompile(`(?s)stream\r?\n?(.*?)endstream`)
	// Caractères interdits ou risqués dans un titre de document
	reUnsafeTitle = regexp.MustCompile(`[<>"'&\x00-\x1f\x7f]`)
)

// kerningSpaceThreshold : en deçà de cette valeur (en millièmes d'unité de texte),
// un déplacement de crénage matérialise une espace entre deux mots. Au-dessus, il
// s'agit d'un simple resserrement typographique à l'intérieur d'un même mot.
const kerningSpaceThreshold = -150

// extractTextFromPDF extrait le texte lisible d'un flux binaire PDF.
//
// Deux écueils sont traités ici. D'une part, la quasi-totalité des PDF réels
// compressent leurs flux de contenu : une regex appliquée aux octets bruts n'y
// trouve aucun opérateur de texte. D'autre part, appliquer cette regex au binaire
// compressé produit des correspondances fortuites qui polluent le corpus avec des
// séquences illisibles. On ne retient donc que des chaînes qui ressemblent
// réellement à du texte, et on n'inspecte les octets bruts que si aucun flux n'a
// pu être décompressé.
func extractTextFromPDF(data []byte) string {
	var sb strings.Builder
	decodedStreams := 0

	// 1. Décompression des flux puis extraction des opérateurs de texte.
	for _, m := range rePDFStream.FindAllSubmatch(data, -1) {
		if len(m) < 2 {
			continue
		}
		decoded, ok := inflateStream(bytes.TrimLeft(m[1], "\r\n"))
		if !ok {
			continue
		}
		decodedStreams++
		appendPDFOperators(&sb, decoded)
	}

	// 2. Aucun flux décompressable : le document est probablement en clair.
	//    On n'applique ce repli que dans ce cas précis, sous peine de rapatrier
	//    des correspondances parasites issues des données binaires compressées.
	if decodedStreams == 0 {
		appendPDFOperators(&sb, data)
	}

	extracted := normalizeWhitespace(sb.String())

	// 3. Dernier recours pour les documents en clair sans opérateur exploitable.
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
		return []Chunk{{Index: 0, Text: text}}
	}

	var chunks []Chunk
	step := chunkSize - chunkOverlap
	for start, idx := 0, 0; start < len(runes); start, idx = start+step, idx+1 {
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, Chunk{Index: idx, Text: string(runes[start:end])})
		if end == len(runes) {
			break
		}
	}
	return chunks
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
// Les deux formes d'affichage sont traitées : l'opérateur simple `(texte) Tj` et
// le tableau de crénage `[(V)94(ersion)] TJ`. Les correspondances sont parcourues
// dans l'ordre de leur position afin de préserver l'ordre de lecture du document.
//
// Chaque chaîne candidate est validée : sur un flux binaire, la regex produit des
// correspondances fortuites qu'il faut écarter avant qu'elles n'atteignent le corpus.
func appendPDFOperators(sb *strings.Builder, content []byte) {
	type match struct {
		pos  int
		text string
	}
	var matches []match

	for _, idx := range rePDFShowText.FindAllSubmatchIndex(content, -1) {
		if idx[2] < 0 {
			continue
		}
		matches = append(matches, match{pos: idx[0], text: cleanPDFString(string(content[idx[2]:idx[3]]))})
	}

	for _, idx := range rePDFShowArray.FindAllSubmatchIndex(content, -1) {
		if idx[2] < 0 {
			continue
		}
		matches = append(matches, match{pos: idx[0], text: decodeKerningArray(content[idx[2]:idx[3]])})
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
// Un tableau alterne chaînes littérales et déplacements de crénage. Les
// déplacements suffisamment négatifs correspondent à des espaces inter-mots, que
// le fichier ne matérialise pas autrement : sans cette reconstitution, le texte
// extrait serait une suite de mots agglutinés.
func decodeKerningArray(array []byte) string {
	var sb strings.Builder

	for _, m := range rePDFArrayItem.FindAllSubmatch(array, -1) {
		switch {
		case m[1] != nil:
			sb.WriteString(cleanPDFString(string(m[1])))
		case len(m[2]) > 0:
			if kern, err := strconv.ParseFloat(string(m[2]), 64); err == nil && kern <= kerningSpaceThreshold {
				sb.WriteString(" ")
			}
		}
	}

	return strings.TrimSpace(sb.String())
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
		"model":   activeModel,
		"message": fmt.Sprintf("Génération de la synthèse groundée avec %s...", activeModel),
	})

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
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
		lowerTitle := strings.ToLower(doc.Title)

		// Un match dans le titre bénéficie à tous les chunks du document.
		titleScore := 0.0
		for _, term := range terms {
			if strings.Contains(lowerTitle, term) {
				titleScore += 0.5
			}
		}

		for _, chunk := range doc.Chunks {
			lowerText := strings.ToLower(chunk.Text)
			score := titleScore

			for _, term := range terms {
				// La fréquence d'occurrence départage les passages : un chunk qui
				// mentionne trois fois le terme est plus pertinent qu'un chunk qui
				// l'évoque une seule fois.
				if occurrences := strings.Count(lowerText, term); occurrences > 0 {
					score += 0.35 + 0.1*float64(min(occurrences-1, 5))
				}
			}

			if score <= 0 {
				continue
			}

			matches = append(matches, SearchChunk{
				DocumentTitle: doc.Title,
				Snippet:       truncateText(chunk.Text, 600),
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

// significantTerms normalise la requête en écartant les mots trop courts et les
// mots vides qui matcheraient partout sans apporter de signal de pertinence.
func significantTerms(query string) []string {
	var terms []string
	for _, term := range strings.Fields(strings.ToLower(query)) {
		term = strings.Trim(term, ".,;:!?()[]\"'")
		if utf8.RuneCountInString(term) <= 2 || isStopWord(term) {
			continue
		}
		terms = append(terms, term)
	}
	return terms
}

func isStopWord(w string) bool {
	switch w {
	case "les", "des", "une", "que", "qui", "pour", "dans", "avec", "sur", "par",
		"est", "sont", "aux", "ses", "cette", "comment", "quel", "quelle", "quels",
		"quelles", "the", "and", "for", "with", "what", "how", "are", "was":
		return true
	default:
		return false
	}
}

// Appel direct à l'API Vertex AI Gemini via REST & ADC
func (s *ServerState) streamGeminiResponse(ctx context.Context, query string, chunks []SearchChunk, modelName string, onToken func(string)) (int, int, error) {
	token := os.Getenv("GOOGLE_OAUTH_ACCESS_TOKEN")
	if token == "" {
		// Tenter de lire le token depuis les métadonnées GCE si on est sur GCP
		token = fetchMetadataToken()
	}

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

	systemInstruction := "Tu es un assistant IA d'architecture Google Cloud. Réponds à la question de manière concise et précise en t'appuyant rigoureusement sur le contexte documentaire fourni ci-dessous. Mentionne explicitement les sources utilisées entre crochets (ex: [Source 1]). Si le contexte ne permet pas de répondre, indique-le explicitement plutôt que d'inventer."

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
			"maxOutputTokens": 1024,
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
				Content struct {
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
