package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/sherifzaher/clone-simplebank/api"
	"github.com/sherifzaher/clone-simplebank/db/sqlc"
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

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("error during run the server %v", err)
		return
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		panic(err)
	}
}
