package db_test

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	db "github/adedaryorh/pooler_Remmitance_Application/db/sqlc"
	"github/adedaryorh/pooler_Remmitance_Application/utils"
	"log"
	"os"
	"testing"
)

var testQuery *db.Store

const testDbName = "testdb"
const sslmode = "?sslmode=disable"

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cant load env config", err)
	}
	connection, err := sql.Open(config.DBdriver, config.DB_source+sslmode)
	if err != nil {
		log.Fatalf("could not connect to %s serer %v", config.DBdriver, err)
	}

	_, err = connection.Exec(fmt.Sprintf("CREATE DATABASE %s;", testDbName))
	if err != nil {
		log.Fatalf("ENCOUNTERED CREATING DB %v", err)
	}

	tconn, err := sql.Open(config.DBdriver, config.DB_source+testDbName+sslmode)
	if err != nil {
		teardown(connection)
		log.Fatalf("Could not conect to DB %v", err)
	}

	driver, err := postgres.WithInstance(tconn, &postgres.Config{})
	if err != nil {
		teardown(connection)
		log.Fatalf("Could not create migrate driver%v", err)
	}

	mig, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "../migrations"),
		config.DBdriver, driver)

	if err != nil {
		teardown(connection)
		log.Fatalf("migration setup failed %v", err)
	}

	if err = mig.Up(); err != nil && err != migrate.ErrNoChange {
		teardown(connection)
		log.Fatalf("migration up failed %v", err)
	}

	testQuery = db.NewStore(tconn)
	code := m.Run()
	tconn.Close()

	teardown(connection)
	//exit the process when u r done running test
	os.Exit(code)
}

func teardown(connection *sql.DB) {
	_, err := connection.Exec(fmt.Sprintf("DROP DATABASE %s WITH (FORCE);", testDbName))

	if err != nil {
		teardown(connection)
		log.Fatalf("failed to drop db %v", err)
	}

	connection.Close()
}
