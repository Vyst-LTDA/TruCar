package schemas

import (
	"time"
)

// --- Estruturas para Dados Agregados ---

type KPI struct {
	TotalVehicles   int     `json:"total_vehicles"`
	TotalDrivers    int     `json:"total_drivers"`
	TotalDistance   float64 `json:"total_distance"`
	TotalFuel       float64 `json:"total_fuel"`
}

type EfficiencyKPIs struct {
	AverageConsumption float64 `json:"average_consumption"`
	AverageCostPerKm   float64 `json:"average_cost_per_km"`
}

type CostByCategory struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

type KmPerDay struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type PodiumDriver struct {
	Position int    `json:"position"`
	Name     string `json:"name"`
	Score    int    `json:"score"`
}

type RecentAlert struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	VehicleName string    `json:"vehicle_name"`
}

type UpcomingMaintenance struct {
	ID            uint      `json:"id"`
	Description   string    `json:"description"`
	DueDate       time.Time `json:"due_date"`
	VehicleName   string    `json:"vehicle_name"`
}

type ActiveGoal struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Progress    float64 `json:"progress"` // 0.0 a 1.0
}

// --- Schemas de Resposta dos Endpoints ---

// ManagerDashboardResponse é o DTO para o dashboard do gestor.
type ManagerDashboardResponse struct {
	KPIs                 KPI                   `json:"kpis"`
	EfficiencyKPIs       EfficiencyKPIs        `json:"efficiency_kpis"`
	CostsByCategory      []CostByCategory      `json:"costs_by_category"`
	KmPerDayLast30Days   []KmPerDay            `json:"km_per_day_last_30_days,omitempty"`
	PodiumDrivers        []PodiumDriver        `json:"podium_drivers,omitempty"`
	RecentAlerts         []RecentAlert         `json:"recent_alerts"`
	UpcomingMaintenances []UpcomingMaintenance `json:"upcoming_maintenances"`
	ActiveGoal           *ActiveGoal           `json:"active_goal"`
}

// DriverDashboardResponse é o DTO para o dashboard do motorista.
type DriverDashboardResponse struct {
	Metrics       interface{} `json:"metrics"`       // Usar interface{} para flexibilidade
	RankingContext interface{} `json:"ranking_context"`
	Achievements  interface{} `json:"achievements"`
}

// VehiclePosition representa a geolocalização de um veículo.
type VehiclePosition struct {
	ID          uint    `json:"id"`
	LicensePlate string `json:"license_plate"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timestamp   time.Time `json:"timestamp"`
}

// DemoResourceLimit representa o uso atual e o limite de um recurso.
type DemoResourceLimit struct {
    Current int `json:"current"`
    Limit   int `json:"limit"`
}

// DemoStatsResponse é o DTO para as estatísticas da conta demo.
type DemoStatsResponse struct {
    Vehicles            DemoResourceLimit `json:"vehicles"`
    Users               DemoResourceLimit `json:"users"`
    Parts               DemoResourceLimit `json:"parts"`
    Clients             DemoResourceLimit `json:"clients"`
    Reports             DemoResourceLimit `json:"reports"`
    Fines               DemoResourceLimit `json:"fines"`
    Documents           DemoResourceLimit `json:"documents"`
    FreightOrders       DemoResourceLimit `json:"freight_orders"`
    MaintenanceRequests DemoResourceLimit `json:"maintenance_requests"`
    FuelLogs            DemoResourceLimit `json:"fuel_logs"`
}
