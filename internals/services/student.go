package services

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/edarha/kafka-golang/internals/models"
	"github.com/edarha/kafka-golang/internals/repositories"
	kafka_pb "github.com/edarha/kafka-golang/pb/kafka"
	"google.golang.org/protobuf/proto"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	KafkaTopic = "student_create"
)

type student struct {
	studentRepo repositories.Student
	producer    sarama.SyncProducer
}

func NewStudent(studentRepo repositories.Student,
	producer sarama.SyncProducer) *student {
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

		produceMessages(s.producer, entity)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// produceMessages will send 'testing 123' to KafkaTopic each second, until receive a os signal to stop e.g. control + c
// by the user in terminal
func produceMessages(producer sarama.SyncProducer, student models.Student) {

	s := kafka_pb.Student{
		UserId: student.UserId,
		Grade:  int32(student.Grade),
	}

	m, err := proto.Marshal(&s)
	if err != nil {
		log.Fatalln("marshal error: ", err.Error())
	}

	message := &sarama.ProducerMessage{Topic: KafkaTopic, Value: sarama.ByteEncoder(m)}

	p, o, err := producer.SendMessage(message)

	if err != nil {
		fmt.Println("Err publish: ", err.Error())
	}

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)

}
