package mercrawl

import (
	"database/sql"
	"os"
)

var connStr = "user=" + os.Getenv("USER") + " dbname=mercrawl" + " sslmode=" + os.Getenv("SSLMODE")
var db *sql.DB

// GetDb returns an instance of Postgresql database connection.
func GetDb() *sql.DB {
	if db != nil {
		return db
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
