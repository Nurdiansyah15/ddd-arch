package http

import (
	authhandler "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/handlers/auth"
	userhandler "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/handlers/user"
	authMiddleware "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/middlewares/auth"
	authRoutes "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/routes/auth"
	userRoutes "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/routes/user"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(
	r *gin.Engine,
	authH *authhandler.AuthHandler,
	userH *userhandler.UserHandler,
	tokenSvc authMiddleware.TokenValidator,
) {
	api := r.Group("/api/v1")

	authRoutes.SetupAuthRoutes(api, authH, tokenSvc)
	userRoutes.SetupUserRoutes(api, userH)

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

	r.Static("/public", "./public")
}
