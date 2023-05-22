package must

import (
	_ "embed"
	"log"
	"time"

	"github.com/edarha/kafka-golang/internals/configs"
	"github.com/edarha/kafka-golang/internals/models"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxDBIdleConns  = 10
	maxDBOpenConns  = 100
	maxConnLifeTime = 30 * time.Minute
)

func ConnectPostgresql(cfg *configs.PostgreSQL) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Error open postgresql", zap.Error(err))
	}

	err = db.Raw("SELECT 1").Error
	if err != nil {
		log.Fatal("Error querying SELECT 1", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error get sql DB", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(maxDBIdleConns)
	sqlDB.SetMaxOpenConns(maxDBOpenConns)
	sqlDB.SetConnMaxLifetime(maxConnLifeTime)
	return db
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Parent{},
	)
}
