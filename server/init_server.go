package server

import (
	"os"

	"github.com/gin-gonic/gin"
	r "github.com/stevensopilidis/dora/server/routes"
)

func InitServer() {
	engine := gin.Default()
	addr := ":8080"
	if os.Getenv("Addr") != "" {
		addr = os.Getenv("Addr")
	}
	r.InitRoutes(engine)
	engine.Run(addr)
}
