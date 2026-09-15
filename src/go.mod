module github.com/cloud-gtm/gcp-ai-foundation-blueprint/examples/04-rag-comparison-chatbot

// L'application n'a aucune dépendance externe : Vertex AI est appelé en REST via
// net/http. Sa seule surface de vulnérabilité est donc la bibliothèque standard,
// et le plancher ci-dessous est le premier patch corrigeant les 31 failles
// remontées par govulncheck sur la ligne 1.22 (net/http, crypto/tls, crypto/x509,
// net/url, encoding/asn1). Le relever est l'unique correctif possible.
go 1.25.13
