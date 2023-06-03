package services

import (
	"context"
	"net/http"

	"github.com/edarha/kafka-golang/internals/models"
	"github.com/edarha/kafka-golang/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type student struct {
	studentRepo repositories.Student
}

func NewStudent(studentRepo repositories.Student) *student {
	return &student{
		studentRepo: studentRepo,
	}
}

func (s *student) Post(c *gin.Context) {
	var student struct {
		UserId string `json:"user_id" biding:"required"`
		Grade  int16  `json:"grade" biding:"required"`
	}

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
