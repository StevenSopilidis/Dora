package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	e "github.com/stevensopilidis/dora/errors"
	r "github.com/stevensopilidis/dora/registry"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/services/:name", getService)
	r.POST("/services/:name", addService)
	r.DELETE("/services/:name", removeService)
}

func getService(c *gin.Context) {
	client := r.GetRedisRegistryClient()
	serviceName := c.Param("name")
	if serviceName == "" {
		res := e.InvalidArgument{
			Message: "Service name must be valid",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}
	err, service := client.Get(c.Request.Context(), serviceName)
	if err != nil {
		if errors.Is(err, &e.ServiceNotFoundError{}) {
			c.JSON(http.StatusNotFound, err)
		} else {
			c.JSON(http.StatusBadRequest, err)
		}
		return
	}
	c.JSON(http.StatusOK, service)
}

func addService(c *gin.Context) {
	client := r.GetRedisRegistryClient()
	var service r.Service
	err := c.ShouldBindJSON(&service)
	if err != nil {
		res := e.InvalidArgument{
			Message: "Must provide a valid service",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}
	serviceName := c.Param("name")
	if serviceName == "" {
		res := e.InvalidArgument{
			Message: "Service name must be valid",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}
	err = client.Append(c.Request.Context(), serviceName, service)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, service)
}

func removeService(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		res := e.InvalidArgument{
			Message: "Service name must be valid",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}
	client := r.GetRedisRegistryClient()
	err := client.Remove(c.Request.Context(), serviceName)
	if err != nil {
		if errors.Is(err, &e.ServiceNotFoundError{}) {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
