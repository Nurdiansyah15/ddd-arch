package auth

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
	"github.com/Nurdiansyah15/ddd-arch/internal/domain/master/user"
)

type LoginUsecase struct {
	UserRepo   user.Repository
	TokenMaker TokenGenerator
	Hasher     user.PasswordHasher
}

func NewLoginUsecase(repo user.Repository, gen TokenGenerator, hasher user.PasswordHasher) *LoginUsecase {
	return &LoginUsecase{UserRepo: repo, TokenMaker: gen, Hasher: hasher}
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
		// jangan bocorkan apakah email ada atau tidak
		return nil, apperror.Unauthorized("invalid credentials")
	}

	if err := u.Authenticate(req.Password, uc.Hasher); err != nil {
		return nil, apperror.Unauthorized("invalid credentials")
	}

	access, err := uc.TokenMaker.GenerateAccess(u.ID)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	refresh, err := uc.TokenMaker.GenerateRefresh(u.ID)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return &LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
