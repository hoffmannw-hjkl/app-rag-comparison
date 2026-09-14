# Multi-stage Dockerfile pour un conteneur minimal et sans shell

# Le binaire étant statique et sans dépendance externe, la stdlib compilée ici
# constitue toute la surface d'attaque de l'image finale : la version du toolchain
# est donc un paramètre de sécurité, pas un simple détail de build. On suit la
# ligne stable courante, qui embarque les correctifs net/http, crypto/tls et
# crypto/x509 exigés par go.mod.
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Les dépendances sont copiées séparément pour bénéficier du cache de couches.
# Pas de `|| true` ici : une résolution de dépendance en échec doit faire échouer
# le build plutôt que produire une image silencieusement incomplète.
COPY go.mod go.su[m] ./
RUN go mod download

COPY . .

# Les tests s'exécutent dans l'image de build : une régression bloque la livraison.
RUN CGO_ENABLED=0 go test ./...

# Binaire statique, sans table de symboles ni informations de débogage.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o rag-chatbot .

# Image finale distroless : ni shell ni gestionnaire de paquets, donc surface
# d'attaque très réduite, et exécution sous un utilisateur non privilégié.
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /app/rag-chatbot .

USER nonroot:nonroot

EXPOSE 8080
ENV PORT=8080

ENTRYPOINT ["/app/rag-chatbot"]
