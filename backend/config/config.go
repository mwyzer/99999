package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName  string
	AppEnv   string
	AppPort  string
	AppURL   string

	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBTimezone        string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime int

	JWTSecret     string
	JWTExpiryHours int

	BcryptCost int

	CORSAllowedOrigins string

	UploadDir          string
	MaxUploadSizeMB    int64
	MaxPhotosPerListing int

	RateLimitLoginPerMinute  int
	RateLimitGlobalPerMinute int

	FrontendURL  string
	PlatformName string
}

func Load() *Config {
	// Load .env if exists (ignore error — env vars might come from system)
	_ = godotenv.Load()

	cfg := &Config{
		AppName:  getEnv("APP_NAME", "PropertyHub"),
		AppEnv:   getEnv("APP_ENV", "development"),
		AppPort:  getEnv("APP_PORT", "8080"),
		AppURL:   getEnv("APP_URL", "http://localhost:8080"),

		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "propertyhub"),
		DBPassword:        getEnv("DB_PASSWORD", "propertyhub_secret"),
		DBName:            getEnv("DB_NAME", "propertyhub"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DBTimezone:        getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		DBMaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 50),
		DBMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 30),

		JWTSecret:      getEnv("JWT_SECRET", "default-secret-change-me"),
		JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 24),

		BcryptCost: getEnvInt("BCRYPT_COST", 12),

		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),

		UploadDir:           getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSizeMB:     int64(getEnvInt("MAX_UPLOAD_SIZE_MB", 5)),
		MaxPhotosPerListing: getEnvInt("MAX_PHOTOS_PER_LISTING", 10),

		RateLimitLoginPerMinute:  getEnvInt("RATE_LIMIT_LOGIN_PER_MINUTE", 5),
		RateLimitGlobalPerMinute: getEnvInt("RATE_LIMIT_GLOBAL_PER_MINUTE", 100),

		FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:5173"),
		PlatformName: getEnv("PLATFORM_NAME", "PropertyHub"),
	}

	return cfg
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode, c.DBTimezone,
	)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		n, err := strconv.Atoi(val)
		if err == nil {
			return n
		}
	}
	return fallback
}
