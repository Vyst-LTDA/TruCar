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
	"go-api/internal/models"
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

	// Repositories
	cacheRepository := repositories.NewRedisCacheRepository(redisClient)
	userRepository := repositories.NewUserRepository(gormDB)
	vehicleRepository := repositories.NewVehicleRepository(gormDB)
	implementRepository := repositories.NewImplementRepository(gormDB)
	journeyRepository := repositories.NewJourneyRepository(gormDB)
	fuelLogRepository := repositories.NewFuelLogRepository(gormDB)
	maintenanceRepository := repositories.NewMaintenanceRepository(gormDB)
	notificationRepository := repositories.NewNotificationRepository(gormDB)
	fineRepository := repositories.NewFineRepository(gormDB)
	partRepository := repositories.NewPartRepository(gormDB)
	inventoryTransactionRepository := repositories.NewInventoryTransactionRepository(gormDB)
	freightOrderRepository := repositories.NewFreightOrderRepository(gormDB)
	documentRepository := repositories.NewDocumentRepository(gormDB)
	organizationRepository := repositories.NewOrganizationRepository(gormDB)

	// Services
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(userRepository)
	vehicleService := services.NewVehicleService(vehicleRepository, cacheRepository)
	implementService := services.NewImplementService(implementRepository)
	journeyService := services.NewJourneyService(journeyRepository, vehicleRepository)
	fuelLogService := services.NewFuelLogService(fuelLogRepository)
	maintenanceService := services.NewMaintenanceService(maintenanceRepository)
	notificationService := services.NewNotificationService(notificationRepository)
	fineService := services.NewFineService(fineRepository, notificationService)
	partService := services.NewPartService(partRepository, inventoryTransactionRepository, notificationService)
	freightOrderService := services.NewFreightOrderService(freightOrderRepository, vehicleRepository, journeyService)
	fileStorageService := storage.NewLocalStorageService("static")
	documentService := services.NewDocumentService(documentRepository, fileStorageService)
	organizationService := services.NewOrganizationService(organizationRepository)

	// Handlers
	userHandler := api.NewUserHandler(userService)
	authHandler := api.NewAuthHandler(authService)
	vehicleHandler := api.NewVehicleHandler(vehicleService)
	implementHandler := api.NewImplementHandler(implementService)
	journeyHandler := api.NewJourneyHandler(journeyService)
	fuelLogHandler := api.NewFuelLogHandler(fuelLogService)
	maintenanceHandler := api.NewMaintenanceHandler(maintenanceService)
	fineHandler := api.NewFineHandler(fineService)
	partHandler := api.NewPartHandler(partService)
	freightOrderHandler := api.NewFreightOrderHandler(freightOrderService)
	documentHandler := api.NewDocumentHandler(documentService)
	adminHandler := api.NewAdminHandler(organizationService, userService, authService)

	router := gin.Default()
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorHandler())

	router.Static("/static", "./static")

	apiV1 := router.Group("/api/v1")
	{
		// Public routes
		routes.RegisterLoginRoutes(authHandler)(apiV1)

		// Authenticated routes
		authRequired := apiV1.Group("/")
		authRequired.Use(middleware.AuthMiddleware(userService))
		{
			// SuperAdmin routes
			superAdminRoutes := authRequired.Group("/admin")
			superAdminRoutes.Use(middleware.AuthorizationMiddleware(models.RoleSuperAdmin))
			{
				routes.RegisterAdminRoutes(adminHandler)(superAdminRoutes)
			}

			// Manager routes
			managerRoutes := authRequired.Group("/")
			managerRoutes.Use(middleware.AuthorizationMiddleware(models.RoleClienteAtivo, models.RoleClienteDemo))
			{
				routes.RegisterUserRoutes(userHandler)(managerRoutes)
				routes.RegisterVehicleRoutes(vehicleHandler)(managerRoutes)
				routes.RegisterImplementRoutes(implementHandler)(managerRoutes)
				routes.RegisterPartRoutes(partHandler)(managerRoutes)
				routes.RegisterDocumentRoutes(documentHandler)(managerRoutes)
				// Add other manager routes here
			}

			// Driver and Manager routes
			driverAndManagerRoutes := authRequired.Group("/")
			driverAndManagerRoutes.Use(middleware.AuthorizationMiddleware(models.RoleDriver, models.RoleClienteAtivo, models.RoleClienteDemo))
			{
				routes.RegisterJourneyRoutes(journeyHandler)(driverAndManagerRoutes)
				routes.RegisterFuelLogRoutes(fuelLogHandler)(driverAndManagerRoutes)
				routes.RegisterMaintenanceRoutes(maintenanceHandler)(driverAndManagerRoutes)
				routes.RegisterFineRoutes(fineHandler)(driverAndManagerRoutes)
				routes.RegisterFreightOrderRoutes(freightOrderHandler)(driverAndManagerRoutes)
			}
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
