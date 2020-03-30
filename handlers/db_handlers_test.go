package handlers_test

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/volatiletech/null"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/pytorchtw/go-db-services-starter/handlers"
	"github.com/pytorchtw/go-db-services-starter/models"
	"github.com/pytorchtw/go-db-services-starter/utils"
)

const (
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres123"
	dbName     = "docker"
)

var testHandler struct {
	db         *sql.DB
	schemaName string
	dbHandler  *handlers.DBHandler
	stop       func()
}

func createTestDatabase(t *testing.T) (*sql.DB, string, func()) {
	connectionString := fmt.Sprintf("port=%d user=%s password=%s dbname=%s sslmode=disable", dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		t.Fatalf("Fail to create database %s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	schemaName := "test" + strconv.FormatInt(rand.Int63(), 10)

	_, err = db.Exec("CREATE SCHEMA " + schemaName)
	if err != nil {
		t.Fatalf("Fail to create schema. %s", err.Error())
	}

	// close db and reconnect with the new test schema set
	db.Close()
	connectionString = fmt.Sprintf("port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s", dbPort, dbUser, dbPassword, dbName, schemaName)
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		t.Fatalf("Fail to create database %s", err.Error())
	}

	return db, schemaName, func() {
		_, err := db.Exec("DROP SCHEMA " + schemaName + " CASCADE")
		if err != nil {
			t.Fatalf("Fail to drop schema. %s", err.Error())
		}
	}
}

func Test_setup(t *testing.T) {
	db, schemaName, dbStopFunc := createTestDatabase(t)
	testHandler.db = db
	testHandler.schemaName = schemaName
	testHandler.stop = dbStopFunc
	h, err := handlers.NewDBHandler(testHandler.db)
	if err != nil {
		t.Fatalf("error creating db handler, %s", err.Error())
	}
	testHandler.dbHandler = h

	log.Println("created test database " + schemaName)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatalf("error getting db instance data, %s", err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+utils.Basepath+"/db_data/migrations", "postgres", driver)
	if err != nil {
		t.Fatalf("error getting db instance data, %s", err.Error())
	}

	err = m.Up()
	if err != nil {
		t.Fatalf("error doing migration steps, %s", err.Error())
	}
	log.Println("completed db migration steps")
}

func Test_CreatePage(t *testing.T) {
	page := models.Page{}
	page.URL = "test url"
	page.Content = null.NewString("test content", true)
	err := testHandler.dbHandler.CreatePage(&page)
	if err != nil {
		t.Fatalf("error creating page, %s", err.Error())
	}
	err = testHandler.dbHandler.CreatePage(&page)
	if err == nil {
		t.Fatalf("error should not be nil, should not be able to create the same page again")
	}
	if page.ID != 1 {
		t.Fatal("id should be 1")
	}
	if page.Content.String != "test content" {
		t.Fatal("error content")
	}
}

func Test_GetPage(t *testing.T) {
	url := "test url"
	page, err := testHandler.dbHandler.GetPage(url)
	if err != nil {
		t.Fatalf("error getting page, %s", err.Error())
	}
	if page.URL != url {
		t.Fatalf("error url")
	}
}

func Test_DeletePage(t *testing.T) {
	url := "test url"
	page, err := testHandler.dbHandler.GetPage(url)
	if err != nil {
		t.Fatalf("error getting page, %s", err.Error())
	}
	if page.URL != url {
		t.Fatalf("error url")
	}

	rowsAff, err := testHandler.dbHandler.DeletePage(url)
	if err != nil {
		t.Fatalf("error deleting page, %s", err.Error())
	}
	if rowsAff != 1 {
		t.Fatalf("affected rows should be 1")
	}

	_, err = testHandler.dbHandler.GetPage(url)
	if err == nil {
		t.Fatalf("should return error when getting deleted page")
	}
}

func Test_shutdown(t *testing.T) {
	testHandler.stop()
	log.Println("dropped test database " + testHandler.schemaName)
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
