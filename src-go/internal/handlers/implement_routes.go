package handlers

import "github.com/gin-gonic/gin"

func RegisterImplementRoutes(router *gin.RouterGroup) {
	router.POST("/implements", CreateImplement)
	router.GET("/implements", GetImplements)
	router.GET("/implements/:id", GetImplementByID)
	router.PUT("/implements/:id", UpdateImplement)
	router.DELETE("/implements/:id", DeleteImplement)
}
