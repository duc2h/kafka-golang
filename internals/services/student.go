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
	student_create = "student_create"
	student_update = "student_update"
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
		Name  string `json:"name" biding:"required"`
		Grade int16  `json:"grade" biding:"required"`
	}

	if c.Bind(&student) == nil {
		entity := models.Student{
			ID:    uuid.NewString(),
			Name:  student.Name,
			Grade: student.Grade,
		}

		err := s.studentRepo.Create(context.Background(), &entity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		produceMessagesCreate(s.producer, entity)

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func (s *student) Put(c *gin.Context) {
	var student struct {
		Name  string `json:"name" biding:"required"`
		Grade int16  `json:"grade" biding:"required"`
	}

	if c.Bind(&student) == nil {
		entity := models.Student{
			Name:  student.Name,
			Grade: student.Grade,
		}

		err := s.studentRepo.Update(context.Background(), c.Param("id"), &entity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		produceMessagesUpdate(s.producer, entity)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// order message by message key. hash key partition.
func produceMessagesCreate(producer sarama.SyncProducer, student models.Student) {
	s := kafka_pb.Student{
		Name:  student.Name,
		Grade: int32(student.Grade),
	}

	m, err := proto.Marshal(&s)
	if err != nil {
		log.Fatalln("marshal error: ", err.Error())
	}

	message := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(s.Grade),
		Topic: student_create,
		Value: sarama.ByteEncoder(m),
	}
	p, o, err := producer.SendMessage(message)

	if err != nil {
		fmt.Println("Err publish: ", err.Error())
	}

	a, _ := sarama.StringEncoder(s.Grade).Encode()

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
	fmt.Println("Key: ", a)

}

// order message by custom partitioner.
func produceMessagesUpdate(producer sarama.SyncProducer, student models.Student) {
	s := kafka_pb.Student{
		Name:  student.Name,
		Grade: int32(student.Grade),
	}

	m, err := proto.Marshal(&s)
	if err != nil {
		log.Fatalln("marshal error: ", err.Error())
	}

	message := &sarama.ProducerMessage{Topic: student_update, Value: sarama.ByteEncoder(m)}
	p, o, err := producer.SendMessage(message)

	if err != nil {
		fmt.Println("Err publish: ", err.Error())
	}

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)

}
