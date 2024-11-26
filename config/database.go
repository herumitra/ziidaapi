package config

import (
	"context" // Impor context
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/herumitra/ziidaapi/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

// SetupDatabase initializes the PostgreSQL and Redis databases.
func SetupDatabase() {
	// PostgreSQL connection
	dsn := "host=localhost user=ziida password=S14n4kC3rd4s dbname=ziida port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL database: %v", err)
	}

	// Automigrate tables
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("failed to migrate PostgreSQL database: %v", err)
	}

	// Redis connection
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Create a context
	ctx := context.Background()

	// Test connection to Redis
	_, err = RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
}
