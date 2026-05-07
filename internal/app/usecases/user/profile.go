package user

import (
	"errors"

	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
	domainuser "github.com/Nurdiansyah15/ddd-arch/internal/domain/master/user"
)

type ProfileUsecase struct {
	Repo domainuser.Repository
}

type ProfileResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func NewProfileUsecase(repo domainuser.Repository) *ProfileUsecase {
	return &ProfileUsecase{Repo: repo}
}

func (uc *ProfileUsecase) Execute(userID int64) (*ProfileResponse, error) {
	u, err := uc.Repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, domainuser.ErrUserNotFound) {
			return nil, apperror.NotFound("user not found", err)
		}
		return nil, apperror.Internal(err)
	}
	return &ProfileResponse{ID: u.ID, Email: u.Email}, nil
}
