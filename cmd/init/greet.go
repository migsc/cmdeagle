package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Get arguments from environment variables
	name := getEnv("ARGS_NAME")
	age := getEnv("ARGS_AGE")

	// Get flags from environment variables
	uppercase := getEnv("FLAGS_UPPERCASE") == "true"
	lowercase := getEnv("FLAGS_LOWERCASE") == "true"
	repeat, _ := strconv.Atoi(getEnv("FLAGS_REPEAT"))

	// Construct base greeting
	greeting := fmt.Sprintf("Hello %s!", name)
	if age != "" {
		greeting += fmt.Sprintf(" You are %s years old.", age)
	}

	// Apply case transformations
	if uppercase {
		greeting = strings.ToUpper(greeting)
	} else if lowercase {
		greeting = strings.ToLower(greeting)
	}

	// Output greeting with repetition
	for i := 0; i < repeat; i++ {
		fmt.Println(greeting)
	}
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}
