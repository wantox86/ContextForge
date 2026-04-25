# ContextForge

<p align="center">
  <img src="build/appicon.png" width="120" alt="ContextForge logo" />
</p>

<p align="center">
  <strong>Generate AI coding instruction files from visual templates — instantly.</strong><br/>
  Auto-detects your tech stack and exports to CLAUDE.md, copilot-instructions.md, and more.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/version-0.1.0-indigo" />
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go" />
  <img src="https://img.shields.io/badge/Wails-v2-orange" />
  <img src="https://img.shields.io/badge/Vue-3-42b883?logo=vue.js" />
  <img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue" />
</p>

---

## What is ContextForge?

ContextForge is an open-source desktop app for developers who use AI coding assistants.
It helps you generate well-structured context/instruction files from curated templates:

| Output file | Used by |
|---|---|
| `CLAUDE.md` | Claude (Anthropic) |
| `.github/copilot-instructions.md` | GitHub Copilot |
| `.cursorrules` *(v0.2.0)* | Cursor |
| `AGENTS.md` *(v0.2.0)* | OpenAI Codex / Agents |

**Key features:**
- 🔍 **Auto stack detection** — scans your project directory, detects Go, Kotlin, Python, Vue, React, Android, Rust, Flutter, and more
- 📋 **Visual rule builder** — enable/disable and edit rules per template
- 👁 **Live preview** — see the exact output before exporting
- ⚡ **Token estimator** — warns you before hitting model context limits
- 🗃 **5 built-in templates** — Spring Boot, Go/Gin, FastAPI, Vue 3, Android Kotlin

---

## Screenshots

> *Coming soon — run `make dev` to see the UI locally.*

---

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.22+ | https://go.dev/dl |
| Node.js | 18+ | https://nodejs.org |
| Wails CLI | v2 | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| gcc / build tools | latest | see platform notes below |

### Platform build dependencies

**macOS**
```bash
xcode-select --install
```

**Linux (Ubuntu/Debian)**
```bash
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.0-dev gcc pkg-config
```

