package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"main-service/config"
)

func main() {
	log.Println("Start migration")
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.PG.Host, cfg.PG.Port, cfg.PG.User, cfg.PG.Password, cfg.PG.DBName)
	db, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		log.Fatalf("Error connectiong: %v\n", err)
	}
	err = goose.Up(db, "./migrations")
	if err != nil {
		log.Fatalf("Error during migration: %v\n", err)
	}
	if err = db.Close(); err != nil {
		log.Fatalf("Error close DB: %v\n", err)
	}

	log.Println("Migration completed successfully")
}
