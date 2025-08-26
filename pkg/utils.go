package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/rasadov/PaymentService/internal/config"
)

func VerifyWebhookSignature(signature string, payload []byte) bool {
	h := hmac.New(sha256.New, []byte(config.GetConfig().DodoWebhookSecret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
