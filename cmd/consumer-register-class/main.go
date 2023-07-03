package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/edarha/kafka-golang/internals/configs"
	"github.com/edarha/kafka-golang/internals/must"
	"github.com/edarha/kafka-golang/internals/repositories"
	"github.com/edarha/kafka-golang/internals/services/consumers"

	"github.com/Shopify/sarama"
)

// Sarama configuration options
var (
	brokers = "localhost:29092,localhost:29093,localhost:29094"
	group   = "group_consumer_class_student"
	topics  = "class_student_register"
)

func main() {
	keepRunning := true
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	ctx, cancel := context.WithCancel(context.Background())

	// init db
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

	// init repo
	classStudentRepo := repositories.NewClassStudentRepo(db)

	// init consumer
	classStudentConsumer := consumers.NewClassStudentConsumer(classStudentRepo)

	reader := consumers.NewReader(config, brokers, group, topics, classStudentConsumer.Handler)

	go func() {

		for {
			if err := reader.Consume(ctx); err != nil {
				return
			}
		}
	}()

	log.Println("server is running...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}
	}
	cancel()

}
