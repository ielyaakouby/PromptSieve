#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
APP_DIR="$( cd "$DIR/../.." && pwd )"
cd "$APP_DIR"

echo "================================================="
echo " Restarting PromptSieve containers               "
echo "================================================="
docker stop promptsieve-frontend promptsieve-backend 2>/dev/null || true
docker rm promptsieve-frontend promptsieve-backend 2>/dev/null || true

docker run -d --name promptsieve-backend -p 8080:8080 promptsieve-backend:latest
docker run -d --name promptsieve-frontend -p 3005:3000 \
  -e NEXT_PUBLIC_API_URL=http://localhost:8080 \
  promptsieve-frontend:latest

echo "PromptSieve API: http://127.0.0.1:8080"
echo "PromptSieve UI:  http://localhost:3005"
