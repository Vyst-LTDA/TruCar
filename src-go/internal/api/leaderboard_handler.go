package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-api/internal/middleware"
	"go-api/internal/services"
)

// LeaderboardHandler gerencia as requisições HTTP para o leaderboard.
type LeaderboardHandler struct {
	service services.LeaderboardService
}

// NewLeaderboardHandler cria uma nova instância de LeaderboardHandler.
func NewLeaderboardHandler(service services.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{service: service}
}

// GetLeaderboard lida com a busca de dados para o leaderboard.
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	orgID, exists := middleware.GetOrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
		return
	}

	leaderboard, err := h.service.GetLeaderboard(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard data"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}
