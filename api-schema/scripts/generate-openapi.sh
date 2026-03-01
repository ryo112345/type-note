#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_ROOT=$(cd -- "${SCRIPT_DIR}/.." && pwd)
OPENAPI_FILE="${PROJECT_ROOT}/generated/openapi.yaml"
TSP_OUTPUT="${PROJECT_ROOT}/tsp-output/@typespec/openapi3/openapi.yaml"

echo "📄 Generating OpenAPI YAML from TypeSpec..."

npx tsp compile typespec

if [ ! -f "${TSP_OUTPUT}" ]; then
  echo "❌ Failed to locate TypeSpec output at ${TSP_OUTPUT}"
  exit 1
fi

mkdir -p "${PROJECT_ROOT}/generated"
cp "${TSP_OUTPUT}" "${OPENAPI_FILE}"

echo "✅ OpenAPI YAML generated at: ${OPENAPI_FILE}"
