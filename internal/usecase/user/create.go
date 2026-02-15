package user

import (
	userdomain "github.com/Nurdiansyah15/ddd-arch/internal/domain/user"
)

type CreateUsecase struct {
	Repo        userdomain.Repository
	UserService *userdomain.UserService // domain service
}

func NewCreateUsecase(repo userdomain.Repository, svc *userdomain.UserService) *CreateUsecase {
	return &CreateUsecase{Repo: repo, UserService: svc}
}

type CreateRequest struct {
	Email    string
	Password string
}

type CreateResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func (uc *CreateUsecase) Execute(req CreateRequest) (*CreateResponse, error) {
	// Domain Service: cek apakah email sudah dipakai
	if err := uc.UserService.CheckEmailAvailability(req.Email); err != nil {
		return nil, err
	}

	hash, err := userdomain.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &userdomain.User{
		Email:        req.Email,
		PasswordHash: hash,
		IsActive:     true,
	}

	if err := uc.Repo.Create(u); err != nil {
		return nil, err
	}

	return &CreateResponse{ID: u.ID, Email: u.Email}, nil
}
