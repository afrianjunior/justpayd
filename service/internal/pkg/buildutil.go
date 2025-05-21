package pkg

type Config struct {
	StoragePath string    `json:"storage_path"`
	ServerPort  string    `json:"server_port"`
	LogLevel    string    `json:"log_level"`
	JWT         JWTConfig `json:"jwt"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string `json:"secret"`
	Expiration int    `json:"expiration"` // time in minutes
}
