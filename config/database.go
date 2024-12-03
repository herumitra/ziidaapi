package config

import (
	"context"
	"log"
	"os"

	"github.com/herumitra/ziidaapi/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

// SetupDatabase initializes the database connection (Postgres & Redis)
func SetupDatabase() {
	// Connect to Postgres
	var err error
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable TimeZone=Asia/Jakarta"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL database: %v", err)
	}

	// Connect to Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       8,
	})

	// Create a Context
	ctx := context.Background()

	// Test the connection to Redis
	// coba := RDB.Ping(ctx).Val()
	// _, err = RDB.Ping(ctx).Result()
	_, err = RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis database: %v", err)
	}

	// Automigrate tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.Branch{},
		&models.UserBranch{},
		&models.Unit{},
		&models.UnitConversion{},
		&models.MemberCategory{},
		&models.Member{},
		&models.ProductCategory{},
		&models.Product{},
		&models.SupplierCategory{},
		&models.Supplier{},
		&models.SupplierProduct{},
	)
	if err != nil {
		log.Fatalf("failed to migrate PostgreSQL database: %v", err)
	}
}
