package repositories

import (
	"context"

	"github.com/edarha/kafka-golang/internals/models"
	"gorm.io/gorm"
)

type Class interface {
	Create(ctx context.Context, entity *models.Class) error
}

type classRepo struct {
	db *gorm.DB
}

func NewClassRepo(db *gorm.DB) Class {
	return &classRepo{
		db: db,
	}
}

func (r *classRepo) Create(ctx context.Context, entity *models.Class) error {
	return r.db.WithContext(ctx).Model(&models.ClassStudent{}).Create(entity).Error
}
