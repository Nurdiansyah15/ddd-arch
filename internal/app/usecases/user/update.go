package user

import (
	"errors"

	"github.com/Nurdiansyah15/ddd-arch/internal/app/domain/master/user"
)

type UpdateUsecase struct {
	Repo user.Repository
}

func NewUpdateUsecase(repo user.Repository) *UpdateUsecase {
	return &UpdateUsecase{Repo: repo}
}

type UpdateRequest struct {
	ID       int64
	Email    *string
	Password *string
	IsActive *bool
}

type UpdateResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

var ErrNotFound = errors.New("user not found")

func (uc *UpdateUsecase) Execute(req UpdateRequest) (*UpdateResponse, error) {
	u, err := uc.Repo.FindByID(req.ID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}

	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Password != nil {
		h, err := user.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = h
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}

	if err := uc.Repo.Update(u); err != nil {
		return nil, err
	}

	return &UpdateResponse{ID: u.ID, Email: u.Email}, nil
}
