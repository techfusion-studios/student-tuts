package database

import (
	"github.com/techfusion/school/student/pkg/data/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type DBEngine interface {
	Connect()
	GetEngine() *gorm.DB
}

type gormDB struct {
	engine *gorm.DB
}

func (db *gormDB) GetEngine() *gorm.DB {
	return db.engine
}

func (db *gormDB) Connect() {
	// Initialize database connection
	conn, err := gorm.Open(sqlite.Open("students.db"), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		TranslateError: false,
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate the Student struct
	conn.AutoMigrate(&models.Student{})

	db.engine = conn
}

func New() DBEngine {
	return &gormDB{}
}
