package config

import (
	"fmt"
	"os"
)

const (
	JwtSigningKeyEnvVar = "JWT_SIGNING_KEY"
	DbUserEnvVar        = "DB_USER"
	DbPasswordEnvVar    = "DB_PASSWORD"
	DbPortEnvVar        = "DB_PORT"
	DbAddressEnvVar     = "DB_ADDRESS"
)

type Config struct {
	JwtSigningKey string
	DbUser        string
	DbPassword    string
	DbPort        string
	DbAddress     string
}

func ConfigFromEnv() (*Config, error) {
	jwtSigningKey := os.Getenv(JwtSigningKeyEnvVar)
	if jwtSigningKey == "" {
		return nil, err(JwtSigningKeyEnvVar)
	}
	dbUser := os.Getenv(DbUserEnvVar)
	if dbUser == "" {
		return nil, err(DbUserEnvVar)
	}
	dbPassword := os.Getenv(DbPasswordEnvVar)
	if dbPassword == "" {
		return nil, err(DbPasswordEnvVar)
	}
	dbPort := os.Getenv(DbPortEnvVar)
	if dbPort == "" {
		return nil, err(DbPortEnvVar)
	}
	dbAddress := os.Getenv(DbAddressEnvVar)
	if dbAddress == "" {
		return nil, err(DbAddressEnvVar)
	}
	return &Config{
		JwtSigningKey: jwtSigningKey,
		DbUser:        dbUser,
		DbPassword:    dbPassword,
		DbPort:        dbPort,
		DbAddress:     dbAddress,
	}, nil
}

func err(envVar string) error {
	return fmt.Errorf("%s must be set as an ENV var", envVar)
}
