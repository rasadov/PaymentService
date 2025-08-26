# Payment Service Documentation

This directory contains comprehensive API documentation for the Payment Service.

## 📁 Documentation Files

| File | Description | Usage |
|------|-------------|-------|
| `API.md` | Complete API documentation in Markdown | Read directly or convert to other formats |
| `api.yaml` | OpenAPI 3.0 specification | Import into API tools (Postman, Insomnia) |
| `index.html` | Interactive Swagger UI viewer | Open in browser for interactive docs |
| `go-docs.txt` | Generated Go package documentation | Reference for Go developers |

## 🚀 Quick Start

### View Interactive Documentation
```bash
# Serve documentation locally
make docs-serve
# Then open http://localhost:8080 in your browser
```

### Generate Documentation
```bash
# Generate all documentation files
make docs-generate
```

### Import into API Tools
- **Postman**: Import `api.yaml` as OpenAPI specification
- **Insomnia**: Import `api.yaml` as OpenAPI specification
- **VS Code**: Use REST Client extension with examples from `API.md`

## 🔧 Available Commands

```bash
make docs-generate  # Generate all documentation
make docs-serve     # Serve docs on localhost:8080
make docs          # Generate and serve (combined)
```

## 📋 API Endpoints Summary

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/checkout` | POST | Create checkout session |
| `/subscriptions` | POST | Get subscription management link |
| `/webhook` | POST | Handle payment webhooks |

## 🔗 External Resources

- [Dodo Payments API Documentation](https://docs.dodopayments.com/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Go Documentation](https://golang.org/doc/)

## 📝 Updating Documentation

When you modify the API:

1. Update the OpenAPI spec in `api.yaml`
2. Update the Markdown docs in `API.md`
3. Add Go doc comments to your code
4. Run `make docs-generate` to refresh generated docs

## 🎯 Best Practices

- **Keep docs in sync** with code changes
- **Use examples** in all documentation
- **Include error cases** and status codes
- **Document authentication** requirements
- **Provide curl examples** for testing
