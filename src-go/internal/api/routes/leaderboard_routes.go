package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/internal/api"
)

// RegisterLeaderboardRoutes registra as rotas de leaderboard.
func RegisterLeaderboardRoutes(handler *api.LeaderboardHandler) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		leaderboard := router.Group("/leaderboard")
		{
			leaderboard.GET("/", handler.GetLeaderboard)
		}
		performance := router.Group("/performance")
		{
			performance.GET("/leaderboard", handler.GetLeaderboard)
		}
	}
}
