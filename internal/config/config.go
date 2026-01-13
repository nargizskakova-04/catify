package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	App        *AppCfg     `json:"app"`
	Repository *Repository `json:"repository"`
}

type AppCfg struct {
	Port int `json:"port"`
}

type Repository struct {
	DBHost      string `json:"db_host"`
	DBPort      int    `json:"db_port"`
	DBUsername  string `json:"db_username"`
	DBPassword  string `json:"db_password"`
	DBName      string `json:"db_name"`
	DBSSLMode   string `json:"db_ssl_mode"`
	MaxConn     int    `json:"max_conn"`
	MaxIdleConn int    `json:"max_idle_conn"`
}

func GetConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err = json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Repository.DBHost = host
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Repository.DBUsername = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Repository.DBPassword = password
	}

	return &cfg, nil
}
