package user

import (
	userdomain "github.com/Nurdiansyah15/ddd-arch/internal/domain/user"
	userrepo "github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/persistence/user"
	useruc "github.com/Nurdiansyah15/ddd-arch/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupUserRoutes(r *gin.RouterGroup, db *sqlx.DB) {
	userRepo := userrepo.NewUserRepositoryPG(db)
	userSvc := userdomain.NewUserService(userRepo)

	createUc := useruc.NewCreateUsecase(userRepo, userSvc)
	listUc := useruc.NewListUsecase(userRepo)
	updateUc := useruc.NewUpdateUsecase(userRepo)
	deleteUc := useruc.NewDeleteUsecase(userRepo)

	h := NewUserHandler(createUc, listUc, updateUc, deleteUc)

	ur := r.Group("/users")
	ur.POST("/", h.Create)
	ur.GET("/", h.List)
	ur.GET(":id", h.Get)
	ur.PUT(":id", h.Update)
	ur.DELETE(":id", h.Delete)
}
