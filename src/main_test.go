package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
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

func TestNormalizeForSearch(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"apostrophe droite", "L'anarchie", "l anarchie"},
		{"apostrophe typographique", "L\u2019anarchie", "l anarchie"},
		{"accents repliés", "Théoriciens de l'État", "theoriciens de l etat"},
		{"trait d'union", "cloud-nat", "cloud nat"},
		{"ligature", "Sœur & cœur", "soeur coeur"},
		{"ponctuation terminale", "Qu'est-ce que l'anarchie ?", "qu est ce que l anarchie"},
		{"espaces multiples", "  deux\t\nmots  ", "deux mots"},
		{"chiffres conservés", "Article 49.3", "article 49 3"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeForSearch(tc.in); got != tc.want {
				t.Errorf("normalizeForSearch(%q) = %q, attendu %q", tc.in, got, tc.want)
			}
		})
	}
}

// Les deux formes d'apostrophe doivent produire la même représentation, sans quoi
// une question saisie au clavier ne peut pas matcher un PDF, qui emploie U+2019.
func TestNormalizeForSearchUnifiesApostrophes(t *testing.T) {
	straight := normalizeForSearch("Qu'est-ce que l'anarchie ?")
	typographic := normalizeForSearch("Qu\u2019est-ce que l\u2019anarchie ?")

	if straight != typographic {
		t.Errorf("les deux apostrophes divergent : %q vs %q", straight, typographic)
	}
}

func TestSignificantTermsSplitsElision(t *testing.T) {
	for _, query := range []string{
		"Qu'est-ce que l'anarchie ?",
		"Qu\u2019est-ce que l\u2019anarchie ?",
		"l'anarchie",
		"ANARCHIE",
	} {
		terms := significantTerms(query)
		if !slices.Contains(terms, "anarchie") {
			t.Errorf("significantTerms(%q) = %v, attendu le terme 'anarchie'", query, terms)
		}
	}
}

// Reproduction du défaut signalé : les documents étaient bien indexés mais la
// question « Qu'est-ce que l'anarchie ? » ne remontait aucun passage, car le
// terme restait collé à son élision et l'apostrophe du PDF différait de celle
// saisie par l'utilisateur.
func TestSearchDocumentsFrenchElision(t *testing.T) {
	s := &ServerState{documents: []Document{
		newDocument(
			"d1",
			"Th\u00e9oriciens politiques",
			"gs://x/1",
			"L\u2019anarchie d\u00e9signe une organisation politique sans autorit\u00e9 centrale. "+
				"Les partisans de l\u2019anarchie pr\u00f4nent l\u2019autogestion.",
			time.Now(),
		),
	}}

	for _, query := range []string{
		"Qu'est-ce que l'anarchie ?",
		"Qu\u2019est-ce que l\u2019anarchie ?",
		"l'anarchie",
		"anarchie",
		"ANARCHIE",
	} {
		if got := s.searchDocuments(query); len(got) == 0 {
			t.Errorf("searchDocuments(%q) n'a remont\u00e9 aucun passage", query)
		}
	}
}

// La question peut être saisie sans accent : le repliage doit fonctionner dans
// les deux sens puisqu'il est appliqué au contenu comme à la requête.
func TestSearchDocumentsIgnoresAccents(t *testing.T) {
	s := &ServerState{documents: []Document{
		newDocument("d1", "Note", "gs://x/1", "La s\u00e9curit\u00e9 p\u00e9rim\u00e9trique repose sur le pare-feu.", time.Now()),
	}}

	if got := s.searchDocuments("securite perimetrique"); len(got) == 0 {
		t.Error("une requête sans accent doit retrouver un contenu accentué")
	}
	if got := s.searchDocuments("s\u00e9curit\u00e9"); len(got) == 0 {
		t.Error("une requête accentuée doit retrouver le même contenu")
	}
}

func TestParseToUnicodeCMap(t *testing.T) {
	cmapContent := []byte(`
/CIDInit /ProcSet findresource begin
12 dict begin
begincmap
1 begincodespacerange
<0000> <FFFF>
endcodespacerange
2 beginbfchar
<0002> <0021>
<0051> <0070>
endbfchar
1 beginbfrange
<0042> <0044> <0061>
endbfrange
endcmap
`)

	cmap := make(map[uint32]string)
	parseToUnicodeCMap(cmapContent, cmap)

	if got := cmap[0x0002]; got != "!" {
		t.Errorf("cmap[0x0002] = %q, attendu '!'", got)
	}
	if got := cmap[0x0051]; got != "p" {
		t.Errorf("cmap[0x0051] = %q, attendu 'p'", got)
	}
	if got := cmap[0x0042]; got != "a" {
		t.Errorf("cmap[0x0042] = %q, attendu 'a'", got)
	}
	if got := cmap[0x0043]; got != "b" {
		t.Errorf("cmap[0x0043] = %q, attendu 'b'", got)
	}
	if got := cmap[0x0044]; got != "c" {
		t.Errorf("cmap[0x0044] = %q, attendu 'c'", got)
	}
}

