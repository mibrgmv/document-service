package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port         string        `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
		IdleTimeout  time.Duration `yaml:"idle_timeout"`
	} `yaml:"server"`

	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string
		Password string
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"postgres"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string
		DB       int `yaml:"db"`
	} `yaml:"redis"`

	JWT struct {
		Secret     string
		Expiration time.Duration `yaml:"expiration"`
	} `yaml:"jwt"`

	Migrations struct {
		Path string `yaml:"path"`
	} `yaml:"migrations"`

	AdminToken string
}

func Load() (*Config, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	data, err := os.ReadFile(filepath.Join(basepath, "config.yaml"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	_ = godotenv.Load()

	cfg.Postgres.User, err = getEnv("POSTGRES_USER")
	cfg.Postgres.Password, err = getEnv("POSTGRES_PASSWORD")
	cfg.Redis.Password, err = getEnv("REDIS_PASSWORD")
	cfg.JWT.Secret, err = getEnv("JWT_SECRET")
	cfg.AdminToken, err = getEnv("ADMIN_TOKEN")

	return cfg, err
}

func getEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("could not get %s from environment", key)
}
