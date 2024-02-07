package main

import (
	"github.com/VadimGossip/calculator/agent/internal/app"
	"time"
)

var configDir = "agent/config"

func main() {
	calculatorAgent := app.NewApp("Calculator Agent", configDir, time.Now())
	calculatorAgent.Run()
}
