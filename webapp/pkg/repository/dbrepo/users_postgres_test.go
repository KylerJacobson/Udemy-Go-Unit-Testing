//go:build integration

package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	"webapp/pkg/data"
	"webapp/pkg/repository"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "users_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepo repository.DatabaseRepo

func TestMain(m *testing.M) {
	// connect to docker; fail if docker not running
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	pool = p

	// set up our docker options, specifying the image and so forth
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// get a resource (docker image)
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until it's ready
	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	// populate the database with empty tables
	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	testRepo = &PostgresDBRepo{DB: testDB}

	// run tests
	code := m.Run()

	// clean up
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping database")
	}
}

func TestPostgresDBRepoInsertUser(t *testing.T) {
	testUser := data.User{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
		Password:  "secret",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := testRepo.InsertUser(testUser)
	if err != nil {
		t.Errorf("insert user returned an error: %s", err)
	}

	if id != 1 {
		t.Errorf("insert user returned wrong id; expected 1, but got %d", id)
	}
}

func TestPostgresDBRepoAllUsers(t *testing.T) {
	users, err := testRepo.AllUsers()
	if err != nil {
		t.Errorf("All users reports and error: %s", err)
	}
	if len(users) != 1 {
		t.Errorf("all users reports wrong size; expected 1 but got %d", len(users))
	}
	testUser := data.User{
		FirstName: "Jack",
		LastName:  "Smith",
		Email:     "Jack@example.com",
		Password:  "secret",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, _ = testRepo.InsertUser(testUser)
	users, err = testRepo.AllUsers()
	if err != nil {
		t.Errorf("All users reports and error: %s", err)
	}
	if len(users) != 2 {
		t.Errorf("all users reports wrong size after insert; expected 2 but got %d", len(users))
	}
}

func TestPostgresDBRepoGetUser(t *testing.T) {
	user, err := testRepo.GetUser(1)
	if err != nil {
		t.Errorf("Error getting user by ID: %s", err)

	}
	if user.Email != "admin@example.com" {
		t.Errorf("wrong email returned by getUser; Expcted admin@example.com but got %s", user.Email)
	}

	_, err = testRepo.GetUser(3)
	if err == nil {
		t.Errorf("no error reported when getting non-existent user")
	}

}

func TestPostgresDBRepoGetUserByEmail(t *testing.T) {
	user, err := testRepo.GetUserByEmail("Jack@example.com")
	if err != nil {
		t.Errorf("Error getting user by email: %s", err)

	}
	if user.ID != 2 {
		t.Errorf("wrong email returned by getUserByEmail; Expcted user 2 but got %d", user.ID)
	}
}

func TestPostgresDBRepoUpdateUser(t *testing.T) {
	user, _ := testRepo.GetUser(2)

	user.FirstName = "Jane"
	user.Email = "jane@example.com"

	err := testRepo.UpdateUser(*user)
	if err != nil {
		t.Errorf("Error updating user %d: %s", 2, err)
	}

	user, _ = testRepo.GetUser(2)
	if user.FirstName != "Jane" || user.Email != "jane@example.com" {
		t.Errorf("expected updated record to have first name jane and email jane@example.com but got %s %s", user.FirstName, user.Email)
	}
}

func TestPostgresDBRepoDeleteUser(t *testing.T) {
	err := testRepo.DeleteUser(2)
	if err != nil {
		t.Errorf("Error deleteing user 2 from database: %s", err)
	}

	_, err = testRepo.GetUser(2)
	if err == nil {
		t.Error("Retrieved user id 2 who should have been deleted")
	}
}

func TestPostgresDBRepoResetPassword(t *testing.T) {
	err := testRepo.ResetPassword(1, "test")
	if err != nil {
		t.Error("Error updating user's password: ", err)
	}

	user, _ := testRepo.GetUser(1)

	matches, err := user.PasswordMatches("test")
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Error("user password does not match the changed password")
	}
}

func TestPostgresDBRepoInsertUserImage(t *testing.T) {
	id, err := testRepo.InsertUserImage(data.UserImage{1, 1, "test.jpg", time.Now(), time.Now()})
	if err != nil {
		t.Error("Error inserting image: ", err)
	}
	if id != 1 {
		t.Errorf("Error inserting image; expected id 1 but got %d", id)
	}

	_, err = testRepo.InsertUserImage(data.UserImage{1, 2, "test.jpg", time.Now(), time.Now()})
	if err == nil {
		t.Errorf("Exepcted error inserting image with userID 2; which should not exist")
	}
}
