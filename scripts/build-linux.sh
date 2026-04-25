#!/usr/bin/env bash
# build-linux.sh — Build ContextForge for Linux (amd64)
# Requires: Go 1.22+, Wails v2, Node.js 18+, gcc
# On Linux: run natively. On Mac: use Docker (see --docker flag).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$PROJECT_DIR/build/bin"
USE_DOCKER=false

for arg in "$@"; do
  case $arg in
    --docker) USE_DOCKER=true ;;
  esac
done

cd "$PROJECT_DIR"

if [ "$USE_DOCKER" = true ]; then
  echo "▶ Building via Docker (cross-compile Mac → Linux amd64)…"
  docker run --rm \
    -v "$(pwd)":/app \
    -w /app \
    -e HOME=/tmp \
    wailsapp/wails:latest \
    wails build -platform linux/amd64 -o ContextForge-linux
else
  echo "▶ Building ContextForge for Linux (amd64) natively…"
  wails build -platform linux/amd64 -clean -o ContextForge
fi

echo ""
echo "✅ Linux build ready in: $BIN_DIR"
ls -lh "$BIN_DIR"
