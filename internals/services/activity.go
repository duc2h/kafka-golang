package services

import (
	"context"
	"fmt"

	"github.com/edarha/kafka-golang/internals/models"
	"github.com/edarha/kafka-golang/internals/repositories"
	"github.com/edarha/kafka-golang/internals/utils"
)

type ClassStudentActivity struct {
	classStudentRepo repositories.ClassStudent
}

func NewClassStudentActivity(repo repositories.ClassStudent) *ClassStudentActivity {
	return &ClassStudentActivity{
		classStudentRepo: repo,
	}
}

func (a *ClassStudentActivity) CheckClassSize(ctx context.Context, classStudent models.ClassStudent) error {
	total, err := a.classStudentRepo.CountByClassId(ctx, classStudent.ClassID)
	if err != nil {

		return err
	}

	if total > 10 {
		return fmt.Errorf("Class size is over 10.")
	}

	return nil
}

func (a *ClassStudentActivity) RegisterStudentClass(ctx context.Context, classStudent models.ClassStudent) error {
	// check the studentClassSize
	total, err := a.classStudentRepo.CountByStudentId(ctx, classStudent.StudentID)
	if err != nil {
		return err
	}

	if total > 3 {
		return utils.STUDENT_CLASS_SIZE_OVER
	}

	return a.classStudentRepo.Create(ctx, &classStudent)
}

func (a *ClassStudentActivity) StudentClassSizeOver(ctx context.Context, classStudent models.ClassStudent) error {
	// get the lasted studentClass by studentId
	deleteClassStudent, err := a.classStudentRepo.GetLastedByStudentId(ctx, classStudent.StudentID)
	if err != nil {
		return err
	}

	// delete it
	err = a.classStudentRepo.DeleteByID(ctx, deleteClassStudent.ID)
	if err != nil {
		return err
	}

	// insert the new one.
	return a.classStudentRepo.Create(ctx, &classStudent)
}
