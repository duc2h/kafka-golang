package main

import (
	"log"

	"github.com/edarha/kafka-golang/internals/configs"
	"github.com/edarha/kafka-golang/internals/must"
	"github.com/edarha/kafka-golang/internals/repositories"
	"github.com/edarha/kafka-golang/internals/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.PostgreSQL{
		Host:     "localhost",
		Port:     "5432",
		Database: "postgres",
		Username: "postgres",
		Password: "admin",
	}

	db := must.ConnectPostgresql(&cfg)

	// migrate database
	if err := must.MigrateDB(db); err != nil {
		log.Fatal("something wrong while migrating database. err: %w", err)
	}

	// init repo
	studentRepo := repositories.NewStudentRepo(db)

	// init service
	studentSvc := services.NewStudent(studentRepo)

	// setup server
	r := gin.Default()
	r.POST("/student", studentSvc.Post)

	r.Run(":8080")
}
