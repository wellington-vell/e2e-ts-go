package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Env(key string) string {
	envPath := filepath.Join("..", "..", ".env")
	_ = godotenv.Load(envPath)

	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Environment variable %s is not set", key))
	}

	return value
}
