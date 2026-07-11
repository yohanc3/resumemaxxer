package db 

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

//Given a set of required db connection parameters, it returns a reference to a database
// handle if connection is established successfully, or an error otherwise.
func GetDB(username, password, name, port string) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=disable",
		username, password, port, name)
	
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error when getting db. %w", err)
	}

	return db, nil

}

func GetDBURL(username, password, name, port string) (string) {

	return fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=disable",
		username, password, port, name)

}
