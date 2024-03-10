package main

import (
	"forum/internal/app"
	"os"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.json"
	}
	app.Run(cfgPath)
}
