package main

import (
	"log"

	"github.com/edarha/kafka-golang/internals/configs"
	"github.com/edarha/kafka-golang/internals/must"
)

func main() {
	cfg := configs.PostgreSQL{
		Host:     "localhost",
		Port:     "5432",
		Database: "postgres",
		Username: "postgres",
		Password: "admin",
	}

	db := must.ConnectPostgresql(&cfg)

	// migrate database
	err := must.MigrateDB(db)
	if err != nil {
		log.Fatal("something wrong while migrating database. err: %w", err)
	}

}
