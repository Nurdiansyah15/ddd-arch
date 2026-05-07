package user

import (
	userhandler "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/handlers/user"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup, h *userhandler.UserHandler) {
	ur := r.Group("/users")
	ur.POST("/", h.Create)
	ur.GET("/", h.List)
	ur.GET(":id", h.Get)
	ur.PUT(":id", h.Update)
	ur.DELETE(":id", h.Delete)
}
