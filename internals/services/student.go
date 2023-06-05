package services

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/edarha/kafka-golang/internals/models"
	"github.com/edarha/kafka-golang/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	kafkaBrokers = []string{"localhost:9092"}
	KafkaTopic   = "student_create"
	enqueued     int
)

type student struct {
	studentRepo repositories.Student
	producer    sarama.AsyncProducer
}

func NewStudent(studentRepo repositories.Student,
	producer sarama.AsyncProducer) *student {
	return &student{
		studentRepo: studentRepo,
		producer:    producer,
	}
}

func (s *student) Post(c *gin.Context) {
	var student struct {
		UserId string `json:"user_id" biding:"required"`
		Grade  int16  `json:"grade" biding:"required"`
	}

	// publish list student to kafka

	if c.Bind(&student) == nil {
		entity := models.Student{
			ID:     uuid.NewString(),
			UserId: student.UserId,
			Grade:  student.Grade,
		}

		err := s.studentRepo.Create(context.Background(), &entity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// produceMessages will send 'testing 123' to KafkaTopic each second, until receive a os signal to stop e.g. control + c
// by the user in terminal
func produceMessages(producer sarama.AsyncProducer, signals chan os.Signal) {
	for {
		time.Sleep(time.Second)
		message := &sarama.ProducerMessage{Topic: KafkaTopic, Value: sarama.StringEncoder("testing 123")}
		select {
		case producer.Input() <- message:
			enqueued++
			log.Println("New Message produced")
		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			return
		}
	}
}
