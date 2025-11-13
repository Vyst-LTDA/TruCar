package db

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"go-api/internal/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-api/internal/config"
	"go-api/internal/models"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.AppConfig.DB_DSN), &gorm.Config{})
	if err != nil {
		logging.Logger.Fatal("Failed to connect to database", zap.Error(err))
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
		logging.Logger.Fatal("Failed to migrate database", zap.Error(err))
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
		logging.Logger.Fatal("Could not connect to Redis", zap.Error(err))
	}

	return rdb
}
