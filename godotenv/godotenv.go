package godotenv

import (
	"os"
	
	"github.com/joho/godotenv"
)

func New() App {
	return &Params{}
}

type Init func() App

func (*Params) Load() string {
	if err := godotenv.Load(".env"); err != nil {
		panic("failed to load .env: " + err.Error())
	}
	
	appMode := os.Getenv("APP_MODE")
	
	if appMode == "" {
		panic("APP MODE environment variable is not set")
	}
	
	return appMode
}
