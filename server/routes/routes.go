package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/services/:name", getService)
	r.POST("/services", addService)
	r.DELETE("/services", removeService)
}

func getService(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Getservice",
	})
}

func addService(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "PostService",
	})
}

func removeService(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "RemoveService",
	})
}
