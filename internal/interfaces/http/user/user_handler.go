package user

import (
	"fmt"
	"net/http"

	useruc "github.com/Nurdiansyah15/ddd-arch/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	CreateUC *useruc.CreateUsecase
	ListUC   *useruc.ListUsecase
	UpdateUC *useruc.UpdateUsecase
	DeleteUC *useruc.DeleteUsecase
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.CreateUC.Execute(useruc.CreateRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) List(c *gin.Context) {
	resp, err := h.ListUC.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	var id int64
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	// reuse list/find usecase via repo inside usecase; direct call to repo is avoided here
	// we can call ProfileUsecase if present but to keep code self-contained, use Repo via CreateUC
	if h.ListUC == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not configured"})
		return
	}
	// find by id using the repo behind createUC
	// access repo from CreateUC
	repo := h.CreateUC.Repo
	u, err := repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email})
}

func (h *UserHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	var id int64
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
		IsActive *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	resp, err := h.UpdateUC.Execute(useruc.UpdateRequest{ID: id, Email: req.Email, Password: req.Password, IsActive: req.IsActive})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	var id int64
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.DeleteUC.Execute(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
