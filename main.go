package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"

	"github.com/sherifzaher/clone-simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB:", err)
		return
	}

	if err = conn.Ping(); err != nil {
		log.Fatal("DB is not live: ", err)
	}
}
