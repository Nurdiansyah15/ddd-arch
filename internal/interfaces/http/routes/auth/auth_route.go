package auth

import (
	"os"
	"time"

	"github.com/Nurdiansyah15/ddd-arch/internal/app/domain/master/user"
	useauth "github.com/Nurdiansyah15/ddd-arch/internal/app/usecases/auth"
	userprofile "github.com/Nurdiansyah15/ddd-arch/internal/app/usecases/user"
	userrepo "github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/persistence/user"
	"github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/token"
	"github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/handlers/auth"
	authMiddleware "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// compile-time assertions: ensure infra TokenService implements usecase contracts
var _ useauth.TokenGenerator = (*token.TokenService)(nil)
var _ useauth.TokenService = (*token.TokenService)(nil)

func SetupAuthRoutes(r *gin.RouterGroup, db *sqlx.DB) {

	userRepo := userrepo.NewUserRepositoryPG(db)

	// Domain Service — business rule yang bukan milik satu entity
	userSvc := user.NewUserService(userRepo)

	// token service (use env JWT_SECRET, default fallback)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}

	tokenSvc := token.NewTokenService(secret, time.Minute*15, secret, time.Hour*24*7)

	loginUc := useauth.NewLoginUsecase(userRepo, tokenSvc)
	registerUc := useauth.NewRegisterUsecase(userRepo, userSvc)
	refreshUc := useauth.NewRefreshUsecase(tokenSvc)
	profileUc := userprofile.NewProfileUsecase(userRepo)

	authHandler := &auth.AuthHandler{LoginUC: loginUc, RegisterUC: registerUc, RefreshUC: refreshUc, ProfileUC: profileUc}

	authRoute := r.Group("/auth")

	authRoute.POST("/login", authHandler.Login)
	authRoute.POST("/register", authHandler.Register)
	authRoute.POST("/refresh", authHandler.Refresh)

	// protected routes
	authRoute.Use(authMiddleware.AuthMiddleware(tokenSvc))
	{
		authRoute.GET("/me", authHandler.GetMe)
	}
}
