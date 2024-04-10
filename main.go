// main.go

package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func getenvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	a := App{}
	a.Initialize(
		getenvOrDefault("APP_DB_USERNAME", "postgres"),
		getenvOrDefault("APP_DB_PASSWORD", ""),
		getenvOrDefault("APP_DB_NAME", "postgres"))

	//List env variables (bash): printenv
	a.Run(getenvOrDefault("APP_PORT", ":8010"))
}
