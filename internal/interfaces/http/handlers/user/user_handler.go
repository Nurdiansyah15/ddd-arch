package user

import (
	"fmt"

	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
	useruc "github.com/Nurdiansyah15/ddd-arch/internal/app/usecases/user"
	"github.com/Nurdiansyah15/ddd-arch/internal/interfaces/http/handlers/respond"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	CreateUC  *useruc.CreateUsecase
	ListUC    *useruc.ListUsecase
	UpdateUC  *useruc.UpdateUsecase
	DeleteUC  *useruc.DeleteUsecase
	ProfileUC *useruc.ProfileUsecase
}

func NewUserHandler(createUC *useruc.CreateUsecase, listUC *useruc.ListUsecase, updateUC *useruc.UpdateUsecase, deleteUC *useruc.DeleteUsecase, profileUC *useruc.ProfileUsecase) *UserHandler {
	return &UserHandler{CreateUC: createUC, ListUC: listUC, UpdateUC: updateUC, DeleteUC: deleteUC, ProfileUC: profileUC}
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
func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, apperror.Validation("invalid request body"))
		return
	}

	resp, err := h.CreateUC.Execute(useruc.CreateRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		respond.Error(c, err)
		return
	}
	respond.Created(c, "user created", resp)
}

// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} []useruc.ListResponseItem
// @Failure 500 {object} gin.H
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	resp, err := h.ListUC.Execute()
	if err != nil {
		respond.Error(c, err)
		return
	}
	respond.OK(c, "users retrieved", resp)
}

// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} useruc.ProfileResponse
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	var id int64
	if _, err := fmt.Sscan(c.Param("id"), &id); err != nil {
		respond.Error(c, apperror.Validation("invalid id"))
		return
	}

	resp, err := h.ProfileUC.Execute(id)
	if err != nil {
		respond.Error(c, err)
		return
	}
	respond.OK(c, "user retrieved", resp)
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
func (h *UserHandler) Update(c *gin.Context) {
	var id int64
	if _, err := fmt.Sscan(c.Param("id"), &id); err != nil {
		respond.Error(c, apperror.Validation("invalid id"))
		return
	}

	var req struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
		IsActive *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, apperror.Validation("invalid request body"))
		return
	}

	resp, err := h.UpdateUC.Execute(useruc.UpdateRequest{ID: id, Email: req.Email, Password: req.Password, IsActive: req.IsActive})
	if err != nil {
		respond.Error(c, err)
		return
	}
	respond.OK(c, "user updated", resp)
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
func (h *UserHandler) Delete(c *gin.Context) {
	var id int64
	if _, err := fmt.Sscan(c.Param("id"), &id); err != nil {
		respond.Error(c, apperror.Validation("invalid id"))
		return
	}

	if err := h.DeleteUC.Execute(id); err != nil {
		respond.Error(c, err)
		return
	}
	c.Status(204)
}
