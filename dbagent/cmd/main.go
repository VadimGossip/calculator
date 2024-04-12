package main

import (
	"github.com/VadimGossip/calculator/dbagent/internal/app"
	"time"
)

var configDir = "api/config"

func main() {
	dbAgent := app.NewApp("Calculator db agent", configDir, time.Now())
	dbAgent.Run()
}
