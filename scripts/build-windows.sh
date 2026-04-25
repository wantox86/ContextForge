#!/usr/bin/env bash
# build-windows.sh — Build ContextForge for Windows (amd64)
# Requires: Go 1.22+, Wails v2, Node.js 18+
# On Windows: run in Git Bash or PowerShell (use build-windows.ps1 instead).
# On Linux/Mac with mingw: cross-compile with --cross flag.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$PROJECT_DIR/build/bin"
USE_NSIS=false
USE_DOCKER=false

for arg in "$@"; do
  case $arg in
    --nsis)   USE_NSIS=true ;;
    --docker) USE_DOCKER=true ;;
  esac
done

cd "$PROJECT_DIR"

if [ "$USE_DOCKER" = true ]; then
  echo "▶ Building via Docker (cross-compile → Windows amd64)…"
  NSIS_FLAG=""
  [ "$USE_NSIS" = true ] && NSIS_FLAG="-nsis"
  docker run --rm \
    -v "$(pwd)":/app \
    -w /app \
    -e HOME=/tmp \
    wailsapp/wails:latest \
    wails build -platform windows/amd64 $NSIS_FLAG -o ContextForge
else
  echo "▶ Building ContextForge for Windows (amd64) natively…"
  NSIS_FLAG=""
  [ "$USE_NSIS" = true ] && NSIS_FLAG="-nsis"
  wails build -platform windows/amd64 $NSIS_FLAG -clean -o ContextForge
fi

echo ""
echo "✅ Windows build ready in: $BIN_DIR"
ls -lh "$BIN_DIR"
