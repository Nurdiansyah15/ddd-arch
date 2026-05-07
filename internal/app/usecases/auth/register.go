package auth

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
	"github.com/Nurdiansyah15/ddd-arch/internal/domain/master/user"
)

type RegisterUsecase struct {
	UserRepo    user.Repository
	UserService *user.UserService
	Hasher      user.PasswordHasher
}

func NewRegisterUsecase(repo user.Repository, svc *user.UserService, hasher user.PasswordHasher) *RegisterUsecase {
	return &RegisterUsecase{UserRepo: repo, UserService: svc, Hasher: hasher}
}

type RegisterRequest struct {
	Email    string
	Password string
}

type RegisterResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func (uc *RegisterUsecase) Execute(req RegisterRequest) (*RegisterResponse, error) {
	if err := uc.UserService.CheckEmailAvailability(req.Email); err != nil {
		return nil, apperror.Conflict("email already registered")
	}

	hash, err := uc.Hasher.Hash(req.Password)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	u := &user.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		IsActive:     true,
	}

	if err := uc.UserRepo.Create(u); err != nil {
		return nil, apperror.Internal(err)
	}

	return &RegisterResponse{ID: u.ID, Email: u.Email}, nil
}