func TestDecodePDFHexStringWithCMap(t *testing.T) {
	cmap := map[uint32]string{
		0x0051: "p",
		0x0053: "r",
		0x0046: "e",
	}

	got := decodePDFHexString("005100530046", cmap)
	if got != "pre" {
		t.Errorf("decodePDFHexString avec CMap = %q, attendu 'pre'", got)
	}
}

func TestDecodePDFHexStringUTF16BE(t *testing.T) {
	// "FEFF" + "0048 0065 006C 006C 006F" -> "Hello"
	got := decodePDFHexString("FEFF00480065006C006C006F", nil)
	if got != "Hello" {
		t.Errorf("decodePDFHexString UTF-16BE = %q, attendu 'Hello'", got)
	}
}

func TestDecodePDFHexStringASCIIHex(t *testing.T) {
	// "48656C6C6F" -> "Hello"
	got := decodePDFHexString("48656C6C6F", nil)
	if got != "Hello" {
		t.Errorf("decodePDFHexString ASCII hex = %q, attendu 'Hello'", got)
	}
}

func TestExtractTextFromPDFWithCMapAndHexArray(t *testing.T) {
	cmapRaw := []byte(`
begincmap
1 beginbfchar
<0001> <0041>
endbfchar
endcmap
`)
	var cmapBuf bytes.Buffer
	zw := zlib.NewWriter(&cmapBuf)
	zw.Write(cmapRaw)
	zw.Close()

	contentRaw := []byte(`
BT
/F1 12 Tf
[<0001> -250 (l'anarchie)] TJ
ET
`)
	var contentBuf bytes.Buffer
	zw2 := zlib.NewWriter(&contentBuf)
	zw2.Write(contentRaw)
	zw2.Close()

	var pdf bytes.Buffer
	pdf.WriteString("%PDF-1.4\n")
	pdf.WriteString("1 0 obj << /Length 100 >> stream\n")
	pdf.Write(cmapBuf.Bytes())
	pdf.WriteString("\nendstream\nendobj\n")
	pdf.WriteString("2 0 obj << /Length 100 >> stream\n")
	pdf.Write(contentBuf.Bytes())
	pdf.WriteString("\nendstream\nendobj\n%%EOF")

	got := extractTextFromPDF(pdf.Bytes())
	if !strings.Contains(got, "A") || !strings.Contains(got, "l'anarchie") {
		t.Errorf("extractTextFromPDF = %q, attendu 'A l'anarchie'", got)
	}
}

