package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Enables file:// driver
	_ "github.com/lib/pq"

	"github.com/yohanc3/resumemaxxer/internal/config"
	"github.com/yohanc3/resumemaxxer/internal/storage/db"
)

// validates DB connection and applies migrations
func main() {
	loadEnvConfig()
	validateDBConnection()
	applyMigrations()
	fmt.Println("DB ready to be used.")
}

func applyMigrations() {

	dbConnectionString := db.GetDBURL(config.Cfg.DBUsername, config.Cfg.DBPassword, config.Cfg.DBDriver, config.Cfg.DBPort)

	migration, err := migrate.New(
		"file://migration_files",
		dbConnectionString,
	)

	if err != nil {
		log.Fatalf("Error when starting migration: %v", err.Error())
	}

	if err := migration.Up(); err != nil {

		// Don't stop execution while we don't have any migration files (which makes
		// this always error out)
		log.Println("Error during migration up: %v", err.Error())
	}

}

func loadEnvConfig() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error when loading env vars. %v", err.Error())
	}

	log.Println("Success loading env vars.")
}

func validateDBConnection() {

	dbConnection, err := db.GetDB(config.Cfg.DBUsername, config.Cfg.DBPassword, config.Cfg.DBDriver, config.Cfg.DBPort)
	if err != nil {
		log.Fatalf("Error when getting DB connection at db setup. %v", err.Error())
	}

	if err := dbConnection.Ping(); err != nil {
		log.Fatalf("Error when pinging db at db setup. %v", err.Error())
	}
}
