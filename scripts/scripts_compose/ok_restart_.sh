#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
APP_DIR="$( cd "$DIR/../.." && pwd )"
cd "$APP_DIR"

echo "================================================="
echo " Restarting PromptSieve                          "
echo "================================================="
docker-compose up -d --force-recreate
echo "PromptSieve API: http://127.0.0.1:8080"
echo "PromptSieve UI:  http://localhost:3005"
