package main

import (
	"github.com/VadimGossip/calculator/api/internal/app"
	"time"
)

var configDir = "api/config"

func main() {
	calculatorApi := app.NewApp("Calculator Api", configDir, time.Now())
	calculatorApi.Run()
}
