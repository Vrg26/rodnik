package main

import (
	"log"
	"rodnik/config"
	"rodnik/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config errors: %s", err)
	}

	app.Run(cfg)
}
