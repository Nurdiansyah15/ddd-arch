package http

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/auth"
	user "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/user"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, db *sqlx.DB) {
	api := r.Group("/api/v1")

	auth.SetupAuthRoutes(api, db)
	user.SetupUserRoutes(api, db)

	// Setup health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	// Setup not found handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})

	// Swagger
	if gin.Mode() != gin.ReleaseMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// router.Static("/exports", "./public/exports")
	r.Static("/public", "./public")
}
