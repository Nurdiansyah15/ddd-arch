package user

import (
	"fmt"
	"net/http"

	useruc "github.com/Nurdiansyah15/ddd-arch/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	CreateUC *useruc.CreateUsecase
	ListUC   *useruc.ListUsecase
	UpdateUC *useruc.UpdateUsecase
	DeleteUC *useruc.DeleteUsecase
}

func NewUserHandler(createUC *useruc.CreateUsecase, listUC *useruc.ListUsecase, updateUC *useruc.UpdateUsecase, deleteUC *useruc.DeleteUsecase) *userHandler {
	return &userHandler{CreateUC: createUC, ListUC: listUC, UpdateUC: updateUC, DeleteUC: deleteUC}
}

// @Summary Create a new user
// @Description Create a new user with the given email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body useruc.CreateRequest true "Create user request"
// @Success 201 {object} useruc.CreateResponse
// @Failure 400 {object} gin.H
// @Router /api/v1/users [post]
func (h *userHandler) Create(c *gin.Context) {
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

// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} []useruc.ListResponseItem
// @Failure 500 {object} gin.H
// @Router /api/v1/users [get]
func (h *userHandler) List(c *gin.Context) {
	resp, err := h.ListUC.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/v1/users/{id} [get]
func (h *userHandler) Get(c *gin.Context) {
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

// @Summary Update a user by ID
// @Description Update a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body useruc.UpdateRequest true "Update user request"
// @Success 200 {object} useruc.UpdateResponse
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/v1/users/{id} [put]
func (h *userHandler) Update(c *gin.Context) {
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

// @Summary Delete a user by ID
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {object} nil
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/v1/users/{id} [delete]
func (h *userHandler) Delete(c *gin.Context) {
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
