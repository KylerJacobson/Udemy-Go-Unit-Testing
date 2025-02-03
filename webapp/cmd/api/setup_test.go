package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

var app application
var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE3MzgxODg0OTQsImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMSJ9.V_zyCNE63WdejK6flEykaZwkvMY5_-YZNQ8I4bP26kQ"

// This will be executed before all the tests
// We can use it to run setup before the tests run
func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}

	app.JWTSecret = "verysecret"
	os.Exit(m.Run())
}
