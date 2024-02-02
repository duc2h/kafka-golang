package services

import (
	"time"

	"github.com/edarha/kafka-golang/internals/models"
	"github.com/edarha/kafka-golang/internals/utils"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type ClassStudentWorkflow struct {
	activity *ClassStudentActivity
}

func NewClassStudentWorkflow(activity *ClassStudentActivity) *ClassStudentWorkflow {
	return &ClassStudentWorkflow{
		activity: activity,
	}
}

func (w *ClassStudentWorkflow) RegisterStudentClass(ctx workflow.Context, classStudent models.ClassStudent) error {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    3,
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failures by default, this is just an example.
		RetryPolicy: retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, w.activity.CheckClassSize, classStudent).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, w.activity.RegisterStudentClass, classStudent).Get(ctx, nil)
	if err != nil {
		if err == utils.STUDENT_CLASS_SIZE_OVER {
			err = workflow.ExecuteActivity(ctx, w.activity.StudentClassSizeOver, classStudent).Get(ctx, nil)
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	return nil
}
