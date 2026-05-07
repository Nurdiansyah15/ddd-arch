package main

import (
	"log"

	_ "github.com/Nurdiansyah15/ddd-arch/docs"

	useauth "github.com/Nurdiansyah15/ddd-arch/internal/app/usecases/auth"
	useruc "github.com/Nurdiansyah15/ddd-arch/internal/app/usecases/user"
	"github.com/Nurdiansyah15/ddd-arch/internal/config"
	domainuser "github.com/Nurdiansyah15/ddd-arch/internal/domain/master/user"
	"github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/crypto"
	"github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/database/mysql"
	"github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/database/postgres"
	userrepo "github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/persistence/user"
	"github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/token"
	authhandler "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/handlers/auth"
	userhandler "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/handlers/user"
	httpRoutes "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/routes"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// @title Api Doc
// @version 1.0
// @description API for App System
// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {
	// Load .env (local only)
	_ = godotenv.Load()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// ── Infrastructure: Database ───────────────────────────────────────────────
	var db *sqlx.DB

	switch cfg.DB.Driver {
	case "mysql":
		db, err = mysql.NewMySQL(cfg.DB)
	case "postgres":
		db, err = postgres.NewPostgres(cfg.DB.DSN())
	default:
		log.Fatal("unsupported db driver")
	}
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// ── Infrastructure: Adapters ───────────────────────────────────────────────
	hasher := crypto.NewBcryptHasher()
	userRepo := userrepo.NewUserRepositoryPG(db)
	tokenSvc := token.NewTokenService(
		cfg.JWT.Secret, cfg.JWT.AccessTTL,
		cfg.JWT.Secret, cfg.JWT.RefreshTTL,
	)

	// ── Domain Services ────────────────────────────────────────────────────────
	userSvc := domainuser.NewUserService(userRepo)

	// ── Usecases ───────────────────────────────────────────────────────────────
	loginUc := useauth.NewLoginUsecase(userRepo, tokenSvc, hasher)
	registerUc := useauth.NewRegisterUsecase(userRepo, userSvc, hasher)
	refreshUc := useauth.NewRefreshUsecase(tokenSvc)
	profileUc := useruc.NewProfileUsecase(userRepo)
	createUc := useruc.NewCreateUsecase(userRepo, userSvc, hasher)
	listUc := useruc.NewListUsecase(userRepo)
	updateUc := useruc.NewUpdateUsecase(userRepo, hasher)
	deleteUc := useruc.NewDeleteUsecase(userRepo)

	// ── Handlers ───────────────────────────────────────────────────────────────
	authH := authhandler.NewAuthHandler(loginUc, registerUc, refreshUc, profileUc)
	userH := userhandler.NewUserHandler(createUc, listUc, updateUc, deleteUc, profileUc)

	// ── HTTP Server ────────────────────────────────────────────────────────────
	r := gin.Default()
	httpRoutes.RegisterRoutes(r, authH, userH, tokenSvc)

	log.Printf("🚀 %s running on :%s", cfg.App.Name, cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatal(err)
	}
}
