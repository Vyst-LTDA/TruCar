package services

import (
	"go-api/internal/models"
	"go-api/internal/repositories"
	"go-api/internal/schemas"
	"sort"
)

// LeaderboardService define a interface para a lógica de negócios do leaderboard.
type LeaderboardService interface {
	GetLeaderboard(orgID uint) (*schemas.LeaderboardResponse, error)
}

type leaderboardService struct {
	userRepo    repositories.UserRepository
	journeyRepo repositories.JourneyRepository
	fineRepo    repositories.FineRepository
}

// NewLeaderboardService cria uma nova instância de LeaderboardService.
func NewLeaderboardService(
	userRepo repositories.UserRepository,
	journeyRepo repositories.JourneyRepository,
	fineRepo repositories.FineRepository,
) LeaderboardService {
	return &leaderboardService{
		userRepo:    userRepo,
		journeyRepo: journeyRepo,
		fineRepo:    fineRepo,
	}
}

// GetLeaderboard calcula e retorna o ranking dos motoristas.
func (s *leaderboardService) GetLeaderboard(orgID uint) (*schemas.LeaderboardResponse, error) {
	// 1. Obter todos os motoristas da organização
	drivers, err := s.userRepo.FindByRoleAndOrganization(orgID, models.RoleDriver)
	if err != nil {
		return nil, err
	}

	// 2. Para cada motorista, calcular a pontuação (lógica de placeholder)
	var rankings []schemas.DriverRanking
	for _, driver := range drivers {
		// Simular uma pontuação
		score := 100.0 - float64(driver.ID) // Lógica de placeholder
		rankings = append(rankings, schemas.DriverRanking{
			DriverID:   driver.ID,
			DriverName: driver.FullName,
			Score:      score,
		})
	}

	// 3. Ordenar os motoristas pela pontuação
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	// 4. Atribuir a posição final
	for i := range rankings {
		rankings[i].Position = i + 1
	}

	return &schemas.LeaderboardResponse{Rankings: rankings}, nil
}
