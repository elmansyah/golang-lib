package godotenv

import (
	"fmt"
	"os"
	
	"github.com/joho/godotenv"
)

type Init func() App

func New() App {
	return &Params{}
}

func (*Params) Load() (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		return "", fmt.Errorf("%w: %w", errLoadEnv, err)
	}
	
	appMode := os.Getenv("APP_MODE")
	
	if appMode == "" {
		return "", fmt.Errorf("%w", errEmptyAppMode)
	}
	
	return appMode, nil
}
