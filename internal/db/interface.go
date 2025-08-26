package db

import "time"

type Storage interface {
	Get(key string) (string, error)
	PutWithExpiration(key string, value string, expiration time.Duration) error
}

func GetConnection() (Storage, error) {
	return &kvStorage{}, nil
}
