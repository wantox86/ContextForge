# Makefile — ContextForge build targets
# Usage: make [target]

BINARY     := ContextForge
BIN_DIR    := build/bin
WAILS      := wails
GO         := go

.PHONY: all dev mac mac-amd64 linux linux-docker windows windows-nsis windows-docker \
        test test-race vet fmt build-frontend clean help

# ── Default ─────────────────────────────────────────────────────────────────
all: mac

# ── Development ─────────────────────────────────────────────────────────────
## dev: Start hot-reload development server
dev:
	$(WAILS) dev

# ── macOS ────────────────────────────────────────────────────────────────────
## mac: Build for macOS arm64 (Apple Silicon)
mac:
	$(WAILS) build -platform darwin/arm64 -clean -o $(BINARY)

## mac-amd64: Build for macOS amd64 (Intel)
mac-amd64:
	$(WAILS) build -platform darwin/amd64 -clean -o $(BINARY)-intel

# ── Linux ────────────────────────────────────────────────────────────────────
## linux: Build for Linux amd64 (run natively on Linux)
linux:
	$(WAILS) build -platform linux/amd64 -clean -o $(BINARY)

## linux-docker: Cross-compile for Linux amd64 via Docker (from Mac/Windows)
linux-docker:
	bash scripts/build-linux.sh --docker

# ── Windows ──────────────────────────────────────────────────────────────────
## windows: Build for Windows amd64 .exe (run natively on Windows)
windows:
	$(WAILS) build -platform windows/amd64 -clean -o $(BINARY)

## windows-nsis: Build Windows .exe + NSIS installer
windows-nsis:
	$(WAILS) build -platform windows/amd64 -nsis -clean -o $(BINARY)

## windows-docker: Cross-compile for Windows amd64 via Docker
windows-docker:
	bash scripts/build-windows.sh --docker

# ── Quality ──────────────────────────────────────────────────────────────────
## test: Run all Go tests
test:
	$(GO) test ./...

## test-race: Run all Go tests with race detector
test-race:
	$(GO) test -race ./...

## vet: Run go vet
vet:
	$(GO) vet ./...

## fmt: Format all Go code
fmt:
	gofmt -w .

## build-frontend: Build frontend only
build-frontend:
	cd frontend && npm run build

# ── Clean ────────────────────────────────────────────────────────────────────
## clean: Remove built binaries
clean:
	rm -rf $(BIN_DIR)/*

# ── Help ─────────────────────────────────────────────────────────────────────
## help: Show available targets
help:
	@echo ""
	@echo "ContextForge — Available make targets:"
	@echo ""
	@grep -E '^## ' Makefile | sed 's/## /  /'
	@echo ""