**Windows**
- Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MSYS2](https://www.msys2.org/)
- Optional: install [NSIS](https://nsis.sourceforge.io/) for the installer builder

---

## Quick Start

```bash
# Clone
git clone https://github.com/wantox86/ContextForge.git
cd ContextForge

# Install Go deps + frontend deps
go mod download
cd frontend && npm install && cd ..

# Run in development mode (hot reload)
make dev
# or: wails dev
```

---

## Building

### Using Make (recommended)

```bash
make help          # list all targets

make mac           # macOS arm64 (Apple Silicon)
make mac-amd64     # macOS amd64 (Intel)
make linux         # Linux amd64  (run on Linux)
make windows       # Windows amd64 .exe  (run on Windows)
make windows-nsis  # Windows .exe + NSIS installer
```

### Using build scripts

```bash
# macOS
bash scripts/build-mac.sh

# Linux — native (run on a Linux machine)
bash scripts/build-linux.sh

# Linux — cross-compile from Mac via Docker
bash scripts/build-linux.sh --docker

# Windows — native (run in Git Bash on Windows)
bash scripts/build-windows.sh

# Windows — with NSIS installer
bash scripts/build-windows.sh --nsis

# Windows — PowerShell
.\scripts\build-windows.ps1 -Nsis
```

### Manual Wails build

```bash
wails build -platform darwin/arm64   -o ContextForge
wails build -platform linux/amd64    -o ContextForge
wails build -platform windows/amd64  -o ContextForge
wails build -platform windows/amd64  -nsis -o ContextForge  # + installer
```

Build output is placed in `build/bin/`.

---

## Releases (CI/CD)

Push a version tag to trigger automated builds on all platforms via GitHub Actions:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The [release workflow](.github/workflows/release.yml) builds on native runners (macOS, Ubuntu, Windows) and attaches the following artifacts to a GitHub Release:

| Artifact | Platform |
|---|---|
| `ContextForge-darwin-arm64.zip` | macOS Apple Silicon |
| `ContextForge-darwin-amd64.zip` | macOS Intel |
| `ContextForge-linux-amd64.tar.gz` | Linux |
| `ContextForge-windows-amd64.exe` | Windows portable |
| `ContextForge-amd64-installer.exe` | Windows NSIS installer |

---

## Architecture

```
ContextForge/
├── main.go                        # Wails entry point
├── app.go                         # Go → Vue bindings (all exposed methods)
├── internal/
│   ├── analyzer/                  # Auto-detect tech stack from project files
│   ├── exporter/                  # Render rules → CLAUDE.md / copilot-instructions.md
│   ├── template/                  # Built-in template JSON + DB seeder
│   ├── tokenizer/                 # Token count estimate + model limit warnings
│   └── storage/                   # SQLite via GORM (models + migrations + CRUD)
├── frontend/
│   ├── src/
│   │   ├── components/            # Vue 3 components (Composition API)
│   │   │   ├── ProjectAnalyzer.vue
│   │   │   ├── RuleBuilder.vue
│   │   │   ├── PreviewPanel.vue
│   │   │   └── ExportPanel.vue
│   │   ├── stores/                # Pinia stores (project, template, export)
│   │   └── types/                 # TypeScript interfaces
│   └── wailsjs/                   # Auto-generated Wails JS bindings
├── assets/templates/              # Built-in template JSON source files
├── build/                         # Wails platform build configs + icons
├── scripts/                       # Build scripts (bash + PowerShell)
└── .github/workflows/             # CI + Release GitHub Actions
```

---

## Built-in Templates

| Template | Stack Tags | Rules |
|---|---|---|
| Spring Boot + Kotlin | `kotlin`, `spring-boot`, `gradle` | 6 |
| Go + Gin | `golang` | 6 |
| Python + FastAPI | `python`, `fastapi` | 6 |
| Vue 3 + TypeScript | `vue`, `nodejs` | 6 |
| Android + Kotlin | `android`, `kotlin` | 6 |

Templates are embedded into the binary and seeded to the local SQLite database on first launch.
Custom templates can be added via the Rule Builder UI.

---

## Stack Detection

ContextForge scans indicator files in your project root:

| File / Folder | Detected Stack |
|---|---|
| `go.mod` | `golang` |
| `pom.xml` | `java`, `maven` |
| `build.gradle.kts` | `kotlin`, `gradle` |
| `build.gradle` | `java`, `gradle` |
| `package.json` | `nodejs` |
| `requirements.txt` / `pyproject.toml` | `python` |
| `Cargo.toml` | `rust` |
| `pubspec.yaml` | `flutter` |
| `android/` folder | `android` |

Secondary scan inspects file contents for framework hints (e.g. `spring-boot-starter` in `pom.xml` → adds `spring-boot` tag).

---

## Token Counter

ContextForge estimates token usage (1 token ≈ 4 characters) and warns you before hitting model limits:

| Model | Warn At | Max |
|---|---|---|
| claude-sonnet | 8,000 | 200,000 |
| github-copilot | 6,000 | 8,000 |
| cursor | 6,000 | 8,000 |
| ollama-8k | 6,000 | 8,192 |
| ollama-32k | 24,000 | 32,768 |

---

## Development

```bash
# Run all Go tests
make test

# Run with race detector
make test-race

# Run go vet
make vet

# Format code
make fmt

# Build frontend only
make build-frontend
```

---

## Roadmap

### v0.1.0 (current)
- [x] Project analyzer with auto stack detection
- [x] 5 built-in templates
- [x] Rule enable/disable + inline edit
- [x] Export to CLAUDE.md and copilot-instructions.md
- [x] Token count estimate with per-model warning
- [x] Mac (.app) + Linux + Windows (.exe) builds

### v0.2.0 (planned)
- [ ] Cursor (`.cursorrules`) export
- [ ] AGENTS.md export
- [ ] Folder picker UI (native dialog)
- [ ] Custom template sharing (import/export JSON)
- [ ] AI-assisted rule suggestion

---

## Contributing

1. Fork the repo
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Commit your changes (no TODO comments — open an issue instead)
4. Push and open a Pull Request

Please follow the conventions in [copilot-instructions.md](copilot-instructions.md).

---

## License

MIT — see [LICENSE](LICENSE) for details.

---

<p align="center">Built with ❤️ using <a href="https://wails.io">Wails</a>, <a href="https://vuejs.org">Vue 3</a>, and <a href="https://go.dev">Go</a></p>

