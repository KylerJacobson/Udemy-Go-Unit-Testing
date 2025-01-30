package main

import (
	"log"
	"os"
	"testing"
	"webapp/pkg/db"
)

var app application

// This will be executed before all the tests
// We can use it to run setup before the tests run
func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"
	app.Session = getSession()
	app.DSN = "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5"
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = db.PostgresConn{DB: conn}

	os.Exit(m.Run())
}
