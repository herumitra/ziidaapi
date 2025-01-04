package config

import (
	context "context"
	log "log"
	os "os"

	models "github.com/herumitra/ziidaapi/models"
	redis "github.com/redis/go-redis/v9"
	postgres "gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

// Setup database connection
func SetupDB() (err error) {
	// Initialize context
	ctx := context.Background()

	// Get environment variable
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")

	// Connect to database Postgres
	// dsn := "host=" + db_host + " user=" + db_user + " password=" + db_pass + " dbname=" + db_name + "  port=" + db_port + " sslmode=disable TimeZone=Asia/Jakarta"
	dsn := "user=" + db_user + " password=" + db_pass + " host=" + db_host + " port=" + db_port + " dbname=" + db_name + "  sslmode=disable TimeZone=Asia/Jakarta"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Migrate the schema
	DB.AutoMigrate(
		&models.User{},
		&models.Branch{},
		&models.UserBranch{},
		&models.Unit{},
		&models.UnitConversion{},
		&models.ProductCategory{},
		&models.Product{},
		&models.MemberCategory{},
		&models.Member{},
		&models.SupplierCategory{},
		&models.Supplier{},
		&models.SupplierProduct{},
	)

	// Check connection Postgres
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL database: %v", err)
	}

	// Connect to database Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     redis_host + ":" + redis_port,
		Password: "",
		DB:       1,
	})

	// Check connection Redis
	_, err = RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis database: %v", err)
	}

	return nil
}
