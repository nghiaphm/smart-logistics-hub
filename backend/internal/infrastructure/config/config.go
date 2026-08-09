package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ProjectName       string
	Version           string
	Environment       string
	ServerHost        string
	ServerPort        string
	MetricsEnabled    bool
	MetricsHost       string
	MetricsPort       string
	S3Enabled         bool
	AWSRegion         string
	S3Bucket          string
	S3Endpoint        string
	S3UsePathStyle    bool
	RedisEnabled      bool
	RedisURI          string
	RedisHost         string
	RedisPort         string
	RedisPassword     string
	RedisDB           int
	RedisCacheTTL     int
	MariaDBEnabled    bool
	MariaDBHost       string
	MariaDBPort       string
	MariaDBUser       string
	MariaDBPassword   string
	MariaDBName       string
	MariaDBURI        string
	FrontendURL       string
	AIServiceURL      string
	KeycloakServerURL string
	KeycloakRealm     string
	KeycloakClientID  string
	DevSkipAuth       bool
}

func LoadConfig() *Config {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	envFile := ".env." + appEnv
	_ = godotenv.Load(envFile)

	return &Config{
		ProjectName:       getEnv("PROJECT_NAME", "Smart Logistics Hub API"),
		Version:           getEnv("VERSION", "1.0.0"),
		Environment:       appEnv,
		ServerHost:        getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:        getEnv("SERVER_PORT", "8000"),
		MetricsEnabled:    getEnvBool("METRICS_ENABLED", true),
		MetricsHost:       getEnv("METRICS_HOST", "0.0.0.0"),
		MetricsPort:       getEnv("METRICS_PORT", "9090"),
		S3Enabled:         getEnvBool("S3_ENABLED", false),
		AWSRegion:         getEnv("AWS_REGION", "ap-southeast-1"),
		S3Bucket:          getEnv("S3_BUCKET", ""),
		S3Endpoint:        getEnv("S3_ENDPOINT", ""),
		S3UsePathStyle:    getEnvBool("S3_USE_PATH_STYLE", false),
		RedisEnabled:      getEnvBool("REDIS_ENABLED", false),
		RedisURI:          getEnv("REDIS_URI", ""),
		RedisHost:         getEnv("REDIS_HOST", "localhost"),
		RedisPort:         getEnv("REDIS_PORT", "6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           getEnvInt("REDIS_DB", 0),
		RedisCacheTTL:     getEnvInt("REDIS_TTL_CACHE", 300),
		MariaDBEnabled:    getEnvBool("MARIADB_ENABLED", true),
		MariaDBHost:       getEnv("MARIADB_HOST", "localhost"),
		MariaDBPort:       getEnv("MARIADB_PORT", "3306"),
		MariaDBUser:       getEnv("MARIADB_USER", "root"),
		MariaDBPassword:   getEnv("MARIADB_PASSWORD", ""),
		MariaDBName:       getEnv("MARIADB_DB_NAME", "smart_logistics"),
		MariaDBURI:        getEnv("MARIADB_URI", ""),
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:3000"),
		AIServiceURL:      getEnv("AI_SERVICE_URL", "http://localhost:5000"),
		KeycloakServerURL: getEnv("KEYCLOAK_SERVER_URL", "http://localhost:8180"),
		KeycloakRealm:     getEnv("KEYCLOAK_REALM", "web-app-project"),
		KeycloakClientID:  getEnv("KEYCLOAK_CLIENT_ID", "fastapi_backend_client"),
		DevSkipAuth:       getEnv("DEV_SKIP_AUTH", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
