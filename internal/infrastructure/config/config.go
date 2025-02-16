package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	devMode     bool
	pgUrl       string
	httpPort    string
	jwtSecret   string
	rdbAddr     string
	rdbPassword string
	rdbDb       int
}

func MustNewConfigWithEnv() *Config {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("error loading .env file: %s", err))
	}
	return &Config{
		devMode:     mustGetEnvAsBool("DEV_MODE"),
		pgUrl:       mustGetEnv("PG_URL"),
		httpPort:    mustGetEnv("HTTP_PORT"),
		jwtSecret:   mustGetEnv("JWT_SECRET"),
		rdbAddr:     mustGetEnv("RDB_ADDR"),
		rdbPassword: mustGetEnv("RDB_PASSWORD"),
		rdbDb:       mustGetEnvAsInt("RDB_DB"),
	}
}

func (c *Config) PgUrl() string {
	return c.pgUrl
}
func (c *Config) HttpPort() string {
	return c.httpPort
}
func (c *Config) JwtSecret() string {
	return c.jwtSecret
}
func (c *Config) RdbAddr() string {
	return c.rdbAddr
}
func (c *Config) RdbPassword() string {
	return c.rdbPassword
}
func (c *Config) RdbDb() int {
	return c.rdbDb
}
func (c *Config) DevMode() bool {
	return c.devMode
}
func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}
	return value
}
func mustGetEnvAsBool(key string) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}
	normalized := strings.TrimSpace(strings.ToLower(value))

	if normalized == "true" || normalized == "1" {
		return true
	} else if normalized == "false" || normalized == "0" {
		return false
	}
	panic(fmt.Sprintf("invalid boolean value for environment variable %s: %s", key, value))
}

func mustGetEnvAsInt(key string) int {
	valueStr := mustGetEnv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s must be an integer, but got '%s'", key, valueStr))
	}
	return value
}
