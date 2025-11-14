package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-api/internal/api"
	"go-api/internal/api/routes"
	"go-api/internal/config"
	"go-api/internal/db"
	"go-api/internal/logging"
	"go-api/internal/middleware"
	"go-api/internal/repositories"
	"go-api/internal/services"
	"go-api/internal/storage"
)

func main() {
	logging.InitLogger()
	defer logging.Logger.Sync()

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

	notificationRepository := repositories.NewNotificationRepository(gormDB)
	notificationService := services.NewNotificationService(notificationRepository)

	fineRepository := repositories.NewFineRepository(gormDB)
	fineService := services.NewFineService(fineRepository, notificationService)
	fineHandler := api.NewFineHandler(fineService)

	partRepository := repositories.NewPartRepository(gormDB)
	inventoryTransactionRepository := repositories.NewInventoryTransactionRepository(gormDB)
	partService := services.NewPartService(partRepository, inventoryTransactionRepository, notificationService)
	partHandler := api.NewPartHandler(partService)

	freightOrderRepository := repositories.NewFreightOrderRepository(gormDB)
	freightOrderService := services.NewFreightOrderService(freightOrderRepository, vehicleRepository, journeyService)
	freightOrderHandler := api.NewFreightOrderHandler(freightOrderService)

	fileStorageService := storage.NewLocalStorageService("static")
	documentRepository := repositories.NewDocumentRepository(gormDB)
	documentService := services.NewDocumentService(documentRepository, fileStorageService)
	documentHandler := api.NewDocumentHandler(documentService)

	router := gin.Default()
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorHandler())

	router.Static("/static", "./static")

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
			routes.RegisterFineRoutes(fineHandler)(authRequired)
			routes.RegisterPartRoutes(partHandler)(authRequired)
			routes.RegisterFreightOrderRoutes(freightOrderHandler)(authRequired)
			routes.RegisterDocumentRoutes(documentHandler)(authRequired)
		}
	}

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.SERVER_PORT,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Fatal("Failed to run server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.Logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logging.Logger.Info("Server exiting")
}
