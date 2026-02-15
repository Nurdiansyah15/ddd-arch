package auth

import (
	"net/http"

	authuc "github.com/Nurdiansyah15/ddd-arch/internal/usecase/auth"
	useruc "github.com/Nurdiansyah15/ddd-arch/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	LoginUC    *authuc.LoginUsecase
	RegisterUC *authuc.RegisterUsecase
	RefreshUC  *authuc.RefreshUsecase
	ProfileUC  *useruc.ProfileUsecase
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.LoginUC.Execute(authuc.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.RegisterUC.Execute(authuc.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	uidRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, ok := uidRaw.(int64)
	if !ok {
		// sometimes numbers come as float64 from jwt
		if f, ok := uidRaw.(float64); ok {
			uid = int64(f)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
			return
		}
	}

	if h.ProfileUC == nil {
		c.JSON(http.StatusOK, gin.H{"id": uid})
		return
	}

	p, err := h.ProfileUC.Execute(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if h.RefreshUC == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh not configured"})
		return
	}

	resp, err := h.RefreshUC.Execute(authuc.RefreshRequest{RefreshToken: req.RefreshToken})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
