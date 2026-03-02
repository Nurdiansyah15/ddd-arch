package user

import "github.com/Nurdiansyah15/ddd-arch/internal/app/domain/master/user"

// ProfileUsecase depends on the domain Repository interface.
type ProfileUsecase struct {
	Repo user.Repository
}

type ProfileResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func NewProfileUsecase(repo user.Repository) *ProfileUsecase {
	return &ProfileUsecase{Repo: repo}
}

func (uc *ProfileUsecase) Execute(userID int64) (*ProfileResponse, error) {
	u, err := uc.Repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return &ProfileResponse{ID: u.ID, Email: u.Email}, nil
}
