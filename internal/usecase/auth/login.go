package auth

import (
	"errors"

	"github.com/Nurdiansyah15/ddd-arch/internal/domain/user"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginUsecase struct {
	UserRepo   user.Repository
	TokenMaker TokenGenerator
}

func NewLoginUsecase(repo user.Repository, gen TokenGenerator) *LoginUsecase {
	return &LoginUsecase{UserRepo: repo, TokenMaker: gen}
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (uc *LoginUsecase) Execute(req LoginRequest) (*LoginResponse, error) {
	u, err := uc.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := u.Authenticate(req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	access, err := uc.TokenMaker.GenerateAccess(u.ID)
	if err != nil {
		return nil, err
	}

	refresh, err := uc.TokenMaker.GenerateRefresh(u.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
