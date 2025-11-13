package db

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-api/internal/config"
	"go-api/internal/models"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.AppConfig.DB_DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.Vehicle{},
		&models.Implement{},
		&models.Journey{},
		&models.FuelLog{},
		&models.MaintenanceRequest{},
		&models.MaintenanceComment{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.REDIS_ADDR,
		Password: config.AppConfig.REDIS_PASSWORD,
		DB:       config.AppConfig.REDIS_DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return rdb
}
