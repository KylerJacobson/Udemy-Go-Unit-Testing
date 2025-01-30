package main

import (
	"os"
	"testing"
)

var app application

// This will be executed before all the tests
// We can use it to run setup before the tests run
func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"
	app.Session = getSession()
	os.Exit(m.Run())
}
