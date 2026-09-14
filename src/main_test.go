package main

import (
	"bytes"
	"compress/zlib"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// truncateText tronquait auparavant sur les octets, ce qui coupait les caractères
// accentués en deux et produisait des séquences UTF-8 invalides dans le JSON.
func TestTruncateTextPreservesUTF8(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"texte court inchangé", "Bonjour", 20, "Bonjour"},
		{"longueur exacte", "abcde", 5, "abcde"},
		{"troncature ascii", "abcdefghij", 5, "abcde..."},
		{"coupure sur accent", "ééééééééé", 4, "éééé..."},
		{"vide", "", 10, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := truncateText(tc.input, tc.maxLen)
			if got != tc.want {
				t.Errorf("truncateText(%q, %d) = %q, attendu %q", tc.input, tc.maxLen, got, tc.want)
			}
			if !utf8.ValidString(got) {
				t.Errorf("truncateText a produit une chaîne UTF-8 invalide : %q", got)
			}
		})
	}
}

func TestTruncateTextNeverEmitsReplacementChar(t *testing.T) {
	// Chaîne francophone tronquée à toutes les longueurs possibles : aucune
	// troncature ne doit générer de caractère de remplacement.
	input := "Le déploiement sécurisé requiert Cloud Armor, l'égress via Cloud NAT et l'accès par IAP."
	for i := 1; i <= utf8.RuneCountInString(input); i++ {
		got := truncateText(input, i)
		if strings.ContainsRune(got, '\uFFFD') {
			t.Fatalf("troncature à %d runes a produit un caractère de remplacement : %q", i, got)
		}
	}
}

func TestChunkTextSplitsWithOverlap(t *testing.T) {
	// Texte plus long que chunkSize pour forcer le découpage.
	long := strings.Repeat("alpha beta gamma delta ", 200)

	chunks := chunkText(long)
	if len(chunks) < 2 {
		t.Fatalf("attendu plusieurs chunks pour un texte de %d runes, obtenu %d",
			utf8.RuneCountInString(long), len(chunks))
	}

	for i, c := range chunks {
		if c.Index != i {
			t.Errorf("chunk %d porte l'index %d", i, c.Index)
		}
		if n := utf8.RuneCountInString(c.Text); n > chunkSize {
			t.Errorf("chunk %d dépasse la taille maximale : %d > %d", i, n, chunkSize)
		}
	}

	// Le chevauchement garantit que la fin d'un chunk se retrouve au début du suivant.
	first := []rune(chunks[0].Text)
	tail := string(first[len(first)-chunkOverlap:])
	if !strings.HasPrefix(chunks[1].Text, tail) {
		t.Errorf("chevauchement absent entre les chunks 0 et 1")
	}
}

func TestChunkTextShortInput(t *testing.T) {
	if got := chunkText(""); got != nil {
		t.Errorf("un texte vide ne doit produire aucun chunk, obtenu %v", got)
	}

	chunks := chunkText("Texte court")
	if len(chunks) != 1 || chunks[0].Text != "Texte court" {
		t.Errorf("un texte court doit produire un unique chunk intact, obtenu %v", chunks)
	}
}

func TestChunkTextPreservesUTF8(t *testing.T) {
	long := strings.Repeat("éàüöç ", 500)
	for _, c := range chunkText(long) {
		if !utf8.ValidString(c.Text) {
			t.Fatalf("chunk %d contient de l'UTF-8 invalide", c.Index)
		}
	}
}

// searchDocuments doit trier par score décroissant et borner le nombre de résultats,
// au lieu de renvoyer les documents dans leur ordre d'insertion.
func TestSearchDocumentsSortsByScore(t *testing.T) {
	s := &ServerState{documents: []Document{
		newDocument("d1", "Document sans rapport", "gs://x/1", "Contenu neutre sur des sujets variés.", time.Now()),
		newDocument("d2", "Guide réseau", "gs://x/2", "Le cloud nat gère la sortie. Le cloud nat est mutualisé. Le cloud nat facture par heure.", time.Now()),
		newDocument("d3", "Note brève", "gs://x/3", "Une mention unique de cloud nat ici.", time.Now()),
	}}

	results := s.searchDocuments("cloud nat")
	if len(results) < 2 {
		t.Fatalf("attendu au moins 2 résultats, obtenu %d", len(results))
	}

	// Le document mentionnant trois fois le terme doit primer sur la mention unique.
	if results[0].DocumentTitle != "Guide réseau" {
		t.Errorf("premier résultat attendu 'Guide réseau', obtenu %q", results[0].DocumentTitle)
	}

	for i := 1; i < len(results); i++ {
		if results[i-1].Score < results[i].Score {
			t.Errorf("résultats non triés : score[%d]=%.2f < score[%d]=%.2f",
				i-1, results[i-1].Score, i, results[i].Score)
		}
	}

	// Aucun document sans occurrence ne doit apparaître.
	for _, r := range results {
		if r.DocumentTitle == "Document sans rapport" {
			t.Error("un document sans occurrence a été retourné")
		}
	}
}

