package services

import (
	"sync"
	"time"

	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
)

// DashboardService define a interface para a lógica de negócios do dashboard.
type DashboardService interface {
	GetManagerDashboard(orgID uint, period string, userRole models.UserRole) (*schemas.ManagerDashboardResponse, error)
	GetDriverDashboard(userID, orgID uint) (*schemas.DriverDashboardResponse, error)
	GetVehiclePositions(orgID uint) ([]schemas.VehiclePosition, error)
	GetDemoStats(orgID uint) (*schemas.DemoStatsResponse, error)
}

type dashboardService struct {
	userRepo        repositories.UserRepository
	vehicleRepo     repositories.VehicleRepository
	vehicleCostRepo repositories.VehicleCostRepository
	journeyRepo     repositories.JourneyRepository
	maintenanceRepo repositories.MaintenanceRepository
	fineRepo        repositories.FineRepository
	notificationRepo repositories.NotificationRepository
	partRepo        repositories.PartRepository
	documentRepo    repositories.DocumentRepository
	fuelLogRepo     repositories.FuelLogRepository
}

// NewDashboardService cria uma nova instância de DashboardService.
func NewDashboardService(
	userRepo repositories.UserRepository,
	vehicleRepo repositories.VehicleRepository,
	vehicleCostRepo repositories.VehicleCostRepository,
	journeyRepo repositories.JourneyRepository,
	maintenanceRepo repositories.MaintenanceRepository,
	fineRepo repositories.FineRepository,
	notificationRepo repositories.NotificationRepository,
	partRepo repositories.PartRepository,
	documentRepo repositories.DocumentRepository,
	fuelLogRepo repositories.FuelLogRepository,
) DashboardService {
	return &dashboardService{
		userRepo:        userRepo,
		vehicleRepo:     vehicleRepo,
		vehicleCostRepo: vehicleCostRepo,
		journeyRepo:     journeyRepo,
		maintenanceRepo: maintenanceRepo,
		fineRepo:        fineRepo,
		notificationRepo: notificationRepo,
		partRepo:        partRepo,
		documentRepo:    documentRepo,
		fuelLogRepo:     fuelLogRepo,
	}
}

// getStartDateFromPeriod converte uma string de período em uma data de início.
func getStartDateFromPeriod(period string) time.Time {
	today := time.Now()
	if period == "last_7_days" {
		return today.AddDate(0, 0, -7)
	}
	if period == "this_month" {
		return time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	}
	// Padrão para "last_30_days"
	return today.AddDate(0, 0, -30)
}

func (s *dashboardService) GetManagerDashboard(orgID uint, period string, userRole models.UserRole) (*schemas.ManagerDashboardResponse, error) {
	startDate := getStartDateFromPeriod(period)

	// Usar um WaitGroup para aguardar a conclusão de todas as goroutines
	var wg sync.WaitGroup
	errChan := make(chan error, 7) // Canal para coletar erros

	var kpis schemas.KPI
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalVehicles, err := s.vehicleRepo.CountByOrganization(orgID, "")
		if err != nil {
			errChan <- err
			return
		}
		kpis.TotalVehicles = int(totalVehicles)

		totalDrivers, err := s.userRepo.CountDriversByOrganization(orgID)
		if err != nil {
			errChan <- err
			return
		}
		kpis.TotalDrivers = int(totalDrivers)

		totalDistance, err := s.journeyRepo.SumDistanceByOrganization(orgID, startDate)
		if err != nil {
			errChan <- err
			return
		}
		kpis.TotalDistance = totalDistance

		totalFuel, err := s.fuelLogRepo.SumFuelByOrganization(orgID, startDate)
		if err != nil {
			errChan <- err
			return
		}
		kpis.TotalFuel = totalFuel
	}()

	var efficiencyKPIs schemas.EfficiencyKPIs
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Implementar busca real
		efficiencyKPIs = schemas.EfficiencyKPIs{AverageConsumption: 6.25, AverageCostPerKm: 1.5}
	}()

	var costsByCategory []schemas.CostByCategory
	wg.Add(1)
	go func() {
		defer wg.Done()
		costs, err := s.vehicleCostRepo.FindByOrganizationID(orgID, &startDate, nil)
		if err != nil {
			errChan <- err
			return
		}
		// Agrupar custos por categoria
		costMap := make(map[models.CostType]float64)
		for _, cost := range costs {
			costMap[cost.CostType] += cost.Amount
		}
		for category, amount := range costMap {
			costsByCategory = append(costsByCategory, schemas.CostByCategory{Category: string(category), Amount: amount})
		}
	}()

	// ... (outras goroutines para alertas, manutenções, etc.)

	wg.Wait()
	close(errChan)

	// Verificar se alguma goroutine retornou erro
	for err := range errChan {
		if err != nil {
			return nil, err // Retornar o primeiro erro encontrado
		}
	}

	response := &schemas.ManagerDashboardResponse{
		KPIs:                 kpis,
		EfficiencyKPIs:       efficiencyKPIs,
		CostsByCategory:      costsByCategory,
		RecentAlerts:         []schemas.RecentAlert{},         // Placeholder
		UpcomingMaintenances: []schemas.UpcomingMaintenance{}, // Placeholder
		ActiveGoal:           nil,                             // Placeholder
	}

	// Adicionar dados premium se o usuário for CLIENTE_ATIVO
	if userRole == models.RoleClienteAtivo {
		// Implementar busca por Km por dia e pódio
		response.KmPerDayLast30Days = []schemas.KmPerDay{} // Placeholder
		response.PodiumDrivers = []schemas.PodiumDriver{}   // Placeholder
	}

	return response, nil
}

