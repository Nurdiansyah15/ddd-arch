package auth

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
	authuc "github.com/Nurdiansyah15/ddd-arch/internal/app/usecases/auth"
	useruc "github.com/Nurdiansyah15/ddd-arch/internal/app/usecases/user"
	"github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/handlers/respond"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	LoginUC    *authuc.LoginUsecase
	RegisterUC *authuc.RegisterUsecase
	RefreshUC  *authuc.RefreshUsecase
	ProfileUC  *useruc.ProfileUsecase
}

func NewAuthHandler(loginUC *authuc.LoginUsecase, registerUC *authuc.RegisterUsecase, refreshUC *authuc.RefreshUsecase, profileUC *useruc.ProfileUsecase) *AuthHandler {
	return &AuthHandler{LoginUC: loginUC, RegisterUC: registerUC, RefreshUC: refreshUC, ProfileUC: profileUC}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, apperror.Validation("invalid request body"))
		return
	}

	resp, err := h.LoginUC.Execute(authuc.LoginRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		respond.Error(c, err)
		return
	}

	respond.OK(c, "login successful", resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, apperror.Validation("invalid request body"))
		return
	}

	resp, err := h.RegisterUC.Execute(authuc.RegisterRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		respond.Error(c, err)
		return
	}

	respond.Created(c, "registration successful", resp)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	uidRaw, exists := c.Get("user_id")
	if !exists {
		respond.Error(c, apperror.Unauthorized("unauthorized"))
		return
	}

	uid, ok := uidRaw.(int64)
	if !ok {
		// jwt kadang menyimpan number sebagai float64
		if f, ok := uidRaw.(float64); ok {
			uid = int64(f)
		} else {
			respond.Error(c, apperror.Internal(nil))
			return
		}
	}

	p, err := h.ProfileUC.Execute(uid)
	if err != nil {
		respond.Error(c, err)
		return
	}

	respond.OK(c, "profile retrieved", p)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, apperror.Validation("invalid request body"))
		return
	}

	resp, err := h.RefreshUC.Execute(authuc.RefreshRequest{RefreshToken: req.RefreshToken})
	if err != nil {
		respond.Error(c, err)
		return
	}

	respond.OK(c, "token refreshed", resp)
}
