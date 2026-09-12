#!/bin/bash
set -e

echo "================================================="
echo " Stopping PromptSieve containers                 "
echo "================================================="
docker stop promptsieve-frontend promptsieve-backend
echo "PromptSieve containers stopped."
echo "UI :3005 and API :8080 are no longer served by this stack."
