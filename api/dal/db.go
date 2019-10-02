package dal

import (
	"database/sql"
	"fmt"
	"sync"

	// postgres dialect
	_ "github.com/lib/pq"
)

var once sync.Once

const (
	// DBUSER - db user
	DBUSER = "admin"
	// DBPASSWORD - db password
	DBPASSWORD = "password"
	// DBNAME - db name
	DBNAME = "graphqldb"
)

// Connect - connects to db
func Connect() (*sql.DB, error) {
	var db *sql.DB
	var err error
	once.Do(func() {
		dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DBUSER, DBPASSWORD, DBNAME)
		db, _ = sql.Open("postgres", dbinfo)
		err = db.Ping()
	})
	return db, err
}

// LogAndQuery - log and query
func LogAndQuery(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	fmt.Println(query)
	return db.Query(query, args...)
}

// MustExec - query than must exec
func MustExec(db *sql.DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		panic(err)
	}
}
