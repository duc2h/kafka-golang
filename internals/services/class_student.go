package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/edarha/kafka-golang/internals/models"
	"github.com/edarha/kafka-golang/internals/repositories"
	kafka_pb "github.com/edarha/kafka-golang/pb/kafka"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

var (
	classStudent_register = "class_student_register"
)

type classStudent struct {
	classStudentRepo repositories.ClassStudent
	producer         sarama.SyncProducer
}

func NewClassStudent(classStudentRepo repositories.ClassStudent, producer sarama.SyncProducer) *classStudent {
	return &classStudent{
		classStudentRepo: classStudentRepo,
		producer:         producer,
	}
}

func (s *classStudent) Register(c *gin.Context) {
	var classStudent struct {
		StudentId string `json:"student_id" biding:"required"`
		ClassId   string `json:"class_id" biding:"required"`
	}

	if c.Bind(&classStudent) == nil {
		entity := models.ClassStudent{
			ClassID:   classStudent.ClassId,
			StudentID: classStudent.StudentId,
		}

		// send event to kafka only.
		produceStudentRegister(s.producer, entity)

		c.JSON(http.StatusOK, gin.H{"message": "ok"})

	}
}

func produceStudentRegister(producer sarama.SyncProducer, classStudent models.ClassStudent) {
	s := kafka_pb.ClassStudent{
		ClassId:   classStudent.ClassID,
		StudentId: classStudent.StudentID,
	}

	m, err := proto.Marshal(&s)
	if err != nil {
		log.Fatalln("marshal error: ", err.Error())
	}

	message := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(classStudent.ClassID),
		Topic: classStudent_register,
		Value: sarama.ByteEncoder(m),
	}
	p, o, err := producer.SendMessage(message)

	if err != nil {
		fmt.Println("Err publish: ", err.Error())
	}

	a, _ := sarama.StringEncoder(classStudent.ClassID).Encode()

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
	fmt.Println("Key: ", a)

}
