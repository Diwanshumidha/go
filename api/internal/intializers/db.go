package initializers

import (
	"go-api/internal/env"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// InitializeDB initializes the singleton database connection
func InitializeDB() {
	once.Do(func() {
		dsn := env.GetString("DB_CONNECTION_STRING", "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable")
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		// Configure connection pool
		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(0)

		dbInstance = db
		log.Println("Database connected successfully")
	})
}

// GetDB returns the singleton database instance
func GetDB() *gorm.DB {
	if dbInstance == nil {
		log.Fatal("Database not initialized. Call InitializeDB() first.")
	}
	return dbInstance
}
