package config

import (
	"os"
	"strings"
)

type Config struct {
	Port           string
	RulesPath      string
	AllowedOrigins []string
	Env            string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rulesPath := os.Getenv("RULES_PATH")
	if rulesPath == "" {
		// Default assumes the server is run from the backend/ module directory,
		// where the rule library lives one level up at ../rules.
		rulesPath = "../rules"
	}

	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	if originsEnv == "" {
		originsEnv = "http://localhost:3000"
	}
	origins := strings.Split(originsEnv, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	return Config{
		Port:           port,
		RulesPath:      rulesPath,
		AllowedOrigins: origins,
		Env:            env,
	}
}
