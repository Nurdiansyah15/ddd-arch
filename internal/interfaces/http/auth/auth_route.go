package auth

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/persistence/user"
	"github.com/Nurdiansyah15/ddd-arch/internal/usecase/auth"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupAuthRoutes(r *gin.RouterGroup, db *sqlx.DB) {

	userRepo := &user.UserRepositoryPG{DB: db}
	loginUc := &auth.LoginUsecase{UserRepo: userRepo}
	registerUc := &auth.RegisterUsecase{UserRepo: userRepo}
	authHandler := &AuthHandler{LoginUC: loginUc, RegisterUC: registerUc}

	authRoute := r.Group("/auth")

	authRoute.POST("/login", authHandler.Login)
	authRoute.POST("/register", authHandler.Register)
}
