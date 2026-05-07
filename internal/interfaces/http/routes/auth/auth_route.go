package auth

import (
	authhandler "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/handlers/auth"
	authMiddleware "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/middlewares/auth"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup, h *authhandler.AuthHandler, tokenSvc authMiddleware.TokenValidator) {
	authRoute := r.Group("/auth")

	authRoute.POST("/login", h.Login)
	authRoute.POST("/register", h.Register)
	authRoute.POST("/refresh", h.Refresh)

	// protected routes
	authRoute.Use(authMiddleware.AuthMiddleware(tokenSvc))
	{
		authRoute.GET("/me", h.GetMe)
	}
}
