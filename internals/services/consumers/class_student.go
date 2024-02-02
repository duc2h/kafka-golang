package consumers

import (
	"context"
	"log"

	"github.com/edarha/kafka-golang/internals/models"
	"github.com/edarha/kafka-golang/internals/repositories"
	"github.com/edarha/kafka-golang/internals/services"
	"github.com/edarha/kafka-golang/internals/utils"
	kafka_pb "github.com/edarha/kafka-golang/pb/kafka"
	"go.temporal.io/sdk/client"

	"google.golang.org/protobuf/proto"
)

type classStudentConsumer struct {
	classStudentRepo repositories.ClassStudent
	workflow         *services.ClassStudentWorkflow
}

func NewClassStudentConsumer(classStudentRepo repositories.ClassStudent,
	workflow *services.ClassStudentWorkflow) *classStudentConsumer {
	return &classStudentConsumer{
		classStudentRepo: classStudentRepo,
		workflow:         workflow,
	}
}

func (c *classStudentConsumer) Handler(msg *Message) error {
	ctx := context.Background()
	classStudent := kafka_pb.ClassStudent{}
	err := proto.Unmarshal(msg.Value, &classStudent)

	if err != nil {
		return err
	}

	log.Printf("Message claimed: class_id = %s, student_id = %s, timestamp = %v, topic = %s, partition = %d, key = %s\n\n", classStudent.ClassId, classStudent.StudentId, msg.Timestamp, msg.Topic, msg.Partition, string(msg.Key))

	// call register workflow here
	// logic1: check class has less than 10 user.
	// if the condition fail => return fail, not do anything

	// logic2: check user has less than 3 class.
	// if the condition fail, remove the lasted class, then insert the new class.

	cc, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer cc.Close()

	options := client.StartWorkflowOptions{
		ID:        "transfer-money-workflow",
		TaskQueue: utils.CLASS_STUDENT_REGISTER_QUEUE,
	}

	we, err := cc.ExecuteWorkflow(ctx, options, c.workflow.RegisterStudentClass, models.ClassStudent{
		ClassID:   classStudent.ClassId,
		StudentID: classStudent.StudentId,
	})

	if err != nil {
		return err
	}

	log.Printf("Execute workflow success ID: %s, RunID: %s", we.GetID(), we.GetRunID())

	return nil
}
