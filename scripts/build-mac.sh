#!/usr/bin/env bash
# build-mac.sh — Build ContextForge for macOS (arm64 + amd64 universal)
# Requires: Go 1.22+, Wails v2, Node.js 18+
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$PROJECT_DIR/build/bin"

cd "$PROJECT_DIR"

echo "▶ Building ContextForge for macOS (arm64)…"
wails build -platform darwin/arm64 -clean -o ContextForge

echo "▶ Building ContextForge for macOS (amd64)…"
wails build -platform darwin/amd64 -o ContextForge-amd64

echo ""
echo "✅ macOS builds ready in: $BIN_DIR"
ls -lh "$BIN_DIR"
