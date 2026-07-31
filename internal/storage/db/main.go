package db 

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/yohanc3/resumemaxxer/internal/config"
)

//Given a set of required db connection parameters, it returns a reference to a database
// handle if connection is established successfully, or an error otherwise.
func GetDB() (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Cfg.DBUsername, config.Cfg.DBPassword, config.Cfg.DBHost, config.Cfg.DBPort, config.Cfg.DBName)
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error when getting db. %w", err)
	}

	return db, nil

}

func GetDBURL() string {

	return fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=disable",
		config.Cfg.DBUsername, config.Cfg.DBPassword, config.Cfg.DBHost, config.Cfg.DBPort, config.Cfg.DBName)

}
