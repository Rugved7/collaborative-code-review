package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	ServiceName string
	ServicePort string

	DatabaseURL       string
	DBMaxConnections  int
	DBMaxIdleConns    int
	DBConnMaxLifeTime time.Duration

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	AWSRegion          string
	AWSS3Bucket        string
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	AuthServicePort      string
	ReviewServicePort    string
	CommentServicePort   string
	WebSocketServicePort string
	AnalysisServicePort  string
}

var AppConfig *Config

func LoadConfig(ServiceName string) *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	jwtAccessExpiry, _ := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	jwtRefreshExpiry, _ := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
	dbMaxLife, _ := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "1h"))

	AppConfig = &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		ServiceName: ServiceName,
		ServicePort: getEnv(ServiceName+"_PORT", "8080"),

		DatabaseURL:       getEnv("DATABASE_URL", ""),
		DBMaxConnections:  getEnvAsInt("DB_MAX_CONNECTIONS", 25),
		DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifeTime: dbMaxLife,

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTAccessExpiry:  jwtAccessExpiry,
		JWTRefreshExpiry: jwtRefreshExpiry,

		AWSRegion:          getEnv("AWS_REGION", "ap-south-1"),
		AWSS3Bucket:        getEnv("AWS_S3_BUCKET", ""),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),

		AuthServicePort:      getEnv("AUTH_SERVICE_PORT", "8001"),
		ReviewServicePort:    getEnv("REVIEW_SERVICE_PORT", "8002"),
		CommentServicePort:   getEnv("COMMENT_SERVICE_PORT", "8003"),
		WebSocketServicePort: getEnv("WEBSOCKET_SERVICE_PORT", "8004"),
		AnalysisServicePort:  getEnv("ANALYSIS_SERVICE_PORT", "8005"),
	}

	if AppConfig.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if AppConfig.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
