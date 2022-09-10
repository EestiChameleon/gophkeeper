package cfg

const (
	DatabaseURI   = "postgresql://localhost:5432/yandex_practicum_db?sslmode=disable"
	CryptoKey     = "secret_123456789" // for test proj ok, but sure it's better to pass it via ENV ;)
	ServerAddress = "localhost:3200"
)
