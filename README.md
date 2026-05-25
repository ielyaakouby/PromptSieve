# PromptSieve

**Protect your data before you prompt.**

PromptSieve is an open-source security and privacy tool that scans AI prompts for sensitive information before they are used with LLMs. It is designed to detect credentials, secrets, private cryptographic material, and similar high-risk data before you send a prompt to ChatGPT, Claude, Gemini, Copilot, Cursor, or another AI system.

PromptSieve scans prompts for sensitive information before they are sent to an AI model.

## Features

- Self-hosted, privacy-first scanning with no prompt retention
- Extensible detection engine with built-in AWS access key, JWT, private key, and password detectors
- Risk scoring with `SAFE`, `REVIEW`, `SENSITIVE`, and `CRITICAL` statuses
- Deterministic, overlap-safe sanitization
- Web UI for paste → scan → review → copy
- HTTP API (`POST /api/v1/scan`)
- CLI for scanning files

## Architecture

- **Frontend:** Next.js (App Router) + React + Tailwind CSS
- **Backend:** Go HTTP API and CLI
- **Scan path:** the browser sends prompt text to your PromptSieve API. Detection runs on the server. The API is stateless and does not persist prompts.

```
Browser UI  ──POST /api/v1/scan──►  Go API  ──►  Detectors → Risk → Sanitizer
```

## Quick Start

Run the stack with Docker Compose:

```bash
./scripts/scripts_compose/ok_start_.sh
./scripts/scripts_compose/ok_restart_.sh
./scripts/scripts_compose/ok_stopt_.sh
```

Or run the Go API and Next.js app on the host:

```bash
./scripts/scripts_compose/ok_start_.sh --local
```

Open the UI at http://localhost:3005. The API listens on http://127.0.0.1:8080.

## Development

```bash
# API
cd promptsieve/backend
go run ./cmd/promptsieve-server

# Web UI
cd promptsieve/frontend
npm install
npm run dev
```

Makefile targets from `promptsieve/`:

```bash
make docker build backend
make docker build frontend
make docker build all

make deploy backend
make deploy frontend
make deploy all

make destroy backend
make destroy frontend
make destroy all

make dev-backend
make dev-frontend
make build-api
make build-cli
make test
```

## Docker

Build the images and start the stack with Compose:

```bash
docker build -t promptsieve-backend:latest -f backend/Dockerfile backend
docker build -t promptsieve-frontend:latest -f frontend/Dockerfile frontend
docker-compose up -d
```

### Run containers directly

You can start the same images with `docker run` instead of Compose. Host ports match the current stack: API `8080`, UI `3005`.

```bash
./scripts/docker/ok_start_.sh
./scripts/docker/ok_restart_.sh
./scripts/docker/ok_stopt_.sh
```

Or run the commands by hand:

```bash
docker build -t promptsieve-backend:latest -f backend/Dockerfile backend
docker build -t promptsieve-frontend:latest -f frontend/Dockerfile frontend

docker run -d --name promptsieve-backend -p 8080:8080 promptsieve-backend:latest
docker run -d --name promptsieve-frontend -p 3005:3000 \
  -e NEXT_PUBLIC_API_URL=http://localhost:8080 \
  promptsieve-frontend:latest
```

`NEXT_PUBLIC_API_URL` is inlined when the Next.js image is built. For localhost the frontend default `http://127.0.0.1:8080` already works when both ports are published.

Stop and remove those containers:

```bash
docker stop promptsieve-frontend promptsieve-backend
docker rm promptsieve-frontend promptsieve-backend
```

## Configuration

Copy `.env.example` when running the frontend outside Compose:

```bash
# Preferred
NEXT_PUBLIC_API_URL=http://127.0.0.1:8080
```

The frontend reads `NEXT_PUBLIC_API_URL` and falls back to `http://127.0.0.1:8080`.

## Detection Rules

Current built-in detectors:

| Detector | Category | Severity |
| --- | --- | --- |
| Email, phone, DOB, address, contextual person name | Personal information | Medium |
| Passwords and sensitive assignments | Credentials | High |
| OpenAI-style API keys | API keys | Critical |
| GitHub tokens | Source control tokens | Critical |
| AWS access key IDs and contextual secret keys | Cloud credentials | Critical |
| JWTs and Bearer tokens | Authentication tokens | High |
| Private keys | Cryptographic material | Critical |
| IBAN (MOD-97), BIC/SWIFT | Financial information | High |
| Payment cards (Luhn), contextual CVV and expiry | Financial information | High / Critical |

The engine is extensible: additional detectors can be registered with the scanner without changing the API contract.

## CLI

```bash
promptsieve scan path/to/file.txt
```

Exit codes:

- `0` SAFE
- `1` REVIEW
- `2` SENSITIVE
- `3` CRITICAL
- `4` usage or argument error

## Security

- Detected secrets are masked in API responses (`maskedValue`). Raw matches are not serialized to JSON.
- Sanitization replaces findings in place so a safe copy can be used with an AI tool.
- Do not treat a clean scan as proof that a prompt is free of every possible secret.

## Privacy

PromptSieve is intended to be self-hosted. The API does not store prompts. Scanning happens on the PromptSieve server you run, not inside the browser. When you run it on localhost, prompt text stays on your machine except for the local HTTP request from the UI to the API.

## Contributing

Improvements to detectors, UX, and documentation are welcome. Keep detection changes conservative and avoid logging raw prompt content.

## License

No project license file is included yet. Add one before publishing the repository if you intend to distribute PromptSieve.