func TestHandleDocumentsGet(t *testing.T) {
	s := &ServerState{
		documents: []Document{
			newDocument("doc-1", "Doc 1", "gs://bucket/1", "Contenu 1", time.Now()),
			newDocument("doc-2", "Doc 2", "gs://bucket/2", "Contenu 2", time.Now()),
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/documents", nil)
	w := httptest.NewRecorder()
	s.handleDocuments(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("handleDocuments GET code = %d, attendu %d", w.Code, http.StatusOK)
	}

	var docs []Document
	if err := json.NewDecoder(w.Body).Decode(&docs); err != nil {
		t.Fatalf("décodage json impossible: %v", err)
	}
	if len(docs) != 2 {
		t.Fatalf("attendu 2 documents, obtenu %d", len(docs))
	}
	if docs[0].ID != "doc-1" || docs[1].ID != "doc-2" {
		t.Errorf("documents reçus inattendus: %+v", docs)
	}
}

func TestHandleDocumentsDeleteSingle(t *testing.T) {
	s := &ServerState{
		documents: []Document{
			newDocument("doc-1", "Doc 1", "gs://bucket/1", "Contenu 1", time.Now()),
			newDocument("doc-2", "Doc 2", "gs://bucket/2", "Contenu 2", time.Now()),
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/documents?id=doc-1", nil)
	w := httptest.NewRecorder()
	s.handleDocuments(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("handleDocuments DELETE single code = %d, attendu %d", w.Code, http.StatusOK)
	}

	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("décodage json impossible: %v", err)
	}
	if resp["status"] != "deleted" || resp["id"] != "doc-1" {
		t.Errorf("réponse inattendue: %+v", resp)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.documents) != 1 || s.documents[0].ID != "doc-2" {
		t.Errorf("état des documents après suppression incorrect: %+v", s.documents)
	}
}

func TestHandleDocumentsDeleteNotFound(t *testing.T) {
	s := &ServerState{
		documents: []Document{
			newDocument("doc-1", "Doc 1", "gs://bucket/1", "Contenu 1", time.Now()),
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/documents?id=non-existant", nil)
	w := httptest.NewRecorder()
	s.handleDocuments(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("handleDocuments DELETE not found code = %d, attendu %d", w.Code, http.StatusNotFound)
	}
}

func TestHandleDocumentsDeleteAll(t *testing.T) {
	s := &ServerState{
		documents: []Document{
			newDocument("doc-1", "Doc 1", "gs://bucket/1", "Contenu 1", time.Now()),
			newDocument("doc-2", "Doc 2", "gs://bucket/2", "Contenu 2", time.Now()),
			newDocument("doc-3", "Doc 3", "gs://bucket/3", "Contenu 3", time.Now()),
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/documents?all=true", nil)
	w := httptest.NewRecorder()
	s.handleDocuments(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("handleDocuments DELETE all code = %d, attendu %d", w.Code, http.StatusOK)
	}

	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("décodage json impossible: %v", err)
	}
	if resp["status"] != "cleared" {
		t.Errorf("réponse inattendue: %+v", resp)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.documents) != 0 {
		t.Errorf("le corpus n'a pas été entièrement vidé: %+v", s.documents)
	}
}

func TestHandleDocumentsDeleteBadRequest(t *testing.T) {
	s := &ServerState{
		documents: []Document{
			newDocument("doc-1", "Doc 1", "gs://bucket/1", "Contenu 1", time.Now()),
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/documents", nil)
	w := httptest.NewRecorder()
	s.handleDocuments(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("handleDocuments DELETE sans paramètre code = %d, attendu %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleUploadReturnsTelemetry(t *testing.T) {
	s := &ServerState{}
	payload := `{"title": "Test Telemetrie", "content": "Contenu riche pour tester la generation de chunks et la telemetrie d'upload."}`
	req := httptest.NewRequest(http.MethodPost, "/api/documents/upload", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleUpload(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("handleUpload code = %d, attendu %d", w.Code, http.StatusOK)
	}

	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("décodage json impossible: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("status attendu 'ok', obtenu %v", resp["status"])
	}
	if cnt, ok := resp["count"].(float64); !ok || cnt != 1 {
		t.Errorf("count attendu 1, obtenu %v", resp["count"])
	}
	if chunks, ok := resp["total_chunks"].(float64); !ok || chunks < 1 {
		t.Errorf("total_chunks attendu >= 1, obtenu %v", resp["total_chunks"])
	}
	if bytes, ok := resp["total_bytes"].(float64); !ok || bytes <= 0 {
		t.Errorf("total_bytes attendu > 0, obtenu %v", resp["total_bytes"])
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.documents) != 1 {
		t.Errorf("document non ajouté au corpus, len = %d", len(s.documents))
	}
}

func TestStoredDocumentSerializationAndRestoration(t *testing.T) {
	orig := newDocument("doc-123",
		"Architecture Haute Disponibilité GCP",
		"gs://bucket/test.pdf",
		"Le déploiement multi-régional garantit une résilience maximale contre les sinistres zonaux.",
		time.Date(2026, 9, 15, 10, 0, 0, 0, time.UTC),
	)

	stored := orig.toStored()
	if stored.ID != orig.ID {
		t.Errorf("ID attendu %q, obtenu %q", orig.ID, stored.ID)
	}
	if stored.Title != orig.Title {
		t.Errorf("Title attendu %q, obtenu %q", orig.Title, stored.Title)
	}
	if len(stored.Chunks) != len(orig.Chunks) {
		t.Fatalf("nombre de chunks attendu %d, obtenu %d", len(orig.Chunks), len(stored.Chunks))
	}

	restored := storedToDocument(stored)
	if restored.ID != orig.ID || restored.Title != orig.Title || restored.Source != orig.Source {
		t.Errorf("champs de base non restaurés correctement: %+v", restored)
	}
	if restored.normTitle != orig.normTitle {
		t.Errorf("normTitle non recalculé: %q != %q", restored.normTitle, orig.normTitle)
	}
	if len(restored.Chunks) != len(orig.Chunks) {
		t.Fatalf("nombre de chunks restaurés différent: %d != %d", len(restored.Chunks), len(orig.Chunks))
	}
	for i := range restored.Chunks {
		if restored.Chunks[i].norm != orig.Chunks[i].norm {
			t.Errorf("chunk %d norm non recalculé: %q != %q", i, restored.Chunks[i].norm, orig.Chunks[i].norm)
		}
	}
}

func TestCorpusSnapshotJSONSerialization(t *testing.T) {
	doc := newDocument("d1", "Test Snapshot", "gs://b/d1", "Texte du document de test pour snapshot.", time.Now())
	snapshot := CorpusSnapshot{
		Version:   1,
		UpdatedAt: time.Now().UTC(),
		Documents: []StoredDocument{doc.toStored()},
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatalf("erreur encodage JSON snapshot: %v", err)
	}

	var decoded CorpusSnapshot
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("erreur décodage JSON snapshot: %v", err)
	}

	if decoded.Version != 1 || len(decoded.Documents) != 1 {
		t.Fatalf("snapshot décodé invalide: %+v", decoded)
	}
	if decoded.Documents[0].ID != "d1" {
		t.Errorf("document ID attendu 'd1', obtenu %q", decoded.Documents[0].ID)
	}
}

func TestGCSOperationsWithMockServer(t *testing.T) {
	storage := make(map[string][]byte)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Non autorisé", http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodPost:
			// Upload URL: /upload/storage/v1/b/{bucket}/o?uploadType=media&name={name}
			name := r.URL.Query().Get("name")
			if name == "" {
				http.Error(w, "nom manquant", http.StatusBadRequest)
				return
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			storage[name] = body
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{"name": name, "size": len(body)})

		case http.MethodGet:
			// Download URL: /storage/v1/b/{bucket}/o/{name}?alt=media
			path := strings.TrimPrefix(r.URL.Path, "/storage/v1/b/test-bucket/o/")
			data, ok := storage[path]
			if !ok {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(data)

		case http.MethodDelete:
			// Delete URL: /storage/v1/b/{bucket}/o/{name}
			path := strings.TrimPrefix(r.URL.Path, "/storage/v1/b/test-bucket/o/")
			delete(storage, path)
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Méthode non supportée", http.StatusMethodNotAllowed)
		}
	}))
	defer ts.Close()

	s := &ServerState{
		gcsBucket:   "test-bucket",
		gcsEndpoint: ts.URL,
	}

	ctx := context.Background()

	// 1. Test Upload
	testData := []byte("contenu secret gcs test")
	if err := s.uploadToGCS(ctx, "test/doc.txt", "text/plain", testData); err != nil {
		t.Fatalf("uploadToGCS a échoué: %v", err)
	}

	// 2. Test Download
	downloaded, err := s.downloadFromGCS(ctx, "test/doc.txt")
	if err != nil {
		t.Fatalf("downloadFromGCS a échoué: %v", err)
	}
	if !bytes.Equal(downloaded, testData) {
		t.Errorf("données téléchargées différentes: %q != %q", string(downloaded), string(testData))
	}

	// 3. Test Delete
	if err := s.deleteFromGCS(ctx, "test/doc.txt"); err != nil {
		t.Fatalf("deleteFromGCS a échoué: %v", err)
	}

	// 4. Test Download après Delete doit renvoyer os.ErrNotExist
	_, err = s.downloadFromGCS(ctx, "test/doc.txt")
	if err == nil {
		t.Fatalf("attendu erreur os.ErrNotExist après suppression, obtenu nil")
	}
}

func TestSaveAndInitCorpusSnapshotWithMockServer(t *testing.T) {
	storage := make(map[string][]byte)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			name := r.URL.Query().Get("name")
			body, _ := io.ReadAll(r.Body)
			storage[name] = body
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{"name": name})
		case http.MethodGet:
			path := strings.TrimPrefix(r.URL.Path, "/storage/v1/b/mock-bucket/o/")
			data, ok := storage[path]
			if !ok {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
	}))
	defer ts.Close()

	// 1. Sauvegarder un corpus
	s1 := &ServerState{
		gcsBucket:   "mock-bucket",
		gcsEndpoint: ts.URL,
		documents: []Document{
			newDocument("d1", "Titre SRE Cloud", "gs://mock/sre", "Pratiques SRE avancées sur Kubernetes et GKE.", time.Now()),
			newDocument("d2", "Titre FinOps GCP", "gs://mock/finops", "Gestion des budgets et alertes de facturation.", time.Now()),
		},
	}

	ctx := context.Background()
	if err := s1.saveCorpusSnapshot(ctx); err != nil {
		t.Fatalf("saveCorpusSnapshot a échoué: %v", err)
	}

	// 2. Démarrer une nouvelle instance (s2) et restaurer le corpus depuis GCS
	s2 := &ServerState{
		gcsBucket:   "mock-bucket",
		gcsEndpoint: ts.URL,
	}
	s2.initCorpusFromGCS(ctx)

	s2.mu.RLock()
	docsCount := len(s2.documents)
	s2.mu.RUnlock()

	if docsCount != 2 {
		t.Fatalf("attendu 2 documents restaurés sur la nouvelle instance, obtenu %d", docsCount)
	}

	// 3. Vérifier que la recherche fonctionne parfaitement sur la nouvelle instance
	results := s2.searchDocuments("FinOps")
	if len(results) == 0 {
		t.Fatalf("aucun résultat de recherche sur l'instance restaurée pour 'FinOps'")
	}
	if results[0].DocumentTitle != "Titre FinOps GCP" {
		t.Errorf("résultat inattendu: %q", results[0].DocumentTitle)
	}
}

func TestInitCorpusFromGCSFallbackWhenNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	}))
	defer ts.Close()

	defaultDoc := newDocument("def-1", "Doc Defaut", "gs://b/def", "Contenu par défaut en mémoire.", time.Now())
	s := &ServerState{
		gcsBucket:   "empty-bucket",
		gcsEndpoint: ts.URL,
		documents:   []Document{defaultDoc},
	}

	ctx := context.Background()
	s.initCorpusFromGCS(ctx)

	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.documents) != 1 || s.documents[0].ID != "def-1" {
		t.Fatalf("le corpus par défaut doit être préservé en cas de 404 sur GCS")
	}
}

