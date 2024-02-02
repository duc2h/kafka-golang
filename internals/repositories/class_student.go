package repositories

import (
	"context"
	"time"

	"github.com/edarha/kafka-golang/internals/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClassStudent interface {
	Create(ctx context.Context, entity *models.ClassStudent) error
	DeleteByID(ctx context.Context, ID string) error

	GetLastedByStudentId(ctx context.Context, studentId string) (*models.ClassStudent, error)
	CountByClassId(ctx context.Context, classId string) (int64, error)
	CountByStudentId(ctx context.Context, studentId string) (int64, error)
}

type classStudentRepo struct {
	db *gorm.DB
}

func NewClassStudentRepo(db *gorm.DB) ClassStudent {
	return &classStudentRepo{
		db: db,
	}
}

func (r *classStudentRepo) Create(ctx context.Context, entity *models.ClassStudent) error {
	entity.ID = uuid.NewString()
	entity.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Model(&models.ClassStudent{}).Create(entity).Error
}

func (r *classStudentRepo) CountByClassId(ctx context.Context, classId string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.ClassStudent{}).Where("class_id = ?", classId).Count(&count)

	return count, query.Error
}

func (r *classStudentRepo) CountByStudentId(ctx context.Context, studentId string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.ClassStudent{}).Where("student_id = ?", studentId).Count(&count)

	return count, query.Error
}

func (r *classStudentRepo) GetLastedByStudentId(ctx context.Context, studentId string) (*models.ClassStudent, error) {
	var result models.ClassStudent
	err := r.db.WithContext(ctx).Model(&models.ClassStudent{}).
		Where("student_id = ?", studentId).
		Order("created_at").
		First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *classStudentRepo) DeleteByID(ctx context.Context, ID string) error {
	return r.db.WithContext(ctx).Where("id = ?", ID).Delete(&models.ClassStudent{}).Error
}
