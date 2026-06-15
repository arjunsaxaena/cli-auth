package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	SessionDuration  time.Duration
	LockDuration     time.Duration
	MaxFailedAttempt int
}

func Load() *Config {

	sessionMinutes := getEnvInt(
		"SESSION_TIMEOUT_MINUTES",
		30,
	)

	lockMinutes := getEnvInt(
		"LOCK_DURATION_MINUTES",
		15,
	)

	maxAttempts := getEnvInt(
		"MAX_FAILED_ATTEMPTS",
		5,
	)

	return &Config{
		SessionDuration:  time.Duration(sessionMinutes) * time.Minute,
		LockDuration:     time.Duration(lockMinutes) * time.Minute,
		MaxFailedAttempt: maxAttempts,
	}
}

func getEnvInt(
	key string,
	defaultValue int,
) int {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return n
}

var App = Load()
