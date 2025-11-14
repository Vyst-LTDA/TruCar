package schemas

// DriverRanking representa os dados de um único motorista no ranking.
type DriverRanking struct {
	Position  int     `json:"position"`
	DriverID  uint    `json:"driver_id"`
	DriverName string `json:"driver_name"`
	Score     float64 `json:"score"`
	// Adicionar outras métricas conforme necessário (distância, consumo, etc.)
}

// LeaderboardResponse é o schema para a resposta do endpoint de leaderboard.
type LeaderboardResponse struct {
	Rankings []DriverRanking `json:"rankings"`
}
