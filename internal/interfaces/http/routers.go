package http

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/auth"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r *gin.Engine, db *sqlx.DB) {
	api := r.Group("/api")

	auth.SetupAuthRoutes(api, db)
}
