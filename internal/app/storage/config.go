package storage

type StorageConfig struct {
	URL string `json:"database_url"`
}

func NewStorageConfig() *StorageConfig {
	return &StorageConfig{}
}
