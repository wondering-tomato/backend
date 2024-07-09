package mysql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type StoreMySQL interface{}
type StoreMySQLImpl struct {
	db *sql.DB
}

var _ StoreMySQL = &StoreMySQLImpl{}

var port = 3306

func New() *StoreMySQLImpl {
	db, err := connect()
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("using store: mysql")
	return &StoreMySQLImpl{
		db: db,
	}
}

func connect() (*sql.DB, error) {
	// Get a database handle.
	db, err := sql.Open("mysql", fmt.Sprintf("exploreuser:test@tcp(db:%d)/explore", port))

	// ## Uncomment the line below to connect to the DB using a localhost name.
	// db, err := sql.Open("mysql", fmt.Sprintf("root:example@tcp(0.0.0.0:%d)/explore", port))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("mysql connected on %d\n", port)
	return db, err
}
