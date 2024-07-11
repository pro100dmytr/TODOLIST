package postgresql

import (
	"database/sql"
	"fmt"
	"log"
)

var Db *sql.DB

func ConnectDB() {
	connStr := "host=localhost port=4050 user=myuser dbname=todolistdb password=4050 sslmode=disable"
	var err error
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Успішно підключено до бази даних!")
}
