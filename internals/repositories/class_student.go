package repositories

import (
	"context"

	"github.com/edarha/kafka-golang/internals/models"
	"gorm.io/gorm"
)

type ClassStudent interface {
	Create(ctx context.Context, entity *models.ClassStudent) error

	Count(ctx context.Context, classId string) (int64, error)
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
	return r.db.WithContext(ctx).Model(&models.ClassStudent{}).Create(entity).Error
}

func (r *classStudentRepo) Count(ctx context.Context, classId string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.ClassStudent{}).Where("class_id = ?", classId).Count(&count)

	return count, query.Error
}
