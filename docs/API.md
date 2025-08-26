# Payment Service API Documentation

A lightweight payment service built for Cloudflare Workers that handles checkout sessions, subscription management, and webhooks.

## Base URL

- **Production**: `https://your-domain.workers.dev`
- **Development**: `http://localhost:8080`

## Authentication

The API uses webhook signature verification for the `/webhook` endpoint. All other endpoints are public but should be secured in production.

## Endpoints

### 1. Create Checkout Session

Creates a new checkout session for a customer to purchase a product.

**Endpoint**: `POST /checkout`

**Request Body**:
```json
{
  "email": "customer@example.com",
  "name": "John Doe",
  "product_id": "prod_123456"
}
```

**Response**:
```json
{
  "url": "https://checkout.dodopayments.com/session_abc123"
}
```

**Status Codes**:
- `200` - Success
- `400` - Invalid request body
- `405` - Method not allowed

---

### 2. Get Subscription Management Link

Retrieves a link for customers to manage their subscription.

**Endpoint**: `POST /subscriptions`

**Request Body**:
```json
{
  "customer_id": "cus_123456789"
}
```

**Response**:
```json
{
  "url": "https://billing.dodopayments.com/manage/cus_123456789"
}
```

**Status Codes**:
- `200` - Success
- `400` - Invalid request body
- `405` - Method not allowed
- `500` - Failed to get subscription management link

---

### 3. Handle Payment Webhook

Processes webhook events from Dodo Payments.

**Endpoint**: `POST /webhook`

**Headers**:
- `webhook-signature` (required) - Webhook signature for verification
- `webhook-id` (required) - Unique webhook identifier

**Request Body**:
```json
{
  "type": "subscription.created",
  "data": {
    "subscription_id": "sub_123456789",
    "status": "active",
    "product_id": "prod_123456",
    "customer": {
      "customer_id": "cus_123456789",
      "email": "customer@example.com",
      "name": "John Doe"
    }
  }
}
```

**Response**:
```json
{
  "message": "Webhook processed"
}
```

**Status Codes**:
- `200` - Webhook processed successfully
- `400` - Invalid payload or missing webhook-id
- `401` - Invalid webhook signature
- `405` - Method not allowed
- `500` - Failed to process webhook

## Webhook Event Types

The service handles the following webhook event types:

- `subscription.created` - New subscription created
- `subscription.updated` - Subscription updated
- `subscription.cancelled` - Subscription cancelled
- `payment.succeeded` - Payment processed successfully
- `payment.failed` - Payment failed

## Data Models

### Customer
```json
{
  "customer_id": "string",
  "email": "string",
  "name": "string"
}
```

### Subscription Statuses
- `active` - Subscription is active
- `cancelled` - Subscription has been cancelled
- `past_due` - Payment is overdue
- `unpaid` - Payment failed

## Error Handling

All endpoints return appropriate HTTP status codes and error messages:

- **400 Bad Request** - Invalid request body or missing required fields
- **401 Unauthorized** - Invalid webhook signature
- **405 Method Not Allowed** - Incorrect HTTP method
- **500 Internal Server Error** - Server-side processing error

## Example Usage

### Creating a Checkout Session

```bash
curl -X POST https://your-domain.workers.dev/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "name": "John Doe",
    "product_id": "prod_123456"
  }'
```

### Getting Subscription Management Link

```bash
curl -X POST https://your-domain.workers.dev/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cus_123456789"
  }'
```

### Testing Webhook (with proper headers)

```bash
curl -X POST https://your-domain.workers.dev/webhook \
  -H "Content-Type: application/json" \
  -H "webhook-signature: your_signature_here" \
  -H "webhook-id: webhook_123" \
  -d '{
    "type": "subscription.created",
    "data": {
      "subscription_id": "sub_123456789",
      "status": "active",
      "product_id": "prod_123456",
      "customer": {
        "customer_id": "cus_123456789",
        "email": "customer@example.com",
        "name": "John Doe"
      }
    }
  }'
```

## Integration Notes

1. **Webhook Verification**: Always verify webhook signatures to ensure authenticity
2. **Idempotency**: Webhook endpoints should handle duplicate events gracefully
3. **Error Handling**: Implement proper error handling for all API calls
4. **Rate Limiting**: Consider implementing rate limiting for production use
5. **CORS**: Configure CORS headers if accessing from web browsers

## Development

To view the interactive API documentation:

1. Install a Swagger/OpenAPI viewer
2. Open `docs/api.yaml` in the viewer
3. Or use the HTML viewer at `docs/index.html`