func TestSearchDocumentsRespectsTopK(t *testing.T) {
	docs := make([]Document, 0, 20)
	for i := 0; i < 20; i++ {
		docs = append(docs, newDocument(
			"d", "Doc terraform", "gs://x/d",
			"Terraform applique l'infrastructure. Terraform gère l'état distant.", time.Now()))
	}
	s := &ServerState{documents: docs}

	if results := s.searchDocuments("terraform"); len(results) > topKChunks {
		t.Errorf("attendu au maximum %d résultats, obtenu %d", topKChunks, len(results))
	}
}

func TestSearchDocumentsEmptyQuery(t *testing.T) {
	s := &ServerState{documents: []Document{
		newDocument("d1", "Doc", "gs://x/1", "Contenu quelconque", time.Now()),
	}}

	// Une requête composée uniquement de mots vides ne doit rien remonter,
	// plutôt que de renvoyer arbitrairement le premier document du corpus.
	if results := s.searchDocuments("les des une"); len(results) != 0 {
		t.Errorf("attendu aucun résultat pour une requête de mots vides, obtenu %d", len(results))
	}
}

// sanitizeFilename doit neutraliser les noms de fichiers permettant une injection
// HTML ou une traversée de répertoire.
func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"rapport.pdf", "rapport.pdf"},
		{"../../etc/passwd", "passwd"},
		{`C:\Windows\system32\config`, "config"},
		{"  espaces.txt  ", "espaces.txt"},
		{"..", ""},
	}

	for _, tc := range tests {
		got := sanitizeFilename(tc.input)
		if got != tc.want {
			t.Errorf("sanitizeFilename(%q) = %q, attendu %q", tc.input, got, tc.want)
		}
	}

	// Quel que soit le nom fourni, aucun caractère permettant une injection HTML
	// ne doit subsister dans le titre affiché par l'interface.
	hostile := []string{
		`<img src=x onerror=alert(1)>.pdf`,
		`"><script>alert(1)</script>.txt`,
		"rapport&co.pdf",
	}
	for _, name := range hostile {
		if got := sanitizeFilename(name); strings.ContainsAny(got, `<>"'&`) {
			t.Errorf("sanitizeFilename(%q) a laissé passer un caractère dangereux : %q", name, got)
		}
	}
}

// extractTextFromPDF doit décompresser les flux FlateDecode, qui représentent la
// quasi-totalité des PDF réels.
func TestExtractTextFromPDFWithFlateStream(t *testing.T) {
	const phrase = "Architecture securisee avec Cloud Armor et Cloud NAT pour la production"

	var compressed bytes.Buffer
	zw := zlib.NewWriter(&compressed)
	if _, err := zw.Write([]byte("BT /F1 12 Tf (" + phrase + ") Tj ET")); err != nil {
		t.Fatalf("écriture du flux compressé : %v", err)
	}
	zw.Close()

	var pdf bytes.Buffer
	pdf.WriteString("%PDF-1.4\n4 0 obj\n<< /Length 120 /Filter /FlateDecode >>\nstream\n")
	pdf.Write(compressed.Bytes())
	pdf.WriteString("\nendstream\nendobj\n%%EOF")

	got := extractTextFromPDF(pdf.Bytes())
	if !strings.Contains(got, "Cloud Armor") {
		t.Errorf("texte compressé non extrait, obtenu %q", got)
	}
}

func TestExtractTextFromPDFUncompressed(t *testing.T) {
	pdf := []byte("%PDF-1.4\nBT (Le budget FinOps est plafonne par Cloud Billing) Tj ET\n%%EOF")

	got := extractTextFromPDF(pdf)
	if !strings.Contains(got, "budget FinOps") {
		t.Errorf("texte non compressé non extrait, obtenu %q", got)
	}
}

// Les PDF produits par TeX/dvips placent le texte dans des tableaux de crénage.
// Ignorer cette forme revenait à n'extraire aucun texte de ces documents.
func TestExtractTextFromPDFKerningArray(t *testing.T) {
	content := `BT /F82 20 Tf [(V)94(ersion)-375(Managemen)31(t)]TJ ET`

	var compressed bytes.Buffer
	zw := zlib.NewWriter(&compressed)
	zw.Write([]byte(content))
	zw.Close()

	var pdf bytes.Buffer
	pdf.WriteString("%PDF-1.4\n1 0 obj\n<< /Filter /FlateDecode >>\nstream\n")
	pdf.Write(compressed.Bytes())
	pdf.WriteString("\nendstream\nendobj\n%%EOF")

	got := extractTextFromPDF(pdf.Bytes())

	// Les crénages faibles (94, 31) restent à l'intérieur d'un mot, le crénage
	// fortement négatif (-375) matérialise l'espace entre les deux mots.
	if !strings.Contains(got, "Version Management") {
		t.Errorf("tableau de crénage mal décodé : attendu 'Version Management', obtenu %q", got)
	}
}

