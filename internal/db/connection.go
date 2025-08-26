//go:build js && wasm

package db

import (
	"time"

	"github.com/rasadov/PaymentService/internal/config"
	"github.com/syumai/workers/cloudflare/kv"
)

type kvStorage struct {
	namespace *kv.Namespace
}

func (k *kvStorage) Get(key string) (string, error) {
	if k.namespace == nil {
		namespace, err := kv.NewNamespace(config.GetConfig().KVNamespace)
		if err != nil {
			return "", err
		}
		k.namespace = namespace
	}
	return k.namespace.GetString(key, nil)
}

func (k *kvStorage) PutWithExpiration(key string, value string, expiration time.Duration) error {
	if k.namespace == nil {
		namespace, err := kv.NewNamespace(config.GetConfig().KVNamespace)
		if err != nil {
			return err
		}
		k.namespace = namespace
	}
	return k.namespace.PutString(key, value, &kv.PutOptions{ExpirationTTL: int(expiration.Seconds())})
}
