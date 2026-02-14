package user

import userdomain "github.com/Nurdiansyah15/ddd-arch/internal/domain/user"

type DeleteUsecase struct {
	Repo userdomain.Repository
}

func NewDeleteUsecase(repo userdomain.Repository) *DeleteUsecase {
	return &DeleteUsecase{Repo: repo}
}

func (uc *DeleteUsecase) Execute(id int64) error {
	return uc.Repo.Delete(id)
}
