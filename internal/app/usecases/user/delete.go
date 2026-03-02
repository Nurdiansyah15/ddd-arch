package user

import "github.com/Nurdiansyah15/ddd-arch/internal/app/domain/master/user"

type DeleteUsecase struct {
	Repo user.Repository
}

func NewDeleteUsecase(repo user.Repository) *DeleteUsecase {
	return &DeleteUsecase{Repo: repo}
}

func (uc *DeleteUsecase) Execute(id int64) error {
	return uc.Repo.Delete(id)
}
