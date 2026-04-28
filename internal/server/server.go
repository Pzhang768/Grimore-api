package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Pzhang768/Grimore-api/internal/handlers"
	"github.com/Pzhang768/Grimore-api/internal/middleware"
	"github.com/Pzhang768/Grimore-api/internal/services"
)

func New(jwtSecret string, db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	handlers.NewAgentHandler().Register(r)

	protected := r.Group("/")
	protected.Use(middleware.Auth(jwtSecret, db))

	handlers.NewTeamHandler(services.NewTeamService(db)).Register(protected)

	return r
}
