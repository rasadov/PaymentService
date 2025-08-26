# Payment Service Makefile

.PHONY: docs docs-serve docs-generate test build clean

# Generate all documentation
docs-generate:
	@echo "🚀 Generating API documentation..."
	@mkdir -p docs
	@go doc -all ./internal/handler > docs/go-docs.txt
	@chmod +x scripts/generate-docs.sh
	@./scripts/generate-docs.sh

# Serve documentation locally
docs-serve:
	@echo "🌐 Serving documentation at http://localhost:8080"
	@cd docs && python3 -m http.server 8080

# Generate and serve documentation
docs: docs-generate docs-serve

# Run tests
test:
	@go test ./...

# Build the application
build:
	@go build -o bin/payment-service main.go

# Clean build artifacts
clean:
	@rm -rf bin/
	@rm -f docs/go-docs.txt

# Help
help:
	@echo "Available commands:"
	@echo "  docs-generate  - Generate all documentation files"
	@echo "  docs-serve     - Serve documentation on localhost:8080"
	@echo "  docs          - Generate and serve documentation"
	@echo "  test          - Run tests"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build artifacts"
