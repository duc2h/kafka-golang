package consumers

import (
	"context"
	"log"

	"github.com/edarha/kafka-golang/internals/models"
	"github.com/edarha/kafka-golang/internals/repositories"
	kafka_pb "github.com/edarha/kafka-golang/pb/kafka"
	"github.com/google/uuid"

	"google.golang.org/protobuf/proto"
)

type classStudentConsumer struct {
	classStudentRepo repositories.ClassStudent
}

func NewClassStudentConsumer(classStudentRepo repositories.ClassStudent) *classStudentConsumer {
	return &classStudentConsumer{
		classStudentRepo: classStudentRepo,
	}
}

func (c *classStudentConsumer) Handler(msg *Message) error {
	ctx := context.Background()
	classStudent := kafka_pb.ClassStudent{}
	err := proto.Unmarshal(msg.Value, &classStudent)

	if err != nil {
		return err
	}

	total, err := c.classStudentRepo.Count(ctx, classStudent.ClassId)
	if err != nil {

		return err
	}

	log.Printf("Message claimed: class_id = %s, student_id = %s, timestamp = %v, topic = %s, partition = %d, key = %s", classStudent.ClassId, classStudent.StudentId, msg.Timestamp, msg.Topic, msg.Partition, string(msg.Key))

	if total < 10 {
		return c.classStudentRepo.Create(ctx, &models.ClassStudent{
			ID:        uuid.NewString(),
			ClassID:   classStudent.ClassId,
			StudentID: classStudent.StudentId,
		})

	}

	// we can manage the error logic here by notification. (not now)

	return nil
}
