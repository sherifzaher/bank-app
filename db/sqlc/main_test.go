package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sherifzaher/clone-simplebank/util"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot open env:", err)
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	defer testDb.Close()
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDb)
	fmt.Println("Connected to db successfully!")
	os.Exit(m.Run())
}
