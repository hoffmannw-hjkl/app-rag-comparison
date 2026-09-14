# Multi-stage Dockerfile pour un conteneur ultra-léger et rapide (<25MB)
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Téléchargement des dépendances
COPY go.mod ./
# Note: go.sum si présent
RUN go mod download || true

COPY . .

# Compilation du binaire statique Go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o rag-chatbot main.go

# Image finale de production minimale
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/rag-chatbot .

EXPOSE 8080
ENV PORT=8080

CMD ["./rag-chatbot"]
