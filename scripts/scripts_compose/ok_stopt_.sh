#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
APP_DIR="$( cd "$DIR/../.." && pwd )"
cd "$APP_DIR"

echo "================================================="
echo " Stopping PromptSieve                            "
echo "================================================="
docker-compose stop
echo "PromptSieve containers stopped."
echo "UI :3005 and API :8080 are no longer served by this stack."
