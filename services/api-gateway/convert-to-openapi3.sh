#!/bin/bash

# Script to convert Swagger 2.0 to OpenAPI 3.0
# Usage: ./convert-to-openapi3.sh

set -e

echo "🔄 Converting Swagger 2.0 to OpenAPI 3.0..."

# Change to docs directory
cd "$(dirname "$0")/docs"

# Check if swagger files exist
if [ ! -f "swagger.yaml" ] || [ ! -f "swagger.json" ]; then
    echo "❌ Error: swagger.yaml or swagger.json not found"
    echo "Please run 'swag init' first to generate Swagger documentation"
    exit 1
fi

# Check if swagger2openapi is installed
if ! command -v swagger2openapi &> /dev/null; then
    echo "📦 Installing swagger2openapi..."
    npm install -g swagger2openapi
fi

# Convert YAML
echo "📄 Converting swagger.yaml to openapi.yaml..."
swagger2openapi swagger.yaml -o openapi.yaml

# Convert JSON
echo "📄 Converting swagger.json to openapi.json..."
swagger2openapi swagger.json -o openapi.json

echo "✅ Conversion complete!"
echo ""
echo "Generated files:"
echo "  - docs/openapi.yaml"
echo "  - docs/openapi.json"
echo ""
echo "Access OpenAPI 3.0 documentation at:"
echo "  - http://localhost:8080/openapi.yaml"
echo "  - http://localhost:8080/openapi.json"
echo "  - http://localhost:8080/api/v1/openapi.yaml"
echo "  - http://localhost:8080/api/v1/openapi.json"
