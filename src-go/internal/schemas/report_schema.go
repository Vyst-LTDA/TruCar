package schemas

// DashboardSummary é o schema para o resumo de dados do dashboard.
type DashboardSummary struct {
	TotalDistance float64 `json:"total_distance"`
	TotalFuel     float64 `json:"total_fuel"`
	TotalCost     float64 `json:"total_cost"`
	// Adicionar outros campos conforme necessário
}

// FleetManagementReport é o schema para o relatório de gerenciamento da frota.
type FleetManagementReport struct {
	// Definir a estrutura com base no `crud.report.get_fleet_management_data`
}

// DriverPerformanceReport é o schema para o relatório de desempenho dos motoristas.
type DriverPerformanceReport struct {
	// Definir a estrutura com base no `crud.report.get_driver_performance_data`
}

// VehicleConsolidatedReport é o schema para o relatório consolidado do veículo.
type VehicleConsolidatedReport struct {
	// Definir a estrutura com base no `crud.report.get_vehicle_consolidated_data`
}
