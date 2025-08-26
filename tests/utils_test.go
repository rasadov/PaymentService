package tests

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/rasadov/PaymentService/internal/config"
	"github.com/rasadov/PaymentService/pkg"
)

func TestWebhookSecretVerification(t *testing.T) {
	testPayload := []byte("test payload")

	h := hmac.New(sha256.New, []byte(config.GetConfig().DodoWebhookSecret))

	h.Write(testPayload)

	testSignature := hex.EncodeToString(h.Sum(nil))

	if !pkg.VerifyWebhookSignature(testSignature, testPayload) {
		t.Errorf("Webhook secret verification failed")
	}
}

func TestWebhookSecretVerificationFailed(t *testing.T) {
	testPayload := []byte("test payload")

	h := hmac.New(sha256.New, []byte("wrong secret"))

	h.Write(testPayload)

	testSignature := hex.EncodeToString(h.Sum(nil))

	if pkg.VerifyWebhookSignature(testSignature, testPayload) {
		t.Errorf("Webhook secret verification failed")
	}
}