// Les codes octaux \050 et \051 doivent produire les parenthèses correspondantes.
func TestCleanPDFStringDecodesEscapes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`con\050figuration\051`, "con(figuration)"},
		{`ligne\nsuivante`, "ligne suivante"},
		{`parenthese \( echappee`, "parenthese ( echappee"},
		{`antislash \\ seul`, `antislash \ seul`},
		{"texte simple", "texte simple"},
	}

	for _, tc := range tests {
		if got := cleanPDFString(tc.input); got != tc.want {
			t.Errorf("cleanPDFString(%q) = %q, attendu %q", tc.input, got, tc.want)
		}
	}
}

// Un flux binaire ne doit jamais produire de pseudo-texte : c'est ce qui polluait
// le corpus avec des séquences illisibles lors des uploads de PDF.
func TestExtractTextFromPDFRejectsBinaryNoise(t *testing.T) {
	binary := make([]byte, 4096)
	for i := range binary {
		binary[i] = byte(i*7 + i/3) // Motif pseudo-aléatoire non textuel
	}

	var pdf bytes.Buffer
	pdf.WriteString("%PDF-1.4\nstream\n")
	pdf.Write(binary)
	pdf.WriteString("\nendstream\n%%EOF")

	got := extractTextFromPDF(pdf.Bytes())
	if !utf8.ValidString(got) {
		t.Errorf("sortie UTF-8 invalide : %q", got)
	}
	if !looksLikeText(got) {
		t.Errorf("du binaire a été remonté dans le corpus : %q", got)
	}
}

func TestLooksLikeText(t *testing.T) {
	if !looksLikeText("Architecture Cloud Armor et Cloud NAT") {
		t.Error("un texte lisible doit être accepté")
	}
	if looksLikeText(string([]byte{0x01, 0x02, 0x03, 0x00, 0x7f, 0x05})) {
		t.Error("une séquence binaire doit être rejetée")
	}
	if looksLikeText("") {
		t.Error("une chaîne vide doit être rejetée")
	}
}

func TestExtractTextFromPDFWithoutTextLayer(t *testing.T) {
	got := extractTextFromPDF([]byte{0x00, 0x01, 0x02, 0xff, 0xfe})
	if got == "" {
		t.Error("extractTextFromPDF ne doit jamais renvoyer une chaîne vide")
	}
	if !utf8.ValidString(got) {
		t.Errorf("résultat UTF-8 invalide : %q", got)
	}
}

func TestNewDocumentComputesMetadata(t *testing.T) {
	content := strings.Repeat("Analyse de l'architecture cloud. ", 100)
	doc := newDocument("id-1", "Titre", "gs://b/o", content, time.Now())

	if doc.NumChunks != len(doc.Chunks) {
		t.Errorf("NumChunks (%d) incohérent avec le nombre de chunks (%d)", doc.NumChunks, len(doc.Chunks))
	}
	if doc.NumChunks == 0 {
		t.Error("aucun chunk produit pour un contenu non vide")
	}
	if utf8.RuneCountInString(doc.Snippet) > 163 { // 160 runes + "..."
		t.Errorf("snippet trop long : %d runes", utf8.RuneCountInString(doc.Snippet))
	}
}

func TestIsKnownModel(t *testing.T) {
	if !isKnownModel(AvailableModels[0].ID) {
		t.Errorf("le modèle %q du catalogue devrait être reconnu", AvailableModels[0].ID)
	}
	if isKnownModel("gemini-inexistant") {
		t.Error("un modèle hors catalogue ne doit pas être accepté")
	}
	if isKnownModel("") {
		t.Error("une chaîne vide ne doit pas être acceptée")
	}
}

func TestSignificantTermsFiltersNoise(t *testing.T) {
	got := significantTerms("Comment sont gérés les budgets, et le WAF ?")

	for _, term := range got {
		if isStopWord(term) {
			t.Errorf("le mot vide %q n'a pas été filtré", term)
		}
		if utf8.RuneCountInString(term) <= 2 {
			t.Errorf("le terme trop court %q n'a pas été filtré", term)
		}
		if strings.ContainsAny(term, ",?") {
			t.Errorf("la ponctuation n'a pas été retirée de %q", term)
		}
	}

	if len(got) == 0 {
		t.Error("tous les termes ont été filtrés, attendu au moins 'budgets' et 'waf'")
	}
}
