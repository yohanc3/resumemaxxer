package main

import (
	"database/sql"
	"errors"
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

	dbConnection, err := db.GetDB(config.Cfg.DBUsername, config.Cfg.DBPassword, config.Cfg.DBName, config.Cfg.DBPort)
	if err != nil {
		log.Fatalf("Error when getting DB connection at db setup. %v", err.Error())
	}

	validateDBConnection(dbConnection)
	applyMigrations()

	log.Println("DB ready to be used.")
}

func applyMigrations() {

	dbConnectionString := db.GetDBURL(config.Cfg.DBUsername, config.Cfg.DBPassword, config.Cfg.DBName, config.Cfg.DBPort)

	migration, err := migrate.New(
		"file://migration_files",
		dbConnectionString,
	)

	if err != nil {
		log.Fatalf("Error when starting migration: %v", err.Error())
	}

	if err := migration.Up(); err != nil {
		
		if errors.Is(err, migrate.ErrNoChange){
			log.Println("No new migrations to apply. Schema is up to date.")
			return
		} 

		log.Fatalf("Error during migration up: %v", err.Error())

	}

	log.Println("Succesfully ran migration up.")

}

func loadEnvConfig() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error when loading env vars. %v", err.Error())
	}

	log.Println("Succesfully loaded env vars.")
}

func validateDBConnection(dbConnection *sql.DB) {

	if err := dbConnection.Ping(); err != nil {
		log.Fatalf("Error when pinging db at db setup. %v", err.Error())
	} else {
		log.Println("Succesfully pinged db at setup.")
	}
}
