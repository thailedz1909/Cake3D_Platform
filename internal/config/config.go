package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	AppEnv         string
	MySQLDSN       string
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	JWTSecret      string
	JWTIssuer      string
	AccessTokenTTL time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	mysqlHost := getEnv("MYSQL_HOST", "localhost")
	mysqlPort := getEnv("MYSQL_PORT", "3306")
	mysqlUser := getEnv("MYSQL_USER", "cake3d")
	mysqlPassword := getEnv("MYSQL_PASSWORD", "cake3d")
	mysqlDatabase := getEnv("MYSQL_DATABASE", "cake3d")
	mysqlParams := getEnv("MYSQL_PARAMS", "charset=utf8mb4&parseTime=True&loc=Local")
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase, mysqlParams)

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	ttlMinutes, _ := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_TTL_MINUTES", "60"))

	appEnv := getEnv("APP_ENV", "debug")
	jwtSecret := getEnv("JWT_SECRET", "change-me-in-production")
	if appEnv == "release" && jwtSecret == "change-me-in-production" {
		panic("JWT_SECRET must be set to a strong secret in release environment")
	}

	return Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppEnv:         appEnv,
		MySQLDSN:       mysqlDSN,
		RedisAddr:      fmt.Sprintf("%s:%s", redisHost, redisPort),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        redisDB,
		JWTSecret:      jwtSecret,
		JWTIssuer:      getEnv("JWT_ISSUER", "cake3d-platform"),
		AccessTokenTTL: time.Duration(ttlMinutes) * time.Minute,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
