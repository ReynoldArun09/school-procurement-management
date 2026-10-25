package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv string
	Port   string
	Mysql  MysqlConfig
}

type MysqlConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("MYSQL_PORT", "8080"),
		Mysql: MysqlConfig{
			Host:            getEnv("MYSQL_HOST", "127.0.0.1"),
			Port:            getEnv("MYSQL_PORT", "3306"),
			User:            strings.TrimSpace(os.Getenv("MYSQL_USER")),
			Password:        os.Getenv("MYSQL_PASSWORD"),
			Database:        strings.TrimSpace(os.Getenv("MYSQL_DATABASE")),
			MaxOpenConns:    25,
			MaxIdleConns:    10,
			ConnMaxLifetime: 5 * time.Minute,
		},
	}

	if cfg.AppEnv != "production" && cfg.AppEnv != "development" && cfg.AppEnv != "test" {
		return nil, fmt.Errorf("Invalid APP_ENV: %q", cfg.AppEnv)
	}

	port, err := strconv.Atoi(cfg.Port)

	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid port: %q", cfg.Port)
	}

	mysqlPort, err := strconv.Atoi(cfg.Mysql.Port)

	if err != nil || mysqlPort < 1 || mysqlPort > 65535 {
		return nil, fmt.Errorf("invalid port: %q", cfg.Mysql.Port)
	}

	if cfg.Mysql.Database == "" {
		return nil, fmt.Errorf("MYSQL_USER is required")
	}

	if cfg.Mysql.User == "" {
		return nil, fmt.Errorf("MYSQL_USER is required")
	}

	if cfg.Mysql.MaxOpenConns < 1 {
		return nil, fmt.Errorf("MaxOpenConns must be greater than zero")
	}

	if cfg.Mysql.MaxIdleConns < 0 ||
		cfg.Mysql.MaxIdleConns > cfg.Mysql.MaxOpenConns {
		return nil, fmt.Errorf(
			"MaxIdleConns must be between zero and MaxOpenConns",
		)
	}

	if cfg.Mysql.MaxOpenConns < 1 {
		return nil, fmt.Errorf("MaxOpenConns must be greater than zero")
	}

	if cfg.Mysql.MaxIdleConns < 0 ||
		cfg.Mysql.MaxIdleConns > cfg.Mysql.MaxOpenConns {
		return nil, fmt.Errorf(
			"MaxIdleConns must be between zero and MaxOpenConns",
		)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		return fallback
	}

	return value
}
