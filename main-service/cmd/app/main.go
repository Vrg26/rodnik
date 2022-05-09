package main

import (
	_ "github.com/lib/pq"
	"log"
	"main-service/config"
	"main-service/internal/app"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Config errors: %s", err)
	}
	app.Run(cfg)
}
