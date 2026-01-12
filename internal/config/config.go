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
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
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
	return &cfg, nil
}
