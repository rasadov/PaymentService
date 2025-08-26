package db

type Storage interface {
	Get(key string) (string, error)
	Put(key string, value string) error
	Delete(key string) error
}

func GetConnection() (Storage, error) {
	return &kvStorage{}, nil
}
