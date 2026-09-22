#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/../../../../src"

echo "🔍 [1/2] Formatting Go files with gofmt..."
gofmt -w .
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
  echo "❌ Unformatted Go files detected: $UNFORMATTED"
  exit 1
fi
echo "✅ gofmt check passed."

echo "🔍 [2/2] Running Go unit test suite..."
go test ./...
echo "✅ All Go unit tests passed."
