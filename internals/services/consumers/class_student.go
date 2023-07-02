package consumers

import (
	"log"

	kafka_pb "github.com/edarha/kafka-golang/pb/kafka"

	"google.golang.org/protobuf/proto"
)

type classStudentConsumer struct {
}

func NewClassStudentConsumer() *classStudentConsumer {
	return &classStudentConsumer{}
}

func (c *classStudentConsumer) Handler(msg *Message) error {
	classStudent := kafka_pb.ClassStudent{}
	err := proto.Unmarshal(msg.Value, &classStudent)

	if err != nil {
		return err
	}

	log.Printf("Message claimed: class_id = %s, student_id = %s, timestamp = %v, topic = %s, partition = %d, key = %s", classStudent.ClassId, classStudent.StudentId, msg.Timestamp, msg.Topic, msg.Partition, string(msg.Key))

	return nil
}
