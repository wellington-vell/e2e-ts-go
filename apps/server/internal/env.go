package internal

import (
	"fmt"
	"os"
)

func Env(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Environment variable %s is not set", key))
	}

	return value
}
