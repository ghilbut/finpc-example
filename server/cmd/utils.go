package main

import (
	"os"
	"path/filepath"

	// external packages
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func getEnvValue(name, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}
	return v
}

func loadEnvs(basepath, name string) {
	envpath := filepath.Join(basepath, name)
	if _, err := os.Stat(envpath); err != nil {
		return
	}
	if err := godotenv.Load(envpath); err != nil {
		log.Fatalf("Error loading %s file", name)
	}
}
