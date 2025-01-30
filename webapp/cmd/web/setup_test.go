package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

var app application

// This will be executed before all the tests
// We can use it to run setup before the tests run
func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"
	app.Session = getSession()
	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run())
}
