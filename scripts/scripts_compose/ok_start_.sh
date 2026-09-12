#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
APP_DIR="$( cd "$DIR/../.." && pwd )"
cd "$APP_DIR"

if [ "${1:-}" = "--local" ]; then
  echo "================================================="
  echo " Starting PromptSieve locally                    "
  echo "================================================="
  cd "$APP_DIR/backend"
  go fmt ./...
  go vet ./...
  cd "$APP_DIR/frontend"
  npm run lint
  cd "$APP_DIR"
  echo "Starting PromptSieve API on :8080 and web UI on :3005"
  (cd "$APP_DIR/backend" && go run ./cmd/promptsieve-server) &
  API_PID=$!
  trap 'kill "$API_PID" 2>/dev/null || true' EXIT
  cd "$APP_DIR/frontend"
  npm run dev -- -p 3005
  exit 0
fi

echo "================================================="
echo " 1. Linting Backend (Go)                         "
echo "================================================="
cd "$APP_DIR/backend"
go fmt ./...
go vet ./...
cd "$APP_DIR"

echo "================================================="
echo " 2. Linting Frontend (Next.js)                   "
echo "================================================="
cd "$APP_DIR/frontend"
npm run lint
cd "$APP_DIR"

echo "================================================="
echo " 3. Building Docker Images                       "
echo "================================================="
echo "Building Backend Image (promptsieve-backend:latest)..."
docker build -t promptsieve-backend:latest -f backend/Dockerfile backend

echo "Building Frontend Image (promptsieve-frontend:latest)..."
docker build -t promptsieve-frontend:latest -f frontend/Dockerfile frontend

echo "================================================="
echo " 4. Starting Containers via Docker Compose       "
echo "================================================="
docker-compose up -d
echo "PromptSieve API: http://127.0.0.1:8080"
echo "PromptSieve UI:  http://localhost:3005"
