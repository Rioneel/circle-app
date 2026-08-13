package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CognoURI string
	CognoUser string
	CognoPassword string
	Port string
}

func Load() (*Config, error) {
	_ = godotenv.Load() 

	cfg := &Config{
		CognoURI: os.Getenv("COGNODB_URI"),
		CognoUser:os.Getenv("COGNODB_USER"),
		CognoPassword:os.Getenv("COGNODB_PASSWORD"),
		Port: os.Getenv("PORT"),
	}

	if cfg.CognoURI == "" || cfg.CognoUser == "" || cfg.CognoPassword == "" {
		return nil, fmt.Errorf("missing required env vars: COGNODB_URI, COGNODB_USER, COGNODB_PASSWORD")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	return cfg, nil
}	