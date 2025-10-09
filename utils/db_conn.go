package utils

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

func ConnectDb(db_url string) (*sql.DB, error ) {
	db_conn , err := sql.Open("postgres",db_url)
	if err != nil {
		return nil, errors.New("Can't connect")
	}
	
    return db_conn , nil
}