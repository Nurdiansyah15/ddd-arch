package user

import userdomain "github.com/Nurdiansyah15/ddd-arch/internal/domain/user"

// ProfileUsecase depends on the domain Repository interface.
type ProfileUsecase struct {
	Repo userdomain.Repository
}

type ProfileResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func NewProfileUsecase(repo userdomain.Repository) *ProfileUsecase {
	return &ProfileUsecase{Repo: repo}
}

func (uc *ProfileUsecase) Execute(userID int64) (*ProfileResponse, error) {
	u, err := uc.Repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return &ProfileResponse{ID: u.ID, Email: u.Email}, nil
}
