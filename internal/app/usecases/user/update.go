package user

import (
	"errors"

	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
	domainuser "github.com/Nurdiansyah15/ddd-arch/internal/domain/master/user"
)

type UpdateUsecase struct {
	Repo   domainuser.Repository
	Hasher domainuser.PasswordHasher
}

func NewUpdateUsecase(repo domainuser.Repository, hasher domainuser.PasswordHasher) *UpdateUsecase {
	return &UpdateUsecase{Repo: repo, Hasher: hasher}
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

func (uc *UpdateUsecase) Execute(req UpdateRequest) (*UpdateResponse, error) {
	u, err := uc.Repo.FindByID(req.ID)
	if err != nil {
		if errors.Is(err, domainuser.ErrUserNotFound) {
			return nil, apperror.NotFound("user not found", err)
		}
		return nil, apperror.Internal(err)
	}

	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Password != nil {
		h, err := uc.Hasher.Hash(*req.Password)
		if err != nil {
			return nil, apperror.Internal(err)
		}
		u.PasswordHash = h
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}

	if err := uc.Repo.Update(u); err != nil {
		return nil, apperror.Internal(err)
	}

	return &UpdateResponse{ID: u.ID, Email: u.Email}, nil
}
