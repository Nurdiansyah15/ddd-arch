package main

import (
	"log"

	_ "github.com/Nurdiansyah15/ddd-arch/docs"

	"github.com/Nurdiansyah15/ddd-arch/internal/config"
	"github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/database/mysql"
	"github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/database/postgres"
	"github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http"
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

	// Init DB
	var db *sqlx.DB

	switch cfg.DB.Driver {
	case "mysql":
		db, err = mysql.NewMySQL(cfg.DB)
	case "postgres":
		db, err = postgres.NewPostgres(cfg.DB.DSN())
	default:
		log.Fatal("unsupported db driver")
	}

	// HTTP server
	r := gin.Default()
	http.RegisterRoutes(r, db)

	log.Printf("🚀 %s running on :%s", cfg.App.Name, cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatal(err)
	}

}
