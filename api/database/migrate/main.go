package main

import (
	"go-api/database/model"
	initializers "go-api/internal/intializers"
	"log"
)

func init() {
	initializers.EnvironmentVariables()
	initializers.InitializeDB()
}

func registerModels() []interface{} {
	return []interface{}{
		&model.User{},
		&model.ShortLink{},
	}
}

func main() {
	log.Printf("🕧 Migrating database models...")
	db := initializers.GetDB()

	for _, model := range registerModels() {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("Migration failed for %T: %v", model, err)
		}
	}

	log.Println("✅ All migrations applied successfully!")
}
