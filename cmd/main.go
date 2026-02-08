package main

import (
	"log"

	"github.com/Nurdiansyah15/ddd-arch/internal/config"
	"github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/database"
	"github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

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
		db, err = database.NewMySQL(cfg.DB)
	case "postgres":
		db, err = database.NewPostgres(cfg.DB.DSN())
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