func (s *dashboardService) GetDriverDashboard(userID, orgID uint) (*schemas.DriverDashboardResponse, error) {
	// TODO: Implementar a lógica para buscar métricas, ranking e conquistas do motorista.

	// Placeholder de retorno
	return &schemas.DriverDashboardResponse{
		Metrics:       "Dados de métricas do motorista",
		RankingContext: "Contexto do ranking",
		Achievements:  "Conquistas do motorista",
	}, nil
}

func (s *dashboardService) GetVehiclePositions(orgID uint) ([]schemas.VehiclePosition, error) {
	vehicles, err := s.vehicleRepo.FindLatestPositionsByOrganization(orgID)
	if err != nil {
		return nil, err
	}

	positions := make([]schemas.VehiclePosition, len(vehicles))
	for i, v := range vehicles {
		positions[i] = schemas.VehiclePosition{
			ID:           v.ID,
			LicensePlate: *v.LicensePlate,
			Latitude:     *v.LastLatitude,
			Longitude:    *v.LastLongitude,
			Timestamp:    *v.LastLocationUpdate,
		}
	}

	return positions, nil
}

func (s *dashboardService) GetDemoStats(orgID uint) (*schemas.DemoStatsResponse, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, 4)

	var stats schemas.DemoStatsResponse

	// Limites (poderiam vir de uma configuração)
	const (
		vehicleLimit = 10
		userLimit    = 5
		partLimit    = 50
		docLimit     = 20
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := s.vehicleRepo.CountByOrganization(orgID, "")
		if err != nil {
			errChan <- err
			return
		}
		stats.Vehicles = schemas.DemoResourceLimit{Current: int(count), Limit: vehicleLimit}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := s.userRepo.CountByOrganization(orgID)
		if err != nil {
			errChan <- err
			return
		}
		stats.Users = schemas.DemoResourceLimit{Current: int(count), Limit: userLimit}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := s.partRepo.CountByOrganization(orgID)
		if err != nil {
			errChan <- err
			return
		}
		stats.Parts = schemas.DemoResourceLimit{Current: int(count), Limit: partLimit}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Este método precisa ser adicionado ao DocumentRepository
		count, err := s.documentRepo.CountByOrganization(orgID)
		if err != nil {
			errChan <- err
			return
		}
		stats.Documents = schemas.DemoResourceLimit{Current: int(count), Limit: docLimit}
	}()

	// Adicionar placeholders para estatísticas mensais
	stats.Reports = schemas.DemoResourceLimit{Current: 0, Limit: 10}
	stats.Fines = schemas.DemoResourceLimit{Current: 0, Limit: 5}
	stats.FreightOrders = schemas.DemoResourceLimit{Current: 0, Limit: 10}
	stats.MaintenanceRequests = schemas.DemoResourceLimit{Current: 0, Limit: 5}
	stats.FuelLogs = schemas.DemoResourceLimit{Current: 0, Limit: 100}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return &stats, nil
}
