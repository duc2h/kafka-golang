package main

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/edarha/kafka-golang/internals/configs"
	"github.com/edarha/kafka-golang/internals/must"
	"github.com/edarha/kafka-golang/internals/repositories"
	"github.com/edarha/kafka-golang/internals/services"
	"github.com/gin-gonic/gin"
)

var (
	kafkaBrokers = []string{"localhost:9092"}
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
	if err := must.MigrateDB(db); err != nil {
		log.Fatal("something wrong while migrating database. err: %w", err)
	}

	// init kafka producer
	producer, err := SetupProducer()
	if err != nil {
		log.Fatal("something wrong while setup producer. err: %w", err)

	}

	// init repo
	studentRepo := repositories.NewStudentRepo(db)

	// init service
	studentSvc := services.NewStudent(studentRepo, producer)

	// setup server
	r := gin.Default()
	r.POST("/student", studentSvc.Post)

	r.Run(":8080")
}

// setupProducer will create a AsyncProducer and returns it
func SetupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	return sarama.NewSyncProducer(kafkaBrokers, config)
}