func TestHandleDocumentsDeleteSyncsGCS(t *testing.T) {
	storage := make(map[string][]byte)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			name := r.URL.Query().Get("name")
			body, _ := io.ReadAll(r.Body)
			storage[name] = body
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			path := strings.TrimPrefix(r.URL.Path, "/storage/v1/b/sync-bucket/o/")
			data, ok := storage[path]
			if !ok {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			w.Write(data)
		}
	}))
	defer ts.Close()

	s := &ServerState{
		gcsBucket:   "sync-bucket",
		gcsEndpoint: ts.URL,
		documents: []Document{
			newDocument("doc-a", "Doc A", "gs://b/a", "Contenu doc A", time.Now()),
			newDocument("doc-b", "Doc B", "gs://b/b", "Contenu doc B", time.Now()),
		},
	}

	// Suppression unitaire
	req := httptest.NewRequest(http.MethodDelete, "/api/documents?id=doc-a", nil)
	w := httptest.NewRecorder()
	s.handleDocuments(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("code = %d, attendu %d", w.Code, http.StatusOK)
	}

	// Vérifier que le snapshot GCS contient maintenant 1 seul document
	snapshotData, ok := storage[corpusSnapshotPath]
	if !ok {
		t.Fatalf("snapshot GCS absent après suppression")
	}
	var snap CorpusSnapshot
	if err := json.Unmarshal(snapshotData, &snap); err != nil {
		t.Fatalf("snapshot JSON invalide: %v", err)
	}
	if len(snap.Documents) != 1 || snap.Documents[0].ID != "doc-b" {
		t.Errorf("snapshot GCS mal synchronisé: %+v", snap.Documents)
	}

	// Suppression totale
	reqAll := httptest.NewRequest(http.MethodDelete, "/api/documents?all=true", nil)
	wAll := httptest.NewRecorder()
	s.handleDocuments(wAll, reqAll)

	if wAll.Code != http.StatusOK {
		t.Fatalf("code = %d, attendu %d", wAll.Code, http.StatusOK)
	}

	snapshotDataAll := storage[corpusSnapshotPath]
	var snapAll CorpusSnapshot
	json.Unmarshal(snapshotDataAll, &snapAll)
	if len(snapAll.Documents) != 0 {
		t.Errorf("snapshot GCS attendu vide après purge, obtenu %d documents", len(snapAll.Documents))
	}
}
