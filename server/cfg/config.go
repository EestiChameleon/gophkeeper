package cfg

const (
	DatabaseURI   = "postgresql://localhost:5432/yandex_practicum_db?sslmode=disable" // local DB.
	CryptoKey     = "secret_123456789"                                                // for test proj ok, but sure it's better to pass it via ENV ;)
	ServerAddress = "localhost:3200"                                                  // also could be passed via ENV, but ok for test proj.
)
