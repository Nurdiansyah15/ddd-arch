package user

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
	"github.com/Nurdiansyah15/ddd-arch/internal/domain/master/user"
)

type DeleteUsecase struct {
	Repo user.Repository
}

func NewDeleteUsecase(repo user.Repository) *DeleteUsecase {
	return &DeleteUsecase{Repo: repo}
}

func (uc *DeleteUsecase) Execute(id int64) error {
	if err := uc.Repo.Delete(id); err != nil {
		return apperror.Internal(err)
	}
	return nil
}
