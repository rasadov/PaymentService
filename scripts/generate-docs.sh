#!/bin/bash

# Generate API documentation for Payment Service

echo "🚀 Generating Payment Service API Documentation..."

# Create docs directory if it doesn't exist
mkdir -p docs

# Generate Go documentation
echo "📝 Generating Go documentation..."
go doc -all ./internal/handler > docs/go-docs.txt

# Check if OpenAPI spec exists
if [ -f "docs/api.yaml" ]; then
    echo "✅ OpenAPI specification found at docs/api.yaml"
else
    echo "❌ OpenAPI specification not found"
fi

# Check if HTML viewer exists
if [ -f "docs/index.html" ]; then
    echo "✅ HTML documentation viewer found at docs/index.html"
else
    echo "❌ HTML documentation viewer not found"
fi

echo ""
echo "📚 Documentation files generated:"
echo "  - docs/API.md           - Markdown documentation"
echo "  - docs/api.yaml         - OpenAPI specification"
echo "  - docs/index.html       - Interactive HTML viewer"
echo "  - docs/go-docs.txt      - Go package documentation"
echo ""
echo "🌐 To view interactive documentation:"
echo "  1. Serve the docs directory with a web server"
echo "  2. Open http://localhost:8080/docs/ in your browser"
echo ""
echo "💡 Quick serve command:"
echo "  cd docs && python3 -m http.server 8080"
echo "  or"
echo "  npx serve docs"
