package user

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
	"github.com/Nurdiansyah15/ddd-arch/internal/domain/master/user"
)

type CreateUsecase struct {
	Repo        user.Repository
	UserService *user.UserService
	Hasher      user.PasswordHasher
}

func NewCreateUsecase(repo user.Repository, svc *user.UserService, hasher user.PasswordHasher) *CreateUsecase {
	return &CreateUsecase{Repo: repo, UserService: svc, Hasher: hasher}
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
	if err := uc.UserService.CheckEmailAvailability(req.Email); err != nil {
		return nil, apperror.Conflict("email already registered")
	}

	hash, err := uc.Hasher.Hash(req.Password)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	u := &user.User{
		Email:        req.Email,
		PasswordHash: hash,
		IsActive:     true,
	}

	if err := uc.Repo.Create(u); err != nil {
		return nil, apperror.Internal(err)
	}

	return &CreateResponse{ID: u.ID, Email: u.Email}, nil
}
