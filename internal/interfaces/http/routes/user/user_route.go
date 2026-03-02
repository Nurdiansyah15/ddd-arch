package user

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/app/domain/master/user"
	useruc "github.com/Nurdiansyah15/ddd-arch/internal/app/usecases/user"
	userrepo "github.com/Nurdiansyah15/ddd-arch/internal/infrastructure/persistence/user"
	userHandler "github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/handlers/user"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupUserRoutes(r *gin.RouterGroup, db *sqlx.DB) {
	userRepo := userrepo.NewUserRepositoryPG(db)
	userSvc := user.NewUserService(userRepo)

	createUc := useruc.NewCreateUsecase(userRepo, userSvc)
	listUc := useruc.NewListUsecase(userRepo)
	updateUc := useruc.NewUpdateUsecase(userRepo)
	deleteUc := useruc.NewDeleteUsecase(userRepo)

	h := userHandler.NewUserHandler(createUc, listUc, updateUc, deleteUc)

	ur := r.Group("/users")
	ur.POST("/", h.Create)
	ur.GET("/", h.List)
	ur.GET(":id", h.Get)
	ur.PUT(":id", h.Update)
	ur.DELETE(":id", h.Delete)
}
