package auth

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUsecase struct {
	UserRepo    user.Repository
	UserService *user.UserService // domain service
}

func NewRegisterUsecase(repo user.Repository, svc *user.UserService) *RegisterUsecase {
	return &RegisterUsecase{UserRepo: repo, UserService: svc}
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
	// Domain Service: cek apakah email sudah dipakai
	if err := uc.UserService.CheckEmailAvailability(req.Email); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &user.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		IsActive:     true,
	}

	if err := uc.UserRepo.Create(user); err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
