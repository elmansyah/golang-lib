package godotenv

import (
	"errors"
)

var (
	errEmptyAppMode = errors.New("APP_MODE environment variable is not set")
	errLoadEnv      = errors.New("failed to load .env file")
)

type Params struct{}

type App interface {
	Load() (string, error)
}
