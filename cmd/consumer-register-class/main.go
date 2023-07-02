package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/edarha/kafka-golang/internals/services/consumers"

	"github.com/Shopify/sarama"
)

// Todo. implement consume logic.

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
	classStudentConsumer := consumers.NewClassStudentConsumer()

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
