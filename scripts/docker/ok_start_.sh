#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
APP_DIR="$( cd "$DIR/../.." && pwd )"
cd "$APP_DIR"

echo "================================================="
echo " Building PromptSieve images                     "
echo "================================================="
docker build -t promptsieve-backend:latest -f backend/Dockerfile backend
docker build -t promptsieve-frontend:latest -f frontend/Dockerfile frontend

echo "================================================="
echo " Starting containers with docker run             "
echo "================================================="
if docker container inspect promptsieve-backend >/dev/null 2>&1; then
  docker start promptsieve-backend
else
  docker run -d --name promptsieve-backend -p 8080:8080 promptsieve-backend:latest
fi

if docker container inspect promptsieve-frontend >/dev/null 2>&1; then
  docker start promptsieve-frontend
else
  docker run -d --name promptsieve-frontend -p 3005:3000 \
    -e NEXT_PUBLIC_API_URL=http://localhost:8080 \
    promptsieve-frontend:latest
fi

echo "PromptSieve API: http://127.0.0.1:8080"
echo "PromptSieve UI:  http://localhost:3005"
