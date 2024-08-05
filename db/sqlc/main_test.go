package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/RenanWinter/bank/util/config"
)

var testQueries *Queries
var testDB *sql.DB
var storeDB Store

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	storeDB = NewStore(testDB)

	os.Exit(m.Run())
}
