package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"go-api/internal/api"
	"go-api/internal/api/routes"
	"go-api/internal/config"
	"go-api/internal/db"
	"go-api/internal/middleware"
	"go-api/internal/repositories"
	"go-api/internal/services"
)

func main() {
	config.LoadConfig()

	gormDB := db.InitDB()
	db.Migrate(gormDB)
	redisClient := db.InitRedis()

	cacheRepository := repositories.NewRedisCacheRepository(redisClient)

	userRepository := repositories.NewUserRepository(gormDB)
	userService := services.NewUserService(userRepository)
	userHandler := api.NewUserHandler(userService)

	authService := services.NewAuthService(userRepository)
	authHandler := api.NewAuthHandler(authService)

	vehicleRepository := repositories.NewVehicleRepository(gormDB)
	vehicleService := services.NewVehicleService(vehicleRepository, cacheRepository)
	vehicleHandler := api.NewVehicleHandler(vehicleService)

	implementRepository := repositories.NewImplementRepository(gormDB)
	implementService := services.NewImplementService(implementRepository)
	implementHandler := api.NewImplementHandler(implementService)

	journeyRepository := repositories.NewJourneyRepository(gormDB)
	journeyService := services.NewJourneyService(journeyRepository, vehicleRepository)
	journeyHandler := api.NewJourneyHandler(journeyService)

	fuelLogRepository := repositories.NewFuelLogRepository(gormDB)
	fuelLogService := services.NewFuelLogService(fuelLogRepository)
	fuelLogHandler := api.NewFuelLogHandler(fuelLogService)

	maintenanceRepository := repositories.NewMaintenanceRepository(gormDB)
	maintenanceService := services.NewMaintenanceService(maintenanceRepository)
	maintenanceHandler := api.NewMaintenanceHandler(maintenanceService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	apiV1 := router.Group("/api/v1")
	{
		routes.RegisterLoginRoutes(authHandler)(apiV1)

		authRequired := apiV1.Group("/")
		authRequired.Use(middleware.AuthMiddleware(userService))
		{
			routes.RegisterUserRoutes(userHandler)(authRequired)
			routes.RegisterVehicleRoutes(vehicleHandler)(authRequired)
			routes.RegisterImplementRoutes(implementHandler)(authRequired)
			routes.RegisterJourneyRoutes(journeyHandler)(authRequired)
			routes.RegisterFuelLogRoutes(fuelLogHandler)(authRequired)
			routes.RegisterMaintenanceRoutes(maintenanceHandler)(authRequired)
		}
	}

	if err := router.Run(":" + config.AppConfig.SERVER_PORT); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
